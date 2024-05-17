package google

import (
	"context"
	"errors"

	"multichat_bot/internal/common/cookie"
)

type googleUserKeyType struct{}

var (
	userKey = googleUserKeyType{}
)

func withPlatformInfo(ctx context.Context, info cookie.PlatformInfo) context.Context {
	return context.WithValue(ctx, userKey, info)
}

func PlatformInfoFromContext(ctx context.Context) (cookie.PlatformInfo, error) {
	val, ok := ctx.Value(userKey).(cookie.PlatformInfo)
	if !ok {
		return cookie.PlatformInfo{}, errors.New("google: Context missing Google User")
	}

	return val, nil
}
