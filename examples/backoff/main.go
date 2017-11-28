package main

import (
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
// Attempt #1, Delay 373.617912ms
// Attempt #2, Delay 668.004225ms
// Attempt #3, Delay 1.220076558s
// Attempt #4, Delay 2.716156336s
// Attempt #5, Delay 6.458431017s
// Repetition process is finished with: <nil>
//

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

	// Repeat op on any error, with 10 retries, with a backoff.
	err := repeat.Repeat(
		// Our op with additional call counter.
		repeat.FnWithCounter(op),
		// Force the repetition to stop in case the previous operation
		// returns nil.
		repeat.StopOnSuccess(),
		// 10 retries max.
		repeat.LimitMaxTries(10),
		// Specify a delay that uses a backoff.
		repeat.WithDelay(
			repeat.FullJitterBackoff(500*time.Millisecond).Set(),
		),
	)

	fmt.Printf("Repetition process is finished with: %v\n", err)
}
