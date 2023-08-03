// Code generated by command: go run main.go -out ../lsh256_amd64.s -stubs ../lsh256_amd64_stubs.go -pkg lsh256. DO NOT EDIT.

package lsh256

func lsh256InitSSE2(ctx *lsh256ContextAsmData)

func lsh256UpdateSSE2(ctx *lsh256ContextAsmData, data []byte, databitlen uint32)

func lsh256FinalSSE2(ctx *lsh256ContextAsmData, hashval []byte)

func lsh256UpdateSSSE3(ctx *lsh256ContextAsmData, data []byte, databitlen uint32)

func lsh256FinalSSSE3(ctx *lsh256ContextAsmData, hashval []byte)

func lsh256InitAVX2(ctx *lsh256ContextAsmData)

func lsh256UpdateAVX2(ctx *lsh256ContextAsmData, data []byte, databitlen uint32)

func lsh256FinalAVX2(ctx *lsh256ContextAsmData, hashval []byte)
