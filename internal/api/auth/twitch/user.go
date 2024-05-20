package twitch

import (
	"encoding/json"
	"errors"
	"io"

	"golang.org/x/oauth2"

	"multichat_bot/internal/domain"
)

type getUserResponseDesc struct {
	Data []userInfo `json:"data"`
}
type userInfo struct {
	ID          string `json:"id,omitempty"`
	Login       string `json:"login,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Description string `json:"description,omitempty"`
	AvatarURL   string `json:"profile_image_url,omitempty"`
	Email       string `json:"email,omitempty"`
}

func getPlatformInfo(reader io.Reader) (userInfo, error) {
	body, err := io.ReadAll(reader)
	if err != nil {
		return userInfo{}, err
	}

	data := new(getUserResponseDesc)
	if err := json.Unmarshal(body, data); err != nil {
		return userInfo{}, err
	}

	if len(data.Data) == 0 {
		return userInfo{}, errors.New("user not found")
	}

	return data.Data[0], nil
}

func convertToDomain(user userInfo, token *oauth2.Token) *domain.PlatformConfig {
	return &domain.PlatformConfig{
		ExpiresIn:     token.Expiry,
		ID:            user.ID,
		Channel:       user.DisplayName,
		AccessToken:   token.AccessToken,
		RefreshToken:  token.RefreshToken,
		DisabledUsers: nil,
		BannedWords:   nil,
		IsJoined:      false,
	}
}
