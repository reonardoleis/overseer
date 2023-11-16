package discord

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func Init(token string, clientId ...string) (*discordgo.Session, error) {
	discordClient, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("discord: error creating discord client", err)
		return nil, err
	}

	err = loadFavorites()
	if err != nil {
		log.Println("discord: error loading favorites", err)
		return nil, err
	}

	setupHandlers(discordClient)

	discordClient.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuildMembers)

	return discordClient, nil
}

func Start(cli *discordgo.Session) error {
	err := cli.Open()
	if err != nil {
		log.Println("discord: error opening connection", err)
		return err
	}

	sc := make(chan os.Signal, 1)
	<-sc

	return nil
}
