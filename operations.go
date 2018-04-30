package repeat

// Operation is the type of function for repetition.
type Operation func(error) error

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

// FnOnSuccess executes operation in case of error is nil.
func FnOnSuccess(op func(error) error) Operation {
	return func(e error) error {
		if e != nil {
			return e
		}

		return op(e)
	}
}

// FnOnError executes operation in case error is NOT nil.
func FnOnError(op func(error) error) Operation {
	return func(e error) error {
		if e == nil {
			return e
		}

		return op(e)
	}
}

// FnHintTemporary hints all operation errors as temporary.
func FnHintTemporary(op func(error) error) Operation {
	return func(e error) error {
		err := op(e)
		switch err.(type) {
		case nil:
		case *TemporaryError:
		case *StopError:
		default:
			err = HintTemporary(err)
		}

		return err
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

// FnS wraps operation with no arguments and return value.
func FnS(op func()) Operation {
	return func(e error) error {
		op()
		return e
	}
}

// FnES wraps operation with no return value.
func FnES(op func(error)) Operation {
	return func(e error) error {
		op(e)
		return e
	}
}

// Nope does nothing.
func Nope(e error) error {
	return e
}
