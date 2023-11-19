package sound

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"

	"github.com/jonas747/dca"
	"github.com/tcolgate/mp3"
)

func LoadSound(path string) ([][]byte, error) {
	fmt.Println(path)
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

func GetDuration(path string) (int, error) {
	t := 0.0

	r, err := os.Open(path)
	if err != nil {
		return 0, err
	}

	d := mp3.NewDecoder(r)
	var f mp3.Frame
	skipped := 0

	for {

		if err := d.Decode(&f, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			return 0, err
		}

		t = t + f.Duration().Seconds()
	}

	return int(math.Ceil(t)), nil
}
