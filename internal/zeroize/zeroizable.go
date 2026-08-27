package zeroize

import (
	"sync"
	"sync/atomic"
)

type ZeroizedFlag struct {
	once  sync.Once
	valid atomic.Bool
}

func (z *ZeroizedFlag) init() {
	z.once.Do(func() {
		z.valid.Store(true)
	})
}

func (z *ZeroizedFlag) SetInvalid() {
	z.init()
	z.valid.Store(false)
}

func (z *ZeroizedFlag) ReturnIfZeroized() error {
	z.init()
	if z.valid.Load() {
		return nil
	}
	return ErrStateCleared
}

func (z *ZeroizedFlag) ThrowIfZeroized() {
	z.init()
	if z.valid.Load() {
		return
	}

	panic(ErrStateCleared)
}
