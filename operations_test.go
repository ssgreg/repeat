package repeat

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestNope(t *testing.T) {
	// Nope should be transparent for all kinds of errors.
	require.NoError(t, Nope(nil))
	require.EqualError(t, Nope(errors.New("kiwi")), "kiwi")
	require.EqualError(t, Nope(HintTemporary(errors.New("kiwi"))), "repeat.temporary: kiwi")
}

func TestFnNope(t *testing.T) {
	// Nope should be transparent for all kinds of errors.
	require.NoError(t, FnNope(Nope)(nil))
	require.EqualError(t, FnNope(Nope)(errors.New("kiwi")), "kiwi")
	require.EqualError(t, FnNope(Nope)(HintTemporary(errors.New("kiwi"))), "repeat.temporary: kiwi")
}

func TestDone(t *testing.T) {
	// Nope should be transparent for all kinds of errors.
	require.NoError(t, Done(nil))
	require.NoError(t, Done(errors.New("kiwi")))
	require.NoError(t, Done(HintTemporary(errors.New("kiwi"))))
}

func TestFnDone(t *testing.T) {
	// Nope should be transparent for all kinds of errors.
	require.NoError(t, FnDone(Nope)(nil))
	require.NoError(t, FnDone(Nope)(errors.New("kiwi")))
	require.NoError(t, FnDone(Nope)(HintTemporary(errors.New("kiwi"))))
}

func TestFnES(t *testing.T) {
	opNil := func(e error) {
	}
	opErr := func(e error) {
		require.EqualError(t, e, "kiwi")
	}
	require.NoError(t, FnES(opNil)(nil))
	require.EqualError(t, FnES(opErr)(errors.New("kiwi")), "kiwi")
}

func TestFnS(t *testing.T) {
	op := func() {
	}
	require.NoError(t, FnS(op)(nil))
	require.EqualError(t, FnS(op)(errors.New("kiwi")), "kiwi")
}

func TestFn(t *testing.T) {
	opNil := func() error {
		return nil
	}
	opErr := func() error {
		return errors.New("kiwi")
	}
	require.NoError(t, Fn(opNil)(nil))
	require.EqualError(t, Fn(opErr)(errors.New("apple")), "kiwi")
}

func TestLimitMaxTries(t *testing.T) {
	fn := LimitMaxTries(5)
	fn(nil)
	fn(nil)
	fn(nil)
	fn(nil)
	assert.True(t, fn(nil) == nil)
	assert.False(t, fn(nil) == nil)
}

func TestStopOnSuccess(t *testing.T) {
	fn := StopOnSuccess()
	assert.True(t, fn(fn(nil)) != nil)
	assert.False(t, fn(fn(errors.New("error"))) == nil)
}

func TestFnOnSuccess_CalledOnNil(t *testing.T) {
	require.EqualError(t, FnOnSuccess(func(e error) error {
		require.NoError(t, e)
		return errors.New("called")
	})(nil), "called")
}

func TestFnOnSuccess_NotCalledOnError(t *testing.T) {
	require.EqualError(t, FnOnSuccess(func(e error) error {
		require.Fail(t, "must not be called")
		return e
	})(errors.New("not called")), "not called")
}

func TestFnOnError_NotCalledOnNil(t *testing.T) {
	require.NoError(t, FnOnError(func(e error) error {
		require.Fail(t, "must not be called")
		return e
	})(nil))
}

func TestFnOnError_CalledOnError(t *testing.T) {
	require.EqualError(t, FnOnError(func(e error) error {
		require.EqualError(t, e, "calling error")
		return errors.New("called")
	})(errors.New("calling error")), "called")
}

func TestFnWithErrorAndCounter(t *testing.T) {
	cc := 0
	op := func(e error, c int) error {
		require.Equal(t, cc, c)
		return e
	}
	fn := FnWithErrorAndCounter(op)

	// tick
	require.NoError(t, fn(nil))
	// tick
	cc++
	require.NoError(t, fn(nil))
	// tick
	cc++
	require.EqualError(t, fn(errors.New("passed")), "passed")
}

func TestFnWithCounter(t *testing.T) {
	cc := 0
	op := func(c int) error {
		require.Equal(t, cc, c)
		return nil
	}
	fn := FnWithCounter(op)

	// tick
	require.NoError(t, fn(nil))
	// tick
	cc++
	require.NoError(t, fn(nil))
}

func TestFnHintTemporary(t *testing.T) {
	fn := FnHintTemporary(Nope)

	// No action on nil
	require.NoError(t, fn(nil))
	// No action on temporary error.
	te := fn(HintTemporary(nil))
	require.EqualError(t, te, "repeat.temporary")
	// No action on stop error.
	se := fn(HintStop(errors.New("stop")))
	require.EqualError(t, se, "repeat.stop: stop")

	// Hint common error as temporary.
	ce := fn(errors.New("common"))
	require.EqualError(t, ce, "repeat.temporary: common")
}

func TestFnHintStop(t *testing.T) {
	fn := FnHintStop(Nope)

	// Hint nil as StopError.
	require.EqualError(t, fn(nil), "repeat.stop")
	// No action on temporary error.
	te := fn(HintTemporary(nil))
	require.EqualError(t, te, "repeat.temporary")
	// No action on stop error.
	se := fn(HintStop(errors.New("stop")))
	require.EqualError(t, se, "repeat.stop: stop")
	// Hint common error as StopError.
	ce := fn(errors.New("common"))
	require.EqualError(t, ce, "repeat.stop: common")
}

func TestFnPanic(t *testing.T) {
	fn := FnPanic(Nope)

	// No action on nil
	require.NoError(t, fn(nil))
	// No action on temporary error.
	te := fn(HintTemporary(nil))
	require.EqualError(t, te, "repeat.temporary")
	// No action on stop error.
	se := fn(HintStop(errors.New("stop")))
	require.EqualError(t, se, "repeat.stop: stop")

	// Panic
	require.Panics(t, func() { fn(errors.New("common")) })
}

func TestFnOnlyOnce(t *testing.T) {
	c := 0
	op := func(e error) error {
		c++
		return e
	}

	fn := FnOnlyOnce(op)
	fn(nil)
	fn(nil)
	fn(nil)

	require.Equal(t, 1, c)
}
