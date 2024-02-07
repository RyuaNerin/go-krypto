package drbg

import (
	"io"
)

func WriteBytes(arr ...[]byte) func(io.Writer) {
	switch len(arr) {
	case 0:
		return nil

	case 1:
		if len(arr[0]) == 0 {
			return nil
		}
		return func(w io.Writer) {
			w.Write(arr[0])
		}

	default:
		return func(w io.Writer) {
			for _, v := range arr {
				if len(v) > 0 {
					w.Write(v)
				}
			}
		}
	}
}
