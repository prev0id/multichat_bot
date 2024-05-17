package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"multichat_bot/internal/config"
	"multichat_bot/internal/database/adapter"
	"multichat_bot/internal/database/async_cache"
	"multichat_bot/internal/domain"
)

type Manager struct {
	db    *adapter.DB
	cache *async_cache.Cache
}

func New(ctx context.Context, cfg config.DB) (*Manager, error) {
	db, err := adapter.New(cfg.DBPath)
	if err != nil {
		return nil, err
	}

	cache := async_cache.New(db.ListUsers, 5*time.Second)

	if err := cache.StartSyncing(ctx); err != nil {
		return nil, fmt.Errorf("db_manager::cache first sync failed: %w", err)
	}

	return &Manager{
		db:    db,
		cache: cache,
	}, nil
}

func (m *Manager) GetUserByChannel(platform domain.Platform, channel string) (domain.User, error) {
	user, err := m.cache.GetByPlatform(platform, channel)
	if err != nil {
		return user, fmt.Errorf("db_manager::get_by_platform get failed: %w", err)
	}

	return user, nil
}

func (m *Manager) GetUserByID(id int64) (domain.User, error) {
	user, err := m.cache.GetByID(id)
	if err != nil {
		return user, fmt.Errorf("db_manager::get_by_id get failed: %w", err)
	}

	return user, nil
}

func (m *Manager) GetUserByChannelOrCreateNew(platform domain.Platform, channel string) (domain.User, error) {
	user, err := m.GetUserByChannel(platform, channel)
	if err == nil {
		return user, nil
	}

	if !errors.Is(err, async_cache.ErrNotFound) {
		return domain.User{}, err
	}

	id, err := m.db.NewUser()
	if err != nil {
		return domain.User{}, fmt.Errorf("db_manager::new_user create failed: %w", err)
	}

	return domain.User{
		ID:        id,
		Platforms: make(map[domain.Platform]*domain.PlatformConfig),
	}, nil
}

func (m *Manager) DeleteUserPlatform(id int64, platform domain.Platform) error {
	return m.db.DeletePlatform(id, platform)
}

func (m *Manager) JoinChannel(id int64, platform domain.Platform) error {
	return m.db.ChangeJoined(id, platform, true)
}

func (m *Manager) LeaveChannel(id int64, platform domain.Platform) error {
	return m.db.ChangeJoined(id, platform, false)
}

func (m *Manager) NewUser() (domain.User, error) {
	id, err := m.db.NewUser()
	if err != nil {
		return domain.User{}, fmt.Errorf("db_manager::new_user create failed: %w", err)
	}

	return domain.User{
		Platforms: make(map[domain.Platform]*domain.PlatformConfig),
		ID:        id,
	}, nil
}

func (m *Manager) UpdatePlatform(id int64, platform domain.Platform, platformConfig *domain.PlatformConfig) error {
	if err := m.db.UpsertPlatform(id, platform, platformConfig); err != nil {
		return fmt.Errorf("db_manager::update db failed: %w", err)
	}

	return nil
}
