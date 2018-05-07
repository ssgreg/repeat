package repeat

import (
	"context"
)

var (
	def = NewRepeater()
)

// Once composes the operations and executes the result once.
//
// It is guaranteed that the first op will be called at least once.
func Once(ops ...Operation) error {
	return def.Once(ops...)
}

// Repeat repeat operations until one of them stops the repetition.
//
// It is guaranteed that the first op will be called at least once.
func Repeat(ops ...Operation) error {
	return def.Repeat(ops...)
}

// FnRepeat is a Repeat operation.
func FnRepeat(ops ...Operation) Operation {
	return def.FnRepeat(ops...)
}

// Compose composes all passed operations into a single one.
func Compose(ops ...Operation) Operation {
	return def.Compose(ops...)
}

// WithContext repeat operations until one of them stops the
// repetition or context will be canceled.
//
// It is guaranteed that the first op will be called at least once.
func WithContext(ctx context.Context) Repeater {
	return Wrap(WrStopOnContextError(ctx))
}

// Repeater represents general package concept.
type Repeater interface {
	Once(...Operation) error
	Repeat(...Operation) error
	Compose(...Operation) Operation
	FnRepeat(...Operation) Operation
}

type stdRepeater struct {
	wop OpWrapper
	c   Operation
	d   Operation
}

// NewRepeater sets up everything to be able to repeat operations.
func NewRepeater() Repeater {
	return &stdRepeater{Forward, Done, Done}
}

// Wrap returns object that wraps all repeating ops with passed OpWrapper.
func Wrap(wop OpWrapper) Repeater {
	return &stdRepeater{wop, Done, Done}
}

// Cpp returns object that calls C (constructor) at first, then ops,
// then D (destructor). D will be called in any case if C returns nil.
//
// Note! Cpp panics if D returns non nil error. Wrap it using Done if
// you log D's error or handle it somehow else.
//
func Cpp(c Operation, d Operation) Repeater {
	return &stdRepeater{Forward, c, FnPanic(d)}
}

// Once composes the operations and executes the result once.
//
// It is guaranteed that the first op will be called at least once.
func (w *stdRepeater) Once(ops ...Operation) error {
	return Cause(w.Compose(ops...)(nil))
}

// Repeat repeat operations until one of them stops the repetition.
//
// It is guaranteed that the first op will be called at least once.
func (w *stdRepeater) Repeat(ops ...Operation) error {
	return Cause(w.FnRepeat(ops...)(nil))
}

// FnRepeat is a Repeat operation.
func (w *stdRepeater) FnRepeat(ops ...Operation) Operation {
	return func(e error) (err error) {
		op := w.Compose(ops...)

		for {
			err = op(e)
			switch err.(type) {
			case nil:
				e = nil
			case *TemporaryError:
				e = err
			case *StopError:
				return err
			default:
				return err
			}
		}
	}
}

// Compose wraps ops with wop and composes all passed operations info
// a single one.
func (w *stdRepeater) Compose(ops ...Operation) Operation {
	return func(e error) (err error) {
		err = w.c(e)
		if err != nil {
			// If C failed with temporary error, stop error or any other
			// error: stop compose with this error.
			return err
		}
		defer func() {
			// Note: handle error using D wrapper.
			_ = w.d(err)
		}()

		for _, op := range ops {
			err = w.wop(op)(e)
			switch err.(type) {
			// Replace last E with nil.
			case nil:
				e = nil
			// Replace last E with new temporary error.
			case *TemporaryError:
				e = err
			// Stop.
			case *StopError:
				return err
			// Stop.
			default:
				return err
			}
		}

		return e
	}
}
