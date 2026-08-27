package krypto

import (
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

// ErrStateCleared is returned when an operation is attempted on a value that has
// already been zeroized.
var ErrStateCleared = zeroize.ErrStateCleared

// Zeroizable is implemented by values that can securely wipe their secret state.
type Zeroizable = zeroize.Zeroizier

// ZeroizeT wipes the value that p points to, resetting it to the zero value of its
// type. It is a no-op when p is nil.
func ZeroizeT[T any](p *T) { zeroize.ZeroizeT(p) }

// Zeroize securely wipes each value in values, clearing any secret state. Values
// implementing Zeroizable are zeroized via their Zeroize method.
func Zeroize(values ...any) { zeroize.Zeroize(values...) }
