package room

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/parsaaes/jasoos-telegram-bot/message"
)

type Room struct {
	ChatID   int64
	State    State
	Members  []*tgbotapi.User
	SendChan chan message.Message
}

func (r Room) Created() {
	createMsg := fmt.Sprintf("A game created by %s.", r.Members[0].String())
	r.SendChan <- message.Message{
		ChatID: r.ChatID,
		Text:   createMsg,
	}
}
