package discord

import (
	"log"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

type playableItemQueue struct {
	guildID string
	queue   chan playableItem
}

type playableItem struct {
	id     int
	alias  string
	buffer [][]byte
}

func (p playableItem) String() string {
	if p.id == 0 {
		return p.alias
	}

	return strconv.Itoa(p.id)
}

func (p *playableItemQueue) add(item playableItem) {
	p.queue <- item
}

func (p *playableItemQueue) audioPlayerWorker(vc *discordgo.VoiceConnection) {
	log.Println("starting audio player worker for guild", p.guildID)
	for {
		item := <-p.queue
		log.Println("playing audio", item.String(), "for guild", p.guildID)
		playAudio(vc, item.buffer)
	}
}

func playAudio(vc *discordgo.VoiceConnection, buffer [][]byte) (err error) {
	time.Sleep(50 * time.Millisecond)
	vc.Speaking(true)

	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	vc.Speaking(false)

	time.Sleep(50 * time.Millisecond)

	return nil
}
