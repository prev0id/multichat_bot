package google

import (
	"context"
	"errors"

	"multichat_bot/internal/domain"
)

type googleUserKeyType struct{}

var (
	userKey = googleUserKeyType{}
)

func withPlatformInfo(ctx context.Context, info *domain.PlatformConfig) context.Context {
	return context.WithValue(ctx, userKey, info)
}

func PlatformInfoFromContext(ctx context.Context) (*domain.PlatformConfig, error) {
	val, ok := ctx.Value(userKey).(*domain.PlatformConfig)
	if !ok {
		return nil, errors.New("google: Context missing Google User")
	}

	return val, nil
}
