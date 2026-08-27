package lsh512

//go:generate go run ./internal/generate_compress.go
//go:generate go run ../internal/generators/goat_fmt/main.go --platform amd64 --namespace lsh512 --filename asm_amd64_sse2    --src src/lsh512_sse2.c    --machine sse2              --symbol StepConstants=stepConstants
//go:generate go run ../internal/generators/goat_fmt/main.go --platform amd64 --namespace lsh512 --filename asm_amd64_ssse3   --src src/lsh512_ssse3.c   --machine ssse3             --symbol StepConstants=stepConstants --symbol                                 BytePermInfo=bytePermInfoSSSE3
//go:generate go run ../internal/generators/goat_fmt/main.go --platform amd64 --namespace lsh512 --filename asm_amd64_avx2    --src src/lsh512_avx2.c    --machine avx2              --symbol StepConstants=stepConstants --symbol MsgWordPermInfo=msgWordPermInfo,BytePermInfo=bytePermInfoAVX2
//go:generate go run ../internal/generators/goat_fmt/main.go --platform amd64 --namespace lsh512 --filename asm_amd64_avx512  --src src/lsh512_avx512.c  --machine avx512f,avx512vl  --symbol StepConstants=stepConstants --symbol MsgWordPermInfo=msgWordPermInfo,BytePermInfo=bytePermInfoAVX2
//go:generate go run ../internal/generators/goat_fmt/main.go --platform arm64 --namespace lsh512 --filename asm_arm64_neon    --src src/lsh512_neon.c                                --symbol StepConstants=stepConstants --symbol GAMMA=gammaNEON,GAMMA_INV=gammaNEONInv
