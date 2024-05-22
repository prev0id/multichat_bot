package client

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type userClient struct {
	yt *youtube.Service
}

func newUserClient(config *oauth2.Config, token *oauth2.Token) (*userClient, error) {
	ctx := context.Background()
	httpClient := config.Client(ctx, token)

	service, err := youtube.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("create youtube user service: %w", err)
	}

	return &userClient{yt: service}, nil
}

func (c *userClient) listMessages(channel, pageToken string) (*youtube.LiveChatMessageListResponse, error) {
	resp, err := c.yt.LiveChatMessages.
		List(channel, []string{partSnippet, partAuthorDetails}).
		PageToken(pageToken).
		Do()

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *userClient) sendMessage(message *youtube.LiveChatMessage) error {
	_, err := c.yt.LiveChatMessages.
		Insert([]string{partSnippet}, message).
		Do()

	if err != nil {
		return fmt.Errorf("send message to youtube for chat [%s]: %w", message.Snippet.LiveChatId, err)
	}

	return nil
}
