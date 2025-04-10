package app

import "time"

type TimerTask[T any] struct {
	ticker   *time.Ticker
	quit     chan struct{}
	interval time.Duration
	task     func(T)
	arg      T
}

func NewTimerTask[T any](interval time.Duration, task func(T), arg T) *TimerTask[T] {
	return &TimerTask[T]{
		ticker:   time.NewTicker(interval),
		quit:     make(chan struct{}),
		interval: interval,
		task:     task,
		arg:      arg,
	}
}

func (gt *TimerTask[T]) Start() {
	go func() {
		for {
			select {
			case <-gt.ticker.C:
				gt.task(gt.arg)
			case <-gt.quit:
				gt.ticker.Stop()
				return
			}
		}
	}()
}

func (gt *TimerTask[T]) Stop() {
	close(gt.quit)
}
