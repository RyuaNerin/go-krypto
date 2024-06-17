// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpu

const cacheLineSize = 32

// HWCAP/HWCAP2 bits.
// These are specific to Linux.
const (
	hwcap_NEON = 1 << 12

	hwcap2_PMULL = 1 << 1
)

func initOptions() {
	options = []option{
		{Name: "pmull", Feature: &ARM.HasPMULL},
		{Name: "neon", Feature: &ARM.HasNEON},
	}
}
