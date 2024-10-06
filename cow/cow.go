package cow

import "sync/atomic"

type Cow[T any] struct {
	p *atomic.Pointer[T]
}

func NewCow[T any](v *T) *Cow[T] {
	ptr := &atomic.Pointer[T]{}
	ptr.Store(v)

	return &Cow[T]{p: ptr}
}

func (c *Cow[T]) Get() *T {
	return c.p.Load()
}

func (c *Cow[T]) Tx(f func(*T) *T) {
	for {
		oldPtr := c.p.Load()
		newPtr := f(oldPtr)

		if oldPtr == newPtr {
			return
		}

		if c.p.CompareAndSwap(oldPtr, newPtr) {
			return
		}
	}
}
