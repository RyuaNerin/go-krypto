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

	// The embedded zeroize guard holds bookkeeping (sync.Once/atomic flags) that
	// is not required to be all-zero after Zeroize. For such types the byte-level
	// check is skipped; the field-level check below verifies the secret state.
	if containsGuard(rv.Type()) {
		et := rv.Type()

		var bad []string
		nonZeroFields(rv, et.Name(), &bad)
		require.Emptyf(t, bad, "non-zero fields after zeroization:\n%s", strings.Join(bad, "\n"))

		return
	}

	if !bytes.Equal(want, b) {
		et := rv.Type()

		var bad []string
		nonZeroFields(rv, et.Name(), &bad)
		require.NotEmpty(t, bad, "padding bytes are not zeroized")

		t.Fatalf("non-zero fields after zeroization:\n%s", strings.Join(bad, "\n"))
	}
}

// containsGuard reports whether t embeds a zeroize.Zeroizable guard anywhere in
// its structure. Such a guard's bookkeeping is not required to be all-zero after
// Zeroize, so the byte-level zeroize check must not be enforced for it.
func containsGuard(t reflect.Type) bool {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return false
	}
	for i := range t.NumField() {
		if t.Field(i).Type == reflect.TypeFor[zeroize.ZeroizedFlag]() {
			return true
		}

		if containsGuard(t.Field(i).Type) {
			return true
		}
	}

	return false
}

func nonZeroFields(rv reflect.Value, path string, out *[]string) {
	switch rv.Kind() {
	case reflect.Slice:
		for i := range rv.Len() {
			nonZeroFields(rv.Index(i), fmt.Sprintf("%s[%d]", path, i), out)
		}
	case reflect.Struct:
		// The embedded zeroize guard holds bookkeeping (sync.Once/atomic flags),
		// not secret state, and is not required to be all-zero after Zeroize.
		if rv.Type() == reflect.TypeFor[zeroize.ZeroizedFlag]() {
			return
		}

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
