package cmac

import (
	"bytes"
	"crypto/aes"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
)

func TestCMAC_SEED(t *testing.T) {
	K := internal.HB(`00112233445566778899aabbccddeeff`)
	M := internal.HB(`000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f`)
	T := [][]byte{
		internal.HB(`91773796cf510124d3593a331b9d7c51`),
		internal.HB(`ac589f018e897633bac105559c737877`),
		internal.HB(`7775acd0e899670e36700354dfe86af0`),
		internal.HB(`abc1a497093fb8b0ac6b4ee45cb8aa97`),
		internal.HB(`70579a5d04f88ea52eabd7c5ed205127`),
		internal.HB(`7ccc1b050f049846e977cce7d90d00c3`),
		internal.HB(`e332380e2edf5152d4d3c9b27cd3c528`),
		internal.HB(`ca0624353c0250c112e29e6c8ad90472`),
		internal.HB(`bf16eed199d2800e75ee96f6521a270b`),
		internal.HB(`9303cba6ee6167a4a61b3083c27a6ca3`),
		internal.HB(`9adea49371bd29c79e0577b1b37ceac6`),
		internal.HB(`3b29905b9671ee8da9a74ed2e1a88387`),
		internal.HB(`fac847a9b65495daefd0c76448956242`),
		internal.HB(`1d96cd54d6c6fc73931115afa2d69a94`),
		internal.HB(`46f9bd96af320c18958b3f3a507375d2`),
		internal.HB(`a3f7db8a0f76764ef0f2e9dfefae95fd`),
		internal.HB(`c1f732b52fb20caab58d5b6c78cbd514`),
		internal.HB(`3de5233b32fd1d5cb391900ae88a6b55`),
		internal.HB(`b08b2062315af6887d451c4bf99f804b`),
		internal.HB(`a5f251c2472819f89c7d3d1434af9fa3`),
		internal.HB(`c9c37c5a03d8ce644b0782de9a196eb4`),
		internal.HB(`0ddab9cb6b8337d180353105f3dc80a3`),
		internal.HB(`a960de18a95fce2e99e997ef56736e6e`),
		internal.HB(`1216199e6eb717fedeb5f65b216f80a7`),
		internal.HB(`af99c263be9f1f59824b1f381ed9882f`),
		internal.HB(`5301fb1ed8a026bce0a018b2850f52ed`),
		internal.HB(`99f519c6c5a1b80b88ff766590d86c7e`),
		internal.HB(`5b5e4eaff09c55abf89aa942cb24d9da`),
		internal.HB(`1e5dcbb348ebfa79f72f48ef53f517ab`),
		internal.HB(`729afbcbe9ba0a70f89cc747142538ba`),
		internal.HB(`a1e421e85473d05ddbd725536a471106`),
		internal.HB(`4a96b7fcc31df8bb541d6a07a244b60a`),
		internal.HB(`85c86ec86d508d8accf907972fda0436`),
	}

	b, err := aes.NewCipher(K)
	if err != nil {
		panic(err)
	}

	h := New(b)

	var dst []byte
	for length := 1; length < len(M); length++ {
		// full write
		h.Reset()
		h.Write(M[:length])
		dst = h.Sum(dst[:0])

		if !bytes.Equal(dst, T[length]) {
			t.Errorf("FAILED %d\n", length)
			return
		}

		// partial write
		h.Reset()
		for idx := range M[:length] {
			h.Write(M[idx : idx+1])
		}
		h.Reset()
		h.Write(M[:length])
		dst = h.Sum(dst[:0])

		if !bytes.Equal(dst, T[length]) {
			t.Errorf("FAILED %d\n", length)
			return
		}
	}
}
