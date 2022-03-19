package seed

import (
	"encoding/hex"
	"testing"
)

type testCase struct {
	key    string
	plain  string
	cipher string
}

var (
	test128 = []testCase{
		{
			"00000000000000000000000000000000",
			"000102030405060708090A0B0C0D0E0F",
			"5EBAC6E0054E166819AFF1CC6D346CDB",
		},
		{
			"000102030405060708090A0B0C0D0E0F",
			"00000000000000000000000000000000",
			"C11F22F20140505084483597E4370F43",
		},
		{
			"4706480851E61BE85D74BFB3FD956185",
			"83A2F8A288641FB9A4E9A5CC2F131C7D",
			"EE54D13EBCAE706D226BC3142CD40D4A",
		},
		{
			"28DBC3BC49FFD87DCFA509B11D422BE7",
			"B41E6BE2EBA84A148E2EED84593C5EC7",
			"9B9B7BFCD1813CB95D0B3618F40F5122",
		},
	}
)

func TestSeed128(t *testing.T) {
	for _, test := range test128 {
		testKey, _ := hex.DecodeString(test.key)
		testPlain, _ := hex.DecodeString(test.plain)
		testCipher, _ := hex.DecodeString(test.cipher)

		b, _ := NewCipher(testKey)

		cipher := make([]byte, 16)
		b.Encrypt(cipher, testPlain)

		plain := make([]byte, 16)
		b.Decrypt(plain, testCipher)

		if !testEq(cipher, testCipher) {
			t.Errorf("Encrypt failed test: %s / answer: %s", hex.EncodeToString(cipher), test.cipher)
		}

		if !testEq(plain, testPlain) {
			t.Errorf("Encrypt failed test: %s / answer: %s", hex.EncodeToString(plain), test.plain)
		}
	}
}

func testEq(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
