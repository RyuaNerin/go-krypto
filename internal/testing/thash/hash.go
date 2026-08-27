package thash

import (
	"fmt"
	"hash"
	"reflect"
)

type (
	HashInfo struct {
		BlockSize         int
		Sizes             []CipherSize
		SkipUnalignedTest bool
	}
	CipherSize struct {
		Name        string
		DigestSize  int // bytes
		KATFilename string
	}

	HashTest struct {
		HashInfo
		New map[int]func() hash.Hash
		Sum map[int]any
	}
)

type (
	SumFunction interface {
		~func([]byte) []byte |
			~func(size int, data []byte) []byte |
			~func([]byte) [20]byte |
			~func(size int, data []byte) [20]byte |
			~func([]byte) [28]byte |
			~func(size int, data []byte) [28]byte |
			~func([]byte) [32]byte |
			~func(size int, data []byte) [32]byte |
			~func([]byte) [48]byte |
			~func(size int, data []byte) [48]byte |
			~func([]byte) [64]byte |
			~func(size int, data []byte) [64]byte
	}
	Digest interface {
		~[]byte |
			~[20]byte |
			~[28]byte | ~[32]byte |
			~[48]byte | ~[64]byte
	}
	SumFunc[T Digest] interface {
		func(size int, data []byte) T | func(data []byte) T
	}
)

func NormalizeSumT[T Digest, F SumFunc[T]](sumFunc F, digestSize int) func([]byte) T {
	switch f := any(sumFunc).(type) {
	case func(size int, data []byte) T:
		return func(data []byte) T { return f(digestSize, data) }
	case func(data []byte) T:
		return f
	default:
		panic(fmt.Sprintf("invalid type %T", sumFunc))
	}
}

func NormalizeSum(sumFunc any, digestSize int) func([]byte) []byte {
	v := reflect.ValueOf(sumFunc)
	ty := v.Type()
	byteSlice := reflect.TypeFor[[]byte]()

	if ty.Kind() != reflect.Func || ty.NumOut() != 1 {
		panic(fmt.Sprintf("unsupported sum function type: %T", sumFunc))
	}

	out := ty.Out(0)
	switch {
	case out.Kind() != reflect.Slice && out.Kind() != reflect.Array:
		fallthrough
	case out.Elem().Kind() != reflect.Uint8:
		panic(fmt.Sprintf("unsupported sum function type: %T", sumFunc))
	}

	var makeArgs func([]byte) []reflect.Value
	switch {
	case ty.NumIn() == 1 && ty.In(0) == byteSlice:
		makeArgs = func(data []byte) []reflect.Value {
			return []reflect.Value{reflect.ValueOf(data)}
		}

	case ty.NumIn() == 2 && ty.In(0).Kind() == reflect.Int && ty.In(1) == byteSlice:
		makeArgs = func(data []byte) []reflect.Value {
			return []reflect.Value{reflect.ValueOf(digestSize), reflect.ValueOf(data)}
		}

	default:
		panic(fmt.Sprintf("unsupported sum function type: %T", sumFunc))
	}

	return func(data []byte) []byte {
		ret := v.Call(makeArgs(data))[0]
		if ret.Kind() == reflect.Slice {
			return ret.Bytes()
		}

		buf := make([]byte, ret.Len())
		reflect.Copy(reflect.ValueOf(buf), ret)
		return buf
	}
}
