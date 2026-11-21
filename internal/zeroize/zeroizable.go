package zeroize

import (
	"sync"
	"sync/atomic"
)

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

type ZeroizedFlag struct {
	_     noCopy
	valid atomic.Bool
	m     sync.Mutex
}

// Marks the state as valid.
func (z *ZeroizedFlag) SetValid() {
	z.m.Lock()
	defer z.m.Unlock()

	z.valid.Store(true)
}

// Marks the state as zeroized and invalid.
func (z *ZeroizedFlag) Zeroize() {
	z.m.Lock()
	defer z.m.Unlock()

	z.valid.Store(false)
}

func (z *ZeroizedFlag) ReturnIfZeroized() error {
	z.m.Lock()
	defer z.m.Unlock()

	if z.valid.Load() {
		return nil
	}
	return ErrStateCleared
}

func (z *ZeroizedFlag) ThrowIfZeroized() {
	z.m.Lock()
	defer z.m.Unlock()

	if z.valid.Load() {
		return
	}

	panic(ErrStateCleared)
}
