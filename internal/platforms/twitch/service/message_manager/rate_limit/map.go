package rate_limit

import (
	"sync"
)

type Map struct {
	values       map[string]checkerUnsafe
	defaultLimit int
	m            sync.Mutex
}

func NewMapChecker(defaultLimit int) *Map {
	return &Map{
		defaultLimit: defaultLimit,
	}
}

func (m *Map) IsLimitExceeded(key string) bool {
	m.m.Lock()
	defer m.m.Unlock()

	rl, isExist := m.values[key]
	if isExist {
		return rl.isLimitExceeded()
	}

	m.values[key] = checkerUnsafe{
		limit: m.defaultLimit,
	}

	return true
}

func (m *Map) UpgradeRateLimit(key string, limit int) {
	m.m.Lock()
	defer m.m.Unlock()

	rl, isExist := m.values[key]
	if !isExist {
		return
	}

	rl.limit = limit
}
