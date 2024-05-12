package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/google/uuid"

	"multichat_bot/internal/domain"
)

func (s *Service) HandleJoin(w http.ResponseWriter, r *http.Request) {
	req, err := prepareRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	service, ok := s.platforms[req.platform]
	if !ok {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	if err := s.db.UpdateUserPlatform(req.uuid, req.platform, req.channel); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := service.Join(req.channel); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) HandleLeave(w http.ResponseWriter, r *http.Request) {
	req, err := prepareRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	service, ok := s.platforms[req.platform]
	if !ok {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	if err := s.db.RemoveUserPlatform(req.uuid, req.platform); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := service.Leave(req.channel); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type request struct {
	channel  string
	platform domain.Platform
	uuid     uuid.UUID
}

type requestBody struct {
	Channel string `json:"channel"`
}

func prepareRequest(r *http.Request) (*request, error) {
	platform, err := getPlatform(r)
	if err != nil {
		return nil, fmt.Errorf("[prepare_request] get platform: %w", err)
	}

	body, err := getBody(r)
	if err != nil {
		return nil, fmt.Errorf("[prepare_request] get body: %w", err)
	}

	userUUID, err := getUUID(r)
	if err != nil {
		return nil, fmt.Errorf("[prepare_request] get user uuid: %w", err)
	}

	return &request{
		uuid:     userUUID,
		channel:  body.Channel,
		platform: platform,
	}, nil
}

func getPlatform(r *http.Request) (domain.Platform, error) {
	rawPlatform := chi.URLParam(r, domain.URLParamPlatform)
	platform, ok := domain.StringToPlatform[rawPlatform]
	if !ok {
		return "", errors.New("unknown platform")
	}

	return platform, nil
}

func getUUID(r *http.Request) (uuid.UUID, error) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("extract claims: %w", err)
	}

	rawUUID, ok := claims[domain.ClaimUserUUID].(string)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("get uuid from claims: %w", err)
	}

	userUUID, err := uuid.Parse(rawUUID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("parse uuid: %w", err)
	}

	return userUUID, nil
}

func getBody(r *http.Request) (*requestBody, error) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("user::prepare_request read body: %w", err)
	}

	body := &requestBody{}
	if err := json.Unmarshal(bytes, body); err != nil {
		return nil, fmt.Errorf("user::prepare_request unmarshal body: %w", err)
	}

	return body, nil
}
