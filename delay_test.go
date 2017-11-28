package repeat

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// var backoff100 = &ConstantBackoff{100 * time.Millisecond}
// var backoff10 = &ConstantBackoff{10 * time.Millisecond}

func TestDelay(t *testing.T) {
	hb := WithDelay(FixedBackoff(10 * time.Millisecond).Set())

	for i := 0; i < 10; i++ {
		InRange(t, GetDelay(t, hb, nil), time.Millisecond*5, time.Millisecond*15)
	}
}

func TestDelayErrorTimeout(t *testing.T) {
	hb := WithDelay(FixedBackoff(10*time.Millisecond).Set(), SetErrorsTimeout(27*time.Millisecond))

	InRange(t, GetDelay(t, hb, nil), time.Millisecond*10, time.Millisecond*15)
	InRange(t, GetDelay(t, hb, nil), time.Millisecond*10, time.Millisecond*15)
	InRange(t, GetDelay(t, hb, nil), time.Millisecond*10, time.Millisecond*15)
	InRange(t, GetDelay(t, hb, errors.New("error")), time.Millisecond*10, time.Millisecond*15)
	InRange(t, GetDelay(t, hb, errors.New("error")), time.Millisecond, time.Millisecond*10)
}

func TestDelayCancel(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	hb := WithDelay(FixedBackoff(100*time.Millisecond).Set(), SetContext(ctx))
	InRange(t, GetDelay(t, hb, nil), time.Millisecond*50, time.Millisecond*150)

	go func() {
		time.Sleep(time.Millisecond * 50)
		cancelFunc()
	}()
	InRange(t, GetDelay(t, hb, errors.New("context canceled")), time.Millisecond*50, time.Millisecond*60)
}

func GetDelay(t *testing.T, o Operation, result error) time.Duration {
	start := time.Now()
	if result == nil {
		assert.NoError(t, o(nil))
	} else {
		assert.Error(t, o(result), result.Error())
	}
	return time.Now().Sub(start)
}

func InRange(t *testing.T, delay time.Duration, min time.Duration, max time.Duration) {
	fmt.Println(delay, min, max)
	require.True(t, delay >= min && delay <= max)
}
