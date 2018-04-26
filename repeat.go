package repeat

import "context"

var (
	def = NewRepeater()
)

// Repeat repeat operations until one of them stops the repetition.
//
// It is guaranteed that the first op will be called at least once.
func Repeat(ops ...Operation) error {
	return def.Repeat(ops...)
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
	Repeat(...Operation) error
	Compose(...Operation) Operation
}

type stdRepeater struct {
	wop OpWrapper
}

// NewRepeater sets up everything to be able to repeat operations.
func NewRepeater() Repeater {
	return &stdRepeater{Forward}
}

// Wrap returns object that wraps all repeating ops with passed OpWrapper.
func Wrap(wop OpWrapper) Repeater {
	return &stdRepeater{wop}
}

// Repeat repeat operations until one of them stops the repetition.
//
// It is guaranteed that the first op will be called at least once.
func (w *stdRepeater) Repeat(ops ...Operation) (err error) {
	op := w.Compose(ops...)
	for {
		err = op(err)
		switch e := err.(type) {
		case nil:
		case *TemporaryError:
		case *StopError:
			return e.Cause
		default:
			return e
		}
	}
}

// Compose wraps ops with wop and composes all passed operations info
// a single one.
func (w *stdRepeater) Compose(ops ...Operation) Operation {
	return func(err error) error {
		for _, op := range ops {
			err = w.wop(op)(err)
			switch e := err.(type) {
			case nil:
			case *TemporaryError:
			case *StopError:
				return e
			default:
				return e
			}
		}
		return err
	}
}
