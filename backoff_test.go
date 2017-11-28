package repeat

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConstantBackoff(t *testing.T) {
	fn77 := FixedBackoffAlgorithm(77)
	assert.EqualValues(t, fn77(), 77)
	assert.EqualValues(t, fn77(), 77)
}

func TestFullJitterBackoffDefaults(t *testing.T) {
	do := &DelayOptions{}
	FullJitterBackoff(time.Second).Set()(do)

	for i := 0; i < 20; i++ {
		c := int64(math.Pow(2, float64(i)))
		InRange(t, do.Backoff(), 0, time.Duration(c)*time.Second)
	}
}

func TestFullJitterBackoff(t *testing.T) {
	do := &DelayOptions{}
	FullJitterBackoff(1).WithMaxDelay(30).Set()(do)

	for i := 0; i < 50; i++ {
		c := int64(math.Pow(2, float64(i)))
		if c > 30 {
			c = 30
		}
		InRange(t, do.Backoff(), 0, time.Duration(c))
	}
}

var floatSecond = float64(time.Second)

func TestExponentialBackoffDefaults(t *testing.T) {
	do := &DelayOptions{}
	ExponentialBackoff(time.Second).Set()(do)

	for i := 0; i < 30; i++ {
		c := math.Pow(2, float64(i))
		InRange(t, do.Backoff(), time.Duration(c*floatSecond), time.Duration(c*floatSecond))
	}
}

func TestExponentialBackoffJitter(t *testing.T) {
	do := &DelayOptions{}
	ExponentialBackoff(time.Second).WithJitter(.5).Set()(do)

	for i := 0; i < 30; i++ {
		c := math.Pow(2, float64(i))
		fi := .5 * c
		InRange(t, do.Backoff(), time.Duration((c-fi)*floatSecond), time.Duration((c+fi)*floatSecond))
	}
}

func TestExponentialBackoffJitterAndMultiplier(t *testing.T) {
	do := &DelayOptions{}
	ExponentialBackoff(time.Second).WithJitter(.1).WithMultiplier(1.74).Set()(do)

	for i := 0; i < 30; i++ {
		c := math.Pow(1.74, float64(i))
		fi := .1 * c
		InRange(t, do.Backoff(), time.Duration((c-fi)*floatSecond), time.Duration((c+fi)*floatSecond))
	}
}

func TestExponentialBackoff(t *testing.T) {
	do := &DelayOptions{}
	ExponentialBackoff(354 * time.Millisecond).WithJitter(.9).WithMultiplier(1.12).WithMaxDelay(5 * time.Second).Set()(do)

	initDelay := float64(354 * time.Millisecond)
	for i := 0; i < 30; i++ {
		c := math.Pow(1.12, float64(i))
		fi := .9 * c
		max := time.Duration((c + fi) * initDelay)
		if max > 5*time.Second {
			max = 5 * time.Second
		}
		InRange(t, do.Backoff(), time.Duration((c-fi)*initDelay), max)
	}
}
