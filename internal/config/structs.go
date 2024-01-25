package config

type Config struct {
	Twitch Twitch `json:"twitch"`
	Api    Api    `json:"api"`
}

type Twitch struct {
	Capabilities []string  `json:"capabilities"`
	Username     string    `json:"username"`
	IRCServer    IRCServer `json:"irc_server"`
	Oauth        Oauth     `json:"oauth"`
}

type IRCServer struct {
	Address  string `json:"address"`
	Protocol string `json:"protocol"`
	Origin   string `json:"origin"`
}

type Oauth struct {
	RefreshToken string   `json:"refresh_token"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Scopes       []string `json:"scopes"`
	RedirectURL  string   `json:"redirect_url"`
}

type Api struct {
	Host string `json:"host"`
	Port string `json:"port"`
}
