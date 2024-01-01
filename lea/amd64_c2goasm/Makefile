EXECUTABLES=clang c2goasm asm2plan9s asmfmt as yasm
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))

C2GOASM=c2goasm

CC=clang
CFLAGS=-O3 -fno-asynchronous-unwind-tables -fno-exceptions -fno-rtti -mno-red-zone -mstackrealign -mllvm -inline-threshold=1000
CFLAGS_ARCH=-arch x86_64 -masm=intel

CFLAGS_SSE2=-msse -msse2
CFLAGS_SSSE3=${CFLAGS_SSE2} -mssse3
CFLAGS_AVX2=${CFLAGS_SSSE3} -mavx -mavx2


all: sse2 avx2 init

sse2: 
	${CC} ${CFLAGS} ${CFLAGS_ARCH} ${CFLAGS_SSE2} -S src/lea_sse2.c -o src/lea_sse2.s
	${C2GOASM} -a -f src/lea_sse2.s lea_amd64_sse2.s

avx2:
	${CC} ${CFLAGS} ${CFLAGS_ARCH} ${CFLAGS_AVX2} -S src/lea_avx2.c -o src/lea_avx2.s
	${C2GOASM} -a -f src/lea_avx2.s lea_amd64_avx2.s

init:
	${CC} ${CFLAGS} ${CFLAGS_ARCH} ${CFLAGS_SSE2} -S src/lea_init.c -o src/lea_init.s
	${C2GOASM} -a -f src/lea_init.s lea_init.s
