// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpu

import "runtime"

// cacheLineSize is used to prevent false sharing of cache lines.
// We choose 128 because Apple Silicon, a.k.a. M1, has 128-byte cache line size.
// It doesn't cost much and is much more future-proof.
const cacheLineSize = 128

func initOptions() {
	options = []option{
		{Name: "pmull", Feature: &ARM64.HasPMULL},
	}
}

func archInit() {
	switch runtime.GOOS {
	case "freebsd":
		readARM64Registers()
	case "linux", "netbsd", "openbsd":
		doinit()
	default:
	}
}

func readARM64Registers() {
	Initialized = true

	parseARM64SystemRegisters(getisar0(), getisar1(), getpfr0())
}

func parseARM64SystemRegisters(isar0, isar1, pfr0 uint64) {
	// ID_AA64ISAR0_EL1
	switch extractBits(isar0, 4, 7) {
	case 1:
	case 2:
		ARM64.HasPMULL = true
	}
}

func extractBits(data uint64, start, end uint) uint {
	return (uint)(data>>start) & ((1 << (end - start + 1)) - 1)
}
