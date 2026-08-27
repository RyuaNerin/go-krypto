package drbg_test

import (
	"crypto/sha256"
	"testing"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/drbg"
	"github.com/RyuaNerin/go-krypto/internal/cavp"
	"github.com/RyuaNerin/go-krypto/internal/testing/tzeroize"
	"github.com/RyuaNerin/go-krypto/lea"
	"github.com/RyuaNerin/go-krypto/lsh512"

	"github.com/stretchr/testify/require"
)

const (
	katHashSHA256noPR      = "10. Hash_DRBG and HMAC_DRBG/Hash_DRBG/Hash_DRBG_(no PR)_SHA-256_KAT"
	katHashSHA256usePR     = "10. Hash_DRBG and HMAC_DRBG/Hash_DRBG/Hash_DRBG_(use PR)_SHA-256_KAT"
	katHmacSHA256noPR      = "10. Hash_DRBG and HMAC_DRBG/HMAC_DRBG/HMAC_DRBG_(no PR)_SHA-256_KAT"
	katHmacSHA256usePR     = "10. Hash_DRBG and HMAC_DRBG/HMAC_DRBG/HMAC_DRBG_(use PR)_SHA-256_KAT"
	katHmacLSH512512noPR   = "10. Hash_DRBG and HMAC_DRBG/HMAC_DRBG/HMAC_DRBG_(no PR)_LSH-512-512_KAT"
	katCtrLEA128noDFnoPR   = "11. CTR_DRBG/CTR_DRBG_(no DF)(no PR)_LEA-128_KAT"
	katCtrLEA128noDFusePR  = "11. CTR_DRBG/CTR_DRBG_(no DF)(use PR)_LEA-128_KAT"
	katCtrLEA128useDFnoPR  = "11. CTR_DRBG/CTR_DRBG_(use DF)(no PR)_LEA-128_KAT"
	katCtrLEA128useDFusePR = "11. CTR_DRBG/CTR_DRBG_(use DF)(use PR)_LEA-128_KAT"
)

func Test_HashDRBG_SHA256_KAT_NoPR(t *testing.T) {
	cavp.TestDrbgHashKAT(t, katHashSHA256noPR, sha256.New)
}

func Test_HashDRBG_SHA256_KAT_PR(t *testing.T) {
	cavp.TestDrbgHashKAT(t, katHashSHA256usePR, sha256.New)
}

func Test_HMACDRBG_SHA256_KAT_NoPR(t *testing.T) {
	cavp.TestDrbgHmacKAT(t, katHmacSHA256noPR, sha256.New)
}

func Test_HMACDRBG_SHA256_KAT_PR(t *testing.T) {
	cavp.TestDrbgHmacKAT(t, katHmacSHA256usePR, sha256.New)
}

func Test_HMACDRBG_LSH512512_KAT_NoPR(t *testing.T) {
	cavp.TestDrbgHmacKAT(t, katHmacLSH512512noPR, lsh512.New)
}

func Test_CTRDRBG_LEA128_KAT_NoDF_NoPR(t *testing.T) {
	cavp.TestDrbgCtrKAT(t, katCtrLEA128noDFnoPR, lea.NewCipher, 16, false)
}

func Test_CTRDRBG_LEA128_KAT_NoDF_PR(t *testing.T) {
	cavp.TestDrbgCtrKAT(t, katCtrLEA128noDFusePR, lea.NewCipher, 16, false)
}

func Test_CTRDRBG_LEA128_KAT_DF_NoPR(t *testing.T) {
	cavp.TestDrbgCtrKAT(t, katCtrLEA128useDFnoPR, lea.NewCipher, 16, true)
}

func Test_CTRDRBG_LEA128_KAT_DF_PR(t *testing.T) {
	cavp.TestDrbgCtrKAT(t, katCtrLEA128useDFusePR, lea.NewCipher, 16, true)
}

func Test_Zeroize(t *testing.T) {
	t.Run("CTR-DRBG", func(t *testing.T) {
		d, _ := drbg.NewCTRDRBG(lea.NewCipher, 16)
		tzeroize.TestZeroize(t, d)
	})
	t.Run("HMAC-DRBG", func(t *testing.T) {
		d, _ := drbg.NewHMACDRBG(sha256.New)
		tzeroize.TestZeroize(t, d)
	})
	t.Run("Hash-DRBG", func(t *testing.T) {
		d, _ := drbg.NewHashDRBG(sha256.New())
		tzeroize.TestZeroize(t, d)
	})
}

func Test_ZeroizeMarksClosed(t *testing.T) {
	type closer interface {
		Close() error
	}

	t.Run("CTR-DRBG", func(t *testing.T) {
		d, err := drbg.NewCTRDRBG(lea.NewCipher, 16)
		require.NoError(t, err)
		d.Zeroize()
		c, ok := d.(closer)
		require.True(t, ok)
		require.ErrorIs(t, c.Close(), krypto.ErrStateCleared)
	})
	t.Run("HMAC-DRBG", func(t *testing.T) {
		d, err := drbg.NewHMACDRBG(sha256.New)
		require.NoError(t, err)
		d.Zeroize()
		c, ok := d.(closer)
		require.True(t, ok)
		require.ErrorIs(t, c.Close(), krypto.ErrStateCleared)
	})
	t.Run("Hash-DRBG", func(t *testing.T) {
		d, err := drbg.NewHashDRBG(sha256.New())
		require.NoError(t, err)
		d.Zeroize()
		c, ok := d.(closer)
		require.True(t, ok)
		require.ErrorIs(t, c.Close(), krypto.ErrStateCleared)
	})
}
