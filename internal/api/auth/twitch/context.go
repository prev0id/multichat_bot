package twitch

import (
	"context"
	"errors"

	"multichat_bot/internal/common/cookie"
)

type userKeyType struct{}

var (
	userKey = userKeyType{}
)

func withPlatformInfo(ctx context.Context, info cookie.PlatformInfo) context.Context {
	return context.WithValue(ctx, userKey, info)
}

func PlatformInfoFromContext(ctx context.Context) (cookie.PlatformInfo, error) {
	val, ok := ctx.Value(userKey).(cookie.PlatformInfo)
	if !ok {
		return cookie.PlatformInfo{}, errors.New("twitch: Context missing Twitch User")
	}

	return val, nil
}