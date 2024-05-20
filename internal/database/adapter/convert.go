package adapter

import (
	"encoding/json"
	"time"

	"multichat_bot/internal/domain"
)

func convertPlatformsToDomain(rows []platformRow) (map[domain.Platform]*domain.PlatformConfig, error) {
	result := make(map[domain.Platform]*domain.PlatformConfig, len(rows))

	for _, row := range rows {
		converted, err := convertPlatformToDomain(row)
		if err != nil {
			return nil, err
		}

		result[domain.StringToPlatform[row.Name]] = converted
	}

	return result, nil
}

func convertPlatformToDomain(row platformRow) (*domain.PlatformConfig, error) {
	expiresIn, err := time.Parse(time.RFC3339, row.ExpiresIn)
	if err != nil {
		return nil, err
	}

	var (
		bannedUsers domain.BannedList
		bannedWords domain.BannedList
	)

	if err := json.Unmarshal([]byte(row.DisabledUsers), &bannedUsers); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(row.BannedWords), &bannedWords); err != nil {
		return nil, err
	}

	return &domain.PlatformConfig{
		IsJoined:      row.IsJoined == 1,
		ID:            row.ID,
		Channel:       row.Channel,
		AccessToken:   row.AccessToken,
		RefreshToken:  row.RefreshToken,
		ExpiresIn:     expiresIn,
		DisabledUsers: bannedUsers,
		BannedWords:   bannedWords,
	}, nil
}

func convertPlatformToDB(id int64, platform domain.Platform, config *domain.PlatformConfig) platformRow {
	bannedUsers, _ := json.Marshal(config.DisabledUsers)
	bannedWords, _ := json.Marshal(config.BannedWords)
	expiresIn := config.ExpiresIn.Format(time.RFC3339)

	return platformRow{
		UserID:        id,
		IsJoined:      0,
		Name:          platform.String(),
		ID:            config.ID,
		Channel:       config.Channel,
		AccessToken:   config.AccessToken,
		RefreshToken:  config.RefreshToken,
		ExpiresIn:     expiresIn,
		DisabledUsers: string(bannedUsers),
		BannedWords:   string(bannedWords),
	}
}
