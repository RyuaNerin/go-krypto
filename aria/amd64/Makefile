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

CFLAGS_SSSE3=-msse -msse2    -mssse3 -mno-avx -mno-avx2

all: aria

clean:
	rm lsh*.s || true
	rm src/lsh*.s || true

aria: aria_s aria_go
aria_s:  aria_ssse3_s
aria_go: aria_ssse3_go

aria_ssse3: aria_ssse3_s aria_ssse3_go

aria_ssse3_s:
	${CC} ${CFLAGS_ARCH} ${CFLAGS} ${CFLAGS_SSSE3} -S ../src/aria.c -o ../src/aria.ssse3.s

aria_ssse3_go:
	${C2GOASM} -a -f ../src/aria.ssse3.s aria_amd64_ssse3.s
