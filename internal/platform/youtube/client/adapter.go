package client

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/youtube/v3"

	"multichat_bot/internal/config"
	"multichat_bot/internal/domain"
)

const (
	typeTextMessageEvent = "textMessageEvent"
)

type Adapter struct {
	client         *client
	messageChannel chan<- *domain.Message
	channelConfigs map[string]channelConfig
	duration       time.Duration
	m              sync.RWMutex
}

type channelConfig struct {
	chatID    string
	pageToken string
}

func NewAdapter(ctx context.Context, cfg config.Youtube, messageChannel chan<- *domain.Message) (*Adapter, error) {
	client, err := newClient(ctx, cfg.APIKey)
	if err != nil {
		return nil, err
	}

	return &Adapter{
		client:         client,
		duration:       time.Second,
		messageChannel: messageChannel,
	}, nil
}

func (a *Adapter) StartListening(ctx context.Context) {
	go func() {
		a.run(ctx)
	}()
}

func (a *Adapter) Join(channelID string) error {
	searchResult, err := a.client.searchLiveStreams(channelID)
	if err != nil {
		return err
	}

	details, err := a.client.videoDetails(searchResult.Id.VideoId)
	if err != nil {
		return err
	}

	a.m.Lock()
	a.channelConfigs[channelID] = channelConfig{
		chatID: details.LiveStreamingDetails.ActiveLiveChatId,
	}
	a.m.Unlock()

	return nil
}

func (a *Adapter) Leave(channelID string) {
	slog.Warn("[youtube] deleting chat for " + channelID)
	a.m.Lock()
	delete(a.channelConfigs, channelID)
	a.m.Unlock()
}

func (a *Adapter) SendMessage(msg *domain.Message, chatID string) error {
	converted := convertDomainMessageToYoutube(msg, chatID)
	return a.client.sendMessage(converted)
}

func (a *Adapter) run(ctx context.Context) {
	ticker := time.NewTicker(a.duration)
	select {
	case <-ticker.C:
		a.listMessages()
	case <-ctx.Done():
		slog.Warn("[youtube] stop listening: " + ctx.Err().Error())
	}
}

func (a *Adapter) listMessages() {
	a.m.RLock()
	defer a.m.RUnlock()

	for _, channel := range a.channelConfigs {
		resp, err := a.client.listMessages(channel.chatID, channel.pageToken)
		if err == nil {
			for _, message := range resp.Items {
				a.messageChannel <- convertYoutubeMessageToDomain(message)
			}
		}

		var googleErr *googleapi.Error
		if !errors.As(err, &googleErr) {
			slog.Error("[youtube] unknown list chat error: " + err.Error())
		}

		switch googleErr.Code {
		case http.StatusForbidden, http.StatusNotFound:
			a.Leave(channel.chatID)
		}
	}
}

func convertYoutubeMessageToDomain(msg *youtube.LiveChatMessage) *domain.Message {
	return &domain.Message{
		From:     msg.AuthorDetails.DisplayName,
		Text:     msg.Snippet.TextMessageDetails.MessageText,
		Channel:  msg.Snippet.AuthorChannelId,
		Platform: domain.YouTube,
	}
}

func convertDomainMessageToYoutube(msg *domain.Message, chatID string) *youtube.LiveChatMessage {
	return &youtube.LiveChatMessage{
		Snippet: &youtube.LiveChatMessageSnippet{
			LiveChatId: chatID,
			Type:       typeTextMessageEvent,
			TextMessageDetails: &youtube.LiveChatTextMessageDetails{
				MessageText: msg.Text,
			},
		},
	}
}
