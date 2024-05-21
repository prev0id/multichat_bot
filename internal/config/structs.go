package config

type Config struct {
	API     API         `json:"api"`
	Twitch  Twitch      `json:"twitch"`
	Cookie  CookieStore `json:"cookie"`
	Youtube Youtube     `json:"youtube"`
	DB      DB          `json:"db"`
	Auth    Auth        `json:"auth"`
}

type Twitch struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type Youtube struct {
	APIKey string `json:"api_key"`
}

type API struct {
	Port   string `json:"port"`
	Secret string `json:"secret"`
}

type DB struct {
	DBPath string `json:"db_path"`
}

type Auth struct {
	Twitch  AuthProvider `json:"twitch"`
	Youtube AuthProvider `json:"youtube"`
	IsProd  bool         `json:"is_prod"`
}

type CookieStore struct {
	SessionKey string `json:"session_key"`
	IsProd     bool   `json:"is_prod"`
}

type AuthProvider struct {
	ClientKey    string   `json:"client_key"`
	ClientSecret string   `json:"client_secret"`
	CallbackURL  string   `json:"callback_url"`
	Scopes       []string `json:"scopes"`
}
