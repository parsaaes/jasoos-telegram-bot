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
	SendChan chan tgbotapi.MessageConfig
}

// Created must be called on the romm creation
func (r Room) Created() {
	// first member of members slice is a user who creates the game
	msg := tgbotapi.NewMessage(
		r.ChatID,
		fmt.Sprintf("A game created by %s.", r.Members[0].String()),
	)

	jmsg := message.Join
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{
		tgbotapi.InlineKeyboardButton{
			Text:         "Join",
			CallbackData: &jmsg,
		},
	})

	r.SendChan <- msg
}
