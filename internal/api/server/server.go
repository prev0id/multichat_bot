package server

import (
	"net/http"

	"multichat_bot/internal/common/apperr"
)

type Middleware func(http.Handler) http.Handler

type Handler func(w http.ResponseWriter, r *http.Request) error

type Server struct {
	mux         *http.ServeMux
	middlewares []Middleware
}

func New() *Server {
	return &Server{
		mux: http.NewServeMux(),
	}
}

func (s *Server) WithMiddleware(middleware ...Middleware) {
	s.middlewares = append(s.middlewares, middleware...)
}

func (s *Server) GetHandler() http.Handler {
	handler := http.Handler(s.mux)
	for idx := len(s.middlewares) - 1; idx >= 0; idx-- {
		handler = s.middlewares[idx](handler)
	}

	return handler
}

func (s *Server) RegisterHandler(pattern string, handler Handler) {
	s.mux.HandleFunc(pattern, wrapIntoHttpHandler(handler))
}

func wrapIntoHttpHandler(f Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			http.Error(w, err.Error(), apperr.HTTPStatus(err))
		}
	}
}
