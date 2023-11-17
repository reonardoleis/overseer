package discord

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/pion/webrtc/v3/pkg/media/oggwriter"
)

func createPionRTPPacket(p *discordgo.Packet) *rtp.Packet {
	return &rtp.Packet{
		Header: rtp.Header{
			Version:        2,
			PayloadType:    0x78,
			SequenceNumber: p.Sequence,
			Timestamp:      p.Timestamp,
			SSRC:           p.SSRC,
		},
		Payload: p.Opus,
	}
}

func handleVoice(c chan *discordgo.Packet, manager *Manager) {
	files := make(map[uint32]media.Writer)
	for p := range c {
		manager.lastInteractionTime = time.Now()
		file, ok := files[p.SSRC]
		if !ok {
			var err error
			file, err = oggwriter.New(fmt.Sprintf("%d.ogg", p.SSRC), 48000, 2)
			if err != nil {
				log.Println("discord: failed to create file from received voice", err)
				return
			}
			files[p.SSRC] = file
		}

		rtp := createPionRTPPacket(p)
		err := file.WriteRTP(rtp)
		if err != nil {
			log.Println("discord: failed to write voice to discord rtp", err)
		}
	}

	for _, f := range files {
		f.Close()
	}
}
