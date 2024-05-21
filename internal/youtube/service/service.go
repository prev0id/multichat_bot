package service

import (
	"context"
	"google.golang.org/api/youtube/v3"
)

type Service struct {
	youtubeService *youtube.Service
}

func New(ctx context.Context) (*Service, error) {
	service, err := youtube.NewService(ctx)
	if err != nil {
		return nil, err
	}

	return &Service{
		youtubeService: service,
	}, nil
}

func (s *Service) List() error {
	listCall := s.youtubeService.LiveChatMessages.List("id", []string{""})
	resp, err := listCall.Do()
	if err != nil {
		return err
	}

	for message := resp.Items {

	}

}
