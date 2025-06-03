package ir

// Convert to normal int (watch out with integer slicing).
func toInt(t_str string) (int, error) {
	var result int

	integer, err := strconv.ParseInt(t_str, 10, 64)
	if err != nil {
		return result, err
	}

	result = int(integer)

	return result, nil
}

func binaryExprToInt(t_lhs string, t_rhs string) (int, int, error) {
	var (
		lhs int
		rhs int
	)

	lhs, err := toInt(t_lhs)
	if err != nil {
		return lhs, rhs, err
	}

	rhs, err = toInt(t_rhs)
	if err != nil {
		return lhs, rhs, err
	}

	return lhs, rhs, nil
}
