package repeat

// LimitMaxTries returns true if attempt number is less then max.
func LimitMaxTries(max int) Operation {
	return FnWithErrorAndCounter(func(e error, c int) error {
		if c < max {
			return e
		}
		return &StopError{e}
	})
}

// StopOnSuccess returns true in case of error is nil.
func StopOnSuccess() Operation {
	return func(e error) error {
		if e != nil {
			return e
		}
		return &StopError{e}
	}
}

// FnWithErrorAndCounter wraps operation and adds call counter.
func FnWithErrorAndCounter(op func(error, int) error) Operation {
	c := 0
	return func(e error) error {
		defer func() { c++ }()
		return op(e, c)
	}
}

// FnWithCounter wraps operation with counter only.
func FnWithCounter(op func(int) error) Operation {
	return FnWithErrorAndCounter(func(_ error, c int) error {
		return op(c)
	})
}

// Fn wraps operation with no arguments.
func Fn(op func() error) Operation {
	return func(_ error) error {
		return op()
	}
}
