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

// IsTemporary checks if passed error is TemporaryError.
func IsTemporary(e error) bool {
	switch e.(type) {
	case *TemporaryError:
		return true
	default:
		return false
	}
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

// IsStop checks if passed error is StopError.
func IsStop(e error) bool {
	switch e.(type) {
	case *StopError:
		return true
	default:
		return false
	}
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
