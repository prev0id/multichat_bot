package google

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dghubble/gologin/v2"
	oauth2Login "github.com/dghubble/gologin/v2/oauth2"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	"multichat_bot/internal/domain"
)

func StateHandler(config gologin.CookieConfig, success http.Handler) http.Handler {
	return oauth2Login.StateHandler(config, success)
}

func LoginHandler(config *oauth2.Config, failure http.Handler) http.Handler {
	return oauth2Login.LoginHandler(config, failure)
}

func CallbackHandler(config *oauth2.Config, success, failure http.Handler) http.Handler {
	success = googleHandler(config, success, failure)
	return oauth2Login.CallbackHandler(config, success, failure)
}

func googleHandler(config *oauth2.Config, success, failure http.Handler) http.Handler {
	if failure == nil {
		failure = gologin.DefaultFailureHandler
	}
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		token, err := oauth2Login.TokenFromContext(ctx)
		if err != nil {
			ctx = gologin.WithError(ctx, err)
			failure.ServeHTTP(w, req.WithContext(ctx))
			return
		}

		slog.Info(fmt.Sprintf("YouTube logged in with token: %+v", token))

		youtubeService, err := youtube.NewService(ctx, option.WithHTTPClient(config.Client(ctx, token)))
		if err != nil {
			ctx = gologin.WithError(ctx, err)
			failure.ServeHTTP(w, req.WithContext(ctx))
			return
		}

		resp, err := youtubeService.Channels.List([]string{"id", "snippet"}).Mine(true).Do()
		if err != nil {
			ctx = gologin.WithError(ctx, err)
			failure.ServeHTTP(w, req.WithContext(ctx))
			return
		}

		if len(resp.Items) == 0 {
			ctx = gologin.WithError(ctx, errors.New("no youtube channel found"))
			failure.ServeHTTP(w, req.WithContext(ctx))
			return
		}

		ctx = withPlatformInfo(ctx, convertToDomain(resp.Items[0], token))
		success.ServeHTTP(w, req.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func convertToDomain(channel *youtube.Channel, token *oauth2.Token) *domain.PlatformConfig {
	return &domain.PlatformConfig{
		ExpiresIn:     token.Expiry,
		ID:            channel.Id,
		Channel:       channel.Snippet.Title,
		AccessToken:   token.AccessToken,
		RefreshToken:  token.RefreshToken,
		DisabledUsers: nil,
		BannedWords:   nil,
		IsJoined:      false,
	}
}
