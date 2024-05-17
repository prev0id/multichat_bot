package domain

type Platform string

func (p Platform) String() string {
	return string(p)
}

const (
	Twitch  Platform = "twitch"
	YouTube Platform = "youtube"
)

var (
	StringToPlatform = map[string]Platform{
		Twitch.String():  Twitch,
		YouTube.String(): YouTube,
	}
	Platforms = []Platform{Twitch, YouTube}
)
