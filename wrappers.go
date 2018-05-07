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

// Forward returns the passed operation.
func Forward(op Operation) Operation {
	return op
}
