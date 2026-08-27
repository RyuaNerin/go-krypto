// Package cavp reads and processes CAVP/KAT style request-response files.
//
// File format:
//
//	COUNT = 0
//	K = 95008AC1...
//	M =
//	Tlen = 64
//	T = DB8A1B27...
//	INVALID
//
//	COUNT = 1
//	...
//
// Records are separated by blank lines. "[Key = Value]" header lines set
// defaults that apply to all following records. Lines starting with '#'
// are comments.
package cavp

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

type convKind int

const (
	convAuto convKind = iota
	convHex
	convInt
	convBool
	convStr
)

type headerKind int

const (
	headerBracket headerKind = iota
	headerPlain
	headerKCDSASVT
	headerECKCDSASVT
)

type fieldInfo struct {
	idx        int
	typ        reflect.Type
	key        string
	conv       convKind
	optional   bool
	header     bool // value must come from a "[Key = Value]" line
	headerKind headerKind
	skip       bool // blank field "_": parsed and validated, then discarded
}

type Reader[T any] struct {
	sc       *bufio.Scanner
	fields   []fieldInfo
	byKey    map[string]*fieldInfo
	defaults map[string]kv
	pending  *string
	line     int
	err      error
}

type kv struct {
	key string
	val string
	has bool // true if the line had '=' (flag lines have has=false)
	hdr bool // came from a "[Key = Value]" line
}

// NewReader creates a Reader for struct type T. Struct fields are mapped
// with the tag `cavp:"Name[,hex|int|bool|string][,optional]"`. The hex
// conversion accepts both byte slices and integer fields.
func NewReader[T any](r io.Reader) *Reader[T] {
	rd := &Reader[T]{
		sc:       bufio.NewScanner(r),
		defaults: make(map[string]kv),
	}
	rd.sc.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	rd.err = rd.parseType()
	return rd
}

func (r *Reader[T]) parseType() error {
	fields, byKey, err := parseFields[T]()
	if err != nil {
		return err
	}
	r.fields, r.byKey = fields, byKey
	return nil
}

func parseFields[T any]() ([]fieldInfo, map[string]*fieldInfo, error) {
	var fields []fieldInfo
	var zero T
	t := reflect.TypeOf(zero)
	if t == nil || t.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("rsp: type parameter must be a struct")
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag, ok := f.Tag.Lookup("cavp")
		if !ok || tag == "-" {
			continue
		}
		blank := f.Name == "_"
		if f.PkgPath != "" && !blank {
			continue // unexported
		}
		parts := strings.Split(tag, ",")
		fi := fieldInfo{idx: i, typ: f.Type, key: strings.TrimSpace(parts[0]), skip: blank}
		if fi.key == "" {
			fi.key = f.Name
		}
		for _, opt := range parts[1:] {
			switch strings.ToLower(strings.TrimSpace(opt)) {
			case "hex":
				fi.conv = convHex
			case "int":
				fi.conv = convInt
			case "bool":
				fi.conv = convBool
			case "string", "str":
				fi.conv = convStr
			case "optional", "omitempty":
				fi.optional = true
			case "header", "hdr":
				fi.header = true
			case "header=plain":
				fi.header = true
				fi.headerKind = headerPlain
			case "header=kcdsa-svt":
				fi.header = true
				fi.headerKind = headerKCDSASVT
			case "header=eckcdsa-svt":
				fi.header = true
				fi.headerKind = headerECKCDSASVT
			case "":
			default:
				return nil, nil, fmt.Errorf("rsp: unknown tag option %q on field %s", opt, f.Name)
			}
		}
		fields = append(fields, fi)
	}
	byKey := make(map[string]*fieldInfo, len(fields))
	for i := range fields {
		k := strings.ToLower(fields[i].key)
		if _, dup := byKey[k]; dup {
			return nil, nil, fmt.Errorf("rsp: duplicate key %q", fields[i].key)
		}
		byKey[k] = &fields[i]
	}
	return fields, byKey, nil
}

// Read returns the next record. It returns io.EOF when the input is exhausted.
func (r *Reader[T]) Read() (T, error) {
	var out T
	if r.err != nil {
		return out, r.err
	}
	rec, err := r.nextRecord()
	if err != nil {
		return out, err
	}

	v := reflect.ValueOf(&out).Elem()
	seen := make(map[string]bool, len(rec))
	occurrences := make(map[string]int, len(rec))
	entries := make([]kv, 0, len(r.defaults)+len(rec))
	for _, e := range r.defaults {
		entries = append(entries, e)
	}
	entries = append(entries, rec...)
	for _, e := range entries {
		fi, ok := r.byKey[strings.ToLower(e.key)]
		if !ok {
			continue
		}
		if fi.header && !e.hdr {
			return out, fmt.Errorf("rsp: line %d: %q must be given as a [%s = ...] header", r.line, e.key, fi.key)
		}
		dst := reflect.New(fi.typ).Elem()
		if !fi.skip {
			dst = v.Field(fi.idx)
		}
		key := strings.ToLower(e.key)
		if err := setField(dst, fi, e, occurrences[key]); err != nil {
			return out, fmt.Errorf("rsp: line %d: %s: %w", r.line, e.key, err)
		}
		occurrences[key]++
		seen[key] = true
	}
	for i := range r.fields {
		fi := &r.fields[i]
		if fi.optional || seen[strings.ToLower(fi.key)] {
			continue
		}
		if fi.header {
			return out, fmt.Errorf("rsp: line %d: missing [%s = ...] header", r.line, fi.key)
		}
		return out, fmt.Errorf("rsp: line %d: missing required field %q", r.line, fi.key)
	}
	return out, nil
}

// ReadAll reads every remaining record.
func (r *Reader[T]) ReadAll() ([]T, error) {
	var all []T
	for {
		rec, err := r.Read()
		if errors.Is(err, io.EOF) {
			return all, nil
		}
		if err != nil {
			return all, err
		}
		all = append(all, rec)
	}
}

func (r *Reader[T]) nextRecord() ([]kv, error) {
	var rec []kv
	for {
		line, ok := r.nextLine()
		if !ok {
			break
		}

		if line == "" {
			if len(rec) > 0 {
				return rec, nil
			}
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		if entries, ok := r.parseSpecialHeader(line); ok {
			if len(rec) > 0 {
				r.pending = &line // header belongs to the next record
				return rec, nil
			}
			for _, e := range entries {
				r.defaults[strings.ToLower(e.key)] = e
			}
			continue
		}
		// "[Key = Value]" header: stays in effect until the same key
		// is declared again.
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			if len(rec) > 0 {
				r.pending = &line // header belongs to the next record
				return rec, nil
			}
			e := parseLine(strings.TrimSpace(line[1 : len(line)-1]))
			e.hdr = true
			r.defaults[strings.ToLower(e.key)] = e
			continue
		}
		rec = append(rec, parseLine(line))
	}
	if err := r.sc.Err(); err != nil {
		return nil, err
	}
	if len(rec) > 0 {
		return rec, nil
	}
	return nil, io.EOF
}

// parseSpecialHeader handles the header layouts used by the Korean KCDSA
// validation files. Unlike the common CAVP form, KCDSA headers are bare
// assignments and EC-KCDSA headers contain positional values.
func (r *Reader[T]) parseSpecialHeader(line string) ([]kv, bool) {
	if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
		inner := strings.TrimSpace(line[1 : len(line)-1])
		if strings.Contains(inner, "=") {
			return nil, false
		}

		var fields []*fieldInfo
		for i := range r.fields {
			fi := &r.fields[i]
			if fi.header && fi.headerKind == headerECKCDSASVT {
				fields = append(fields, fi)
			}
		}
		values := splitCommaFields(inner)
		if len(fields) == 0 || len(values) != len(fields) {
			return nil, false
		}
		entries := make([]kv, len(fields))
		for i, fi := range fields {
			entries[i] = kv{key: fi.key, val: values[i], has: true, hdr: true}
		}
		return entries, true
	}

	parts := splitCommaFields(line)
	entries := make([]kv, 0, len(parts))
	for _, part := range parts {
		e := parseLine(part)
		if !e.has {
			return nil, false
		}
		fi, ok := r.byKey[strings.ToLower(e.key)]
		if !ok || !fi.header || (fi.headerKind != headerPlain && fi.headerKind != headerKCDSASVT) {
			return nil, false
		}
		e.hdr = true
		entries = append(entries, e)
	}
	return entries, len(entries) > 0
}

func splitCommaFields(s string) []string {
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func (r *Reader[T]) nextLine() (string, bool) {
	if r.pending != nil {
		line := *r.pending
		r.pending = nil
		return line, true
	}
	if !r.sc.Scan() {
		return "", false
	}
	r.line++
	return strings.TrimSpace(strings.TrimRight(r.sc.Text(), "\r")), true
}

// Headers returns the "[Key = Value]" settings currently in effect.
func (r *Reader[T]) Headers() map[string]string {
	m := make(map[string]string, len(r.defaults))
	for _, e := range r.defaults {
		m[e.key] = e.val
	}
	return m
}

func parseLine(line string) kv {
	if before, after, ok := strings.Cut(line, "="); ok {
		return kv{
			key: strings.TrimSpace(before),
			val: strings.TrimSpace(after),
			has: true,
		}
	}
	return kv{key: line}
}

func setField(dst reflect.Value, fi *fieldInfo, e kv, occurrence int) error { //nolint:gocyclo
	conv := fi.conv
	if conv == convAuto {
		switch {
		case dst.Kind() == reflect.Slice && dst.Type().Elem().Kind() == reflect.Uint8:
			conv = convHex
		case isByteSliceCollection(dst.Type()):
			conv = convHex
		case dst.Kind() == reflect.Bool:
			conv = convBool
		case dst.Kind() == reflect.String:
			conv = convStr
		default:
			conv = convInt
		}
	}

	switch conv {
	case convHex:
		val := e.val
		if len(e.val)%2 != 0 {
			val = "0" + e.val
		}

		b, err := hex.DecodeString(strings.Join(strings.Fields(val), ""))
		if err != nil {
			return err
		}

		switch {
		case dst.Kind() == reflect.Slice && dst.Type().Elem().Kind() == reflect.Uint8:
			dst.SetBytes(b)
		case dst.Kind() == reflect.Array && isByteSliceCollection(dst.Type()):
			if occurrence >= dst.Len() {
				return fmt.Errorf("too many values, want at most %d", dst.Len())
			}
			dst.Index(occurrence).SetBytes(b)
		case dst.Kind() == reflect.Slice && isByteSliceCollection(dst.Type()):
			dst.Set(reflect.Append(dst, reflect.ValueOf(b)))
		case isSignedInteger(dst.Kind()):
			n, err := strconv.ParseInt(e.val, 16, dst.Type().Bits())
			if err != nil {
				return err
			}
			dst.SetInt(n)
		case isUnsignedInteger(dst.Kind()):
			n, err := strconv.ParseUint(e.val, 16, dst.Type().Bits())
			if err != nil {
				return err
			}
			dst.SetUint(n)
		default:
			return fmt.Errorf("hex requires a []byte or integer field")
		}

	case convInt:
		if e.val == "" {
			return fmt.Errorf("empty integer value")
		}
		switch dst.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(e.val, 0, dst.Type().Bits())
			if err != nil {
				return err
			}
			dst.SetInt(n)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			n, err := strconv.ParseUint(e.val, 0, dst.Type().Bits())
			if err != nil {
				return err
			}
			dst.SetUint(n)
		default:
			return fmt.Errorf("int requires an integer field")
		}

	case convBool:
		if dst.Kind() != reflect.Bool {
			return fmt.Errorf("bool requires a bool field")
		}
		if !e.has || e.val == "" {
			dst.SetBool(true) // bare flag line such as VALID / INVALID
			return nil
		}
		b, err := strconv.ParseBool(e.val)
		if err != nil {
			return err
		}
		dst.SetBool(b)

	case convStr:
		if dst.Kind() != reflect.String {
			return fmt.Errorf("string requires a string field")
		}
		dst.SetString(e.val)
	}
	return nil
}

func isByteSliceCollection(t reflect.Type) bool {
	return (t.Kind() == reflect.Array || t.Kind() == reflect.Slice) &&
		t.Elem().Kind() == reflect.Slice &&
		t.Elem().Elem().Kind() == reflect.Uint8
}

func isSignedInteger(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	default:
		return false
	}
}

func isUnsignedInteger(kind reflect.Kind) bool {
	switch kind {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}
