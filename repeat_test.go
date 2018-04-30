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

func TestWrap(t *testing.T) {
	c := 0

	wr := func(op Operation) Operation {
		c++
		return func(e error) error {
			return op(e)
		}
	}

	assert.NoError(t, Wrap(wr).Compose(Nope, Nope)(nil))
	assert.Equal(t, 2, c)
}

func TestCpp_C_D(t *testing.T) {
	c := 0

	cd := func(e error) error {
		c++
		return e
	}

	assert.NoError(t, Cpp(cd, cd).Compose(Nope)(nil))
	assert.Equal(t, 2, c)
}

func TestCpp_C_NoD(t *testing.T) {
	c := 0

	cd := func(e error) error {
		c++
		return e
	}

	assert.EqualError(t, Cpp(cd, cd).Compose(Nope)(errGolden), errGolden.Error())
	assert.Equal(t, 1, c)
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

	assert.EqualError(t, Cpp(cd, cd).Compose(errOp)(nil), errGolden.Error())
	assert.Equal(t, 2, c)
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

	assert.Panics(t, func() {
		Cpp(cd, cd).Compose(errOp)(nil)
	})
	assert.Equal(t, 2, c)
}
