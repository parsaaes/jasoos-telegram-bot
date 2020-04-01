package engine

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/parsaaes/jasoos-telegram-bot/message"
	"github.com/parsaaes/jasoos-telegram-bot/room"
	"github.com/sirupsen/logrus"
)

// Engine is a game engine
type Engine struct {
	Bot      *tgbotapi.BotAPI
	RoomList map[int64]*room.Room
	SendChan chan tgbotapi.Chattable
}

// New creates a new game engine
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

// Run runs the game engine
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
			e.handleCallback(update.CallbackQuery)
		}

		if update.Message != nil {
			e.handleMessage(update.Message)
		}
	}
}

func (e *Engine) handleMessage(msg *tgbotapi.Message) {
	if !msg.Chat.IsGroup() && !msg.Chat.IsSuperGroup() {
		return
	}

	switch msg.Command() {
	case message.New:
		if _, ok := e.RoomList[msg.Chat.ID]; !ok {
			r := &room.Room{
				ChatID: msg.Chat.ID,
				State:  room.Join,
				Members: []*room.Member{
					&room.Member{
						Name: msg.From.String(),
						ID:   msg.From.ID,
					},
				},
				SendChan: e.SendChan,
			}
			e.RoomList[msg.Chat.ID] = r

			r.Created()
		}
	}
}

func (e *Engine) handleCallback(callback *tgbotapi.CallbackQuery) {
	r, ok := e.RoomList[callback.Message.Chat.ID]
	if ok {
		switch callback.Data {
		case message.Join:
			r.Joined(callback.From, callback.Message)
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
