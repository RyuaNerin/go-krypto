package lea

type nonCipherContext struct {
	ctx leaContext
}

func (ctx *nonCipherContext) BlockSize() int {
	return ctx.ctx.BlockSize()
}

func (ctx *nonCipherContext) Encrypt(dst, src []byte) {
	ctx.ctx.Encrypt(dst, src)
}

func (ctx *nonCipherContext) Decrypt(dst, src []byte) {
	ctx.ctx.Decrypt(dst, src)
}
