package ast

// Access a cell by the name of the column.
type AccessorName struct {
	Name string
}

// Access a cell by the position of a column.
type AccessorPositional struct {
	Index int64
}
