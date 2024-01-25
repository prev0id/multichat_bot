package rate_limit

import (
	"sync"
	"time"
)

const windowDuration = 30 * time.Second

type Checker struct {
	m sync.Mutex

	rl checkerUnsafe
}

func NewChecker(limit int) *Checker {
	return &Checker{
		rl: checkerUnsafe{limit: limit},
	}
}

func (w *Checker) IsLimitExceeded() bool {
	w.m.Lock()
	defer w.m.Unlock()

	return w.rl.isLimitExceeded()
}

type checkerUnsafe struct {
	endAt   time.Time
	limit   int
	current int
}

func (w *checkerUnsafe) isLimitExceeded() bool {
	now := time.Now()

	if now.After(w.endAt) {
		w.current = 1
		w.endAt = now.Add(windowDuration)
		return true
	}

	if w.current < w.limit {
		w.current++
		return true
	}

	return false
}
