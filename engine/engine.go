package engine

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/parsaaes/jasoos-telegram-bot/message"
	"github.com/parsaaes/jasoos-telegram-bot/room"
	"github.com/sirupsen/logrus"
)

type Engine struct {
	Bot      *tgbotapi.BotAPI
	RoomList map[int64]*room.Room
	SendChan chan tgbotapi.Chattable
}

func New(token string) (*Engine, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return &Engine{}, err
	}

	return &Engine{
		Bot:      bot,
		RoomList: make(map[int64]*room.Room),
		SendChan: make(chan tgbotapi.Chattable),
	}, nil
}

func (e *Engine) Run() {
	bot := e.Bot

	bot.Debug = true

	go e.Sender()

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		logrus.Fatalf("engine: cannot get updates channel: %s", err.Error())
	}

	for update := range updates {

		if update.CallbackQuery != nil {
			callback := update.CallbackQuery

			r, ok := e.RoomList[callback.Message.Chat.ID]
			if ok {
				switch callback.Data {
				case message.Join:
					r.Joined(callback.From, callback.Message)
				}

			}
		}

		if update.Message != nil {
			switch update.Message.Command() {
			case message.New:
				if _, ok := e.RoomList[update.Message.Chat.ID]; !ok {
					r := &room.Room{
						ChatID: update.Message.Chat.ID,
						State:  room.Join,
						Members: []*room.Member{
							&room.Member{
								Name: update.Message.From.String(),
								ID:   update.Message.From.ID,
							},
						},
						SendChan: e.SendChan,
					}
					e.RoomList[update.Message.Chat.ID] = r

					r.Created()
				}
			}
		}
	}
}

// Sender ranges over send channel and sends messages
func (e *Engine) Sender() {
	for msg := range e.SendChan {
		if _, err := e.Bot.Send(msg); err != nil {
			logrus.Errorf("engine: cannot send message: %s", err.Error())
		}
	}
}
