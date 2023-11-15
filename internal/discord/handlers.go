package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	handlers = []func(*discordgo.Session, *discordgo.MessageCreate){
		handleMessageCreation,
	}
)

func setupHandlers(s *discordgo.Session) {
	for _, handler := range handlers {
		s.AddHandler(handler)
	}
}

func handleMessageCreation(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// check if the message is "!airhorn"
	if strings.HasPrefix(m.Content, "!airhorn") {

		// Find the channel that the message came from.
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			// Could not find channel.
			return
		}

		// Find the guild for that channel.
		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			// Could not find guild.
			return
		}

		buff, err := loadSound()
		if err != nil {
			fmt.Println("Error loading sound: ", err)
			return
		}

		// Look for the message sender in that guild's current voice states.
		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				err = playSound(s, buff, g.ID, vs.ChannelID)
				if err != nil {
					fmt.Println("Error playing sound:", err)
				}

				return
			}
		}
	}
}
