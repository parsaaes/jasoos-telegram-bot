package engine

import (
	"encoding/json"
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

	go e.SendHandler()

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		logrus.Fatalf("engine: cannot get updates channel: %s", err.Error())
	}

	for update := range updates {
		if update.Message != nil {
			if update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup() {

				updateRoom, roomExists := e.RoomList[update.Message.Chat.ID]

				if update.CallbackQuery != nil {
					if roomExists {
						var callback message.Callback
						if err := json.Unmarshal([]byte(update.CallbackQuery.Data), &callback); err != nil {
							logrus.Errorf("engine: callback: cannot unmarshal update callback data: %s", err.Error())
						}

						switch callback.Type {
						case message.JoinCallbackType:
							var join message.JoinCallback

							if err := json.Unmarshal([]byte(update.CallbackQuery.Data), &join); err != nil {
								logrus.Errorf("engine: callback: cannot unmarshal update callback data (join): %s", err.Error())
							}

							updateRoom.Joined(join)

							break

						}
					}
				}

				switch update.Message.Command() {
				case message.New:
					if !roomExists {
						r := &room.Room{
							ChatID: update.Message.Chat.ID,
							State:  room.Join,
							Members: []*room.Member{&room.Member{
								Name: update.Message.From.String(),
								ID:   update.Message.From.ID,
							}},
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
}

func (e *Engine) SendHandler() {
	for msg := range e.SendChan {
		if _, err := e.Bot.Send(msg); err != nil {
			logrus.Errorf("engine: cannot send message: %s", err.Error())
		}
	}
}
