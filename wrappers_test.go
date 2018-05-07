package repeat

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	errGolden = errors.New("golden")
)

func TestWrStopOnContextError_CallReturnOp(t *testing.T) {
	opNilError := func(e error) error {
		require.NoError(t, e)
		return e
	}

	require.NoError(t, WrStopOnContextError(context.Background())(opNilError)(nil))

	opError := func(e error) error {
		require.EqualError(t, Cause(e), errGolden.Error())
		return e
	}

	require.Error(t, WrStopOnContextError(context.Background())(opError)(HintTemporary(errGolden)))
}

func TestWrStopOnContextError_CancelOp(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	op := func(e error) error {
		require.Fail(t, "should be never called")
		return nil
	}

	cancel()
	require.EqualError(t, WrStopOnContextError(ctx)(op)(nil), "repeat.stop: context canceled")
}

func TestWrStopOnContextError_CancelOpWithError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	op := func(e error) error {
		require.Fail(t, "should be never called")
		return nil
	}

	cancel()
	require.EqualError(t, WrStopOnContextError(ctx)(op)(HintTemporary(errGolden)), "repeat.stop: golden")
	require.EqualError(t, WrStopOnContextError(ctx)(op)(HintStop(errGolden)), "repeat.stop: golden")
	require.EqualError(t, WrStopOnContextError(ctx)(op)(errGolden), "golden")
}

func TestWrStopOnContextError_SuccessIfDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	called := false
	op := func(e error) error {
		time.Sleep(time.Millisecond * 10)
		called = true
		return nil
	}

	go func() {
		time.Sleep(time.Millisecond * 5)
		cancel()
	}()

	require.NoError(t, WrStopOnContextError(ctx)(op)(HintTemporary(errGolden)))
	require.True(t, called)
}

func TestForward(t *testing.T) {
	op := func(e error) error { return nil }

	require.Nil(t, Forward(nil))
	require.Equal(t, reflect.ValueOf(op).Pointer(), reflect.ValueOf(op).Pointer())
}
