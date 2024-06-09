// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpu

import (
	"io"
	"os"
)

const (
	_AT_HWCAP  = 16
	_AT_HWCAP2 = 26

	procAuxv = "/proc/self/auxv"

	uintSize = int(32 << (^uint(0) >> 63))
)

func readFile(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var size int
	if info, err := f.Stat(); err == nil {
		size64 := info.Size()
		if int64(int(size64)) == size64 {
			size = int(size64)
		}
	}
	size++ // one byte for final read at EOF

	// If a file claims a small size, read at least 512 bytes.
	// In particular, files in Linux's /proc claim size 0 but
	// then do not work right if read in small pieces,
	// so an initial read of 1 byte would not work correctly.
	if size < 512 {
		size = 512
	}

	data := make([]byte, 0, size)
	for {
		n, err := f.Read(data[len(data):cap(data)])
		data = data[:len(data)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return data, err
		}

		if len(data) >= cap(data) {
			d := append(data[:cap(data)], 0)
			data = d[:len(data)]
		}
	}
}

// For those platforms don't have a 'cpuid' equivalent we use HWCAP/HWCAP2
// These are initialized in cpu_$GOARCH.go
// and should not be changed after they are initialized.
var (
	hwCap  uint
	hwCap2 uint
)

func readHWCAP() error {
	// For Go 1.21+, get auxv from the Go runtime.
	if a := getAuxv(); len(a) > 0 {
		for len(a) >= 2 {
			tag, val := a[0], uint(a[1])
			a = a[2:]
			switch tag {
			case _AT_HWCAP:
				hwCap = val
			case _AT_HWCAP2:
				hwCap2 = val
			}
		}
		return nil
	}

	buf, err := readFile(procAuxv)
	if err != nil {
		// e.g. on android /proc/self/auxv is not accessible, so silently
		// ignore the error and leave Initialized = false. On some
		// architectures (e.g. arm64) doinit() implements a fallback
		// readout and will set Initialized = true again.
		return err
	}
	bo := hostByteOrder()
	for len(buf) >= 2*(uintSize/8) {
		var tag, val uint
		switch uintSize {
		case 32:
			tag = uint(bo.Uint32(buf[0:]))
			val = uint(bo.Uint32(buf[4:]))
			buf = buf[8:]
		case 64:
			tag = uint(bo.Uint64(buf[0:]))
			val = uint(bo.Uint64(buf[8:]))
			buf = buf[16:]
		}
		switch tag {
		case _AT_HWCAP:
			hwCap = val
		case _AT_HWCAP2:
			hwCap2 = val
		}
	}
	return nil
}
