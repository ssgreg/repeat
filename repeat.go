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
	opw  OpWrapper
	copw OpWrapper
}

// NewRepeater sets up everything to be able to repeat operations.
func NewRepeater() Repeater {
	return NewRepeaterExt(Forward, Forward)
}

// Wrap returns object that wraps all repeating ops with passed OpWrapper.
func Wrap(opw OpWrapper) Repeater {
	return NewRepeaterExt(opw, Forward)
}

// WrapOnce returns object that wraps all repeating ops combined into a single
// op with passed OpWrapper calling it once.
func WrapOnce(copw OpWrapper) Repeater {
	return NewRepeaterExt(Forward, copw)
}

// NewRepeaterExt returns object that wraps all ops with with the given opw
// and wraps composed operation with the given copw.
func NewRepeaterExt(opw, copw OpWrapper) Repeater {
	return &stdRepeater{opw, copw}
}

// Cpp returns object that calls C (constructor) at first, then ops,
// then D (destructor). D will be called in any case if C returns nil.
//
// Note! Cpp panics if D returns non nil error. Wrap it using Done if
// you log D's error or handle it somehow else.
//
func Cpp(c, d Operation) Repeater {
	return NewRepeaterExt(Forward, WrWith(c, func(e error) error {
		_ = FnPanic(d)(e)

		return e
	}))
}

// With returns object that calls C (constructor) at first, then ops,
// then D (destructor). D will be called in any case if C returns nil.
//
// Note! D is able to hide original error an return nil or return error
// event if the original error is nil.
func With(c, d Operation) Repeater {
	return NewRepeaterExt(Forward, WrWith(c, d))
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
			switch typedError := err.(type) {
			case nil:
				e = nil
			case *TemporaryError:
				e = err
			case *StopError:
				switch typedError.Cause {
				case nil:
					return nil
				default:
					return err
				}
			default:
				return err
			}
		}
	}
}

// Compose wraps ops with wop and composes all passed operations info
// a single one.
func (w *stdRepeater) Compose(ops ...Operation) Operation {
	return w.copw(func(e error) (err error) {
		for _, op := range ops {
			err = w.opw(op)(e)
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
	})
}
