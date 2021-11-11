package gophercloud

import (
	"context"
	"sync"
	"time"
)

type mergeContext struct {
	child  context.Context
	parent context.Context
	ch     chan struct{}

	sync.Mutex
	err error
}

func idempotentlyClose(ch chan struct{}) {
	select {
	case <-ch:
	default:
		close(ch)
	}
}

// MergeContext returns a copy of ctx that also obeys a second parent.
// MergeContext always returns a non-nil context.Context.
func MergeContext(child, parent context.Context) (context.Context, context.CancelFunc) {
	if child == nil {
		if parent == nil {
			return context.WithCancel(context.Background())
		}
		return context.WithCancel(parent)
	}
	if parent == nil {
		return context.WithCancel(child)
	}

	ch := make(chan struct{})
	cancelCh := make(chan struct{})

	go func() {
		select {
		case <-child.Done():
		case <-parent.Done():
		case <-cancelCh:
		}
		close(ch)
	}()

	return &mergeContext{
		child:  child,
		parent: parent,
		ch:     ch,
	}, func() { idempotentlyClose(cancelCh) }
}

// Value returns Context's value for the key, or parent's if nil.
func (ctx *mergeContext) Value(key interface{}) interface{} {
	if v := ctx.child.Value(key); v != nil {
		return v
	}
	return ctx.parent.Value(key)
}

// Done returns a channel that is closed when either Context's or parent's is.
func (ctx *mergeContext) Done() <-chan struct{} {
	return ctx.ch
}

// Err returns Context's Err(), or parent's if nil. After Err returns a
// non-nil error, successive calls to Err return the same error.
func (ctx *mergeContext) Err() error {
	ctx.Lock()
	defer ctx.Unlock()

	if ctx.err == nil {
		if err := ctx.child.Err(); err != nil {
			ctx.err = err
		} else {
			ctx.err = ctx.parent.Err()
		}
	}
	return ctx.err
}

// Deadline returns the closest deadline, if any.
func (ctx *mergeContext) Deadline() (deadline time.Time, ok bool) {
	if d1, ok := ctx.child.Deadline(); ok {
		if d2, ok := ctx.parent.Deadline(); ok {
			if d1.Before(d2) {
				return d1, ok
			} else {
				return d2, ok
			}
		}
		return d1, ok
	}
	return ctx.parent.Deadline()
}
