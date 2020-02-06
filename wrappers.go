package repeat

import (
	"context"
)

// OpWrapper is the type of function for repetition.
type OpWrapper func(Operation) Operation

// WrStopOnContextError stops an operation in case of context error.
func WrStopOnContextError(ctx context.Context) OpWrapper {
	return func(op Operation) Operation {
		return func(e error) error {
			if ctx.Err() != nil {
				switch e.(type) {
				case nil:
					return HintStop(ctx.Err())
				case *StopError:
					return e
				case *TemporaryError:
					return HintStop(e)
				default:
					return e
				}
			}

			return op(e)
		}
	}
}

// WrWith returns wrapper that calls C (constructor) at first, then ops,
// then D (destructor). D will be called in any case if C returns nil.
func WrWith(c, d Operation) OpWrapper {
	return func(op Operation) Operation {
		return func(e error) (err error) {
			err = c(e)
			if err != nil {
				// If C failed with temporary error, stop error or any other
				// error: stop compose with this error.
				return err
			}
			defer func() {
				// Note: handle error using D wrapper.
				err = d(err)
			}()

			return op(e)
		}
	}
}

// Forward returns the passed operation.
func Forward(op Operation) Operation {
	return op
}
