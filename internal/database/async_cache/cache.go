package async_cache

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"

	"multichat_bot/internal/domain"
	"multichat_bot/internal/domain/logger"
)

type list func() ([]*domain.User, error)
type PlatformList map[string]*domain.User

type Cache struct {
	list            list
	users           map[uuid.UUID]*domain.User
	byPlatform      map[domain.Platform]PlatformList
	refreshDuration time.Duration
	m               sync.RWMutex
}

func New(list list, refresh time.Duration) *Cache {
	return &Cache{
		refreshDuration: refresh,
		list:            list,
		byPlatform:      make(map[domain.Platform]PlatformList),
		users:           make(map[uuid.UUID]*domain.User),
	}
}

func (c *Cache) GetByPlatform(platform domain.Platform, channel string) (domain.User, error) {
	c.m.RLock()
	defer c.m.RUnlock()

	list, ok := c.byPlatform[platform]
	if !ok {
		return domain.User{}, errors.New("platform not found")
	}

	user, ok := list[channel]
	if !ok {
		return domain.User{}, errors.New("user not found")
	}

	return *user, nil
}

func (c *Cache) GetByUUID(userUUID uuid.UUID) (domain.User, error) {
	c.m.RLock()
	defer c.m.RUnlock()

	user, ok := c.users[userUUID]
	if !ok {
		if !ok {
			return domain.User{}, fmt.Errorf("user with the uuid %s not exists", userUUID.String())
		}
	}

	return *user, nil
}

func (c *Cache) UpdatePlatformByUUID(userID uuid.UUID, platform domain.Platform, value string) error {
	c.m.Lock()
	defer c.m.Unlock()

	user, ok := c.users[userID]
	if !ok {
		return errors.New("user not found")
	}

	user.Platforms[platform] = value

	return nil
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

	newByPlatform := make(map[domain.Platform]PlatformList, len(domain.StringToPlatform))
	for _, platform := range domain.StringToPlatform {
		newByPlatform[platform] = make(PlatformList, len(newUsers))
	}

	users := make(map[uuid.UUID]*domain.User, len(newUsers))
	for _, user := range newUsers {
		for platform, userName := range user.Platforms {
			newByPlatform[platform][userName] = user
		}
		users[user.UUID] = user
	}

	c.m.Lock()
	c.users = users
	c.byPlatform = newByPlatform
	c.m.Unlock()

	return nil
}
