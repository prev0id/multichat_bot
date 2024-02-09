package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"multichat_bot/internal/common/apperr"
)

type twitchJoinRequest struct {
	Chat string `json:"chat"`
}

type twitchLeaveRequest struct {
	Chat string `json:"chat"`
}

func (s *Service) TwitchLeave(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return apperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	request := &twitchLeaveRequest{}
	if err := json.Unmarshal(body, request); err != nil {
		return apperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	err = s.twitch.LeaveChat(request.Chat)
	return processResponseError(w, err)
}

func (s *Service) TwitchJoin(w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return apperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	request := &twitchJoinRequest{}
	if err := json.Unmarshal(body, request); err != nil {
		return apperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	err = s.twitch.JoinChat(ctx, request.Chat)
	return processResponseError(w, err)
}

func processResponseError(w http.ResponseWriter, err error) error {
	if err == nil {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	if appErr, ok := apperr.GetAppError(err); ok {
		return appErr
	}

	return err
}
