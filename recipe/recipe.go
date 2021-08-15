package recipe

type Output struct {
	Type  string
	Value string
}

type Argument struct {
	Type  string
	Value string
}

type Operation struct {
	Name      string
	Arguments []Argument
}

type Recipe struct {
	Output Output
	Pipe   []Operation
}
