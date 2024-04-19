package async_cache

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"multichat_bot/internal/domain"
	"multichat_bot/internal/domain/logger"
)

type list func(ctx context.Context) ([]*domain.User, error)
type PlatformList map[string]*domain.User

type Cache struct {
	m sync.RWMutex

	refreshDuration time.Duration
	list            list

	users      []*domain.User
	byPlatform map[domain.Platform]PlatformList
}

func New(list list, refresh time.Duration) *Cache {
	return &Cache{
		refreshDuration: refresh,
		list:            list,
		byPlatform:      make(map[domain.Platform]PlatformList),
	}
}

func (c *Cache) ListPlatform(platform domain.Platform) PlatformList {
	c.m.RLock()
	defer c.m.Unlock()

	return c.byPlatform[platform]
}

func (c *Cache) StartSyncing(ctx context.Context) error {
	slog.Info("async_cache::sync start syncing")
	if err := c.update(ctx); err != nil {
		return err
	}

	ticker := time.NewTicker(c.refreshDuration)
	go func() {
		select {

		case <-ctx.Done():
			slog.Error("async_cache::sync end ", slog.String(logger.Error, ctx.Err().Error()))
			ticker.Stop()
			return

		case <-ticker.C:
			if err := c.update(ctx); err != nil {
				slog.Error("async_cache::sync update failed", slog.String(logger.Error, err.Error()))
			}

		}
	}()

	return nil
}

func (c *Cache) update(ctx context.Context) error {
	newUsers, err := c.list(ctx)
	if err != nil {
		return err
	}

	newByPlatform := make(map[domain.Platform]PlatformList, len(domain.AllPlatforms))
	for _, platform := range domain.AllPlatforms {
		newByPlatform[platform] = make(PlatformList, len(newUsers))
	}

	for _, user := range newUsers {
		for platform, userName := range user.Platforms {
			newByPlatform[platform][userName] = user
		}
	}

	c.m.Lock()
	c.users = newUsers
	c.byPlatform = newByPlatform
	c.m.Unlock()

	return nil
}
