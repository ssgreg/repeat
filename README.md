# Repeat

![Alt text](https://user-images.githubusercontent.com/1574981/33310621-48c78416-d433-11e7-9a80-36d2381901d0.png "repeat")
[![GoDoc](https://godoc.org/github.com/ssgreg/repeat?status.svg)](https://godoc.org/github.com/ssgreg/repeat)
[![Build Status](https://travis-ci.org/ssgreg/repeat.svg?branch=master)](https://travis-ci.org/ssgreg/repeat)
[![Go Report Status](https://goreportcard.com/badge/github.com/ssgreg/repeat)](https://goreportcard.com/report/github.com/ssgreg/repeat)
[![GoCover](https://gocover.io/_badge/github.com/ssgreg/repeat)](https://gocover.io/github.com/ssgreg/repeat)

Go implementation of different backoff strategies useful for retrying operations and heartbeating.

## Examples

### Backoff

Let's imagine that we need to do a REST call on remote server but it could fail with a bunch of different issues. We can repeat failed operation using exponential backoff policies.

> *Exponential backoff* is an algorithm that uses feedback to multiplicatively decrease the rate of some process, in order to gradually find an acceptable rate.

The example below tries to repeat operation 10 times using a full jitter backoff. [See algorithm details here.](https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/)

```go
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
```

The example of output:

```
Attempt #0, Delay 0s
Attempt #1, Delay 373.617912ms
Attempt #2, Delay 668.004225ms
Attempt #3, Delay 1.220076558s
Attempt #4, Delay 2.716156336s
Attempt #5, Delay 6.458431017s
Repetition process is finished with: <nil>
```

### Backoff with timeout

The example below is almost the same as the previous one. It adds one important feature - possibility to cancel operation repetition using context's timeout.

```go
    // A context with cancel.
    // Repetition will be cancelled in 3 seconds.
    ctx, cancelFunc := context.WithCancel(context.Background())
    go func() {
        time.Sleep(3 * time.Second)
        cancelFunc()
    }()

    // Repeat op on any error, with 10 retries, with a backoff.
    err := repeat.Repeat(
        ...
        // Specify a delay that uses a backoff.
        repeat.WithDelay(
            repeat.FullJitterBackoff(500*time.Millisecond).Set(),
            repeat.SetContext(ctx),
        ),
        ...
    )
```

The example of output:

```
Attempt #0, Delay 0s
Attempt #1, Delay 358.728046ms
Attempt #2, Delay 845.361787ms
Attempt #3, Delay 61.527485ms
Repetition process is finished with: context canceled
```

### Heartbeating

Let's imagine we need to periodically report execution progress to remote server. The example below repeats the operation each second until it will be cancelled using passed context.

```go
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
```

The example of output:

```
Attempt #0, Delay 0s
Attempt #1, Delay 1.001129426s
Attempt #2, Delay 1.000155727s
Attempt #3, Delay 1.001131014s
Attempt #4, Delay 1.000500428s
Attempt #5, Delay 1.0008985s
Attempt #6, Delay 1.000417057s
Repetition process is finished with: context canceled
```

### Heartbeating with error timeout

The example below is almost the same as the previous one but it will be cancelled using special error timeout. This timeout resets each time the operations return nil.

```go
    // An example operation that do heartbeat.
    // It fails 5 times after 3 successfull tries.
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
```

The example of output:

```
Attempt #0, Delay 0s
Attempt #1, Delay 1.001634616s
Attempt #2, Delay 1.004912408s
Attempt #3, Delay 1.001021358s
Attempt #4, Delay 1.001249459s
Attempt #5, Delay 1.004320833s
Repetition process is finished with: can't connect to a server
```
