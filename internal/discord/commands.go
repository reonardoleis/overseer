package discord

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	openai "github.com/sashabaranov/go-openai"

	"github.com/reonardoleis/overseer/internal/ai"
	"github.com/reonardoleis/overseer/internal/database"
	"github.com/reonardoleis/overseer/internal/database/models"
	"github.com/reonardoleis/overseer/internal/database/repository"
	"github.com/reonardoleis/overseer/internal/functions"
	"github.com/reonardoleis/overseer/internal/prompts"
	"github.com/reonardoleis/overseer/internal/sound"
	"github.com/reonardoleis/overseer/internal/utils"
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

func chatgpt(
	s *discordgo.Session,
	m *discordgo.MessageCreate,
	prompt string,
	useContext bool,
) error {
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
		manager.addChatGptContext(
			ai.MessageContext{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
			ai.MessageContext{
				Role:    openai.ChatMessageRoleAssistant,
				Content: text,
			},
		)
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

func leave(_ *discordgo.Session, m *discordgo.MessageCreate) error {
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

	_, err := s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf("Looping is now %t", manager.audioQueue.loop),
	)
	if err != nil {
		log.Println("discord: error sending message:", err)
		return err
	}

	return nil
}

func chatgpttts(s *discordgo.Session, m *discordgo.MessageCreate, prompt string) error {
	messageResp, err := s.ChannelMessageSend(m.ChannelID, "Generating audio...")
	if err != nil {
		log.Println("discord: error sending message:", err)
		return err
	}

	text, err := ai.Generate(prompt, []ai.MessageContext{}, 500)
	if err != nil {
		log.Println("discord: error generating text:", err)
		return err
	}

	audio, err := ai.TTS(text)
	if err != nil {
		log.Println("discord: error generating audio:", err)
		return err
	}

	file, err := os.Create(utils.GetPath(utils.AUDIOS_PATH, "aitts_"+uuid.New().String()+".mp3"))
	if err != nil {
		log.Println("discord: error creating file:", err)
		return err
	}

	_, err = io.Copy(file, audio)
	if err != nil {
		log.Println("discord: error copying audio to file:", err)
		return err
	}

	buff, err := sound.LoadSound(file.Name())
	if err != nil {
		log.Println("discord: error loading sound:", err)
		return err
	}

	file.Seek(0, 0)
	editContent := "Generating audio... Done!\nContent:\n\n" + text
	editedMessage := &discordgo.MessageEdit{
		Content: &editContent,
		ID:      messageResp.ID,
		Channel: m.ChannelID,
		Files: []*discordgo.File{
			{
				Name:        file.Name(),
				ContentType: "audio/mpeg",
				Reader:      file,
			},
		},
	}

	_, err = s.ChannelMessageEditComplex(editedMessage)
	if err != nil {
		log.Println("discord: error editing message:", err)
		return err
	}

	err = file.Close()
	if err != nil {
		log.Println("discord: error closing file:", err)
		return err
	}

	manager := getManager(m.GuildID)
	manager.audioQueue.add(&playableItem{
		buffer: buff,
	})

	return nil
}

func image(s *discordgo.Session, m *discordgo.MessageCreate, prompt string) error {
	messageResp, err := s.ChannelMessageSend(m.ChannelID, "Generating image...")
	if err != nil {
		log.Println("discord: error sending message:", err)
		return err
	}

	optimizedPrompt, err := ai.Generate(
		"Optmize this prompt for image generation, answer with the optimized prompt only: "+prompt,
		[]ai.MessageContext{},
		500,
	)
	if err != nil {
		log.Println("discord: error generating optimized prompt:", err)
		return err
	}

	image, err := ai.Image(optimizedPrompt)
	if err != nil {
		log.Println("discord: error generating image:", err)
		return err
	}

	editContent := "Generating image... Done!"
	editedMessage := &discordgo.MessageEdit{
		Content: &editContent,
		ID:      messageResp.ID,
		Channel: m.ChannelID,
		Files: []*discordgo.File{
			{
				Name:        strings.ReplaceAll(prompt, " ", "_") + ".png",
				ContentType: "image/png",
				Reader:      image,
			},
		},
	}

	_, err = s.ChannelMessageEditComplex(editedMessage)
	if err != nil {
		log.Println("discord: error editing message:", err)
		return err
	}

	return nil
}

func fncreate(s *discordgo.Session, m *discordgo.MessageCreate, name, code string) error {
	if valid := functions.Validate(code); !valid {
		_, err := s.ChannelMessageSend(m.ChannelID, "Invalid function")
		if err != nil {
			log.Println("discord: error sending message:", err)
			return err
		}
	}

	if strings.Contains(name, "\n") {
		code = strings.Split(name, "\n")[1] + " " + code
		name = strings.Split(name, "\n")[0]
	}

	existed, err := repository.
		Prepare(m.GuildID, m.Author.ID).
		CreateFunction(&models.Function{
			Name: name,
			Code: functions.Code(name, code),
		})
	if err != nil {
		log.Println("discord: error creating function:", err)
		return err
	}

	if existed {
		s.ChannelMessageSend(m.ChannelID, "A function with the same name already exists")
		return nil
	}

	_, err = s.ChannelMessageSend(m.ChannelID, "Function created")
	if err != nil {
		log.Println("discord: error sending message:", err)
		return err
	}

	return nil
}

func fnrun(s *discordgo.Session, m *discordgo.MessageCreate, name string, args []string) error {
	function, exists, err := repository.Prepare(m.GuildID, m.Author.ID).GetFunction(name)
	if err != nil {
		log.Println("discord: error finding function: ", err)
		return err
	}

	if !exists {
		s.ChannelMessageSend(m.ChannelID, "Function not found")
		return nil
	}

	output, err := functions.Run(function.Code, args)
	if err != nil {
		log.Println("discord: error running function:", err)
		return err
	}

	_, err = s.ChannelMessageSend(m.ChannelID, output)
	if err != nil {
		log.Println("discord: error sending message:", err)
		return err
	}

	return nil
}

func magic8(s *discordgo.Session, m *discordgo.MessageCreate, prompt string) error {
	messageResp, err := s.ChannelMessageSend(m.ChannelID, "Generating text...")
	if err != nil {
		log.Println("discord: error sending message:", err)
		return err
	}

	coinflip := rand.Intn(2)
	var work string
	if coinflip == 0 {
		work = "yes"
	} else {
		work = "no"
	}

	text, err := ai.Generate(fmt.Sprintf(prompts.Magic8, work, prompt), []ai.MessageContext{}, 500)
	if err != nil {
		log.Println("discord: error generating text:", err)
		return err
	}

	_, err = s.ChannelMessageEdit(m.ChannelID, messageResp.ID, text)
	if err != nil {
		log.Println("discord: error editing message:", err)
		return err
	}

	return nil
}

func analyze(s *discordgo.Session, m *discordgo.MessageCreate) error {
	messages, err := s.ChannelMessages(m.ChannelID, 10, m.ID, "", "")
	if err != nil {
		log.Println("discord: error getting messages:", err)
		return err
	}

	if len(messages) == 0 {
		return nil
	}

	message := ""
	for _, mm := range messages {
		if mm.Author.ID == m.Author.ID &&
			!strings.HasPrefix(mm.Content, "!") &&
			len(mm.Content) > 10 {
			message = mm.Content
			break
		}
	}

	if message == "" {
		_, err = s.ChannelMessageSend(m.ChannelID, "No suitable message found")
		if err != nil {
			log.Println("discord: error sending message:", err)
			return err
		}

		return nil
	}

	response, err := s.ChannelMessageSend(m.ChannelID, "Analyzing text...")
	if err != nil {
		log.Println("discord: error sending message:", err)
		return err
	}

	generated, err := ai.Generate(fmt.Sprintf(prompts.Analyze, message), []ai.MessageContext{}, 500)
	if err != nil {
		log.Println("discord: error generating text:", err)
		return err
	}

	_, err = s.ChannelMessageEdit(m.ChannelID, response.ID, generated)
	if err != nil {
		log.Println("discord: error sending message:", err)
		return err
	}

	return nil
}

func ping(s *discordgo.Session, m *discordgo.MessageCreate) error {
  _, err := s.ChannelMessageSend(m.ChannelID, "pong")
  if err != nil {
    log.Println("discord: error sending message: ", err)
    return err
  }

  return nil
}
