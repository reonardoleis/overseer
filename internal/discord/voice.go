package discord

import (
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	isPlayingLock = sync.Mutex{}
)

func playSound(vc *discordgo.VoiceConnection, buffer [][]byte) (err error) {
	isPlayingLock.Lock()

	time.Sleep(50 * time.Millisecond)
	vc.Speaking(true)

	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	vc.Speaking(false)

	time.Sleep(50 * time.Millisecond)

	isPlayingLock.Unlock()
	return nil
}
