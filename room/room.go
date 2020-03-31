package room

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/parsaaes/jasoos-telegram-bot/message"
)

type Member struct {
	Name string
	ID   int
}

type Room struct {
	ChatID   int64
	State    State
	Members  []*Member
	SendChan chan tgbotapi.Chattable
}

// Created must be called on the room creation
func (r Room) Created() {
	// first member of members slice is a user who creates the game
	msg := tgbotapi.NewMessage(
		r.ChatID,
		fmt.Sprintf("A game created by %s.", r.Members[0].Name),
	)

	jmsg := message.Join
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{
		{
			Text:         "Join",
			CallbackData: &jmsg,
		},
	})

	r.SendChan <- msg
}

func (r Room) Joined(join message.JoinCallback) {
	r.Members = append(r.Members, &Member{
		Name: join.Username,
		ID:   join.UserID,
	})

	msg := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    r.ChatID,
			MessageID: join.MessageID,
		},
		Text: fmt.Sprintf("%s\n %s has joined.", join.LastMessage, join.Username),
	}

	r.SendChan <- msg
}
