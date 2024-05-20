package domain

type TemplateSettingsData struct {
	Name          string
	Error         string
	DisabledUsers []string
	BannedWords   []string
	IsJoined      bool
}

func NewTemplateSettingsData(platform Platform, config *PlatformConfig, err error) *TemplateSettingsData {
	data := &TemplateSettingsData{Name: platform.String()}

	if err != nil {
		data.Error = err.Error()
	}

	if config != nil {
		data.DisabledUsers = config.DisabledUsers
		data.BannedWords = config.BannedWords
		data.IsJoined = config.IsJoined
	}

	return data
}
