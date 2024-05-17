package async_cache

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"multichat_bot/internal/domain"
	"multichat_bot/internal/domain/logger"
)

var (
	ErrNotFound = errors.New("not found")
)

type list func() ([]*domain.User, error)
type PlatformList map[string]*domain.User

type Cache struct {
	list            list
	byPlatform      map[domain.Platform]PlatformList
	users           map[int64]*domain.User
	refreshDuration time.Duration
	m               sync.RWMutex
}

func New(list list, refresh time.Duration) *Cache {
	return &Cache{
		refreshDuration: refresh,
		list:            list,
		byPlatform:      make(map[domain.Platform]PlatformList),
		users:           make(map[int64]*domain.User),
	}
}

func (c *Cache) GetByPlatform(platform domain.Platform, channelID string) (domain.User, error) {
	c.m.RLock()
	defer c.m.RUnlock()

	list, ok := c.byPlatform[platform]
	if !ok {
		return domain.User{}, ErrNotFound
	}

	user, ok := list[channelID]
	if !ok {
		return domain.User{}, ErrNotFound
	}

	return *user, nil
}
func (c *Cache) GetByID(id int64) (domain.User, error) {
	c.m.RLock()
	defer c.m.RUnlock()

	user, ok := c.users[id]
	if !ok {
		return domain.User{}, ErrNotFound
	}

	return *user, nil
}
func (c *Cache) StartSyncing(ctx context.Context) error {
	slog.Info("async_cache::sync start syncing")
	if err := c.update(); err != nil {
		return err
	}

	ticker := time.NewTicker(c.refreshDuration)
	go func() {
		for {
			select {
			case <-ctx.Done():
				slog.Error("async_cache::sync end ", slog.String(logger.Error, ctx.Err().Error()))
				ticker.Stop()
				return

			case <-ticker.C:
				if err := c.update(); err != nil {
					slog.Error("async_cache::sync update failed", slog.String(logger.Error, err.Error()))
				}
			}
		}
	}()

	return nil
}

func (c *Cache) update() error {
	newUsers, err := c.list()
	if err != nil {
		return err
	}

	newByPlatform := make(map[domain.Platform]PlatformList, len(domain.Platforms))
	for _, platform := range domain.Platforms {
		newByPlatform[platform] = make(PlatformList, len(newUsers))
	}

	users := make(map[int64]*domain.User, len(newUsers))
	for _, user := range newUsers {
		for platform, platformConfig := range user.Platforms {
			newByPlatform[platform][platformConfig.ID] = user
		}
		users[user.ID] = user
	}

	c.m.Lock()
	c.users = users
	c.byPlatform = newByPlatform
	c.m.Unlock()

	return nil
}
