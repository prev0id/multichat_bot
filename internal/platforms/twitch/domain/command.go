package domain

const (
	IRCCommandJoin       = "JOIN"
	IRCCommandPart       = "PART"
	IRCCommandNotice     = "NOTICE"
	IRCCommandClearChat  = "CLEARCHAT"
	IRCCommandHostTarget = "HOSTTARGET"
	IRCCommandPrivmsg    = "PRIVMSG"
	IRCCommandPing       = "PING"
	IRCCommandCap        = "CAP"
	IRCCommandUserState  = "USERSTATE"
	IRCCommandRoomState  = "ROOMSTATE"
	IRCCommand001        = "001"
	IRCCommand366        = "366"
)

type Message struct {
	Tags       map[string]string
	Parameters string
	RawMessage string
	RawSource  string
	Command    *Command
}

type Command struct {
	Name                string
	Channel             string
	IsCapRequestEnabled bool
	RawCommand          string

	BotCommand       string
	BotCommandParams string
	RawBotCommand    string
}

type Source struct {
	Nick      string
	Host      string
	RawSource string
}
