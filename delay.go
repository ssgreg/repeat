package repeat

import (
	"context"
	"time"
)

// SetErrorsTimeout specifies the maximum timeout for repetition
// in case of error. This timeout is reset each time when the
// repetition operation is successfully completed.
//
// Default value is maximum time.Duration value.
func SetErrorsTimeout(t time.Duration) func(*DelayOptions) {
	return func(do *DelayOptions) {
		do.ErrorsTimeout = t
	}
}

// SetContext allows to set a context instead of default one.
func SetContext(ctx context.Context) func(*DelayOptions) {
	return func(do *DelayOptions) {
		do.Context = ctx
	}
}

// WithDelay constructs HeartbeatPredicate.
func WithDelay(options ...func(hb *DelayOptions)) Operation {
	do := applyOptions(applyOptions(&DelayOptions{}, defaultOptions()), options)

	shift := func() time.Time {
		return time.Now().Add(do.ErrorsTimeout)
	}

	deadline := shift()

	return func(e error) error {
		// Shift the deadline in case of success.
		if e == nil {
			deadline = shift()
		}

		delayT := time.NewTimer(do.Backoff())
		defer delayT.Stop()
		deadlineT := time.NewTimer(deadline.Sub(time.Now()))
		defer deadlineT.Stop()

		select {
		case <-do.Context.Done():
			// Let out caller know that the op is cancelled.
			return do.Context.Err()

		case <-deadlineT.C:
			// The reason of a deadline is the previous error. Let our
			// caller to take care of it.
			return Cause(e)

		case <-delayT.C:
			return e
		}
	}
}

// DelayOptions holds parameters for a heartbeat process.
type DelayOptions struct {
	ErrorsTimeout time.Duration
	Backoff       func() time.Duration
	Context       context.Context
}

func defaultOptions() []func(hb *DelayOptions) {
	return []func(do *DelayOptions){
		SetContext(context.Background()),
		SetErrorsTimeout(1<<63 - 1),
		FixedBackoff(time.Second).Set(),
	}
}

func applyOptions(do *DelayOptions, options []func(*DelayOptions)) *DelayOptions {
	for _, o := range options {
		o(do)
	}
	return do
}
