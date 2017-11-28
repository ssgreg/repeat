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

// Output:
//
// Attempt #0, Delay 0s
// Attempt #1, Delay 1.001634616s
// Attempt #2, Delay 1.004912408s
// Attempt #3, Delay 1.001021358s
// Attempt #4, Delay 1.001249459s
// Attempt #5, Delay 1.004320833s
// Repetition process is finished with: can't connect to a server
//

func main() {

	// An example operation that do heartbeat.
	// It fails 5 times after 3 successful tries.
	var last time.Time
	op := func(c int) error {
		printInfo(c, &last)
		if c > 3 && c < 8 {
			return repeat.HintTemporary(errors.New("can't connect to a server"))
		}
		return nil
	}

	err := repeat.Repeat(
		// Heartbeating op.
		repeat.FnWithCounter(op),
		// Delay with fixed backoff and error timeout.
		repeat.WithDelay(
			repeat.FixedBackoff(time.Second).Set(),
			repeat.SetErrorsTimeout(3*time.Second),
		),
	)

	fmt.Printf("Repetition process is finished with: %v\n", err)
}
