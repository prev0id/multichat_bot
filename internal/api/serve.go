package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"multichat_bot/internal/api/auth"
	"multichat_bot/internal/api/page"
	"multichat_bot/internal/api/user"
	"multichat_bot/internal/config"
)

func Serve(cfg config.API, userService *user.Service, pageService *page.Service, authService *auth.Service) error {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(3 * time.Second))
	router.Use(middleware.Heartbeat("/ping"))

	router.NotFound(pageService.Handle404)

	// pages
	router.Group(func(r chi.Router) {
		r.Get("/", pageService.HandleRoot)
		r.Get("/settings", pageService.HandleSetting)

		// static
		r.Get("/assets/icon.png", pageService.HandleIcon)
		r.Get("/css/main.min.css", pageService.HandleCSS)
		r.Get("/js/htmx.min.js", pageService.HandleJS)
	})

	// auth
	router.Group(func(r chi.Router) {
		r.Post("/auth/logout", authService.Logout)
		r.Get("/auth/{platform}/callback", authService.CallBack)
		r.Get("/auth/{platform}/login", authService.Login)
		r.Get("/auth/{platform}/delete", authService.DeletePlatform)
	})

	// private
	router.Group(func(r chi.Router) {
		r.Post("/user/{platform}/join", userService.HandleJoin)
		r.Post("/user/{platform}/leave", userService.HandleLeave)

		r.Post("/user/{platform}/ban/add/word", userService.HandleAddBanWord)
		r.Post("/user/{platform}/ban/remove/word/{word}", userService.HandleRemoveBanWord)

		r.Post("/user/{platform}/ban/add/user", userService.HandleAddBanUser)
		r.Post("/user/{platform}/ban/remove/user/{user}", userService.HandleRemoveBanUser)

	})

	server := &http.Server{
		Addr:              cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 3 * time.Second,
	}

	return server.ListenAndServe()
}
