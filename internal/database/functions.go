package database

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/reonardoleis/overseer/internal/functions"
)

var (
	ErrFunctionAlreadyExists = errors.New("function already exists")
)

func CreateFunction(name, code string) error {
	locks[FUNCTIONS].Lock()
	defer locks[FUNCTIONS].Unlock()

	if strings.Contains(name, "\n") {
		code = strings.Split(name, "\n")[1] + " " + code
		name = strings.Split(name, "\n")[0]
	}

	if _, err := os.Stat("functions/" + name + ".js"); !errors.Is(err, os.ErrNotExist) {
		return ErrFunctionAlreadyExists
	}

	file, err := os.Create("functions/" + name + ".js")
	if err != nil {
		log.Println("discord: error creating function file", err)
		return err
	}

	code = functions.Code(name, code)
	_, err = file.WriteString(code)
	if err != nil {
		log.Println("discord: error writing to function file", err)
		return err
	}

	err = file.Close()
	if err != nil {
		log.Println("discord: error closing function file", err)
		return err
	}

	return nil
}
