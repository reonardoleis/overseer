package discord

import (
	"log"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

type playableItemQueue struct {
	guildID string
	loop    bool
	current *playableItem
	queue   chan *playableItem
}

type playableItem struct {
	id     int
	alias  string
	buffer [][]byte
	skip   bool
}

func (p playableItem) String() string {
	if p.id == 0 {
		return p.alias
	}

	return strconv.Itoa(p.id)
}

func (p *playableItemQueue) add(item *playableItem) {
	p.queue <- item
}

func (p *playableItemQueue) skip() {
	p.current.skip = true
}

func (p *playableItemQueue) audioPlayerWorker(vc *discordgo.VoiceConnection) {
	log.Println("starting audio player worker for guild", p.guildID)
	for {
		item := <-p.queue
		log.Println("playing audio", item.String(), "for guild", p.guildID)
		p.current = item
		p.playAudio(vc, item)
	}
}

func (p *playableItemQueue) playAudio(vc *discordgo.VoiceConnection, pi *playableItem) (err error) {
	time.Sleep(50 * time.Millisecond)
	vc.Speaking(true)

	for {
		for _, buff := range pi.buffer {
			if pi.skip {
				break
			}

			vc.OpusSend <- buff
		}

		if !p.loop {
			break
		}
	}

	vc.Speaking(false)
	time.Sleep(50 * time.Millisecond)

	return nil
}
