package discord

import (
	"log"
	"math/rand"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/reonardoleis/overseer/internal/ai"
	"github.com/reonardoleis/overseer/internal/database"
	"github.com/reonardoleis/overseer/internal/sound"
	"github.com/reonardoleis/overseer/internal/utils"
)

func audio(s *discordgo.Session, m *discordgo.MessageCreate, idOrAlias string) error {
	var path string
	id, err := strconv.Atoi(idOrAlias)

	if err != nil {
		filename, exists := database.GetFavorite(idOrAlias)
		if !exists {
			s.ChannelMessageSend(m.ChannelID, "Invalid alias")
			return err
		}

		path = "audios/" + filename
	} else {
		filename, err := utils.GetFileByPosition("audios", id)
		if err != nil {
			log.Println("discord: error getting audio filename:", err)
			return err
		}

		path = "audios/" + filename
	}

	buf, err := sound.LoadSound(path)
	if err != nil {
		log.Println("discord: error loading sound:", err)
		return err
	}

	manager := getManager(m.GuildID)
	manager.audioQueue.add(playableItem{
		buffer: buf,
		id:     id,
		alias:  idOrAlias,
	})

	return nil
}

func favoritecreate(s *discordgo.Session, m *discordgo.MessageCreate, audioId, alias string) error {
	id, err := strconv.Atoi(audioId)
	if err != nil {
		log.Println("discord: error converting audio ID to int:", err)
		return err
	}

	if _, exists := database.GetFavorite(alias); exists {
		s.ChannelMessageSend(m.ChannelID, "Alias already exists")
		return nil
	}

	filename, err := utils.GetFileByPosition("audios", id)
	if err != nil {
		log.Println("discord: error getting audio filename:", err)
		return err
	}

	err = database.CreateFavorite(filename, alias)
	if err != nil {
		log.Println("discord: error creating favorite:", err)
		return err
	}

	s.ChannelMessageSend(m.ChannelID, "Favorite created")

	return nil
}

func join(s *discordgo.Session, m *discordgo.MessageCreate) (*discordgo.VoiceConnection, error) {
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return nil, err
	}

	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		return nil, err
	}

	var channelId string
	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {
			channelId = vs.ChannelID
			break
		}
	}

	vc, err := s.ChannelVoiceJoin(g.ID, channelId, false, false)
	if err != nil {
		return nil, err
	}

	manager := getManager(m.GuildID)
	manager.setVC(vc)

	go manager.audioPlayerWorker()

	return vc, nil
}

func favoritelist(s *discordgo.Session, m *discordgo.MessageCreate) error {
	formattedFavorites := database.GetFormattedFavorites()

	_, err := s.ChannelMessageSend(m.ChannelID, formattedFavorites)
	if err != nil {
		log.Println("discord: error sending message:", err)
		return err
	}

	return nil
}

func randomaudios(s *discordgo.Session, m *discordgo.MessageCreate, n string) error {
	count, err := utils.CountFolderFiles("audios")
	if err != nil {
		log.Println("discord: error counting folder files:", err)
		return err
	}

	numberOfAudios, err := strconv.Atoi(n)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Invalid number of audios")
		return err
	}

	for i := 0; i < numberOfAudios; i++ {
		random := rand.Intn(count)

		err = audio(s, m, strconv.Itoa(random))
		if err != nil {
			log.Println("discord: error playing audio:", err)
			continue
		}

	}

	return nil
}

func chatgpt(s *discordgo.Session, m *discordgo.MessageCreate, prompt string) error {
	text, err := ai.Generate(prompt)
	if err != nil {
		log.Println("discord: error generating text:", err)
		return err
	}

	_, err = s.ChannelMessageSend(m.ChannelID, text)
	if err != nil {
		log.Println("discord: error sending message:", err)
		return err
	}

	return nil
}
