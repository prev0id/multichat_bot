package adapter

import (
	"context"
	"time"

	"multichat_bot/internal/config"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres" // import postgres Dialect
	"github.com/redis/go-redis/v9"
)

const (
	prefixUser    = "user:"
	prefixTwitch  = "twitch:"
	prefixYoutube = "youtube:"

	dialect   = "postgres"
	tableUser = "user"

	columnKey     = "key"
	columnTwitch  = "twitch"
	columnYoutube = "youtube"
)

type Adapter struct {
	client *redis.Client
}

func New(cfg config.Redis) (*Adapter, error) {
	opts, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, err
	}

	return &Adapter{
		client: redis.NewClient(opts),
	}, nil
}

func (a *Adapter) SetPlatformToUser(ctx context.Context, uuid, platform, value string) error {
	if err := a.client.HSet(ctx, userKey(uuid), platform, value).Err(); err != nil {
		return err
	}

	return nil
}

func (a *Adapter) SetConnectionBetweenPlatforms(ctx context.Context, key, value string) error {
	if err := a.client.Set(ctx, key, value, time.Duration(-1)).Err(); err != nil {
		return err
	}

	return nil
}

func (a *Adapter) GetPlatfromToUser(ctx context.Context, uuid, platform string) (string, error) {
	value, err := a.client.HGet(ctx, userKey(uuid), platform).Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

func (a *Adapter) IsUserExists(ctx context.Context, uuid string) (bool, error) {
	fieldsLen, err := a.client.Exists(ctx, userKey(uuid)).Result()
	if err != nil {
		return false, err
	}

	return fieldsLen > 0, nil
}

func userKey(uuid string) string {
	return prefixUser + uuid
}

func twitchKey(chat string) string {
	return prefixTwitch + chat
}

func youtubeKey(chatID string) string {
	return prefixYoutube + chatID
}

func ListAllUsers() error {
	query, params, err := goqu.Dialect(dialect).
		From(tableUser).
		Select(
			goqu.C(columnKey),
			goqu.C(columnTwitch),
			goqu.C(columnYoutube),
		).ToSQL()

	if err != nil {
		return err
	}
}
