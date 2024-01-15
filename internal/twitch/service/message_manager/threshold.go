package message_manager

import (
	"time"
)

type windowRateLimit struct {
	endAt    time.Time
	duration time.Duration
	capacity int
	current  int
}

func (w *windowRateLimit) isSendAllowed() bool {
	now := time.Now()

	if now.After(w.endAt) {
		w.current = 1
		w.endAt = now.Add(w.duration)
		return true
	}

	if w.current < w.capacity {
		w.current++
		return true
	}

	return false
}
