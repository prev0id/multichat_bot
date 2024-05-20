package twitch

import (
	"context"
	"errors"

	"multichat_bot/internal/domain"
)

type userKeyType struct{}

var (
	userKey = userKeyType{}
)

func withPlatformInfo(ctx context.Context, info *domain.PlatformConfig) context.Context {
	return context.WithValue(ctx, userKey, info)
}

func PlatformInfoFromContext(ctx context.Context) (*domain.PlatformConfig, error) {
	val, ok := ctx.Value(userKey).(*domain.PlatformConfig)
	if !ok {
		return nil, errors.New("twitch: Context missing Twitch User")
	}

	return val, nil
}
