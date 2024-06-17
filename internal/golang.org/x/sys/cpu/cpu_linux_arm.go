// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpu

func doinit() {
	ARM.HasNEON = isSet(hwCap, hwcap_NEON)
	ARM.HasPMULL = isSet(hwCap2, hwcap2_PMULL)
}

func isSet(hwc uint, value uint) bool {
	return hwc&value != 0
}
