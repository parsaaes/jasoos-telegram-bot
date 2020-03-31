package message

const (
	// New creates new room for game
	New = "new"
	// Join joins the sender to the room
	Join = "join"
)

type Message struct {
	ChatID int64
	Text   string
}
