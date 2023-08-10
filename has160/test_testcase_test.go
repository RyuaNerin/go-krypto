package has160

import (
	"strings"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
)

func Test_HAS160(t *testing.T) { testGo(t, testCases) }

// TTAS.KO-12.0011/R2
var testCases = []testCase{
	{
		M:  internal.HB(``),
		MD: internal.HB(`30 79 64 ef 34 15 1d 37 c8 04 7a de c7 ab 50 f4 ff 89 76 2d`),
	},
	{
		M:  []byte(`a`),
		MD: internal.HB(`48 72 bc bc 4c d0 f0 a9 dc 7c 2f 70 45 e5 b4 3b 6c 83 0d b8`),
	},
	{
		M:  []byte(`abc`),
		MD: internal.HB(`97 5e 81 04 88 cf 2a 3d 49 83 84 78 12 4a fc e4 b1 c7 88 04`),
	},
	{
		M:  []byte(`message digest`),
		MD: internal.HB(` 23 38 db c8 63 8d 31 22 5f 73 08 62 46 ba 52 9f 96 71 0b c6`),
	},
	{
		M:  []byte(`abcdefghijklmnopqrstuvwxyz`),
		MD: internal.HB(`59 61 85 c9 ab 67 03 d0 d0 db b9 87 02 bc 0f 57 29 cd 1d 3c`),
	},
	{
		M:  []byte(`12345678901234567890123456789012345678901234567890123456789012345678901234567890`),
		MD: internal.HB(`07 f0 5c 8c 07 73 c5 5c a3 a5 a6 95 ce 6a ca 4c 43 89 11 b5`),
	},
	{
		M:  []byte(strings.Repeat("a", 1_000_000)),
		MD: internal.HB(`d6 ad 6f 06 08 b8 78 da 9b 87 99 9c 25 25 cc 84 f4 c9 f1 8d`),
	},
	{
		M:  []byte(`ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789`),
		MD: internal.HB(`cb 5d 7e fb ca 2f 02 e0 fb 71 67 ca bb 12 3a f5 79 57 64 e5`),
	},
}
