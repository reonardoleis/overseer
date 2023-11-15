package discord

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/reonardoleis/overseer/internal/sound"
)

func handleMessageCreation(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!airhorn") {
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			return
		}

		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			return
		}

		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				buff, err := sound.LoadSound("ding.mp3")
				if err != nil {
					log.Println("discord: error loading sound:", err)
					return
				}

				err = playSound(s, buff, g.ID, vs.ChannelID)
				if err != nil {
					log.Println("discord: error playing sound:", err)
					return
				}

			}
		}
	}
}
