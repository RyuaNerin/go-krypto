package internal

import (
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBitsToBytes(t *testing.T) {
	for bits := range 0xFF_FFFF {
		answer := int(math.Ceil(float64(bits) / 8))

		require.Equal(t, answer, BitsToBytes(bits))
	}
}

func TestDigitLen(t *testing.T) {
	testCases := []struct {
		x        *big.Int
		e        int
		expected int
	}{
		// x == 2 ** e
		{HI(`0x100000000000000000000000000000000000000`), 8, 19},                           // log(0x100000000000000000000000000000000000000)                           / log(2**8) = 19
		{HI(`0x10000000000000000000000000000000000000000000000000000000000000000`), 8, 32}, // log(0x10000000000000000000000000000000000000000000000000000000000000000) / log(2**8) = 32
		{HI(`0xffffffffffffffffffffffff99def836146bc9b1b4d22831`), 8, 24},                  // log(0xffffffffffffffffffffffff99def836146bc9b1b4d22831)                  / log(2**8) = 23.99999999999999999999999999999909194074497683942334210799692027429193921791246144182659925527070974
		{HI(`0x1000000000000000000000000000013e974e72f8a6922031d2603cfe0d7`), 8, 30},       // log(0x1000000000000000000000000000013e974e72f8a6922031d2603cfe0d7)       / log(2**8) = 29.00000000000000000000000000000000000270146066501192228995704198641991460168842359846283006852399854
	}
	for idx, tc := range testCases {
		result := DigitLen(tc.x)
		require.Equal(t, tc.expected, result, "failed test case %d", idx)
	}
}

func TestBigEqual(t *testing.T) {
	require.True(t, BigEqual(big.NewInt(5), big.NewInt(5)))
	require.False(t, BigEqual(big.NewInt(5), big.NewInt(-5)))
	require.False(t, BigEqual(big.NewInt(5), nil))
	require.True(t, BigEqual(nil, nil))
}
