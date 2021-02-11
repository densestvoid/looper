// Package looper defines looping functions and interfaces to be looped
package looper

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

////////// Looper Suite //////////

type LooperSuite struct {
	suite.Suite
	looper Looper
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *LooperSuite) SetupSuite() {
	s.looper = NewLooper(
		func() error {
			time.Sleep(time.Millisecond)
			return nil
		},
		func(error) bool { return false },
	)
}

func (s *LooperSuite) SetupTest() {
	s.ctx, s.cancel = context.WithTimeout(context.Background(), 0)
}

func (s *LooperSuite) TearDownTest() {
	s.cancel()
}

func (s *LooperSuite) TearDownSuite() {}

func (s *LooperSuite) TestLoop() {
	// Setup

	// Verification
	s.Assert().IsType(
		context.DeadlineExceeded,
		Loop(
			s.ctx,
			s.looper,
		),
	)
}

func (s *LooperSuite) TestLoopInBackground() {
	// Setup

	// Verification
	s.Assert().IsType(
		context.DeadlineExceeded,
		<-LoopInBackground(
			s.ctx,
			s.looper,
			nil,
		),
	)
}

func TestLooperSuite(t *testing.T) {
	suite.Run(t, new(LooperSuite))
}

////////// Interval Looper Suite //////////

type IntervalLooperSuite struct {
	suite.Suite
	looper IntervalLooper
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *IntervalLooperSuite) SetupSuite() {
	s.looper = NewIntervalLooper(
		func() error {
			time.Sleep(time.Millisecond)
			return nil
		},
		func(error) bool { return false },
		time.Millisecond,
	)
}

func (s *IntervalLooperSuite) SetupTest() {
	s.ctx, s.cancel = context.WithTimeout(context.Background(), 0)
}

func (s *IntervalLooperSuite) TearDownTest() {
	s.cancel()
}

func (s *IntervalLooperSuite) TearDownSuite() {}

func (s *IntervalLooperSuite) TestLoopOnInterval() {
	// Setup

	// Verification
	s.Assert().IsType(
		context.DeadlineExceeded,
		LoopOnInterval(
			s.ctx,
			s.looper,
		),
	)
}

func (s *IntervalLooperSuite) TestLoopOnIntervalInBackground() {
	// Setup

	// Verification
	s.Assert().IsType(
		context.DeadlineExceeded,
		<-LoopOnIntervalInBackground(
			s.ctx,
			s.looper,
			nil,
		),
	)
}

func TestIntervalLooperSuite(t *testing.T) {
	suite.Run(t, new(IntervalLooperSuite))
}

////////// Conditional Looper Suite //////////

type ConditionalLooperSuite struct {
	suite.Suite
	looper ConditionalLooper
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *ConditionalLooperSuite) SetupSuite() {
	s.looper = NewConditionalLooper(
		func() error {
			time.Sleep(time.Millisecond)
			return nil
		},
		func(error) bool { return false },
		func() bool { return false },
	)
}

func (s *ConditionalLooperSuite) SetupTest() {
	s.ctx, s.cancel = context.WithTimeout(context.Background(), 0)
}

func (s *ConditionalLooperSuite) TearDownTest() {
	s.cancel()
}

func (s *ConditionalLooperSuite) TearDownSuite() {}

func (s *ConditionalLooperSuite) TestLoopOnCondition() {
	// Setup

	// Verification
	s.Assert().IsType(
		context.DeadlineExceeded,
		LoopOnCondition(
			s.ctx,
			s.looper,
		),
	)
}

func (s *ConditionalLooperSuite) TestLoopOnConditionInBackground() {
	// Setup

	// Verification
	s.Assert().IsType(
		context.DeadlineExceeded,
		<-LoopOnConditionInBackground(
			s.ctx,
			s.looper,
			nil,
		),
	)
}

func TestConditionalLooperSuite(t *testing.T) {
	suite.Run(t, new(ConditionalLooperSuite))
}

////////// Conditional Interval Looper Suite //////////

type ConditionalIntervalLooperSuite struct {
	suite.Suite
	looper ConditionalIntervalLooper
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *ConditionalIntervalLooperSuite) SetupSuite() {
	s.looper = NewConditionalIntervalLooper(
		func() error {
			time.Sleep(time.Millisecond)
			return nil
		},
		func(error) bool { return false },
		time.Millisecond,
		func() bool { return false },
	)
}

func (s *ConditionalIntervalLooperSuite) SetupTest() {
	s.ctx, s.cancel = context.WithTimeout(context.Background(), 0)
}

func (s *ConditionalIntervalLooperSuite) TearDownTest() {
	s.cancel()
}

func (s *ConditionalIntervalLooperSuite) TearDownSuite() {}

func (s *ConditionalIntervalLooperSuite) TestLoopOnConditionalInterval() {
	// Setup

	// Verification
	s.Assert().IsType(
		context.DeadlineExceeded,
		LoopOnConditionalInterval(
			s.ctx,
			s.looper,
		),
	)
}

func (s *ConditionalIntervalLooperSuite) TestLoopOnConditionalIntervalInBackground() {
	// Setup

	// Verification
	s.Assert().IsType(
		context.DeadlineExceeded,
		<-LoopOnConditionalIntervalInBackground(
			s.ctx,
			s.looper,
			nil,
		),
	)
}

func TestConditionalIntervalLooperSuite(t *testing.T) {
	suite.Run(t, new(ConditionalIntervalLooperSuite))
}

////////// Background Tests //////////

func TestBackground(t *testing.T) {
	// Setup
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()
	chErr := <-Background(
		ctx,
		func() error {
			time.Sleep(time.Millisecond)
			return nil
		},
		&wg,
	)
	wg.Wait()

	// Verification
	assert.IsType(
		t,
		context.DeadlineExceeded,
		chErr,
	)
}
