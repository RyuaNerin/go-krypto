// Package tzeroize provides test helpers for Zeroize implementations.
package tzeroize

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unsafe"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"

	"github.com/stretchr/testify/require"
)

func TestZeroize(t testing.TB, v any) {
	t.Helper()

	_, ok := v.(krypto.Zeroizable)
	require.Truef(t, ok, "%T does not implement krypto.Zeroizable", v)

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Pointer {
		require.Falsef(t, rv.IsNil(), "%T is nil", v)
		rv = rv.Elem()
	}

	var b []byte
	switch rv.Kind() {
	case reflect.Slice:
		b = unsafe.Slice((*byte)(rv.UnsafePointer()), uintptr(rv.Len())*rv.Type().Elem().Size())
	default:
		require.True(t, rv.CanAddr(), "fnNew returns a value that cannot be addressed")
		b = unsafe.Slice((*byte)(unsafe.Pointer(rv.UnsafeAddr())), rv.Type().Size())
	}

	want := make([]byte, len(b))

	require.NotEqual(t, want, b, "fnNew returns a value that is already zeroized")
	zeroize.Zeroize(v)

	if !bytes.Equal(want, b) {
		et := rv.Type()

		var bad []string
		nonZeroFields(rv, et.Name(), &bad)
		require.NotEmpty(t, bad, "padding bytes are not zeroized")

		t.Fatalf("non-zero fields after zeroization:\n%s", strings.Join(bad, "\n"))
	}
}

func nonZeroFields(rv reflect.Value, path string, out *[]string) {
	switch rv.Kind() {
	case reflect.Slice:
		for i := range rv.Len() {
			nonZeroFields(rv.Index(i), fmt.Sprintf("%s[%d]", path, i), out)
		}
	case reflect.Struct:
		for i := range rv.NumField() {
			nonZeroFields(rv.Field(i), path+"."+rv.Type().Field(i).Name, out)
		}
	case reflect.Array:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			if !rv.IsZero() {
				*out = append(*out, fmt.Sprintf("%s (byte array) = %v", path, rv))
			}
			return
		}
		for i := range rv.Len() {
			nonZeroFields(rv.Index(i), fmt.Sprintf("%s[%d]", path, i), out)
		}
	default:
		if !rv.IsZero() {
			*out = append(*out, fmt.Sprintf("%s = %v", path, rv))
		}
	}
}
