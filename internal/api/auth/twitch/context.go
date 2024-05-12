package twitch

import (
	"context"
	"errors"
)

type userKeyType struct{}

var (
	userKey = userKeyType{}
)

func withUser(ctx context.Context, user *UserInfo) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func UserFromContext(ctx context.Context) (*UserInfo, error) {
	val, ok := ctx.Value(userKey).(*UserInfo)
	if !ok {
		return nil, errors.New("twitch: Context missing Twitch User")
	}

	return val, nil
}
