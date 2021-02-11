// Package looper defines looping functions and interfaces to be looped
package looper

import (
	"context"
	"sync"
	"time"
)

// Loop calls Looper.Function until the context is done.
// Canceling the context allows the current call to finish
func Loop(ctx context.Context, l Looper) error {
	for {
		// Attempt to not dominate the go scheduler
		ch := make(chan struct{})
		close(ch)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ch:
			if err := l.Function(); err != nil {
				if l.HandleError(err) {
					return err
				}
			}
		}
	}
}

// LoopInBackground calls Background on the Loop.Function
func LoopInBackground(ctx context.Context, l Looper, wg *sync.WaitGroup) <-chan error {
	return Background(
		ctx,
		func() error {
			return Loop(ctx, l)
		},
		wg,
	)
}

// LoopOnInterval calls Looper.Function on ticker with
// duration Looper.Interval until the context is done.
// Canceling the context allows the current call to finish
func LoopOnInterval(ctx context.Context, l IntervalLooper) error {
	ticker := time.NewTicker(l.Interval())
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := l.Function(); err != nil {
				if l.HandleError(err) {
					return err
				}
			}
		}
	}
}

// LoopOnIntervalInBackground calls Background on LooopOnInterval.Function
func LoopOnIntervalInBackground(ctx context.Context, l IntervalLooper, wg *sync.WaitGroup) <-chan error {
	return Background(
		ctx,
		func() error {
			return LoopOnInterval(ctx, l)
		},
		wg,
	)
}

// LoopOnCondition calls Looper.Function when condition is satisfied until
// the context is done. Canceling the context allows the current call to finish
func LoopOnCondition(ctx context.Context, l ConditionalLooper) error {
	for {
		// Attempt to not dominate the go scheduler
		ch := make(chan struct{})
		close(ch)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ch:
			if !l.Condition() {
				continue
			} else if err := l.Function(); err != nil {
				if l.HandleError(err) {
					return err
				}
			}
		}
	}
}

// LoopOnConditionInBackground calls Background on LoopOnCondition.Function
func LoopOnConditionInBackground(ctx context.Context, l ConditionalLooper, wg *sync.WaitGroup) <-chan error {
	return Background(
		ctx,
		func() error {
			return LoopOnCondition(ctx, l)
		},
		wg,
	)
}

// LoopOnConditionalInterval calls Looper.Function on ticker with
// duration Looper.Interval when condition is satisfied until the
// context is done. Canceling the context allows the current call to finish
func LoopOnConditionalInterval(ctx context.Context, l ConditionalIntervalLooper) error {
	ticker := time.NewTicker(l.Interval())
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if !l.Condition() {
				continue
			}
			if err := l.Function(); err != nil {
				if l.HandleError(err) {
					return err
				}
			}
		}
	}
}

// LoopOnConditionalIntervalInBackground calls Background on LoopOnConditionalInterval.Function
func LoopOnConditionalIntervalInBackground(ctx context.Context, l IntervalLooper, wg *sync.WaitGroup) <-chan error {
	return Background(
		ctx,
		func() error {
			return LoopOnInterval(ctx, l)
		},
		wg,
	)
}

// Background calls function in a go routine and returns the error on the returned channel
func Background(ctx context.Context, function Function, wg *sync.WaitGroup) <-chan error {
	if wg != nil {
		wg.Add(2)
	}

	funcErrChan := make(chan error, 1)

	go func() {
		if wg != nil {
			defer wg.Done()
		}

		defer close(funcErrChan)
		funcErrChan <- function()
	}()

	errChan := make(chan error)

	go func() {
		if wg != nil {
			defer wg.Done()
		}

		defer close(errChan)
		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
			<-funcErrChan
		case err := <-funcErrChan:
			errChan <- err
		}
	}()

	return errChan
}
