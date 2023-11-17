package utils

import (
	"errors"
	"os"
)

func GetFileByPosition(basePath string, position int) (string, error) {
	files, err := os.ReadDir("audios")
	if err != nil {
		return "", err
	}

	if len(files) < position {
		return "", errors.New("invalid file position")
	}

	file := files[position-1]

	return file.Name(), nil
}

func CountFolderFiles(basePath string) (int, error) {
	files, err := os.ReadDir(basePath)
	if err != nil {
		return 0, err
	}

	return len(files), nil
}

func DeleteFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

type BasePath string

const (
	AUDIOS_PATH BasePath = "audios"
	TTS_PATH    BasePath = "tts"
)

func GetPath(base BasePath, path string) string {
	return string(base) + "/" + path
}
