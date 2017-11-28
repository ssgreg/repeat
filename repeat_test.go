package repeat

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testOperation(max int, result error) Operation {
	c := 0
	return func(e error) error {
		defer func() { c++ }()
		if c >= max && result == nil {
			return &StopError{e}
		}
		if c >= max {
			fmt.Println(result)
			return result
		}
		return e
	}
}

func TestRepeat(t *testing.T) {
	op10 := func(e error, c int) error {
		defer func() { c++ }()
		if c == 0 {
			assert.NoError(t, e)
		} else {
			assert.Error(t, e, "temporary")
		}
		if c >= 10 {
			return nil
		}
		return &TemporaryError{errors.New("temporary")}
	}
	assert.NoError(t, Repeat(FnWithErrorAndCounter(op10), testOperation(10, nil)))
	assert.EqualError(t, Repeat(FnWithErrorAndCounter(op10), testOperation(5, nil)), "temporary")
	assert.EqualError(t, Repeat(FnWithErrorAndCounter(op10), testOperation(5, errors.New("external"))), "external")
}

func TestCompose(t *testing.T) {
	op10 := func(e error, c int) error {
		defer func() { c++ }()
		if c == 0 {
			assert.NoError(t, e)
		} else {
			assert.Error(t, e, "temporary")
		}
		if c >= 10 {
			return nil
		}
		return &TemporaryError{errors.New("temporary")}
	}
	assert.NoError(t, Repeat(Compose(FnWithErrorAndCounter(op10), testOperation(10, nil))))
	assert.EqualError(t, Repeat(Compose(FnWithErrorAndCounter(op10), testOperation(5, nil))), "temporary")
	assert.EqualError(t, Repeat(Compose(FnWithErrorAndCounter(op10), testOperation(5, errors.New("external")))), "external")
}
