package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"multichat_bot/internal/domain"
)

func (a *Auth) HandleLogin(w http.ResponseWriter, r *http.Request, platform domain.Platform, config *domain.PlatformConfig) error {
	user, err := a.getUserOnLogin(r, platform, config)
	if err != nil {
		return fmt.Errorf("get user on login failed: %w", err)
	}

	if err := a.cookie.SaveAccessToken(w, user.AccessToken); err != nil {
		return fmt.Errorf("save access token failed: %w", err)
	}

	if err := a.db.UpdatePlatform(user.ID, platform, config); err != nil {
		return fmt.Errorf("update platform failed: %w", err)
	}

	return nil
}

func (a *Auth) getUserOnLogin(r *http.Request, platform domain.Platform, config *domain.PlatformConfig) (domain.User, error) {
	if user, ok := a.IsLoggedIn(r); ok {
		return user, nil
	}

	user, err := a.getUnauthorized(platform, config)
	return user, err
}

func (a *Auth) getUnauthorized(platform domain.Platform, config *domain.PlatformConfig) (domain.User, error) {
	if user, ok := a.db.GetUserByChannel(platform, config.ID); ok {
		return user, nil
	}

	token := generateSecureToken()

	id, err := a.db.NewUser(token)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		Platforms:   make(map[domain.Platform]*domain.PlatformConfig),
		ID:          id,
		AccessToken: token,
	}, nil
}

func generateSecureToken() string {
	b := make([]byte, tokenLength)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
