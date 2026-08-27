package lea

//go:generate go run ../internal/generators/goat_fmt/main.go --platform amd64 --namespace lea --filename asm_amd64_sse2   --src src/lea_sse2.c   --machine sse2
//go:generate go run ../internal/generators/goat_fmt/main.go --platform amd64 --namespace lea --filename asm_amd64_avx2   --src src/lea_avx2.c   --machine avx2
//go:generate go run ../internal/generators/goat_fmt/main.go --platform amd64 --namespace lea --filename asm_amd64_avx512 --src src/lea_avx512.c --machine avx512f

//go:generate go run ../internal/generators/goat_fmt/main.go --platform arm64 --namespace lea --filename asm_arm64_neon --src src/lea_neon.c
