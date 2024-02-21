//go:build arm64 && !purego
// +build arm64,!purego

package aria

var newCipher = newCipherAsm
