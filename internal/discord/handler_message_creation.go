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

	if strings.HasPrefix(m.Content, "!") && len(m.Content) > 1 {
		command, args := parseArguments(m.Content)

		commandInfo := getCommandInfo(command)
		if commandInfo == nil {
			s.ChannelMessageSend(m.ChannelID, "Invalid command")
			return
		}

		if expected, valid := commandInfo.validateArguments(args); !valid {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Invalid arguments, expected %d got %d", expected, len(args)))
			return
		}

		if !managerExists(m.GuildID, true) {
			createManager(m.GuildID, true)
		}

		switch command {
    case "ping":
      ping(s, m)
    case "join":
			join(s, m)
		case "audio":
			audio(s, m, args[0])
		case "favoritecreate":
			favoritecreate(s, m, args[0], args[1])
		case "favoritelist":
			favoritelist(s, m)
		case "randomaudios":
			randomaudios(s, m, args[0])
		case "chatgpt":
			chatgpt(s, m, strings.Join(args[1:], " "), args[0] == "usectx")
		case "skip":
			skip(s, m)
		case "help":
			help(s, m)
		case "leave":
			leave(s, m)
		case "loop":
			loop(s, m)
		case "chatgpttts":
			chatgpttts(s, m, strings.Join(args, " "))
		case "image":
			image(s, m, strings.Join(args, " "))
		case "fncreate":
			fncreate(s, m, args[0], strings.Join(args[1:], " "))
		case "fnrun":
			fnrun(s, m, args[0], args[1:])
		case "magic8":
			magic8(s, m, strings.Join(args, " "))
		case "analyze":
			analyze(s, m)
		}
	}
}
