package repeat

// Operation is the type of function for repetition.
type Operation func(error) error

// LimitMaxTries returns true if attempt number is less then max.
func LimitMaxTries(max int) Operation {
	return FnWithErrorAndCounter(func(e error, c int) error {
		if c < max {
			return e
		}

		return HintStop(e)
	})
}

// StopOnSuccess returns true in case of error is nil.
func StopOnSuccess() Operation {
	return func(e error) error {
		if e != nil {
			return e
		}

		return HintStop(e)
	}
}

// FnOnSuccess executes operation in case of error is nil.
func FnOnSuccess(op Operation) Operation {
	return func(e error) error {
		if e != nil {
			return e
		}

		return op(e)
	}
}

// FnOnError executes operation in case error is NOT nil.
func FnOnError(op Operation) Operation {
	return func(e error) error {
		if e == nil {
			return e
		}

		return op(e)
	}
}

// FnHintTemporary hints all operation errors as temporary.
func FnHintTemporary(op Operation) Operation {
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

// FnHintStop hints all operation errors as StopError.
func FnHintStop(op Operation) Operation {
	return func(e error) error {
		err := op(e)
		switch err.(type) {
		case *TemporaryError:
		case *StopError:
		default:
			err = HintStop(err)
		}

		return err
	}
}

// FnPanic panics if op returns any error other than nil, TemporaryError
// and StopError.
func FnPanic(op Operation) Operation {
	return func(e error) error {
		err := op(e)
		switch err.(type) {
		case nil:
		case *TemporaryError:
		case *StopError:
		default:
			panic(err)
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

// Nope does nothing, returns input error.
func Nope(e error) error {
	return e
}

// FnNope does not call pass op, returns input error.
func FnNope(op Operation) Operation {
	return func(e error) error {
		return Nope(e)
	}
}

// Done does nothing, returns nil.
func Done(e error) error {
	return nil
}

// FnDone returns nil even if wrapped op returns an error.
func FnDone(op Operation) Operation {
	return func(e error) error {
		return Done(op(e))
	}
}

// FnOnlyOnce executes op only once permanently.
func FnOnlyOnce(op Operation) Operation {
	once := false
	return func(e error) error {
		if once {
			return e
		}

		once = true
		return op(e)
	}
}
