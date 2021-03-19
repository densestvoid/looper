// Package looper defines looping functions and interfaces to be looped
package looper

import (
	"context"
	"runtime"
	"time"
)

// Function signature of function that gets called in a loop
type Function func() error

// HandleError signature of function that handles the Function error.
// Returns true if error should cause the loop to exit
type HandleError func(error) bool

// Condition returns true if the requirements to run the function are met
type Condition func() bool

// Looper defines types that should constantly run in a loop
type Looper interface {
	Function() error
	HandleError(error) bool
}

// IntervalLooper defines types the should run in a loop on a ticker
type IntervalLooper interface {
	Looper
	Interval() time.Duration
}

// ConditionalLooper defines types the should run in a loop
// when a condition is satisfied
type ConditionalLooper interface {
	Looper
	Condition() bool
}

func watchConditional(ctx context.Context, condition Condition) <-chan struct{} {
	match := make(chan struct{})

	go func() {
		for {
			runtime.Gosched()

			select {
			case <-ctx.Done():
				return
			default:
			}

			if !condition() {
				continue
			}

			select {
			case match <- struct{}{}:
			default:
			}
		}
	}()

	return match
}

// ChannelLooper defines types that should run in a loop
// whenever a channel can be read
type ChannelLooper interface {
	Looper
	Channel() <-chan struct{}
}

// NewLooper creates a new LooperInterval using the arguments
// in the interface functions calls
func NewLooper(function Function, handleErr HandleError) Looper {
	return &looper{
		function:  function,
		handleErr: handleErr,
	}
}

// NewIntervalLooper creates a new LooperInterval using the arguments
// in the interface functions calls
func NewIntervalLooper(function Function, handleErr HandleError, interval time.Duration) IntervalLooper {
	return &looper{
		function:  function,
		handleErr: handleErr,
		interval:  interval,
	}
}

// NewConditionalLooper creates a new LooperInterval using the arguments
// in the interface functions calls
func NewConditionalLooper(function Function, handleErr HandleError, condition Condition) ConditionalLooper {
	return &looper{
		function:  function,
		handleErr: handleErr,
		condition: condition,
	}
}

// NewConditionalIntervalLooper creates a new LooperInterval using
// the arguments in the interface functions calls
func NewConditionalIntervalLooper(
	function Function,
	handleErr HandleError,
	interval time.Duration,
	condition Condition,
) ConditionalIntervalLooper {
	return &looper{
		function:  function,
		handleErr: handleErr,
		interval:  interval,
		condition: condition,
	}
}

type looper struct {
	function  Function
	handleErr HandleError
	interval  time.Duration
	condition Condition
}

func (l *looper) Function() error {
	return l.function()
}

func (l *looper) HandleError(err error) bool {
	return l.handleErr(err)
}

func (l *looper) Interval() time.Duration {
	return l.interval
}

func (l *looper) Condition() bool {
	return l.condition()
}
