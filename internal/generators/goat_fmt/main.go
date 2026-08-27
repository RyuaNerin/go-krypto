package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

type stringSlice []string

func (s *stringSlice) String() string { return strings.Join(*s, ",") }
func (s *stringSlice) Set(v string) error {
	*s = append(*s, v)
	return nil
}

type symbolMapping struct {
	from string
	to   string
}

// cflags is the fixed, predefined list of extra clang options always
// passed to goat via -e.
var cflags = []string{
	// goat v0.2.1 does not reset its assembly parser from a data section on
	// Clang's shorthand .text directive. At -O1 and above, constant pools can
	// precede the first function. Emitting explicit .section .text.<name>
	// directives keeps goat's section tracking in sync without disabling
	// optimization.
	"-ffunction-sections",
	"-fno-asynchronous-unwind-tables",
	"-fno-exceptions",
	"-fno-rtti",
	"-mno-red-zone",
	"-fno-builtin",
	"-mllvm",
	"-inline-threshold=1000",
	"-fno-stack-protector",
	"-fno-jump-tables",
	"-fomit-frame-pointer",
	"-funroll-loops",
}

// cflagsArm64 are extra clang options for arm64. Go reserves R18 (platform
// register) and R28 (the current goroutine, g); generated code must not
// allocate them.
var cflagsArm64 = []string{
	"-ffixed-x18",
	"-ffixed-x28",
}

var (
	packageRe             = regexp.MustCompile(`(?m)^package\s+\S+`)
	localAssemblySymbolRe = regexp.MustCompile(`([[:alpha:]_][[:alnum:]_]*)<>(?:[+-](?:0[xX][[:xdigit:]]+|[[:digit:]]+))?\(SB\)`)
	dataDefinitionRe      = regexp.MustCompile(`(?m)^[ \t]*(?:DATA|GLOBL)[ \t]+([[:alpha:]_][[:alnum:]_]*)<>(?:[+-](?:0[xX][[:xdigit:]]+|[[:digit:]]+))?\(SB\)[^\r\n]*(?:\r?\n(?:\r?\n)?|$)`)
	noasmBuildTagRe       = regexp.MustCompile(`(^|[[:space:](])!noasm([[:space:])]|$)`)

	// brokenRipLeaRe matches a raw-byte-encoded LEA instruction goat emits
	// for a RIP-relative reference at a nonzero offset from a symbol (e.g.
	// `leaq StepConstants+64(%rip), %rcx`). goat v0.2.1 transcribes the
	// literal bytes it disassembled from the unlinked .o file instead of
	// emitting a real Go instruction, so the encoded displacement is a
	// relocation placeholder (always 0), not the actual offset shown in
	// its own comment. Every such instruction still carries an accurate
	// descriptive comment, which fixBrokenRipLeaqInstructions uses to
	// reconstruct a working replacement.
	brokenRipLeaRe = regexp.MustCompile(`(?m)^([ \t]*)(?:(?:LONG|WORD|BYTE|QUAD)\s+\$0x[0-9a-fA-F]+;\s*)*(?:LONG|WORD|BYTE|QUAD)\s+\$0x[0-9a-fA-F]+\s*//\s*leaq\s+([[:alpha:]_][[:alnum:]_]*)(?:\+([0-9]+))?\(%rip\),\s*%(\w+)\s*$`)

	// amd64AttRegisters maps AT&T 64-bit general-purpose register names
	// (as they appear in goat's disassembly comments) to their Go asm
	// mnemonics.
	amd64AttRegisters = map[string]string{
		"rax": "AX", "rbx": "BX", "rcx": "CX", "rdx": "DX",
		"rsi": "SI", "rdi": "DI", "rbp": "BP", "rsp": "SP",
		"r8": "R8", "r9": "R9", "r10": "R10", "r11": "R11",
		"r12": "R12", "r13": "R13", "r14": "R14", "r15": "R15",
	}
)

// addNoasmBuildTags extends goat's generic !noasm build constraint with this
// project's global and package-specific assembly opt-out tags. Only go:build
// directives are changed so occurrences in generated comments remain intact.
func addNoasmBuildTags(data []byte, packageName string) []byte {
	const buildDirective = "//go:build"
	replacement := "!noasm && !noasm_krypto && !noasm_krypto_" + packageName

	lines := strings.SplitAfter(string(data), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, buildDirective) {
			lines[i] = noasmBuildTagRe.ReplaceAllStringFunc(line, func(match string) string {
				before, after, _ := strings.Cut(match, "!noasm")
				return before + replacement + after
			})
		}
	}
	return []byte(strings.Join(lines, ""))
}

// fixBrokenRipLeaqInstructions rewrites goat's broken raw-byte encoding of
// RIP-relative LEA-with-offset instructions (see brokenRipLeaRe) into real
// Go assembly using the `symbol<>+offset(SB)` local-symbol form. That form
// is understood by the existing collectAssemblySymbols/rewriteAssemblySymbols
// pipeline, which resolves it to the mapped Go package symbol exactly like
// any other local assembly symbol goat translates correctly.
func fixBrokenRipLeaqInstructions(data []byte) ([]byte, error) {
	var firstErr error
	fixed := brokenRipLeaRe.ReplaceAllFunc(data, func(m []byte) []byte {
		match := brokenRipLeaRe.FindSubmatch(m)
		indent, symbol, offset, attReg := match[1], match[2], match[3], match[4]

		goReg, ok := amd64AttRegisters[strings.ToLower(string(attReg))]
		if !ok {
			if firstErr == nil {
				firstErr = fmt.Errorf("fixBrokenRipLeaqInstructions: unrecognized register %%%s in %q", attReg, string(m))
			}
			return m
		}

		ref := string(symbol) + "<>"
		if len(offset) > 0 {
			ref += "+" + string(offset)
		}
		ref += "(SB)"

		return fmt.Appendf(nil, "%sLEAQ %s, %s", indent, ref, goReg)
	})
	if firstErr != nil {
		return nil, firstErr
	}
	return fixed, nil
}

// toolModDir is the directory containing the go.mod that declares goat and
// goimports as tool dependencies (via the `tool` directive). `go tool <name>`
// resolves against the module found from the current working directory, so
// commands must be run with this directory as their working directory. This
// lets us invoke tools that were `go get -tool`'d into go.mod without
// requiring them to be present on PATH.
func toolModDir() (string, error) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("could not determine source file location")
	}
	return filepath.Dir(thisFile), nil
}

// runGoTool runs `go tool <name> args...` with its working directory set to
// the module that declares the tool dependency, so it resolves correctly
// even if the tool binary isn't on PATH.
func runGoTool(name string, args ...string) error {
	dir, err := toolModDir()
	if err != nil {
		return err
	}
	fullArgs := append([]string{"tool", name}, args...)
	cmd := exec.Command("go", fullArgs...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	src := flag.String("src", "", "path to the input C source file")
	outdir := flag.String("outdir", "", "output directory to copy generated files into (default: current directory)")
	namespace := flag.String("namespace", "", "package name to set in the generated .go file")
	filename := flag.String("filename", "", "base filename (without extension) for the output .s/.go files; defaults to the source file's base name")
	platform := flag.String("platform", "", "target GOARCH for cross-compilation, passed to goat -t (e.g. arm64, amd64)")
	var machine stringSlice
	var symbols stringSlice
	flag.Var(&machine, "machine", "goat machine-option (-m) value, repeatable (e.g. --machine avx2)")
	flag.Var(&symbols, "symbol", "rewrite generated local assembly symbols to Go package variables or constants (from=to), repeatable or comma-separated")
	flag.Parse()

	if *src == "" {
		fmt.Fprintln(os.Stderr, "missing required flag: --src")
		os.Exit(1)
	}
	if *namespace == "" {
		fmt.Fprintln(os.Stderr, "missing required flag: --namespace")
		os.Exit(1)
	}

	resolvedOutdir := *outdir
	if resolvedOutdir == "" {
		resolvedOutdir = "."
		log.Printf("--outdir not specified, defaulting to current directory")
	}

	if err := run(*src, resolvedOutdir, *namespace, *filename, *platform, machine, symbols); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(src, outdir, namespace, filename, platform string, machine, symbols []string) error {
	mappings, err := parseSymbolMappings(symbols)
	if err != nil {
		return err
	}

	log.Printf("checking for clang and go binaries")
	if _, err := exec.LookPath("clang"); err != nil {
		return fmt.Errorf("clang not found: %w", err)
	}
	if _, err := exec.LookPath("go"); err != nil {
		return fmt.Errorf("go not found: %w", err)
	}
	log.Printf("clang and go found")

	// resolve src to an absolute path since goat will run with a
	// different working directory (the tool module's directory).
	absSrc, err := filepath.Abs(src)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path for %s: %w", src, err)
	}
	absOutdir, err := filepath.Abs(outdir)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path for %s: %w", outdir, err)
	}

	tmpDir, err := os.MkdirTemp("", "goat-tmp-")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)
	log.Printf("created temp directory: %s", tmpDir)

	// goat src -o tmpDir -O3 -v -m <machine> -t <platform> -e <cflag> ...
	args := []string{absSrc, "-o", tmpDir, "-O3", "-v"}
	for _, m := range machine {
		args = append(args, "-m", m)
	}
	if platform != "" {
		args = append(args, "-t", platform)
	}
	for _, f := range cflags {
		args = append(args, "-e", f)
	}
	if platform == "arm64" {
		for _, f := range cflagsArm64 {
			args = append(args, "-e", f)
		}
	}

	log.Printf("running: go tool goat %s", strings.Join(args, " "))
	if err := runGoTool("goat", args...); err != nil {
		return fmt.Errorf("goat execution failed: %w", err)
	}
	log.Printf("goat finished successfully")

	if err := os.MkdirAll(outdir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	log.Printf("output directory ready: %s", outdir)

	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		return fmt.Errorf("failed to read temp directory: %w", err)
	}
	log.Printf("found %d entries in temp directory", len(entries))

	type generatedFile struct {
		name string
		ext  string
		data []byte
	}
	// goat embeds the absolute source path in comments (e.g.
	// "// source: /abs/path/to/src.c"). Compute the relative path once
	// so the generated files don't leak the local filesystem layout.
	relSrc := ""
	if rel, relErr := filepath.Rel(absOutdir, absSrc); relErr == nil {
		relSrc = rel
	} else {
		log.Printf("could not compute relative path for %s: %v", absSrc, relErr)
	}

	var generated []generatedFile
	assemblySymbols := make(map[string]bool)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		ext := filepath.Ext(name)
		if ext != ".s" && ext != ".go" {
			log.Printf("skipping %s (not .s or .go)", name)
			continue
		}

		data, err := os.ReadFile(filepath.Join(tmpDir, name))
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", name, err)
		}

		// Rewrite any occurrence of the absolute source path to the
		// precomputed relative path so the generated file doesn't leak
		// the local filesystem layout.
		if relSrc != "" {
			data = []byte(strings.ReplaceAll(string(data), absSrc, relSrc))
		}

		data = addNoasmBuildTags(data, namespace)

		// Rewrite the package declaration in generated .go files to
		// use the requested namespace.
		if ext == ".go" {
			log.Printf("rewriting package declaration in %s to %q", name, namespace)
			data = packageRe.ReplaceAll(data, []byte("package "+namespace))
		} else {
			fixed, fixErr := fixBrokenRipLeaqInstructions(data)
			if fixErr != nil {
				return fmt.Errorf("failed to fix broken RIP-relative LEA instructions in %s: %w", name, fixErr)
			}
			data = fixed

			for symbol := range collectAssemblySymbols(data) {
				assemblySymbols[symbol] = true
			}
			data = rewriteAssemblySymbols(data, mappings)
		}
		generated = append(generated, generatedFile{name: name, ext: ext, data: data})
	}

	if err := validateSymbolsUsed(mappings, assemblySymbols); err != nil {
		return err
	}
	if err := validateSymbolTargets(absOutdir, namespace, renamedAssemblySymbols(assemblySymbols, mappings)); err != nil {
		return err
	}

	copied := 0
	for _, file := range generated {
		// Optionally rename the output file (keeping its extension).
		outName := file.name
		if filename != "" {
			outName = filename + file.ext
		}

		dst := filepath.Join(outdir, outName)
		if err := os.WriteFile(dst, file.data, 0o644); err != nil {
			return fmt.Errorf("failed to write %s: %w", dst, err)
		}
		log.Printf("copied %s -> %s", file.name, dst)
		copied++

		if file.ext == ".go" {
			if err := formatGoFile(dst); err != nil {
				return fmt.Errorf("failed to format %s: %w", dst, err)
			}
		}
	}

	log.Printf("done: %d file(s) copied to %s", copied, outdir)
	return nil
}

func parseSymbolMappings(values []string) ([]symbolMapping, error) {
	var mappings []symbolMapping
	for _, value := range values {
		for _, item := range strings.Split(value, ",") {
			from, to, ok := strings.Cut(item, "=")
			from = strings.TrimSpace(from)
			to = strings.TrimSpace(to)
			if !ok || from == "" || to == "" {
				return nil, fmt.Errorf("invalid symbol mapping %q (expected from=to)", item)
			}
			mappings = append(mappings, symbolMapping{from: from, to: to})
		}
	}
	return mappings, nil
}

// validateSymbolTargets ensures every local assembly symbol exists in the
// output package after applying --symbol mappings. Test files are excluded:
// their declarations are unavailable to normal package assembly builds.
func validateSymbolTargets(dir, packageName string, symbols map[string]bool) error {
	if len(symbols) == 0 {
		return nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read output package directory %s: %w", dir, err)
	}

	goSymbols := make(map[string]bool)
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".go" || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			return fmt.Errorf("failed to parse Go file %s: %w", path, err)
		}
		if file.Name.Name != packageName {
			continue
		}
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || (genDecl.Tok != token.VAR && genDecl.Tok != token.CONST) {
				continue
			}
			for _, spec := range genDecl.Specs {
				valueSpec, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for _, name := range valueSpec.Names {
					goSymbols[name.Name] = true
				}
			}
		}
	}

	var missing []string
	for symbol := range symbols {
		if !goSymbols[symbol] {
			missing = append(missing, symbol)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return fmt.Errorf("assembly symbol(s) not found as package-level variable(s) or constant(s) in package %q: %s", packageName, strings.Join(missing, ", "))
	}
	return nil
}

// collectAssemblySymbols returns every file-local symbol referenced from a
// generated assembly file. Comments are ignored so a commented-out reference
// cannot satisfy strict symbol validation.
func collectAssemblySymbols(data []byte) map[string]bool {
	symbols := make(map[string]bool)
	for _, line := range strings.Split(string(data), "\n") {
		if commentStart := strings.Index(line, "//"); commentStart >= 0 {
			line = line[:commentStart]
		}
		for _, match := range localAssemblySymbolRe.FindAllStringSubmatch(line, -1) {
			symbols[match[1]] = true
		}
	}
	return symbols
}

// sourceToTarget reduces the symbol mappings to a single from→to lookup,
// keeping the first mapping for any duplicate source (first-mapping-wins).
func sourceToTarget(mappings []symbolMapping) map[string]string {
	bySource := make(map[string]string, len(mappings))
	for _, mapping := range mappings {
		if _, exists := bySource[mapping.from]; !exists {
			bySource[mapping.from] = mapping.to
		}
	}
	return bySource
}

// rewriteAssemblySymbols redirects goat's file-local C symbols to Go package
// variables or constants according to the supplied mappings. goat v0.2.2
// emits DATA/GLOBL definitions for C data symbols; mapped symbols already have
// storage in Go, so discard those definitions instead of duplicating the data
// in every generated assembly implementation.
func rewriteAssemblySymbols(data []byte, mappings []symbolMapping) []byte {
	bySource := sourceToTarget(mappings)
	data = dataDefinitionRe.ReplaceAllFunc(data, func(line []byte) []byte {
		match := dataDefinitionRe.FindSubmatch(line)
		if _, mapped := bySource[string(match[1])]; mapped {
			return nil
		}
		return line
	})
	return localAssemblySymbolRe.ReplaceAllFunc(data, func(reference []byte) []byte {
		match := localAssemblySymbolRe.FindSubmatch(reference)
		to, ok := bySource[string(match[1])]
		if !ok {
			return reference
		}
		return []byte(strings.Replace(string(reference), string(match[1])+"<>", "·"+to, 1))
	})
}

// renamedAssemblySymbols maps every referenced assembly symbol to the Go name
// its mapping renames it to, so the targets can be validated against the
// output package. Symbols without a mapping keep their own name.
func renamedAssemblySymbols(symbols map[string]bool, mappings []symbolMapping) map[string]bool {
	bySource := sourceToTarget(mappings)
	renamed := make(map[string]bool, len(symbols))
	for symbol := range symbols {
		// Mapped symbols are renamed to their target; unmapped ones keep
		// their own name.
		target := symbol
		if to, ok := bySource[symbol]; ok {
			target = to
		}
		renamed[target] = true
	}
	return renamed
}

func validateSymbolsUsed(mappings []symbolMapping, assemblySymbols map[string]bool) error {
	var unused []string
	for _, mapping := range mappings {
		if !assemblySymbols[mapping.from] {
			unused = append(unused, mapping.from)
		}
	}
	if len(unused) > 0 {
		return fmt.Errorf("symbol mapping source(s) not found in generated assembly: %s", strings.Join(unused, ", "))
	}
	return nil
}

// formatGoFile formats a .go file in place using goimports resolved via
// `go tool goimports` (declared as a tool dependency in go.mod), which
// also fixes up imports. Falls back to `go fmt` if that fails.
func formatGoFile(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path for %s: %w", path, err)
	}

	log.Printf("running go tool goimports on %s", absPath)
	if err := runGoTool("goimports", "-w", absPath); err == nil {
		return nil
	}
	log.Printf("go tool goimports failed (%v), falling back to go fmt", err)

	log.Printf("running go fmt on %s", absPath)
	cmd := exec.Command("go", "fmt", absPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
