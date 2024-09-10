module github.com/RyuaNerin/go-krypto

// The go directive sets the minimum version of Go required to use this module.
// Before Go 1.21, the directive was advisory only;
// now it is a mandatory requirement: Go toolchains refuse to use modules declaring newer Go versions.
go 1.20

require (
	github.com/RyuaNerin/elliptic2 v1.0.0
	github.com/RyuaNerin/testingutil v0.1.1
)
