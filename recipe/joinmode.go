package recipe

type JoinMode int

//go:generate stringer -type=JoinMode
const (
	Replace JoinMode = iota
	Join
)
