package room

type State int

const (
	Join State = iota
	Discuss
	Vote
	End
)
