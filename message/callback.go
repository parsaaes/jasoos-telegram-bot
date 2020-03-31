package message

const (
	JoinCallbackType = "join"
)

type Callback struct {
	Type string `json:"type"`
}

type JoinCallback struct {
	Type        string `json:"type"`
	Username    string `json:"name"`
	UserID      int    `json:"id"`
	MessageID   int    `json:"message_id"`
	LastMessage string `json:"last_message"`
}
