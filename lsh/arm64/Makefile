EXECUTABLES=goat clang ar objdump
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))

GOAT=goat
GOAT_FLAGS=-O3 --extra-option -fno-builtin

all: lsh256_neon lsh512_neon

lsh256_neon:
	${GOAT} ${GOAT_FLAGS} src/lsh256_neon.c

lsh512_neon:
	${GOAT} ${GOAT_FLAGS} src/lsh512_neon.c
