package client

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/youtube/v3"

	"multichat_bot/internal/config"
	"multichat_bot/internal/domain"
)

const (
	typeTextMessageEvent = "textMessageEvent"
)

type Adapter struct {
	client         *serverClient
	messageChannel chan<- *domain.Message
	channelConfigs map[string]*channelConfig
	oauth2Config   *oauth2.Config
	duration       time.Duration
	m              sync.RWMutex
}

type channelConfig struct {
	service   *userClient
	chatID    string
	pageToken string
	id        string
}

func NewAdapter(ctx context.Context, cfg config.Youtube, auth config.AuthProvider, messageChannel chan<- *domain.Message) (*Adapter, error) {
	client, err := newClient(ctx, cfg.APIKey)
	if err != nil {
		return nil, err
	}

	return &Adapter{
		client:         client,
		duration:       5 * time.Second,
		messageChannel: messageChannel,
		channelConfigs: make(map[string]*channelConfig),
		oauth2Config: &oauth2.Config{
			ClientID:     auth.ClientKey,
			ClientSecret: auth.ClientSecret,
			Scopes:       auth.Scopes,
			Endpoint:     google.Endpoint,
		},
	}, nil
}

func (a *Adapter) StartListening(ctx context.Context) {
	go func() {
		a.run(ctx)
	}()
}

func (a *Adapter) Join(cfg *domain.PlatformConfig) error {
	slog.Info(fmt.Sprintf("YouTube joining channel %s(%s)", cfg.Channel, cfg.ID))

	searchResult, err := a.client.searchLiveStreams(cfg.ID)
	if err != nil {
		return fmt.Errorf("search live streams: %w", err)
	}

	details, err := a.client.videoDetails(searchResult.Id.VideoId)
	if err != nil {
		return fmt.Errorf("get video details: %w", err)
	}

	service, err := newUserClient(a.oauth2Config, cfg.Token())
	if err != nil {
		return fmt.Errorf("new user client: %w", err)
	}

	chatID := details.LiveStreamingDetails.ActiveLiveChatId

	slog.Info(fmt.Sprintf("YouTube joind channel %s, chatID %s", cfg.Channel, chatID))

	a.m.Lock()
	a.channelConfigs[cfg.ID] = &channelConfig{
		chatID:  chatID,
		service: service,
		id:      cfg.ID,
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

func (a *Adapter) SendMessage(msg *domain.Message, cfg *domain.PlatformConfig) error {
	a.m.RLock()
	defer a.m.RUnlock()

	channel, ok := a.channelConfigs[cfg.ID]
	if !ok {
		return errors.New("no channel config for " + cfg.ID)
	}

	return channel.service.sendMessage(convertDomainMessageToYoutube(msg, channel.chatID))
}

func (a *Adapter) run(ctx context.Context) {
	ticker := time.NewTicker(a.duration)
	for {
		select {
		case <-ticker.C:
			a.listMessages()
		case <-ctx.Done():
			slog.Warn("[youtube] stop listening: " + ctx.Err().Error())
		}
	}
}

func (a *Adapter) listMessages() {
	a.m.RLock()
	defer a.m.RUnlock()

	for _, channel := range a.channelConfigs {
		resp, err := channel.service.listMessages(channel.chatID, channel.pageToken)
		if err == nil {
			for _, message := range resp.Items {
				a.messageChannel <- convertYoutubeMessageToDomain(message, channel.id)
			}
			channel.pageToken = resp.NextPageToken
			continue
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

func convertYoutubeMessageToDomain(msg *youtube.LiveChatMessage, channelID string) *domain.Message {
	return &domain.Message{
		From:     msg.AuthorDetails.DisplayName,
		Text:     msg.Snippet.TextMessageDetails.MessageText,
		Channel:  channelID,
		Platform: domain.YouTube,
	}
}

func convertDomainMessageToYoutube(msg *domain.Message, chatID string) *youtube.LiveChatMessage {
	return &youtube.LiveChatMessage{
		Snippet: &youtube.LiveChatMessageSnippet{
			LiveChatId: chatID,
			Type:       typeTextMessageEvent,
			TextMessageDetails: &youtube.LiveChatTextMessageDetails{
				MessageText: msg.FormatedText(),
			},
		},
	}
}
