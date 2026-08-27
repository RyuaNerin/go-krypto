package lsh256

//go:generate go run ./internal/generate_compress.go
//go:generate go run ../internal/generators/goat_fmt/main.go --platform amd64 --namespace lsh256 --filename asm_amd64_sse2   --src src/lsh256_sse2.c   --machine sse2             --symbol StepConstants=stepConstants --symbol GammaMask0=gammaMask0,GammaMask1=gammaMask1,GammaMask2=gammaMask2,GammaMask3=gammaMask3,GammaMask4=gammaMask4,GammaMask5=gammaMask5
//go:generate go run ../internal/generators/goat_fmt/main.go --platform amd64 --namespace lsh256 --filename asm_amd64_ssse3  --src src/lsh256_ssse3.c  --machine ssse3            --symbol StepConstants=stepConstants --symbol BytePermInfo_L=bytePermInfo,BytePermInfo_R=bytePermInfoR
//go:generate go run ../internal/generators/goat_fmt/main.go --platform amd64 --namespace lsh256 --filename asm_amd64_avx2   --src src/lsh256_avx2.c   --machine avx2             --symbol StepConstants=stepConstants --symbol BytePermInfo=bytePermInfo,MsgWordPermInfo=msgWordPermInfo
//go:generate go run ../internal/generators/goat_fmt/main.go --platform amd64 --namespace lsh256 --filename asm_amd64_avx512 --src src/lsh256_avx512.c --machine avx512f,avx512vl --symbol StepConstants=stepConstants --symbol BytePermInfo=bytePermInfo,MsgWordPermInfo=msgWordPermInfo
//go:generate go run ../internal/generators/goat_fmt/main.go --platform arm64 --namespace lsh256 --filename asm_arm64_neon   --src src/lsh256_neon.c                              --symbol StepConstants=stepConstants --symbol GAMMA=gammaNEON,GAMMA_INV=gammaNEONInv
