package zeroize_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/zeroize"

	"github.com/stretchr/testify/require"
)

type state struct {
	zeroize.ZeroizedFlag

	Data []byte
}

func (s *state) Zeroize() {
	s.ZeroizedFlag.Zeroize()
	zeroize.Zeroize(s.Data)
	s.Data = nil
}

func TestZeroize(t *testing.T) {
	b := []byte{9, 8, 7}
	zeroize.Zeroize(b)
	require.Equal(t, []byte{0, 0, 0}, b)
}

func TestZeroizeBigInt(t *testing.T) {
	printBits := func(desc string, v *big.Int) {
		bits := v.Bits()
		fmt.Println(desc, bits[:cap(bits)])
	}

	v := new(big.Int).Lsh(big.NewInt(1), 256) // v = 0x100...000
	printBits("bits after Lsh:", v)           //
	v.Sub(v, big.NewInt(1))                   // v = 0x0ff...fff
	printBits("bits after Sub:", v)           //
	v.SetInt64(1)                             // v = 0x        1
	printBits("bits after SetInt64:", v)      //

	zeroize.Zeroize(v)
	printBits("bits after zeroize:", v)

	bits := v.Bits()
	for i, w := range bits[:cap(bits)] {
		require.Zero(t, w, "word %d not cleared", i)
	}
}

func TestGuardLifecycle(t *testing.T) {
	s := &state{Data: []byte{1, 2, 3}}
	s.SetValid()

	require.NotPanics(t, func() { s.ThrowIfZeroized() })

	s.Zeroize()
	require.Nil(t, s.Data)
	require.Panics(t, func() { s.ThrowIfZeroized() })
}

func TestReturnIfZeroized(t *testing.T) {
	s := &state{Data: []byte{1, 2, 3}}
	s.SetValid()

	require.NotPanics(t, func() { s.ThrowIfZeroized() })
	require.NoError(t, s.ReturnIfZeroized())

	s.Zeroize()
	require.ErrorIs(t, s.ReturnIfZeroized(), zeroize.ErrStateCleared)
}

func TestZero(t *testing.T) {
	var x int = 42
	zeroize.ZeroizeT(&x)
	require.Zero(t, x)
}
