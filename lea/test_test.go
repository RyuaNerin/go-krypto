package lea

import (
	"bytes"
	"encoding/hex"
	"testing"
)

type testCase struct {
	Key    []byte
	Plain  []byte
	Secure []byte
}

func testAll(t *testing.T, f func(*testing.T, int)) {
	tests := []struct {
		name    string
		keySize int
	}{
		{"128", 128},
		{"196", 196},
		{"256", 256},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			f(t, test.keySize)
		})
	}
}

func testEncryptGo(t *testing.T, testCases []testCase) {
	dst := make([]byte, BlockSize)

	var ctx leaContext
	for _, tc := range testCases {
		err := ctx.initContext(tc.Key)
		if err != nil {
			t.Error(err)
		}

		leaEnc1Go(&ctx, dst, tc.Plain)
		if !bytes.Equal(dst, tc.Secure) {
			t.Errorf("encrypt failed.\nresult: %s\nanswer: %s", hex.EncodeToString(dst), hex.EncodeToString(tc.Secure))
		}
	}
}

func testDecryptGo(t *testing.T, testCases []testCase) {
	dst := make([]byte, BlockSize)

	var ctx leaContext
	for _, tc := range testCases {
		err := ctx.initContext(tc.Key)
		if err != nil {
			t.Error(err)
		}

		leaDec1Go(&ctx, dst, tc.Secure)
		if !bytes.Equal(dst, tc.Plain) {
			t.Errorf("encrypt failed.\nresult: %s\nanswer: %s", hex.EncodeToString(dst), hex.EncodeToString(tc.Plain))
		}
	}
}
