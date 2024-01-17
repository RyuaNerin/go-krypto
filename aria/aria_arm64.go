//go:build arm64 && !purego
// +build arm64,!purego

package aria

//go:noescape
func __encKeySetup_NEON(rk *byte, mk *byte, keyBytes uint64)

//go:noescape
func __decKeySetup_NEON(rk *byte, rounds uint64)

//go:noescape
func __process_NEON(dst, src, rk *byte, round uint64)

func init() {
	newCipher = newCipherAsm
}

func (ctx *ariaContextAsm) init(key []byte) {
	__encKeySetup_NEON(toBytePtr(ctx.ctx.ek[:]), toBytePtr(key), uint64(len(key)))

	ctx.ctx.dk = ctx.ctx.ek
	__decKeySetup_NEON(toBytePtr(ctx.ctx.dk[:]), uint64(ctx.ctx.rounds))
}

func (ctx *ariaContextAsm) process(dst, src, rk []byte) {
	__process_NEON(toBytePtr(dst), toBytePtr(src), toBytePtr(rk), uint64(ctx.ctx.rounds))
}
