package ir

// Access a cell by the name of the column.
type AccessorName struct {
	Name string
}

func InitAccessorName(t_name string) (AccessorName, error) {
	return AccessorName{t_name}, nil
}

// Access a cell by the position of a column.
type AccessorPositional struct {
	Index int64
}

func InitAccessorPositional(t_index int64) (AccessorPositional, error) {
	return AccessorPositional{t_index}, nil
}
