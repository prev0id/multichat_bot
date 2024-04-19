package client

import (
	"context"
	"errors"
	"net/http"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	"multichat_bot/internal/common/apperr"
)

const (
	partSnippet              = "snippet"
	partContentDetails       = "contentDetails"
	partStatistics           = "statistics"
	partLiveStreamingDetails = "liveStreamingDetails"

	eventTypeLive = "live"

	typeVideo = "video"
)

type Client struct {
	yt *youtube.Service
}

func New(ctx context.Context, apiKey string) (*Client, error) {
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	return &Client{
		yt: service,
	}, nil
}

func (c *Client) SearchChannelID(username string) (string, error) {
	listCall := c.yt.Channels.List([]string{partContentDetails}).ForUsername(username)

	resp, err := listCall.Do()
	if err != nil {
		return "", err
	}

	if len(resp.Items) != 1 {
		return "", apperr.WithHTTPStatus(errors.New("youtube channel not found"), http.StatusNotFound)
	}

	return resp.Items[0].Id, nil
}

func (c *Client) SearchLiveStreams(channelID string) (*youtube.SearchListResponse, error) {
	listCall := c.yt.Search.List([]string{partSnippet}).
		ChannelId(channelID).
		EventType(eventTypeLive).
		Type(typeVideo)

	resp, err := listCall.Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) GetVideosWithLiveStreamingDetails(videoIDs ...string) (*youtube.VideoListResponse, error) {
	listCall := c.yt.Videos.List([]string{
		partSnippet,
		partStatistics,
		partContentDetails,
		partLiveStreamingDetails,
	}).Id(videoIDs...)

	resp, err := listCall.Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}
