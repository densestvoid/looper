// Package looper defines looping functions and interfaces to be looped
package looper

import (
	"context"
	"time"
)

func Loop(ctx context.Context, looper Looper) error {
	switch l := looper.(type) {
	case ConditionalIntervalLooper:
		return loopOnConditionalInterval(ctx, l)
	case IntervalLooper:
		return loopOnInterval(ctx, l)
	case ConditionalLooper:
		return loopOnCondition(ctx, l)
	default:
		return loop(ctx, looper)
	}
}

func loop(ctx context.Context, looper Looper) error {
	for {
		// Attempt to not dominate the go scheduler
		ch := make(chan struct{})
		close(ch)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ch:
			if err := looper.Function(); err != nil {
				if looper.HandleError(err) {
					return err
				}
			}
		}
	}
}

func loopOnInterval(ctx context.Context, looper IntervalLooper) error {
	ticker := time.NewTicker(looper.Interval())
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := looper.Function(); err != nil {
				if looper.HandleError(err) {
					return err
				}
			}
		}
	}
}

func loopOnCondition(ctx context.Context, looper ConditionalLooper) error {
	for {
		// Attempt to not dominate the go scheduler
		ch := make(chan struct{})
		close(ch)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ch:
			if !looper.Condition() {
				continue
			}

			if err := looper.Function(); err != nil {
				if looper.HandleError(err) {
					return err
				}
			}
		}
	}
}
