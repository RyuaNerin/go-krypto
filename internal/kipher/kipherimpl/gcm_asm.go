//go:build gc && (arm64 || amd64) && !noasm && !noasm_krypto && !noasm_krypto_kipher

package kipherimpl

import (
	"runtime"
	"unsafe"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/memory"

	"golang.org/x/sys/cpu"
)

// The amd64 GHASH code also uses PSHUFB (SSSE3) and PINSRQ (SSE4.1).
var SupportsGFMUL = (cpu.X86.HasPCLMULQDQ && cpu.X86.HasSSSE3 && cpu.X86.HasSSE41) || cpu.ARM64.HasPMULL

func init() {
	if SupportsGFMUL {
		if runtime.GOARCH == "amd64" {
			gcmStateID = 1
		} else {
			gcmStateID = 2
		}
		gcmInit = GCMInitAsm
		gcmDeriveCounter = gcmDeriveCounterAsm
		gcmUpdate = gcmUpdateAsm
		gcmAuth = GCMAuthAsm
		gcmFinish = gcmFinishAsm
	} else {
		gcmInit = GCMInitGo
		gcmDeriveCounter = GCMDeriveCounterGo
		gcmUpdate = GCMUpdateGo
		gcmAuth = GCMAuthGo
		gcmFinish = GCMFinishGo
	}
}

//go:noescape
func _gcmFinish(productTable *[16]GCMFieldElement, tagMask, T unsafe.Pointer, pLen, dLen uint64)

//go:noescape
func _gcmInit(productTable *[16]GCMFieldElement, ks *[GCMBlockSize]byte)

//go:noescape
func _gcmData(productTable *[16]GCMFieldElement, data []byte, T unsafe.Pointer)

func GCMInitAsm(g *GCM, cipher krypto.Block) {
	var key [GCMBlockSize]byte
	defer clear(key[:])
	cipher.Encrypt(key[:], key[:])

	_gcmInit(&g.ProductTable, &key)
}

func gcmUpdateAsm(g *GCM, y *GCMFieldElement, blocks []byte) {
	_gcmData(&g.ProductTable, blocks, unsafe.Pointer(y))
}

func gcmDeriveCounterAsm(g *GCM, counter *[GCMBlockSize]byte, nonce []byte) {
	if len(nonce) == GCMStandardNonceSize {
		// Init counter to nonce||1
		copy(counter[:], nonce)
		counter[GCMBlockSize-1] = 1
	} else {
		var tagMask [GCMTagSize]byte

		// Otherwise counter = GHASH(nonce)
		_gcmData(&g.ProductTable, nonce, unsafe.Pointer(counter))
		_gcmFinish(&g.ProductTable, memory.Pointer(tagMask[:]), unsafe.Pointer(counter), uint64(len(nonce)), uint64(0))
	}
}

func GCMAuthAsm(g *GCM, out, ciphertext, additionalData []byte, tagMask *[GCMTagSize]byte) {
	out = out[:GCMTagSize]
	clear(out)

	_gcmData(&g.ProductTable, additionalData, memory.Pointer(out))
	_gcmData(&g.ProductTable, ciphertext, memory.Pointer(out))

	_gcmFinish(&g.ProductTable, unsafe.Pointer(tagMask), memory.Pointer(out), uint64(len(ciphertext)), uint64(len(additionalData)))
}

func gcmFinishAsm(g *GCM, out []byte, y *GCMFieldElement, ciphertextLen, additionalDataLen uint64, tagMask *[GCMTagSize]byte) {
	copy(out, (*[GCMBlockSize]byte)(unsafe.Pointer(y))[:])
	_gcmFinish(&g.ProductTable, unsafe.Pointer(tagMask), memory.Pointer(out), ciphertextLen, additionalDataLen)
}
