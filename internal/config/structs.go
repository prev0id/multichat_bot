package config

const (
	TwitchClient = "twitch"
)

type Application struct {
	Clients     map[string]Client `json:"clients"`
	Credentials Credentials       `json:"credentials"`
}

type Client struct {
	Address  string `json:"address"`
	Protocol string `json:"protocol"`
	Origin   string `json:"origin"`
}

type Credentials struct {
	UserName string `json:"user_name"`
	Token    string `json:"token"`
}
