package auth

import (
	"log/slog"
	"net/http"

	"multichat_bot/internal/common/cookie"
	"multichat_bot/internal/database"
	"multichat_bot/internal/domain"
)

const (
	tokenLength = 256
)

type Auth struct {
	db     *database.Manager
	cookie *cookie.Store
}

func NewAuthService(db *database.Manager, cookieStore *cookie.Store) *Auth {
	return &Auth{db: db, cookie: cookieStore}
}

func (a *Auth) IsLoggedIn(r *http.Request) (domain.User, bool) {
	accessToken := a.cookie.GetAccessToken(r)
	if accessToken == "" {
		return domain.User{}, false
	}

	fromDB, err := a.db.GetUserByAccessToken(accessToken)
	if err != nil {
		slog.Error("auth.isLoggedIn", "error", err)
		return domain.User{}, false
	}

	return fromDB, true
}

func (a *Auth) HandleLogout(w http.ResponseWriter) {
	a.cookie.DestroyAppSession(w)
}

func (a *Auth) HandleDelete(w http.ResponseWriter, user domain.User) error {
	if err := a.db.DeleteUser(user.ID); err != nil {
		return err
	}

	a.cookie.DestroyAppSession(w)
	return nil
}
