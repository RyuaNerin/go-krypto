package cavp

import (
	"encoding/hex"
	"strings"
)

func upperHex(src []byte) string {
	return strings.ToUpper(hex.EncodeToString(src))
}
