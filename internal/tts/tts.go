package tts

import (
	htgotts "github.com/hegedustibor/htgo-tts"
)

func Speak(text string) (string, error) {
	speech := htgotts.Speech{Folder: "tts", Language: "pt-BR"}
	return speech.Save(text)
}
