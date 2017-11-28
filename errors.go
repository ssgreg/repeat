package repeat

// TemporaryError allows not to stop repetitions process right now.
//
// This error never returns to the caller as is, only wrapped error.
type TemporaryError struct {
	Cause error
}

func (e *TemporaryError) Error() string {
	return e.Cause.Error()
}

// HintTemporary makes a TemporaryError.
func HintTemporary(e error) error {
	return &TemporaryError{e}
}

// StopError allows to stop repetition process without specifying a
// separate error.
//
// This error never returns to the caller as is, only wrapped error.
type StopError struct {
	Cause error
}

func (e *StopError) Error() string {
	return e.Cause.Error()
}

// HintStop makes a StopError.
func HintStop(e error) error {
	return &StopError{e}
}

// Cause extracts the cause error from TemporaryError and StopError
// or return the passed one.
func Cause(err error) error {
	switch e := err.(type) {
	case *TemporaryError:
		return e.Cause
	case *StopError:
		return e.Cause
	default:
		return err
	}
}
