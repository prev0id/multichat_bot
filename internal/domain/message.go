package domain

type Message struct {
	From     string
	Text     string
	Channel  string
	Platform Platform
}
