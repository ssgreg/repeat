package repeat

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

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

func TestFnHintTemporary(t *testing.T) {
	op := func(e error) error {
		return e
	}
	fn := FnHintTemporary(op)

	// No action on nil
	require.NoError(t, fn(nil))
	// No action on temporary error.
	te := fn(HintTemporary(nil))
	require.True(t, IsTemporary(te))
	require.NoError(t, Cause(te))
	// No action on stop error.
	se := fn(HintStop(errors.New("stop")))
	require.True(t, IsStop(se))
	require.EqualError(t, Cause(se), "stop")

	// Hints common error as temporary.
	ce := fn(errors.New("common"))
	require.True(t, IsTemporary(ce))
	require.EqualError(t, Cause(ce), "common")
}
