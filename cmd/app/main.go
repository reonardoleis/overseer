package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/reonardoleis/overseer/internal/ai"
	"github.com/reonardoleis/overseer/internal/database"
	"github.com/reonardoleis/overseer/internal/discord"
)

func main() {
	err := godotenv.Overload(".env")
	if err != nil {
		log.Println("no .env file found")
	}

	ai.Init(os.Getenv("OPENAI_KEY"))

	err = database.Connect()
	if err != nil {
		log.Println("error connecting to database", err)
		return
	}

	cli, err := discord.Init(os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Println("error initializing discord client", err)
		return
	}

	err = discord.Start(cli)
	if err != nil {
		log.Println("error starting discord client", err)
		return
	}
}
