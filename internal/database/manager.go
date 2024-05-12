package database

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

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

	cache := async_cache.New(db.ListUsers, time.Second)

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
		return domain.User{}, fmt.Errorf("db_manager::get_by_platform get failed: %w", err)
	}
	return user, nil
}

func (m *Manager) GetUserByUUID(userUUID uuid.UUID) (domain.User, error) {
	user, err := m.cache.GetByUUID(userUUID)
	if err != nil {
		return domain.User{}, fmt.Errorf("db_manager::get_by_uuid get failed: %w", err)
	}

	return user, nil
}

func (m *Manager) UpdateUserPlatform(userUUID uuid.UUID, platform domain.Platform, channel string) error {
	return m.updatePlatform(userUUID, platform, channel)
}

func (m *Manager) RemoveUserPlatform(userUUID uuid.UUID, platform domain.Platform) error {
	return m.updatePlatform(userUUID, platform, "")
}

func (m *Manager) updatePlatform(userUUID uuid.UUID, platform domain.Platform, channel string) error {
	if err := m.db.UpdateUserPlatform(userUUID.String(), platform, channel); err != nil {
		return fmt.Errorf("db_manager::update db failed: %w", err)
	}

	err := m.cache.UpdatePlatformByUUID(userUUID, platform, channel)
	if err != nil {
		return fmt.Errorf("db_manager::update cache failed: %w", err)
	}

	return nil
}
