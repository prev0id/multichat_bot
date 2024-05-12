package client

import (
	"context"
	"fmt"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const (
	partSnippet              = "snippet"
	partContentDetails       = "contentDetails"
	partStatistics           = "statistics"
	partLiveStreamingDetails = "liveStreamingDetails"
	partAuthorDetails        = "authorDetails"

	eventTypeLive = "live"

	typeVideo = "video"
)

type client struct {
	yt *youtube.Service
}

func newClient(ctx context.Context, apiKey string) (*client, error) {
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	return &client{
		yt: service,
	}, nil
}

func (c *client) listMessages(channel, pageToken string) (*youtube.LiveChatMessageListResponse, error) {
	resp, err := c.yt.LiveChatMessages.
		List(channel, []string{partSnippet, partAuthorDetails}).
		PageToken(pageToken).
		Do()

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *client) sendMessage(message *youtube.LiveChatMessage) error {
	_, err := c.yt.LiveChatMessages.
		Insert([]string{partSnippet}, message).
		Do()

	if err != nil {
		return err
	}

	return nil
}

func (c *client) searchLiveStreams(channelID string) (*youtube.SearchResult, error) {
	resp, err := c.yt.Search.
		List([]string{partSnippet}).
		ChannelId(channelID).
		EventType(eventTypeLive).
		Type(typeVideo).
		Do()

	if err != nil {
		return nil, err
	}

	if len(resp.Items) == 0 {
		return nil, fmt.Errorf("no active live streams found")
	}

	return resp.Items[0], nil
}

func (c *client) videoDetails(videoID string) (*youtube.Video, error) {
	resp, err := c.yt.Videos.
		List([]string{
			partSnippet,
			partStatistics,
			partContentDetails,
			partLiveStreamingDetails,
		}).
		Id(videoID).
		Do()

	if err != nil {
		return nil, err
	}

	if len(resp.Items) == 0 {
		return nil, fmt.Errorf("no active live streams found for ID: %s", videoID)
	}

	return resp.Items[0], nil
}
