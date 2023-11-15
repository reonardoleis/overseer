package discord

import (
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
