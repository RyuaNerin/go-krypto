EXECUTABLES=goat clang ar objdump
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))

GOAT=goat
GOAT_FLAGS=-O3

export PATH := /opt/homebrew/opt/binutils/bin/:/opt/homebrew/opt/llvm/bin:$(PATH)
export LDFLAGS := -L/opt/homebrew/opt/llvm/lib
export CPPFLAGS := -I/opt/homebrew/opt/llvm/include

all: neon

neon:
	${GOAT} ${GOAT_FLAGS} src/lea_arm64_neon.c
