package services

import (
	"fmt"
	"log"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/maikpro/mangadownloader/models"
	"github.com/maikpro/mangadownloader/util"
)

func SendChapter(chapter models.Chapter) error {
	bot, err := getBot()
	if err != nil {
		log.Fatal(err)
		return err
	}

	sendMessage(bot, fmt.Sprintf("One Piece chapter: %d - %s - Pages: %d", chapter.Number, chapter.Name, len(chapter.Pages)))

	var mediaGroup []interface{}
	for _, page := range chapter.Pages {
		inputMediaPhoto := tgbotapi.NewInputMediaPhoto(page.Url)
		mediaGroup = append(mediaGroup, inputMediaPhoto)
	}

	err = sendMediaGroup(bot, mediaGroup)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func sleep(ms uint) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func getBot() (*tgbotapi.BotAPI, error) {
	telegramAPITokenString, err := getToken()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	bot, err := tgbotapi.NewBotAPI(telegramAPITokenString)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// bot.Debug = true
	// log.Printf("Authorized on account %s", bot.Self.UserName)

	// check for Telegram chat updates with goroutine
	go func() error {
		update := tgbotapi.NewUpdate(0)
		update.Timeout = 60

		updates, err := bot.GetUpdatesChan(update)
		if err != nil {
			log.Fatal(err)
			return err
		}

		for update := range updates {
			if update.Message == nil { // ignore any non-Message updates
				continue
			}

			if update.Message.Text != "" {
				messageText := update.Message.Text
				log.Printf("chatId: [%d] Username: [%s] message: %s", update.Message.Chat.ID, update.Message.From.UserName, messageText)
				listenForCommands(bot, messageText)
			}
		}
		return nil
	}()

	return bot, nil
}

func listenForCommands(bot *tgbotapi.BotAPI, messageText string) {
	switch messageText {
	case "!img":
		sendImage(bot, "./upload/image.jpg")
	case "!health":
		sendMessage(bot, "I'm alive! 🥳💚")
	default:
		sendMessage(bot, "no comment found!")
	}
}

func getToken() (string, error) {
	telegramAPITokenString, err := util.GetEnvString("TELEGRAM_BOT_TOKEN")
	if err != nil {
		return "", err
	}

	return telegramAPITokenString, nil
}

func getChatId() (int64, error) {
	chatIdString, err := util.GetEnvString("TELEGRAM_CHAT_ID")
	if err != nil {
		return 0, err
	}

	chatId, err := strconv.ParseInt(chatIdString, 10, 64)
	if err != nil {
		return 0, err
	}

	return chatId, nil
}

func sendMessage(bot *tgbotapi.BotAPI, text string) error {
	chatId, err := getChatId()
	if err != nil {
		log.Fatal(err)
		return err
	}

	msg := tgbotapi.NewMessage(chatId, text)
	bot.Send(msg)
	sleep(5500)
	return nil
}

func sendImage(bot *tgbotapi.BotAPI, fullPath string) error {
	log.Println("Sending Image...")

	chatId, err := getChatId()
	if err != nil {
		log.Fatal(err)
		return err
	}

	photoConfig := tgbotapi.NewPhotoUpload(chatId, fullPath)
	msg, err := bot.Send(photoConfig)
	if err != nil {
		log.Fatal("Error sending photo:", err)
	}
	log.Printf("Photo sent with message ID: %d", msg.MessageID)
	sleep(5500)
	return nil
}

func sendMediaGroup(bot *tgbotapi.BotAPI, mediaGroup []interface{}) error {
	chatId, err := getChatId()
	if err != nil {
		log.Fatal(err)
		return err
	}
	chunkSize := 10
	for i := 0; i < len(mediaGroup); i += chunkSize {
		end := i + chunkSize
		if end > len(mediaGroup) {
			end = len(mediaGroup)
		}
		chunk := mediaGroup[i:end]
		mediaGroupConfig := tgbotapi.NewMediaGroup(chatId, chunk)
		_, err := bot.Send(mediaGroupConfig)
		if err != nil {
			log.Fatal("Error sending mediaGroup:", err)
			return err
		}
		sleep(5500)
	}
	return nil
}
