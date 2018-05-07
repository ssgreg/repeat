package repeat

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompose_NilInNilOut(t *testing.T) {
	require.NoError(t, Compose(func(e error) error {
		require.NoError(t, e)
		return e
	})(nil))
}

func TestCompose_ErrInErrOut(t *testing.T) {
	require.EqualError(t, Compose(func(e error) error {
		require.EqualError(t, e, "oil")
		return e
	})(errors.New("oil")), "oil")
}

func TestCompose_TemporaryErrInTemporaryErrOut(t *testing.T) {
	require.EqualError(t, Compose(func(e error) error {
		require.EqualError(t, e, "repeat.temporary: cat")
		return HintTemporary(e)
	})(HintTemporary(errors.New("cat"))), "repeat.temporary: cat")
}

func TestCompose_ErrInStopErrOut(t *testing.T) {
	require.EqualError(t, Compose(func(e error) error {
		return HintStop(e)
	})(errors.New("bob")), "repeat.stop: bob")
}

func TestCompose_NillOverridesTemporaryError(t *testing.T) {
	require.NoError(t, Compose(
		func(e error) error {
			return HintTemporary(e)
		},
		Done,
	)(errors.New("ann")))
}

func TestCompose_OverrideTemporaryError(t *testing.T) {
	require.EqualError(t, Compose(
		func(e error) error {
			return HintTemporary(e)
		},
		func(e error) error {
			return HintTemporary(errors.New("pong"))
		},
	)(errors.New("ping")), "repeat.temporary: pong")
}

func TestOnce_NonTemporaryErrOut(t *testing.T) {
	require.EqualError(t, Once(func(e error) error {
		return HintTemporary(errors.New("zed"))
	}), "zed")
}

func TestOnce_NonStopErrOut(t *testing.T) {
	require.NoError(t, Once(func(e error) error {
		return HintStop(nil)
	}))
}

func TestOnce_ErrOut(t *testing.T) {
	require.EqualError(t, Once(func(e error) error {
		return errors.New("aim")
	}), "aim")
}

func TestFnRepeat_NilInStopErrWithNilOut(t *testing.T) {
	require.EqualError(t, FnRepeat(
		func(e error) error {
			require.NoError(t, e)
			return e
		},
		StopOnSuccess(),
	)(nil), "repeat.stop")
}

func TestFnRepeat_ErrInErrOut(t *testing.T) {
	require.EqualError(t, FnRepeat(
		func(e error) error {
			require.EqualError(t, e, "oil")
			return e
		},
	)(errors.New("oil")), "oil")
}

func TestFnRepeat_TemporaryErrInSameStopErrOut(t *testing.T) {
	require.EqualError(t, FnRepeat(
		LimitMaxTries(1),
		func(e error) error {
			require.EqualError(t, e, "repeat.temporary: cat")
			return HintTemporary(e)
		},
	)(HintTemporary(errors.New("cat"))), "repeat.stop: cat")
}

func TestFnRepeat_ErrInStopErrOut(t *testing.T) {
	require.EqualError(t, FnRepeat(func(e error) error {
		return HintStop(e)
	})(errors.New("bob")), "repeat.stop: bob")
}

func TestRepeat_WithNoErrors(t *testing.T) {
	cn := 0
	require.NoError(t, Repeat(
		LimitMaxTries(3),
		FnWithErrorAndCounter(func(e error, c int) error {
			defer func() { cn++ }()
			require.Equal(t, cn, c, "should be equal on every")

			switch {
			case c < 3:
				require.NoError(t, e, "no error on every call")
			default:
				require.Fail(t, "cant be here, only three tries")
			}

			return nil
		}),
	))
	require.Equal(t, 3, cn)
}

func TestRepeat_WithTemporaryErrors(t *testing.T) {
	cn := 0
	require.EqualError(t, Repeat(
		LimitMaxTries(3),
		// Should be called three times until LimitMaxTries stops the execution.
		FnWithErrorAndCounter(func(e error, c int) error {
			defer func() { cn++ }()
			require.Equal(t, cn, c, "should be equal on every")

			switch c {
			case 0:
				require.NoError(t, e,
					"no error on first call (started with no error)")
			case 1, 2:
				require.EqualError(t, e, "repeat.temporary: my temporary",
					"the same error func returns")
			default:
				require.Fail(t, "cant be here")
			}

			return HintTemporary(errors.New("my temporary"))
		}),
	), "my temporary")
	require.Equal(t, 3, cn)
}

func TestRepeat_WithErrors(t *testing.T) {
	cn := 0
	require.EqualError(t, Repeat(
		FnWithErrorAndCounter(func(e error, c int) error {
			defer func() { cn++ }()
			require.Equal(t, cn, c, "should be equal on every")

			if c == 2 {
				return errors.New("my real")
			}

			return HintTemporary(errors.New("my temporary"))
		}),
	), "my real")
	require.Equal(t, 3, cn)
}

func TestRepeat_WithStopErrors(t *testing.T) {
	cn := 0
	require.EqualError(t, Repeat(
		FnWithErrorAndCounter(func(e error, c int) error {
			defer func() { cn++ }()
			require.Equal(t, cn, c, "should be equal on every")

			if c == 2 {
				return HintStop(errors.New("my real"))
			}

			return HintTemporary(errors.New("my temporary"))
		}),
	), "my real")
	require.Equal(t, 3, cn)
}

func TestWrap(t *testing.T) {
	c := 0

	wr := func(op Operation) Operation {
		c++
		return func(e error) error {
			return op(e)
		}
	}

	require.NoError(t, Wrap(wr).Compose(Nope, Nope)(nil))
	require.Equal(t, 2, c, "wr called two times according to number of ops in Compose")
}

func TestCpp_C_D(t *testing.T) {
	c := 0

	cd := func(e error) error {
		c++
		return e
	}

	require.NoError(t, Cpp(cd, cd).Compose(Nope)(nil))
	require.Equal(t, 2, c)
}

func TestCpp_C_NoD(t *testing.T) {
	c := 0

	cd := func(e error) error {
		c++
		return e
	}

	require.EqualError(t, Cpp(cd, cd).Compose(Nope)(errGolden), errGolden.Error())
	require.Equal(t, 1, c)
}

func TestCpp_C_ErrOP_D(t *testing.T) {
	c := 0

	cd := func(e error) error {
		c++
		return e
	}

	errOp := func(error) error {
		return errGolden
	}

	require.EqualError(t, Cpp(cd, FnDone(cd)).Compose(errOp)(nil), errGolden.Error())
	require.Equal(t, 2, c)
}

func TestCpp_C_PanicOP_D(t *testing.T) {
	c := 0

	cd := func(e error) error {
		c++
		return e
	}

	errOp := func(error) error {
		panic(errGolden)
	}

	require.Panics(t, func() {
		Cpp(cd, cd).Compose(errOp)(nil)
	})
	require.Equal(t, 2, c)
}

func TestCpp_C_Op_ErrorD(t *testing.T) {
	c := 0

	cc := func(e error) error {
		c++
		return e
	}

	dd := func(error) error {
		c++
		return errGolden
	}

	require.Panics(t, func() {
		Cpp(cc, dd).Compose(Nope)(nil)
	})
	require.Equal(t, 2, c)
}

func TestCpp_TransparentC(t *testing.T) {
	c := 0

	cc := func(e error) error {
		c++
		require.EqualError(t, Cause(e), errGolden.Error())
		return nil
	}

	dd := func(e error) error {
		c++
		require.NoError(t, Cause(e))
		return nil
	}

	op := func(e error) error {
		c++
		require.EqualError(t, Cause(e), errGolden.Error())
		return nil
	}

	require.NoError(t, Cause(
		Cpp(cc, dd).Compose(op)(HintTemporary(errGolden)),
	))
	require.Equal(t, 3, c)
}

func TestRepeatWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	require.EqualError(t, WithContext(ctx).Once(Nope), "context canceled")
}
