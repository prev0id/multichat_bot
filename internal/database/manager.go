package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"multichat_bot/internal/database/async_cache"
	"multichat_bot/internal/database/domain"
	global_domain "multichat_bot/internal/domain"
)

var (
	errUserNotFound = errors.New("user not found")
)

type db interface {
	ListUsers(ctx context.Context) ([]*global_domain.User, error)

	SetPlatformToUser(ctx context.Context, uuid string, platform domain.PlatformKey) error
	GetPlatfromToUser(ctx context.Context, uuid string, platform domain.PlatformName) (string, error)
	SetPlatfromConnection(ctx context.Context, from, to domain.PlatformKey)
}

type Manager struct {
	db    db
	cache *async_cache.Cache
}

func New(ctx context.Context, db db, refresh time.Duration) (*Manager, error) {
	cache := async_cache.New(db.ListUsers, refresh)

	if err := cache.StartSyncing(ctx); err != nil {
		return nil, fmt.Errorf("db_manager::cache first sync failed: %s", err.Error())
	}

	return &Manager{
		db:    db,
		cache: cache,
	}, nil
}

func (m *Manager) GetUserByPlatform(platform global_domain.Platform, username string) (*global_domain.User, error) {
	list := m.cache.ListPlatform(platform)
	if list == nil {
		return nil, errUserNotFound
	}

	user, ok := list[username]
	if !ok {
		return nil, errUserNotFound
	}

	return user, nil
}

func (m *Manager) AddTwitchChatToUser(ctx context.Context, uuid, chat string) error {
	if err := m.db.SetPlatformToUser(ctx, uuid, domain.PlatformKey{Name: domain.Twitch, Value: chat}); err != nil {
		return err
	}

	return nil
}

func (m *Manager) AddYoutubeChatID(ctx context.Context, uuid, chatID string) error {
	if err := m.db.SetPlatformToUser(ctx, uuid, domain.PlatformKey{Name: domain.Youtube, Value: chatID}); err != nil {
		return err
	}

	return nil
}

func (m *Manager) GetTwitchChat(ctx context.Context, uuid string) (string, error) {
	chat, err := m.db.GetPlatfromToUser(ctx, uuid, domain.Twitch)
	if err != nil {
		return "", err
	}

	return chat, nil
}

func (m *Manager) GetYoutubeChatID(ctx context.Context, uuid string) (string, error) {
	chatID, err := m.db.GetPlatfromToUser(ctx, uuid, domain.Youtube)
	if err != nil {
		return "", err
	}

	return chatID, nil
}
