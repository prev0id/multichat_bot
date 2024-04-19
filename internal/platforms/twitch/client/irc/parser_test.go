package irc

//func Test_parseSingleMessage(t *testing.T) {
//	t.Parallel()
//
//	tests := []struct {
//		name     string
//		input    string
//		expected *domain.Message
//	}{
//		{
//			name:  "with tags",
//			input: "@badges=staff/1,broadcaster/1,turbo/1;color=#FF0000;display-name=PetsgomOO;emote-only=1;emotes=33:0-7;flags=0-7:A.6/P.6,25-36:A.1/I.2;id=c285c9ed-8b1b-4702-ae1c-c64d76cc74ef;mod=0;room-id=81046256;subscriber=0;turbo=0;tmi-sent-ts=1550868292494;user-id=81046256;user-type=staff :petsgomoo!petsgomoo@petsgomoo.tmi.twitch.tv PRIVMSG #petsgomoo :DansGame",
//			expected: &domain.Message{
//				tags: map[string]string{
//					"badges":       "staff/1,broadcaster/1,turbo/1",
//					"color":        "#FF0000",
//					"display-name": "PetsgomOO",
//					"emote-only":   "1",
//					"emotes":       "33:0-7",
//					"flags":        "0-7:A.6/P.6,25-36:A.1/I.2",
//					"id":           "c285c9ed-8b1b-4702-ae1c-c64d76cc74ef",
//					"mod":          "0",
//					"room-id":      "81046256",
//					"subscriber":   "0",
//					"turbo":        "0",
//					"tmi-sent-ts":  "1550868292494",
//					"user-id":      "81046256",
//					"user-type":    "staff",
//				},
//				source:     "petsgomoo!petsgomoo@petsgomoo.tmi.twitch.tv",
//				command:    "PRIVMSG #petsgomoo",
//				Parameters: "DansGame",
//			},
//		},
//		{
//			name:  "without tags",
//			input: ":lovingt3s!lovingt3s@lovingt3s.tmi.twitch.tv PRIVMSG #lovingt3s :!dilly",
//			expected: &domain.Message{
//				tags:       map[string]string{},
//				source:     "lovingt3s!lovingt3s@lovingt3s.tmi.twitch.tv",
//				command:    "PRIVMSG #lovingt3s",
//				Parameters: "!dilly",
//			},
//		},
//		{
//			name:  "ping msg",
//			input: "PING :tmi.twitch.tv",
//		},
//	}
//	for _, tt := range tests {
//		tt := tt
//		t.Run(tt.name, func(t *testing.T) {
//			t.Parallel()
//
//			actual := parseSingleMessage(tt.input)
//
//			assert.Equal(t, tt.expected, actual)
//		})
//	}
//}
