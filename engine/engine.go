package engine

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/parsaaes/jasoos-telegram-bot/room"
	"github.com/sirupsen/logrus"
	"log"
)

type Engine struct {
	Bot      *tgbotapi.BotAPI
	RoomList []*room.Room
}

func New(token string) (*Engine, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return &Engine{}, err
	}

	return &Engine{
		Bot:      bot,
		RoomList: make([]*room.Room, 1),
	}, nil
}

func (e *Engine) Run() {
	bot := e.Bot

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		logrus.Fatalf("engine: cannot get updates channel: %s", err.Error())
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		if _, err := bot.Send(msg); err != nil {
			logrus.Errorf("engine: cannot send message: %s", err.Error())
		}
	}
}
