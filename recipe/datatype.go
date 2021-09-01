package recipe

import "fmt"

type DataType int

const (
	Column DataType = iota
	Variable
	Literal
	Placeholder
	Header
)

func (d *DataType) String() string {
	switch *d {
	case Column:
		return "Column"
	case Variable:
		return "Variable"
	case Literal:
		return "Literal"
	case Placeholder:
		return "Placeholder"
	case Header:
		return "Header"
	default:
		return fmt.Sprintf("Unknown datatype [%d]", d)
	}
}
