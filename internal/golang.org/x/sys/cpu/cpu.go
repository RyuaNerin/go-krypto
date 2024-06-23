// https://cs.opensource.google/go/x/sys/+/refs/tags/v0.16.0:cpu/cpu.go

// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build 386 || amd64 || amd64p32 || arm || arm64
// +build 386 amd64 amd64p32 arm arm64

// Package cpu implements processor feature detection for
// various CPU architectures.
package cpu

import (
	"os"
	"strings"
)

// Initialized reports whether the CPU features were initialized.
//
// For some GOOS/GOARCH combinations initialization of the CPU features depends
// on reading an operating specific file, e.g. /proc/self/auxv on linux/arm
// Initialized will report false if reading the file fails.
var Initialized bool

// CacheLinePad is used to pad structs to avoid false sharing.
type CacheLinePad struct{ _ [cacheLineSize]byte }

// X86 contains the supported CPU features of the
// current X86/AMD64 platform. If the current platform
// is not X86/AMD64 then all feature flags are false.
//
// X86 is padded to avoid false sharing. Further the HasAVX
// and HasAVX2 are only set if the OS supports XMM and YMM
// registers in addition to the CPUID feature bit being set.
var X86 struct {
	_            CacheLinePad
	HasAVX       bool // Advanced vector extension
	HasAVX2      bool // Advanced vector extension 2
	HasOSXSAVE   bool // OS supports XSAVE/XRESTOR for saving/restoring XMM registers.
	HasPCLMULQDQ bool // PCLMULQDQ instruction - most often used for AES-GCM
	HasSSE2      bool // Streaming SIMD extension 2 (always available on amd64)
	HasSSSE3     bool // Supplemental streaming SIMD extension 3
	_            CacheLinePad
}

// ARM64 contains the supported CPU features of the
// current ARMv8(aarch64) platform. If the current platform
// is not arm64 then all feature flags are false.
var ARM64 struct {
	_        CacheLinePad
	HasPMULL bool // Polynomial multiplication instruction set
	_        CacheLinePad
}

// ARM contains the supported CPU features of the current ARM (32-bit) platform.
// All feature flags are false if:
//  1. the current platform is not arm, or
//  2. the current operating system is not Linux.
var ARM struct {
	_        CacheLinePad
	HasNEON  bool // NEON instruction set
	HasPMULL bool // Polynomial multiplication instruction set
	_        CacheLinePad
}

func init() {
	archInit()
	initOptions()
	processOptions()
}

// options contains the cpu debug options that can be used in GODEBUG.
// Options are arch dependent and are added by the arch specific initOptions functions.
// Features that are mandatory for the specific GOARCH should have the Required field set
// (e.g. SSE2 on amd64).
var options []option

// Option names should be lower case. e.g. avx instead of AVX.
type option struct {
	Name      string
	Feature   *bool
	Specified bool // whether feature value was specified in GODEBUG
	Enable    bool // whether feature should be enabled
	Required  bool // whether feature is mandatory and can not be disabled
}

func processOptions() {
	env := os.Getenv("GODEBUG")
field:
	for env != "" {
		field := ""
		i := strings.IndexByte(env, ',')
		if i < 0 {
			field, env = env, ""
		} else {
			field, env = env[:i], env[i+1:]
		}
		if len(field) < 4 || field[:4] != "cpu." {
			continue
		}
		i = strings.IndexByte(field, '=')
		if i < 0 {
			print("GODEBUG sys/cpu: no value specified for \"", field, "\"\n")
			continue
		}
		key, value := field[4:i], field[i+1:] // e.g. "SSE2", "on"

		var enable bool
		switch value {
		case "on":
			enable = true
		case "off":
			enable = false
		default:
			print("GODEBUG sys/cpu: value \"", value, "\" not supported for cpu option \"", key, "\"\n")
			continue field
		}

		if key == "all" {
			for i := range options {
				options[i].Specified = true
				options[i].Enable = enable || options[i].Required
			}
			continue field
		}

		for i := range options {
			if options[i].Name == key {
				options[i].Specified = true
				options[i].Enable = enable
				continue field
			}
		}

		print("GODEBUG sys/cpu: unknown cpu feature \"", key, "\"\n")
	}

	for _, o := range options {
		if !o.Specified {
			continue
		}

		if o.Enable && !*o.Feature {
			print("GODEBUG sys/cpu: can not enable \"", o.Name, "\", missing CPU support\n")
			continue
		}

		if !o.Enable && o.Required {
			print("GODEBUG sys/cpu: can not disable \"", o.Name, "\", required CPU feature\n")
			continue
		}

		*o.Feature = o.Enable
	}
}
