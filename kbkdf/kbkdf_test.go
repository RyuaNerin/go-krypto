package kbkdf_test

import (
	"crypto/sha256"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/cavp"
	"github.com/RyuaNerin/go-krypto/kbkdf"
	"github.com/RyuaNerin/go-krypto/lea"

	"github.com/stretchr/testify/require"
)

const (
	katCmacCtrLEA128               = "19. KBKDF/CMAC_KDF/CTR/KDF_(CTR)_CMAC_LEA-128_KAT"
	katCmacDpLEA128noCTR           = "19. KBKDF/CMAC_KDF/DP/KDF_(DP)(no CTR)_CMAC_LEA-128_KAT"
	katCmacDpLEA128useCTR          = "19. KBKDF/CMAC_KDF/DP/KDF_(DP)(use CTR)_CMAC_LEA-128_KAT"
	katCmacFbLEA128noCTR           = "19. KBKDF/CMAC_KDF/FB/KDF_(FB)(no CTR)_CMAC_LEA-128_KAT"
	katCmacFbLEA128useCTR          = "19. KBKDF/CMAC_KDF/FB/KDF_(FB)(use CTR)(allowed Null IV)_CMAC_LEA-128_KAT"
	katCmacFbLEA128useCTRrequireIV = "19. KBKDF/CMAC_KDF/FB/KDF_(FB)(use CTR)(Not allowed Null IV)_CMAC_LEA-128_KAT"

	katHmacCtrSHA256               = "19. KBKDF/HMAC_KDF/CTR/KDF_(CTR)_HMAC_SHA2-256_KAT"
	katHmacDpSHA256noCTR           = "19. KBKDF/HMAC_KDF/DP/KDF_(DP)(no CTR)_HMAC_SHA2-256_KAT"
	katHmacDpSHA256useCTR          = "19. KBKDF/HMAC_KDF/DP/KDF_(DP)(use CTR)_HMAC_SHA2-256_KAT"
	katHmacFbSHA256noCTR           = "19. KBKDF/HMAC_KDF/FB/KDF_(FB)(no CTR)_HMAC_SHA2-256_KAT"
	katHmacFbSHA256useCTR          = "19. KBKDF/HMAC_KDF/FB/KDF_(FB)(use CTR)(allowed Null IV)_HMAC_SHA2-256_KAT"
	katHmacFbSHA256useCTRrequireIV = "19. KBKDF/HMAC_KDF/FB/KDF_(FB)(use CTR)(Not allowed Null IV)_HMAC_SHA2-256_KAT"
)

//////////////////////////////////////////////////
// MARK: CMAC

func Test_KBKDF_CMAC_CTR_LEA128_KAT(t *testing.T) {
	cavp.TestKbkdfKat(t, katCmacCtrLEA128, cavp.KBKDFModeCounter, kbkdf.NewCMACPRF(lea.NewCipher))
}

func Test_KBKDF_CMAC_DP_LEA128_KAT_NoCTR(t *testing.T) {
	cavp.TestKbkdfKat(t, katCmacDpLEA128noCTR, cavp.KBKDFModeDoublePipeline, kbkdf.NewCMACPRF(lea.NewCipher))
}

func Test_KBKDF_CMAC_DP_LEA128_KAT_UseCTR(t *testing.T) {
	cavp.TestKbkdfKat(t, katCmacDpLEA128useCTR, cavp.KBKDFModeDoublePipeline, kbkdf.NewCMACPRF(lea.NewCipher))
}

func Test_KBKDF_CMAC_FB_LEA128_KAT_NoCTR(t *testing.T) {
	cavp.TestKbkdfKat(t, katCmacFbLEA128noCTR, cavp.KBKDFModeFeedback, kbkdf.NewCMACPRF(lea.NewCipher))
}

func Test_KBKDF_CMAC_FB_LEA128_KAT_UseCTR(t *testing.T) {
	cavp.TestKbkdfKat(t, katCmacFbLEA128useCTR, cavp.KBKDFModeFeedback, kbkdf.NewCMACPRF(lea.NewCipher))
}

func Test_KBKDF_CMAC_FB_LEA128_KAT_UseCTR_RequireIV(t *testing.T) {
	cavp.TestKbkdfKat(t, katCmacFbLEA128useCTRrequireIV, cavp.KBKDFModeFeedback, kbkdf.NewCMACPRF(lea.NewCipher))
}

//////////////////////////////////////////////////
// MARK: HMAC

func Test_KBKDF_HMAC_CTR_SHA256_KAT(t *testing.T) {
	cavp.TestKbkdfKat(t, katHmacCtrSHA256, cavp.KBKDFModeCounter, kbkdf.NewHMACPRF(sha256.New))
}

func Test_KBKDF_HMAC_DP_SHA256_KAT_NoCTR(t *testing.T) {
	cavp.TestKbkdfKat(t, katHmacDpSHA256noCTR, cavp.KBKDFModeDoublePipeline, kbkdf.NewHMACPRF(sha256.New))
}

func Test_KBKDF_HMAC_DP_SHA256_KAT_UseCTR(t *testing.T) {
	cavp.TestKbkdfKat(t, katHmacDpSHA256useCTR, cavp.KBKDFModeDoublePipeline, kbkdf.NewHMACPRF(sha256.New))
}

func Test_KBKDF_HMAC_FB_SHA256_KAT_NoCTR(t *testing.T) {
	cavp.TestKbkdfKat(t, katHmacFbSHA256noCTR, cavp.KBKDFModeFeedback, kbkdf.NewHMACPRF(sha256.New))
}

func Test_KBKDF_HMAC_FB_SHA256_KAT_UseCTR(t *testing.T) {
	cavp.TestKbkdfKat(t, katHmacFbSHA256useCTR, cavp.KBKDFModeFeedback, kbkdf.NewHMACPRF(sha256.New))
}

func Test_KBKDF_HMAC_FB_SHA256_KAT_UseCTR_RequireIV(t *testing.T) {
	cavp.TestKbkdfKat(t, katHmacFbSHA256useCTRrequireIV, cavp.KBKDFModeFeedback, kbkdf.NewHMACPRF(sha256.New))
}

func Test_KBKDF_CounterLimits(t *testing.T) {
	prf := kbkdf.NewHMACPRF(sha256.New) // 32-byte output
	key := make([]byte, 16)

	// Without a counter every block of the counter mode would be the same.
	require.Panics(t, func() { kbkdf.CounterMode(nil, prf, key, nil, nil, 0, 64) })

	// A 1-byte counter numbers at most 255 blocks.
	require.NotPanics(t, func() { kbkdf.CounterMode(nil, prf, key, nil, nil, 1, 255*32) })
	require.Panics(t, func() { kbkdf.CounterMode(nil, prf, key, nil, nil, 1, 255*32+1) })
	require.Panics(t, func() { kbkdf.FeedbackMode(nil, prf, key, nil, nil, nil, 1, 255*32+1) })
	require.Panics(t, func() { kbkdf.DoublePipeMode(nil, prf, key, nil, nil, 1, 255*32+1) })

	// The feedback and double-pipeline modes do not need a counter.
	require.NotPanics(t, func() { kbkdf.FeedbackMode(nil, prf, key, nil, nil, nil, 0, 256*32) })
	require.NotPanics(t, func() { kbkdf.DoublePipeMode(nil, prf, key, nil, nil, 0, 256*32) })
}
