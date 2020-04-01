package message

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

const (
	// New creates new room for game
	New = "new"
	// Join joins the sender to the room
	Join = "join"
	// Vote votes the member for being an spy
	Vote = "vote"
)

type Response struct {
	Message tgbotapi.Message
	Error   error
}

type Request struct {
	Chattable tgbotapi.Chattable
	Report    chan Response
}
