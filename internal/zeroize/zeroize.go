package zeroize

import (
	"errors"
	"math/big"
	"reflect"
)

var ErrStateCleared = errors.New("state cleared")

type Zeroizier interface {
	Zeroize()
}

func ZeroizeT[T any](p *T) {
	if p == nil {
		return
	}

	var zero T
	*p = zero
}

func Zeroize(values ...any) {
	for _, b := range values {
		if b == nil {
			continue
		}

		switch v := b.(type) {
		case Zeroizier:
			v.Zeroize()

		case *big.Int:
			if v != nil {
				bits := v.Bits()
				clear(bits[:cap(bits)])
				v.SetInt64(0)
			}

		case []byte:
			clear(v)
		case []uint16:
			clear(v)
		case []uint32:
			clear(v)
		case []uint64:
			clear(v)
		case []uint:
			clear(v)
		case []int8:
			clear(v)
		case []int16:
			clear(v)
		case []int32:
			clear(v)
		case []int64:
			clear(v)
		case []int:
			clear(v)

		default:
			rv := reflect.ValueOf(b)
			switch rv.Kind() {
			case reflect.Pointer:
				if !rv.IsNil() {
					rv.Elem().SetZero()
				}
			case reflect.Slice:
				for i, n := 0, rv.Len(); i < n; i++ {
					rv.Index(i).SetZero()
				}
			}
		}
	}
}
