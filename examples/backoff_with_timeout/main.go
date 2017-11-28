package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ssgreg/repeat"
)

func printInfo(attempt int, last *time.Time) {
	tmp := *last
	*last = time.Now()
	if attempt == 0 {
		tmp = *last
	}
	fmt.Printf("Attempt #%d, Delay %v\n", attempt, last.Sub(tmp))
}

// Example of output:
//
// Attempt #0, Delay 0s
// Attempt #1, Delay 358.728046ms
// Attempt #2, Delay 845.361787ms
// Attempt #3, Delay 61.527485ms
// Repetition process is finished with: context canceled
//

func backoff(ctx context.Context) repeat.Operation {
	return repeat.Compose(
		// Force the repetition to stop in case the previous operation
		// returns nil.
		repeat.StopOnSuccess(),
		// 10 retries max.
		repeat.LimitMaxTries(10),
		// Specify a delay that uses a backoff.
		repeat.WithDelay(
			repeat.FullJitterBackoff(500*time.Millisecond).Set(),
			repeat.SetContext(ctx),
		),
	)
}

func main() {

	// An example operation that do some useful stuff.
	// It fails five first times.
	var last time.Time
	op := func(c int) error {
		printInfo(c, &last)
		if c < 5 {
			return repeat.HintTemporary(errors.New("can't connect to a server"))
		}
		return nil
	}

	// A context with cancel.
	// Repetition will be cancelled in 3 seconds.
	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		time.Sleep(3 * time.Second)
		cancelFunc()
	}()

	// Repeat op on any error, with 10 retries, with a backoff.
	err := repeat.Repeat(repeat.FnWithCounter(op), backoff(ctx))

	fmt.Printf("Repetition process is finished with: %v\n", err)
}
