package nagus

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"image/png"
	"os"
	"math/rand"
	"bytes"
)

type NagusBot struct {
	Bot 	*tgbotapi.BotAPI
	Images 	*ImageManager
	Config 	NagusConfig
}

func (b *NagusBot) Main() error {
	b.Bot.Debug = true
	log.Printf("Authorized on account %s", b.Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.Bot.GetUpdatesChan(u)

	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s} %s", update.Message.From.UserName, update.Message.Text)
		buf := bytes.Buffer{}
		buf.WriteString("./")
		buf.WriteString(string(rand.Int63()))
		buf.WriteString(".png")
		out, err := os.Create(buf.String())
		img := b.Images.WriteToImage(update.Message.Text)
		err = png.Encode(out, img)
		if err != nil {
			log.Println(err)
		}
		out.Sync()
		out.Close()
		photoConf := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, buf.String())
		photoConf.ReplyToMessageID = update.Message.MessageID
		_,err = b.Bot.Send(photoConf)
		if err != nil {
			log.Println(err)
		}
		os.Remove(buf.String())
	}
	return nil
}

func BuildBot(config NagusConfig, confFolder string) (NagusBot, error) {
	bot, err := tgbotapi.NewBotAPI(config.ApiKey)
	if err != nil {
		return NagusBot{
			Config:config,
			Bot: nil,
		}, err
	}
	return NagusBot{
		Config:config,
		Bot: bot,
		Images: NewImageManager(confFolder),
	}, nil
}
