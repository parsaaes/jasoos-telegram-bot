package room

// State represents state of the playing room
type State int

const (
	// Join users can join
	Join State = iota
	// Discuss users discuss to find an spy
	Discuss
	// Vote users vote for an spy
	Vote
	// End there is no game in this room
	End
	// CreatorBlocked creator didn't started the bot
	CreatorBlocked
)
