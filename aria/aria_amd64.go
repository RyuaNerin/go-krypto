package aria

var (
	newCipher = newCipherGo
)

func init() {
	if hasSSSE3 {
		newCipher = newCipherAsm
	}
}
