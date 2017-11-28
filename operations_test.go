package repeat

import (
	"errors"
	"testing"

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
