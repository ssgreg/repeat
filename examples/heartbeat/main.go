package main

import (
	"context"
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

// Output:
//
// Attempt #0, Delay 0s
// Attempt #1, Delay 1.001129426s
// Attempt #2, Delay 1.000155727s
// Attempt #3, Delay 1.001131014s
// Attempt #4, Delay 1.000500428s
// Attempt #5, Delay 1.0008985s
// Attempt #6, Delay 1.000417057s
// Repetition process is finished with: context canceled
//

func main() {

	// An example operation that do heartbeat.
	var last time.Time
	op := func(c int) error {
		printInfo(c, &last)
		return nil
	}

	// A context with cancel.
	// Repetition will be cancelled in 7 seconds.
	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		time.Sleep(7 * time.Second)
		cancelFunc()
	}()

	err := repeat.Repeat(
		// Heartbeating op.
		repeat.FnWithCounter(op),
		// Delay with fixed backoff and context.
		repeat.WithDelay(
			repeat.FixedBackoff(time.Second).Set(),
			repeat.SetContext(ctx),
		),
	)

	fmt.Printf("Repetition process is finished with: %v\n", err)
}
