package repeat

// Operation is the type of function for repetition.
type Operation func(error) error

// Repeat repeat operations until one of them stops the repetition.
//
// It is guaranteed that the first op will be called at least once.
func Repeat(ops ...Operation) error {
	var err error

	op := ComposeSlice(ops)
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

// Compose binds all passed operations into a single one.
func Compose(ops ...Operation) Operation {
	return ComposeSlice(ops)
}

// ComposeSlice binds all passed operations into a single one.
func ComposeSlice(ops []Operation) Operation {
	return func(err error) error {
		for _, op := range ops {
			err = op(err)
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
	// Nothing to do here.
}
