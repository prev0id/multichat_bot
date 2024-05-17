package twitch

import (
	"net/http"

	"github.com/dghubble/gologin/v2"
	oauth2Login "github.com/dghubble/gologin/v2/oauth2"
	"golang.org/x/oauth2"
)

const (
	userEndpoint string = "https://api.twitch.tv/helix/users"
)

func StateHandler(config gologin.CookieConfig, success http.Handler) http.Handler {
	return oauth2Login.StateHandler(config, success)
}

func LoginHandler(config *oauth2.Config, failure http.Handler) http.Handler {
	return oauth2Login.LoginHandler(config, failure)
}

func CallbackHandler(config *oauth2.Config, success, failure http.Handler) http.Handler {
	success = twitchHandler(config, success, failure)
	return oauth2Login.CallbackHandler(config, success, failure)
}

func twitchHandler(config *oauth2.Config, success, failure http.Handler) http.Handler {
	if failure == nil {
		failure = gologin.DefaultFailureHandler
	}

	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token, err := oauth2Login.TokenFromContext(ctx)
		if err != nil {
			ctx = gologin.WithError(ctx, err)
			failure.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		req, err := http.NewRequest(http.MethodGet, userEndpoint, nil)
		if err != nil {
			ctx = gologin.WithError(ctx, err)
			failure.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		req.Header.Add("Client-ID", config.ClientID)
		req.Header.Add("Authorization", "Bearer "+token.AccessToken)

		resp, err := http.DefaultClient.Do(req.WithContext(ctx))
		if err != nil {
			ctx = gologin.WithError(ctx, err)
			failure.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		info, err := getPlatformInfo(resp.Body)
		if err != nil {
			ctx = gologin.WithError(ctx, err)
			failure.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		converted := convertToCookie(info, token)

		success.ServeHTTP(w, r.WithContext(withPlatformInfo(ctx, converted)))
	}

	return http.HandlerFunc(fn)
}
