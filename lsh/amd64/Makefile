EXECUTABLES=clang c2goasm asm2plan9s asmfmt as yasm
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))

C2GOASM=c2goasm

CC=clang
CFLAGS_ARCH=-arch x86_64 -masm=intel
CFLAGS= \
	-O3 \
	-fno-asynchronous-unwind-tables \
	-fno-exceptions \
	-fno-rtti \
	-mno-red-zone \
	-mstackrealign \
	-fno-builtin \
	-mllvm \
	-inline-threshold=1000 \
	-fno-stack-protector \
	-fno-jump-tables

CFLAGS_SSE2= -msse -msse2 -mno-ssse3 -mno-avx -mno-avx2
CFLAGS_SSSE3=-msse -msse2    -mssse3 -mno-avx -mno-avx2
CFLAGS_AVX2= -msse -msse2    -mssse3    -mavx    -mavx2

all: lsh256 lsh512

clean:
	rm lsh*.s || true
	rm src/lsh*.s || true

lsh256: lsh256_s lsh256_go
lsh256_s:  lsh256_sse2_s  lsh256_ssse3_s  lsh256_avx2_s
lsh256_go: lsh256_sse2_go lsh256_ssse3_go lsh256_avx2_go

lsh256_sse2:  lsh256_sse2_s  lsh256_sse2_go
lsh256_ssse3: lsh256_ssse3_s lsh256_ssse3_go
lsh256_avx2:  lsh256_avx2_s  lsh256_avx2_go

lsh256_sse2_s:
	${CC} ${CFLAGS_ARCH} ${CFLAGS} ${CFLAGS_SSE2}  -S src/lsh256_sse2.c  -o src/lsh256_sse2.s
lsh256_ssse3_s:
	${CC} ${CFLAGS_ARCH} ${CFLAGS} ${CFLAGS_SSSE3} -S src/lsh256_ssse3.c -o src/lsh256_ssse3.s
lsh256_avx2_s:
	${CC} ${CFLAGS_ARCH} ${CFLAGS} ${CFLAGS_AVX2}  -S src/lsh256_avx2.c  -o src/lsh256_avx2.s

lsh256_sse2_go:
	${C2GOASM} -a -f src/lsh256_sse2.s  lsh256_amd64_sse2.s
lsh256_ssse3_go:
	${C2GOASM} -a -f src/lsh256_ssse3.s lsh256_amd64_ssse3.s
lsh256_avx2_go:
	${C2GOASM} -a -f src/lsh256_avx2.s  lsh256_amd64_avx2.s

lsh512: lsh512_s lsh512_go
lsh512_s:  lsh512_sse2_s  lsh512_ssse3_s  lsh512_avx2_s
lsh512_go: lsh512_sse2_go lsh512_ssse3_go lsh512_avx2_go

lsh512_sse2:  lsh512_sse2_s  lsh512_sse2_go
lsh512_ssse3: lsh512_ssse3_s lsh512_ssse3_go
lsh512_avx2:  lsh512_avx2_s  lsh512_avx2_go

lsh512_sse2_s:
	${CC} ${CFLAGS_ARCH} ${CFLAGS} ${CFLAGS_SSE2}  -S src/lsh512_sse2.c  -o src/lsh512_sse2.s
lsh512_ssse3_s:
	${CC} ${CFLAGS_ARCH} ${CFLAGS} ${CFLAGS_SSSE3} -S src/lsh512_ssse3.c -o src/lsh512_ssse3.s
lsh512_avx2_s:
	${CC} ${CFLAGS_ARCH} ${CFLAGS} ${CFLAGS_AVX2}  -S src/lsh512_avx2.c  -o src/lsh512_avx2.s

lsh512_sse2_go: lsh512_sse2_s
	${C2GOASM} -a -f src/lsh512_sse2.s  lsh512_amd64_sse2.s
lsh512_ssse3_go: lsh512_ssse3_s
	${C2GOASM} -a -f src/lsh512_ssse3.s lsh512_amd64_ssse3.s
lsh512_avx2_go: lsh512_avx2_s
	${C2GOASM} -a -f src/lsh512_avx2.s  lsh512_amd64_avx2.s
