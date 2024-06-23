// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpu

import (
	"strings"
	"syscall"
)

// HWCAP/HWCAP2 bits. These are exposed by Linux.
const (
	hwcap_PMULL = 1 << 4
)

// linuxKernelCanEmulateCPUID reports whether we're running
// on Linux 4.11+. Ideally we'd like to ask the question about
// whether the current kernel contains
// https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/commit/?id=77c97b4ee21290f5f083173d957843b615abbff2
// but the version number will have to do.
func linuxKernelCanEmulateCPUID() bool {
	var un syscall.Utsname
	syscall.Uname(&un)
	var sb strings.Builder
	for _, b := range un.Release[:] {
		if b == 0 {
			break
		}
		sb.WriteByte(byte(b))
	}
	major, minor, _, ok := parseRelease(sb.String())
	return ok && (major > 4 || major == 4 && minor >= 11)
}

func doinit() {
	if err := readHWCAP(); err != nil {
		// We failed to read /proc/self/auxv. This can happen if the binary has
		// been given extra capabilities(7) with /bin/setcap.
		//
		// When this happens, we have two options. If the Linux kernel is new
		// enough (4.11+), we can read the arm64 registers directly which'll
		// trap into the kernel and then return back to userspace.
		//
		// But on older kernels, such as Linux 4.4.180 as used on many Synology
		// devices, calling readARM64Registers (specifically getisar0) will
		// cause a SIGILL and we'll die. So for older kernels, parse /proc/cpuinfo
		// instead.
		//
		// See golang/go#57336.
		if linuxKernelCanEmulateCPUID() {
			readARM64Registers()
		} else {
			readLinuxProcCPUInfo()
		}
		return
	}

	ARM64.HasPMULL = isSet(hwCap, hwcap_PMULL)
}

func isSet(hwc uint, value uint) bool {
	return hwc&value != 0
}
