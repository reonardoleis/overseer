package discord

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func playSound(vc *discordgo.VoiceConnection, buffer [][]byte) (err error) {
	time.Sleep(250 * time.Millisecond)
	vc.Speaking(true)

	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	vc.Speaking(false)

	time.Sleep(250 * time.Millisecond)

	return nil
}
