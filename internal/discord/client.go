package discord

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

type Client struct {
	*discordgo.Session
}

func Init(token string, clientId ...string) (*Client, error) {
	err := createCommands()
	if err != nil {
		log.Println("discord: error creating commands", err)
		return nil, err
	}

	discordClient, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("discord: error creating discord client", err)
		return nil, err
	}

	setupHandlers(discordClient)

	discordClient.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuildMembers)

	return &Client{
		discordClient,
	}, nil
}

func (c *Client) Start() error {
	err := c.Open()
	if err != nil {
		log.Println("discord: error opening connection", err)
		return err
	}

	sc := make(chan os.Signal, 1)
	<-sc

	return nil
}
