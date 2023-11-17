package discord

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/reonardoleis/overseer/internal/ai"
	"github.com/reonardoleis/overseer/internal/database"
	"github.com/reonardoleis/overseer/internal/sound"
	"github.com/reonardoleis/overseer/internal/utils"
	openai "github.com/sashabaranov/go-openai"
)

func audio(s *discordgo.Session, m *discordgo.MessageCreate, idOrAlias string) error {
	var path string
	id, err := strconv.Atoi(idOrAlias)

	var filename string
	if err != nil {
		var exists bool
		filename, exists = database.GetFavorite(idOrAlias)
		if !exists {
			s.ChannelMessageSend(m.ChannelID, "Invalid alias")
			return err
		}
	} else {
		filename, err = utils.GetFileByPosition("audios", id)
		if err != nil {
			log.Println("discord: error getting audio filename:", err)
			return err
		}
	}

	path = utils.GetPath(utils.AUDIOS_PATH, filename)

	buf, err := sound.LoadSound(path)
	if err != nil {
		log.Println("discord: error loading sound:", err)
		return err
	}

	manager := getManager(m.GuildID)
	manager.audioQueue.add(&playableItem{
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

func chatgpt(s *discordgo.Session, m *discordgo.MessageCreate, prompt string, useContext bool) error {
	gptContext := []ai.MessageContext{}
	var manager *Manager
	if useContext {
		log.Println("using context")
		manager = getManager(m.GuildID)
		gptContext = manager.getChatGptContext()
	}

	text, err := ai.Generate(prompt, gptContext)
	if err != nil {
		log.Println("discord: error generating text:", err)
		return err
	}

	if useContext {
		manager.addChatGptContext(ai.MessageContext{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		})

		manager.addChatGptContext(ai.MessageContext{
			Role:    openai.ChatMessageRoleAssistant,
			Content: text,
		})
	}

	_, err = s.ChannelMessageSend(m.ChannelID, text)
	if err != nil {
		log.Println("discord: error sending message:", err)
		return err
	}

	return nil
}

func skip(s *discordgo.Session, m *discordgo.MessageCreate) error {
	manager := getManager(m.GuildID)
	manager.audioQueue.skip()

	_, err := s.ChannelMessageSend(m.ChannelID, "Skipping")
	if err != nil {
		log.Println("discord: error sending message:", err)
		return err
	}

	return nil
}

func help(s *discordgo.Session, m *discordgo.MessageCreate) error {
	_, err := s.ChannelMessageSend(m.ChannelID, commandHelp())
	if err != nil {
		log.Println("discord: error sending message:", err)
		return err
	}

	return nil
}

func leave(s *discordgo.Session, m *discordgo.MessageCreate) error {
	manager := getManager(m.GuildID)

	err := manager.vc.Disconnect()
	if err != nil {
		log.Println("discord: error disconnecting from voice channel:", err)
		return err
	}

	removeManager(m.GuildID)

	return nil
}

func loop(s *discordgo.Session, m *discordgo.MessageCreate) error {
	manager := getManager(m.GuildID)
	manager.audioQueue.loop = !manager.audioQueue.loop

	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Looping is now %t", manager.audioQueue.loop))
	if err != nil {
		log.Println("discord: error sending message:", err)
		return err
	}

	return nil
}

func chatgpttts(s *discordgo.Session, m *discordgo.MessageCreate, prompt string) error {
	_, err := s.ChannelMessageSend(m.ChannelID, "WIP")
	if err != nil {
		log.Println("discord: error sending message:", err)
		return err
	}

	return nil
}
