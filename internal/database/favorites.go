package database

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

var (
	favorites     = make(map[string]string)
	favoritesLock = sync.Mutex{}
)

func GetFavorite(alias string) (string, bool) {
	v, ok := favorites[alias]
	return v, ok
}

func LoadFavorites() error {
	_, err := os.Stat("favorites.json")
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create("favorites.json")
			if err != nil {
				log.Println("discord: error creating favorites.json", err)
				return err
			}

			_, err = file.WriteString("{}")
			if err != nil {
				log.Println("discord: error writing to favorites.json", err)
				return err
			}

			return file.Close()
		} else {
			log.Println("discord: error reading favorites.json", err)
			return err
		}
	}

	content, err := os.ReadFile("favorites.json")
	if err != nil {
		log.Println("discord: error reading favorites.json", err)
		return err
	}

	err = json.Unmarshal(content, &favorites)
	if err != nil {
		log.Println("discord: error unmarshalling favorites.json", err)
		return err
	}

	return nil
}

func CreateFavorite(filename string, alias string) error {
	favoritesLock.Lock()
	defer favoritesLock.Unlock()

	favorites[alias] = filename

	content, err := json.Marshal(favorites)
	if err != nil {
		log.Println("discord: error marshalling favorites.json", err)
		return err
	}

	err = os.WriteFile("favorites.json", content, 0644)
	if err != nil {
		log.Println("discord: error writing to favorites.json", err)
		return err
	}

	return nil
}

func GetFormattedFavorites() string {
	s := "\n**Favorites**\n\n"
	for k, _ := range favorites {
		s += k + "\n"
	}

	return s
}
