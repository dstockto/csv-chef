package recipe

type DataType int

//go:generate stringer -type=DataType

const (
	Column DataType = iota
	Variable
	Literal
	Placeholder
	Header
)
