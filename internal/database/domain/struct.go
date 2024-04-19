package domain

type PlatformName string

const (
	Youtube = PlatformName("youtube")
	Twitch  = PlatformName("twitch")
)

type PlatformKey struct {
	Name  PlatformName
	Value string
}

func (p PlatformKey) Key() string {
	return string(p.Name) + ":" + p.Value
}
