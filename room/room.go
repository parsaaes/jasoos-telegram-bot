package room

import (
	"fmt"
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/parsaaes/jasoos-telegram-bot/message"
)

// JoinDuration in seconds
const JoinDuration = 30

// Member represents a member of a game room
type Member struct {
	Name string
	ID   int
}

// Room is a place for playing
type Room struct {
	ChatID   int64
	State    State
	Members  []*Member
	SendChan chan tgbotapi.Chattable

	Words []string
}

func joinKeyboard() tgbotapi.InlineKeyboardMarkup {
	jmsg := message.Join
	return tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{
		{
			Text:         "Join",
			CallbackData: &jmsg,
		},
	})
}

// Created must be called on the room creation
func (r *Room) Created() {
	// first member of members slice is a user who creates the game
	msg := tgbotapi.NewMessage(
		r.ChatID,
		fmt.Sprintf("A game created by %s.", r.Members[0].Name),
	)
	msg.ReplyMarkup = joinKeyboard()

	r.SendChan <- msg

	go r.CountToStart()
	r.State = Join
}

// CountToStart conts 30 seconds then start the game
func (r *Room) CountToStart() {
	tick := time.NewTicker(1 * time.Second)
	count := 0

	for range tick.C {
		if count%10 == 0 {
			r.SendChan <- tgbotapi.NewMessage(
				r.ChatID,
				fmt.Sprintf("%d sec left to join", JoinDuration-count),
			)
		}
		count++
		if count == JoinDuration {
			tick.Stop()
			break
		}
	}

	r.SendChan <- tgbotapi.NewMessage(
		r.ChatID,
		"Let's play",
	)

	r.Inform()
}

// Inform is automatically called after the join duration and send the chosen word to members
func (r *Room) Inform() {
	spy := rand.Intn(len(r.Members))
	index := rand.Intn(len(r.Words))
	word := r.Words[index]

	for i := range r.Members {
		word := word

		if i == spy {
			word = "Spy"
		}

		msg := tgbotapi.NewMessage(
			int64(r.Members[i].ID),
			word,
		)
		r.SendChan <- msg
	}
}

// Joined must be called when a new member joined
func (r *Room) Joined(from *tgbotapi.User, base *tgbotapi.Message) {
	if r.State != Join {
		return
	}

	if from == nil || base == nil {
		return
	}

	r.Members = append(r.Members, &Member{
		Name: from.String(),
		ID:   from.ID,
	})

	msg := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    r.ChatID,
			MessageID: base.MessageID,
		},
		Text: fmt.Sprintf("%s\n\r- %s has joined.", base.Text, from.String()),
	}

	keyboard := joinKeyboard()
	msg.ReplyMarkup = &keyboard

	r.SendChan <- msg
}
