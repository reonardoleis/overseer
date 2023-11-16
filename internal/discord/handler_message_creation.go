package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func handleMessageCreation(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!") {
		command, args := parseArguments(m.Content)

		argumentInfo := commandsArgc[command]
		if argumentInfo == nil {
			s.ChannelMessageSend(m.ChannelID, "Invalid command")
			return
		}

		if expected, valid := argumentInfo.validateArguments(args); !valid {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Invalid arguments, expected %d got %d", expected, len(args)))
			return
		}

		switch command {
		case "join":
			joinVoiceChannel(s, m)
		case "audio":
			playAudio(s, m, args[0])
		case "favoritecreate":
			favoriteAudio(s, m, args[0], args[1])
		case "favoritelist":
			getFavorites(s, m)
		case "randomaudios":
			playRandomAudios(s, m, args[0])
		}
	}
}
