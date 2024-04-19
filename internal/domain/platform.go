package domain

type Platform string

const (
	Twitch  Platform = "twitch"
	YouTube Platform = "youtube"
)

var (
	AllPlatforms = []Platform{
		Twitch,
		YouTube,
	}
)
