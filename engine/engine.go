package engine

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/parsaaes/jasoos-telegram-bot/message"
	"github.com/parsaaes/jasoos-telegram-bot/room"
	"github.com/sirupsen/logrus"
	"log"
)

type Engine struct {
	Bot      *tgbotapi.BotAPI
	RoomList map[int64]*room.Room
	SendChan chan message.Message
}

func New(token string) (*Engine, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return &Engine{}, err
	}

	return &Engine{
		Bot:      bot,
		RoomList: make(map[int64]*room.Room),
		SendChan: make(chan message.Message),
	}, nil
}

func (e *Engine) Run() {
	bot := e.Bot

	bot.Debug = true

	go e.SendHandler()

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

		if update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup() {
			switch update.Message.Command() {
			case message.New:
				if _, exists := e.RoomList[update.Message.Chat.ID]; !exists {
					r := &room.Room{
						ChatID:   update.Message.Chat.ID,
						State:    room.Join,
						Members:  []*tgbotapi.User{update.Message.From},
						SendChan: e.SendChan,
					}
					e.RoomList[update.Message.Chat.ID] = r

					r.Created()
				}
				break
			}
		}
	}
}

func (e *Engine) SendHandler() {
	for msg := range e.SendChan {
		msg := tgbotapi.NewMessage(msg.ChatID, msg.Text)

		if _, err := e.Bot.Send(msg); err != nil {
			logrus.Errorf("engine: cannot send message: %s", err.Error())
		}
	}
}
