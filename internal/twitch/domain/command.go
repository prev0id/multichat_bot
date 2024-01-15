package domain

const (
	IRCCommandJoin            = "JOIN"
	IRCCommandPart            = "PART"
	IRCCommandNotice          = "NOTICE"
	IRCCommandClearChat       = "CLEARCHAT"
	IRCCommandHostTarget      = "HOSTTARGET"
	IRCCommandPrivmsg         = "PRIVMSG"
	IRCCommandPing            = "PING"
	IRCCommandCap             = "CAP"
	IRCCommandGlobalUserState = "GLOBALUSERSTATE"
	IRCCommandUserState       = "USERSTATE"
	IRCCommandRoomState       = "ROOMSTATE"
	IRCCommandReconnect       = "RECONNECT"
	IRCCommand421             = "421"
	IRCCommand001             = "001"
	IRCCommand002             = "002"
	IRCCommand003             = "003"
	IRCCommand004             = "004"
	IRCCommand353             = "353"
	IRCCommand366             = "366"
	IRCCommand372             = "372"
	IRCCommand375             = "375"
	IRCCommand376             = "376"

	IRCCommandClearMessage = "CLEARMSG"
	IRCCommandUserNotice   = "USERNOTICE"
)

type Message struct {
	Tags       map[string]string
	Parameters string
	RawMessage string
	Source     *Source
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
