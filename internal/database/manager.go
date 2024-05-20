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

func (m *Manager) GetUserByChannel(platform domain.Platform, channel string) (domain.User, bool) {
	user, ok := m.cache.GetByPlatform(platform, channel)
	if !ok {
		return domain.User{}, false
	}

	return user, true
}

func (m *Manager) GetUserByAccessToken(accessToken string) (domain.User, error) {
	users, err := m.db.ListUsers()
	if err != nil {
		return domain.User{}, err
	}

	for _, user := range users {
		if user.AccessToken == accessToken {
			return *user, nil
		}
	}

	return domain.User{}, errors.New("user not found")
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

func (m *Manager) NewUser(token string) (int64, error) {
	id, err := m.db.NewUser(token)
	if err != nil {
		return 0, fmt.Errorf("db_manager::new_user create failed: %w", err)
	}

	return id, nil
}

func (m *Manager) UpdatePlatform(id int64, platform domain.Platform, platformConfig *domain.PlatformConfig) error {
	if err := m.db.UpsertPlatform(id, platform, platformConfig); err != nil {
		return fmt.Errorf("db_manager::update db failed: %w", err)
	}

	return nil
}

func (m *Manager) UpdateBannedUsers(id int64, platform domain.Platform, bannedUsers domain.BannedList) error {
	return m.db.UpdateBannedUsers(id, platform, bannedUsers)
}

func (m *Manager) UpdateBannedWords(id int64, platform domain.Platform, bannedWords domain.BannedList) error {
	return m.db.UpdateBannedWords(id, platform, bannedWords)
}
