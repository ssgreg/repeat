package repeat

import (
	"math/rand"
	"time"
)

// FixedBackoffAlgorithm implements backoff with a fixed delay.
func FixedBackoffAlgorithm(delay time.Duration) func() time.Duration {
	return func() time.Duration {
		return delay
	}
}

// FixedBackoffBuilder is an option builder.
type FixedBackoffBuilder struct {
	// Delay specifyes fixed delay value.
	Delay time.Duration
}

// Set creates a Delay' option.
func (s *FixedBackoffBuilder) Set() func(*DelayOptions) {
	return func(do *DelayOptions) {
		do.Backoff = FixedBackoffAlgorithm(s.Delay)
	}
}

// FixedBackoff create a builder for Delay's option.
func FixedBackoff(delay time.Duration) *FixedBackoffBuilder {
	return &FixedBackoffBuilder{Delay: delay}
}

// FullJitterBackoffAlgorithm implements caped exponential backoff
// with jitter. Algorithm is fast because it does not use floating
// point arithmetics.
//
// Details:
// https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
//
// Example (BaseDelay=1, maxDelay=30):
// Call			Delay
// -------      ----------------
// 1            random [0...1]
// 2            random [0...2]
// 3            random [0...4]
// 4            random [0...8]
// 5            random [0...16]
// 6            random [0...30]
// 7            random [0...30]
//
func FullJitterBackoffAlgorithm(baseDelay time.Duration, maxDelay time.Duration) func() time.Duration {
	rnd := rand.New(rand.NewSource(int64(time.Now().Unix())))
	delay := baseDelay

	return func() time.Duration {
		defer func() {
			delay = delay << 1
			if delay > maxDelay {
				delay = maxDelay
			}
		}()
		return time.Duration(rnd.Int63n(int64(delay)))
	}
}

// FullJitterBackoffBuilder is an option builder.
type FullJitterBackoffBuilder struct {
	// MaxDelay specifies maximum value of a delay calculated by the
	// algorithm.
	//
	// Default value is maximum time.Duration value.
	MaxDelay time.Duration

	// BaseDelay specifies base of an exponent.
	BaseDelay time.Duration
}

// WithMaxDelay allows to set MaxDelay.
//
// MaxDelay specifies the maximum value of a delay calculated by the
// algorithm.
//
// Default value is maximum time.Duration value.
func (s *FullJitterBackoffBuilder) WithMaxDelay(d time.Duration) *FullJitterBackoffBuilder {
	s.MaxDelay = d
	return s
}

// WithBaseDelay allows to set BaseDelay.
//
// BaseDelay specifies base of an exponent.
func (s *FullJitterBackoffBuilder) WithBaseDelay(d time.Duration) *FullJitterBackoffBuilder {
	s.BaseDelay = d
	return s
}

// Set creates a Delay' option.
func (s *FullJitterBackoffBuilder) Set() func(*DelayOptions) {
	return func(do *DelayOptions) {
		do.Backoff = FullJitterBackoffAlgorithm(s.BaseDelay, s.MaxDelay)
	}
}

// FullJitterBackoff create a builder for Delay's option.
func FullJitterBackoff(baseDelay time.Duration) *FullJitterBackoffBuilder {
	return (&FullJitterBackoffBuilder{}).
		WithBaseDelay(baseDelay).
		WithMaxDelay(1<<63 - 1)
}

// ExponentialBackoffAlgorithm implements classic caped exponential backoff.
//
// Example (initialDelay=1, maxDelay=30, Multiplier=2, Jitter=0.5):
// Attempt		Delay
// -------      --------------------------
// 0             1 + random [-0.5...0.5]
// 1             2 + random [-1...1]
// 2             4 + random [-2...2]
// 3             8 + random [-4...4]
// 4            16 + random [-8...8]
// 5            32 + random [-16...16]
// 6            64 + random [-32...32] = 30
//
func ExponentialBackoffAlgorithm(initialDelay time.Duration, maxDelay time.Duration, multiplier float64, jitter float64) func() time.Duration {
	rnd := rand.New(rand.NewSource(int64(time.Now().Unix())))
	nextDelay := float64(initialDelay)
	limit := float64(maxDelay)

	return func() time.Duration {
		delay := nextDelay
		nextDelay = nextDelay * multiplier

		// Fix delay according to jitter.
		delta := delay * jitter
		delay = delay - delta + (2 * delta * rnd.Float64())

		// Fix delay limits.
		if delay > limit {
			delay = limit
		}

		return time.Duration(delay)
	}
}

// ExponentialBackoffBuilder is an option builder.
type ExponentialBackoffBuilder struct {
	// MaxDelay specifies maximum value of a delay calculated by the
	// algorithm.
	//
	// Default value is maximum time.Duration value.
	MaxDelay time.Duration

	// InitialDelay specifies an initial delay for the algorithm.
	//
	// Default value is equal to 1 second.
	InitialDelay time.Duration

	// Miltiplier specifies a multiplier for the last calculated
	// or specified delay.
	//
	// Default value is 2.
	Multiplier float64

	// Jitter specifies randomization factor [0..1].
	//
	// Default value is 0.
	Jitter float64

	nextDelay float64
	maxDelay  float64
	rnd       *rand.Rand
}

// WithMaxDelay allows to set MaxDelay.
//
// MaxDelay specifies the maximum value of a delay calculated by the
// algorithm.
//
// Default value is maximum time.Duration value.
func (s *ExponentialBackoffBuilder) WithMaxDelay(d time.Duration) *ExponentialBackoffBuilder {
	s.MaxDelay = d
	return s
}

// WithInitialDelay allows to set InitialDelay.
//
// InitialDelay specifies an initial delay for the algorithm.
func (s *ExponentialBackoffBuilder) WithInitialDelay(d time.Duration) *ExponentialBackoffBuilder {
	s.InitialDelay = d
	return s
}

// WithMultiplier allows to set Multiplier.
//
// Miltiplier specifies a multiplier for the last calculated
//
// Default value is 2.
func (s *ExponentialBackoffBuilder) WithMultiplier(m float64) *ExponentialBackoffBuilder {
	s.Multiplier = m
	return s
}

// WithJitter allows to set Jitter.
//
// Jitter specifies randomization factor [0..1].
//
// Default value is 0.
func (s *ExponentialBackoffBuilder) WithJitter(j float64) *ExponentialBackoffBuilder {
	s.Jitter = j
	return s
}

// Set creates a Delay' option.
func (s *ExponentialBackoffBuilder) Set() func(*DelayOptions) {
	return func(do *DelayOptions) {
		do.Backoff = ExponentialBackoffAlgorithm(s.InitialDelay, s.MaxDelay, s.Multiplier, s.Jitter)
	}
}

// ExponentialBackoff create a builder for Delay's option.
func ExponentialBackoff(initialDelay time.Duration) *ExponentialBackoffBuilder {
	return (&ExponentialBackoffBuilder{}).
		WithInitialDelay(initialDelay).
		WithMaxDelay(1<<63 - 1).
		WithMultiplier(2).
		WithJitter(0)
}
