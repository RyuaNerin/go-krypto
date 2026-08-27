// Package unaligned provides helpers for testing values at deliberately
// misaligned addresses.
package unaligned

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

const isInTestEnv = "KRYPTO_UNALIGNED_TEST"

type Value[T any] struct {
	Value T

	storage []byte
	source  T
}

func Copy[T any](source T, alignment, residue uintptr) *Value[T] {
	if alignment == 0 {
		panic(msgAlignmentMustBePositive)
	}
	if residue >= alignment {
		panic(msgResidueMustBeLessThanAlignment)
	}

	sourceValue := reflect.ValueOf(source)
	if sourceValue.Kind() != reflect.Pointer || sourceValue.IsNil() {
		panic(msgSourceMustBeNonNilPointer)
	}

	elemType := sourceValue.Type().Elem()
	size := elemType.Size()
	storage := make([]byte, size+alignment-1)
	base := unsafe.Pointer(&storage[0])
	offset := (residue + alignment - uintptr(base)%alignment) % alignment
	destination := unsafe.Add(base, offset)

	sourceBytes := unsafe.Slice((*byte)(sourceValue.UnsafePointer()), size)
	copy(storage[offset:offset+size], sourceBytes)

	value := reflect.NewAt(elemType, destination).Interface().(T)
	return &Value[T]{Value: value, storage: storage, source: source}
}

func (v *Value[T]) KeepAlive() {
	runtime.KeepAlive(v)
}

func Test(t *testing.T, alignment uintptr, test func(residue uintptr)) {
	t.Helper()

	require.GreaterOrEqual(t, alignment, uintptr(2), "unaligned: test alignment must be at least 2")

	testName := t.Name()
	if os.Getenv(isInTestEnv) == testName {
		for shiftIdx := uintptr(1); shiftIdx < alignment; shiftIdx++ {
			test(shiftIdx)
		}
		return
	}

	nameParts := strings.Split(testName, "/")
	for idx := range nameParts {
		nameParts[idx] = "^" + regexp.QuoteMeta(nameParts[idx]) + "$"
	}

	cmd := exec.Command(
		os.Args[0],
		"-test.run="+strings.Join(nameParts, "/"),
		"-test.count=1",
	)
	cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", isInTestEnv, testName))

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "unaligned context subprocess failed. output:\n%s", string(output))
}
