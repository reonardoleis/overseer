package sound

import (
	"io"
	"log"

	"github.com/jonas747/dca"
)

func LoadSound(path string) ([][]byte, error) {
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"

	encodeSession, err := dca.EncodeFile(path, options)
	if err != nil {
		log.Println("discord: error encoding file:", err)
		return nil, err
	}

	defer encodeSession.Cleanup()

	var opuslen int16

	var buffer [][]byte

	for {
		opusFrame, err := encodeSession.OpusFrame()
		if err != nil {
			if err != io.EOF {
				log.Println("discord: error reading opus frame:", err)
			}
			break
		}

		buffer = append(buffer, opusFrame)

		opuslen += int16(len(opusFrame))
	}

	return buffer, nil
}
