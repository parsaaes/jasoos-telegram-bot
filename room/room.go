package room

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/parsaaes/jasoos-telegram-bot/message"
)

// JoinDuration in seconds
const JoinDuration = 60

// DiscussionDuration in minutes
const DiscussionDuration = 3

// VoteDuration in minutes
const VoteDuration = 3

// Member represents a member of a game room
type Member struct {
	Name string
	ID   int
}

// Room is a place for playing
type Room struct {
	ChatID                int64
	Title                 string
	State                 State
	Members               []*Member
	SendChan              chan message.Request
	BotStartedWarningSent map[int]bool

	Words []string
	Spy   string
	Done  chan<- struct{}

	Votes map[string]int
}

func New(msg *tgbotapi.Message,
	sendChan chan message.Request, doneChan chan<- struct{},
	words []string) *Room {
	return &Room{
		ChatID: msg.Chat.ID,
		Title:  msg.Chat.Title,
		State:  Join,
		Members: []*Member{
			{
				Name: msg.From.String(),
				ID:   msg.From.ID,
			},
		},
		SendChan:              sendChan,
		BotStartedWarningSent: make(map[int]bool),

		Words: words,
		Done:  doneChan,

		Votes: make(map[string]int),
	}
}

func (r *Room) joinKeyboard() tgbotapi.InlineKeyboardMarkup {
	jmsg := message.Join

	return tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{
		{
			Text:         "Join",
			CallbackData: &jmsg,
		},
	})
}

func (r *Room) voteKeyboard() tgbotapi.InlineKeyboardMarkup {
	buttons := make([]tgbotapi.InlineKeyboardButton, 0)

	for _, member := range r.Members {
		vmsg := fmt.Sprintf("%s %s", message.Vote, member.Name)

		buttons = append(buttons, tgbotapi.InlineKeyboardButton{
			Text:         member.Name,
			CallbackData: &vmsg,
		})
	}

	return tgbotapi.NewInlineKeyboardMarkup(buttons)
}

// Created must be called on the room creation
func (r *Room) Created() {
	if !(r.State == Join || r.State == CreatorBlocked) {
		return
	}

	// check if user has started the bot
	if err := r.welcome(r.creator().ID, r.Title); err != nil {
		// only send the warning on the first attempt
		if _, ok := r.BotStartedWarningSent[r.creator().ID]; !ok {
			r.botStartedWarning(r.creator().ID, r.creator().Name)
		}

		r.State = CreatorBlocked

		return
	}

	msg := tgbotapi.NewMessage(
		r.ChatID,
		fmt.Sprintf("A game created by %s.", r.creator().Name),
	)
	msg.ReplyMarkup = r.joinKeyboard()

	r.SendChan <- message.Request{
		Chattable: msg,
	}

	go r.countToStart()

	r.State = Join
}

func (r *Room) botStartedWarning(id int, name string) {
	errorMessage := tgbotapi.NewMessage(
		r.ChatID,
		fmt.Sprintf("%s did you start the bot? 🤔", name),
	)

	r.SendChan <- message.Request{
		Chattable: errorMessage,
	}

	r.BotStartedWarningSent[id] = true
}

// countToStart counts 30 seconds then start the game
func (r *Room) countToStart() {
	count := 0

	reportChan := make(chan message.Response)
	r.SendChan <- message.Request{
		Chattable: tgbotapi.NewMessage(
			r.ChatID,
			fmt.Sprintf("%d sec left to join", JoinDuration),
		),
		Report: reportChan,
	}

	id := (<-reportChan).Message.MessageID

	tick := time.NewTicker(1 * time.Second)

	for range tick.C {
		count++
		if count%10 == 0 {
			r.SendChan <- message.Request{
				Chattable: tgbotapi.EditMessageTextConfig{
					BaseEdit: tgbotapi.BaseEdit{
						ChatID:    r.ChatID,
						MessageID: id,
					},
					Text: fmt.Sprintf("%d sec left to join", JoinDuration-count),
				},
			}
		}

		if count == JoinDuration {
			tick.Stop()
			break
		}
	}

	r.SendChan <- message.Request{
		Chattable: tgbotapi.NewMessage(
			r.ChatID,
			"Let's play",
		),
	}

	r.Inform()
}

// Inform is automatically called after the join duration and send the chosen word to members
func (r *Room) Inform() {
	spy := rand.Intn(len(r.Members))
	index := rand.Intn(len(r.Words))
	word := r.Words[index]

	r.Spy = r.Members[spy].Name

	for i := range r.Members {
		word := word

		if i == spy {
			word = "Spy"
		}

		msg := tgbotapi.NewMessage(
			int64(r.Members[i].ID),
			word,
		)
		r.SendChan <- message.Request{
			Chattable: msg,
		}
	}

	r.CountToVote()
}

// CountToVote count 3 minutes then start the voting phase
func (r *Room) CountToVote() {
	tick := time.NewTicker(1 * time.Minute)
	count := 0

	r.SendChan <- message.Request{
		Chattable: tgbotapi.NewMessage(
			r.ChatID,
			fmt.Sprintf("%d min left to discuss", DiscussionDuration),
		),
	}

	for range tick.C {
		count++
		r.SendChan <- message.Request{
			Chattable: tgbotapi.NewMessage(
				r.ChatID,
				fmt.Sprintf("%d min left to discuss", DiscussionDuration-count),
			),
		}

		if count == DiscussionDuration {
			tick.Stop()
			break
		}
	}

	voteMsg := tgbotapi.NewMessage(
		r.ChatID,
		"Let's vote",
	)
	voteMsg.ReplyMarkup = r.voteKeyboard()

	r.SendChan <- message.Request{
		Chattable: voteMsg,
	}

	r.CountToEnd()
}

// CountToEnd count 3 minutes then end the game
func (r *Room) CountToEnd() {
	tick := time.NewTicker(1 * time.Minute)
	count := 0

	r.SendChan <- message.Request{
		Chattable: tgbotapi.NewMessage(
			r.ChatID,
			fmt.Sprintf("%d min left to vote", VoteDuration),
		),
	}

	for range tick.C {
		count++
		r.SendChan <- message.Request{
			Chattable: tgbotapi.NewMessage(
				r.ChatID,
				fmt.Sprintf("%d min left to vote", VoteDuration-count),
			),
		}

		if count == VoteDuration {
			tick.Stop()
			break
		}
	}

	maxName := ""
	maxVote := -1

	for name, vote := range r.Votes {
		if vote > maxVote {
			maxVote = vote
			maxName = name
		}
	}

	endMsg := tgbotapi.NewMessage(
		r.ChatID,
		fmt.Sprintf("%s is the Spy and you select %s as Spy with %d votes", r.Spy, maxName, maxVote),
	)

	r.SendChan <- message.Request{
		Chattable: endMsg,
	}

	close(r.Done)
}

// Voted must be called when a member vote for an spy
// nolint: interfacer
func (r *Room) Voted(from *tgbotapi.User, base *tgbotapi.Message, target string) {
	if from == nil || base == nil {
		return
	}

	r.Votes[target]++

	voteMsg := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    r.ChatID,
			MessageID: base.MessageID,
		},
		Text: fmt.Sprintf("%s\n\r- %s vote for %s.", base.Text, from.String(), target),
	}

	keyboard := r.voteKeyboard()
	voteMsg.ReplyMarkup = &keyboard

	r.SendChan <- message.Request{
		Chattable: voteMsg,
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

	// check to not join one person twice
	for _, member := range r.Members {
		if member.Name == from.String() {
			return
		}
	}

	// check if user has started the bot
	if err := r.welcome(from.ID, base.Chat.Title); err != nil {
		if _, ok := r.BotStartedWarningSent[from.ID]; !ok {
			r.botStartedWarning(from.ID, from.String())
		}

		return
	}

	r.Members = append(r.Members, &Member{
		Name: from.String(),
		ID:   from.ID,
	})

	joinedMsg := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    r.ChatID,
			MessageID: base.MessageID,
		},
		Text: fmt.Sprintf("%s\n\r- %s has joined.", base.Text, from.String()),
	}

	keyboard := r.joinKeyboard()
	joinedMsg.ReplyMarkup = &keyboard

	r.SendChan <- message.Request{
		Chattable: joinedMsg,
	}
}

// welcome sends a welcome message to check If use has started the bot
func (r *Room) welcome(id int, title string) error {
	welcomeMsg := tgbotapi.NewMessage(
		int64(id),
		fmt.Sprintf("Successfully joined a game in %s. 😀", title),
	)

	joinReportChan := make(chan message.Response)
	r.SendChan <- message.Request{
		Chattable: welcomeMsg,
		Report:    joinReportChan,
	}

	if report := <-joinReportChan; report.Error != nil {
		return errors.New("cannot send welcome message")
	}

	return nil
}

func (r *Room) creator() *Member {
	return r.Members[0]
}
