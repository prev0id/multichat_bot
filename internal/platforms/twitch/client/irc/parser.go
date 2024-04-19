package irc

import (
	"strings"

	"multichat_bot/internal/platforms/twitch/domain"
)

func parse(msg string) []*domain.Message {
	messages := strings.Split(msg, "\r\n")

	result := make([]*domain.Message, 0, len(messages))
	for _, message := range messages {
		result = append(result, parseSingleMessage(message))
	}

	return result
}

func parseSingleMessage(msg string) *domain.Message {
	if msg == "" {
		return nil
	}

	var (
		rawTags       string
		rawSource     string
		rawCommand    string
		rawBotCommand string

		rawMessage = msg
	)

	idx := 0
	if msg[idx] == '@' {
		endIdx := strings.IndexByte(msg, ' ')
		rawTags = msg[idx+1 : endIdx]
		msg = msg[endIdx+1:]
	}

	if msg[idx] == ':' {
		endIdx := strings.IndexByte(msg, ' ')
		rawSource = msg[idx+1 : endIdx]
		msg = msg[endIdx+1:]
	}

	endIdx := strings.IndexByte(msg, ':')
	if endIdx == -1 {
		endIdx = len(msg)
	}
	rawCommand = strings.TrimSpace(msg[idx:endIdx])

	if endIdx != len(msg) {
		idx = endIdx + 1
		rawBotCommand = msg[idx:]
	}

	return &domain.Message{
		Tags:       parseTags(rawTags),
		Command:    parseCommands(rawCommand, rawBotCommand),
		Parameters: rawBotCommand,
		RawSource:  rawSource,
		RawMessage: rawMessage,
	}
}

func parseTags(rawTags string) map[string]string {
	separated := strings.Split(rawTags, ";")
	parsedTags := make(map[string]string, len(separated))

	for _, rawTag := range separated {
		keyValue := strings.Split(rawTag, "=")
		if len(keyValue) != 2 {
			continue
		}

		parsedTags[keyValue[0]] = keyValue[1]
	}

	return parsedTags
}

func parseCommands(rawCommand, rawBotCommand string) *domain.Command {
	command := &domain.Command{}

	parseCommand(rawCommand, command)
	parseBotCommand(rawBotCommand, command)

	return command
}

func parseCommand(rawCommand string, command *domain.Command) {
	commandParts := strings.Split(rawCommand, " ")

	command.Name = commandParts[0]
	command.RawCommand = rawCommand

	switch commandParts[0] {
	case domain.IRCCommandJoin,
		domain.IRCCommandPart,
		domain.IRCCommandNotice,
		domain.IRCCommandClearChat,
		domain.IRCCommandHostTarget,
		domain.IRCCommandPrivmsg,
		domain.IRCCommandUserState,
		domain.IRCCommandRoomState,
		domain.IRCCommand001:

		command.Channel = commandParts[1]

	case domain.IRCCommandCap:
		if commandParts[2] == "ACK" {
			command.IsCapRequestEnabled = true
		}
	}
}

func parseBotCommand(rawBotCommand string, command *domain.Command) {
	commandParts := strings.TrimSpace(rawBotCommand)
	command.RawBotCommand = commandParts

	paramIdx := strings.Index(commandParts, " ")

	if paramIdx == -1 {
		command.BotCommand = commandParts
		return
	}

	command.BotCommand = commandParts[0:paramIdx]
	command.BotCommandParams = strings.TrimSpace(commandParts[paramIdx:])
}
