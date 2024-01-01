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

all: lea

clean:
	rm lea*.s || true
	rm src/lea*.s || true

lea: lea_s lea_go
lea_s:  lea_sse2_s  lea_avx2_s
lea_go: lea_sse2_go lea_avx2_go

lea_sse2: lea_sse2_s lea_sse2_go
lea_avx2: lea_avx2_s lea_avx2_go

lea_sse2_s: 
	${CC} ${CFLAGS_ARCH} ${CFLAGS} ${CFLAGS_SSE2}  -S src/lea_sse2.c -o src/lea_sse2.s
lea_avx2_s:
	${CC} ${CFLAGS_ARCH} ${CFLAGS} ${CFLAGS_AVX2}  -S src/lea_avx2.c -o src/lea_avx2.s

lea_sse2_go:
	${C2GOASM} -a -f src/lea_sse2.s lea_sse2.s
lea_avx2_go:
	${C2GOASM} -a -f src/lea_avx2.s lea_avx2.s
