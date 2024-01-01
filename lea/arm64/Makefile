EXECUTABLES=goat clang ar objdump
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))

GOAT=goat
GOAT_FLAGS=-O3 --extra-option -fno-builtin

all: lea_neon

lea_neon:
	${GOAT} ${GOAT_FLAGS} src/lea_neon.c
