package domain

type User struct {
	Platforms   map[Platform]*PlatformConfig
	AccessToken string
	ID          int64
}
