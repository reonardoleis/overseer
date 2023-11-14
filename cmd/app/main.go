package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/reonardoleis/overseer/internal/discord"
)

func main() {
	err := godotenv.Overload(".env")
	if err != nil {
		log.Println("no .env file found")
	}

	cli, err := discord.Init(os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Println("error initializing discord client", err)
		return
	}

	err = cli.Start()
	if err != nil {
		log.Println("error starting discord client", err)
		return
	}
}
