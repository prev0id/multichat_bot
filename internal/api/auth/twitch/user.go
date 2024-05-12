package twitch

import (
	"encoding/json"
	"errors"
	"io"
)

type UserInfo struct {
	ID          string `json:"id,omitempty"`
	Login       string `json:"login,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Description string `json:"description,omitempty"`
	AvatarURL   string `json:"profile_image_url,omitempty"`
	Email       string `json:"email,omitempty"`
}

func getUserInfo(reader io.Reader) (*UserInfo, error) {
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	data := new(getUserResponseDesc)
	if err := json.Unmarshal(body, data); err != nil {
		return nil, err
	}

	if len(data.Data) == 0 {
		return nil, errors.New("user not found")
	}

	return data.Data[0], nil
}

type getUserResponseDesc struct {
	Data []*UserInfo `json:"data"`
}
