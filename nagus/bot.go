package nagus

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

type NagusBot struct {
	Bot* tgbotapi.BotAPI
	Config NagusConfig
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

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		b.Bot.Send(msg)
	}
	return nil
}

func BuildBot(config NagusConfig) (NagusBot, error) {
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
	}, nil
}
