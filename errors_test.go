package repeat

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStopError(t *testing.T) {
	e := HintStop(nil)
	require.False(t, IsStop(nil))
	require.False(t, IsStop(errors.New("test")))
	require.False(t, IsStop(HintTemporary(nil)))
	require.True(t, IsStop(e))
	require.Nil(t, Cause(e))

	require.EqualError(t, e, "repeat.stop")
	require.EqualError(t, HintStop(errors.New("internal")), "repeat.stop: internal")
}

func TestTemporaryError(t *testing.T) {
	e := HintTemporary(nil)
	require.False(t, IsTemporary(nil))
	require.False(t, IsTemporary(errors.New("test")))
	require.False(t, IsTemporary(HintStop(nil)))
	require.True(t, IsTemporary(e))
	require.Nil(t, Cause(e))

	require.EqualError(t, e, "repeat.temporary")
	require.EqualError(t, HintTemporary(errors.New("internal")), "repeat.temporary: internal")
}
