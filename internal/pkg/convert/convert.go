package convert

import "strconv"

func StringToInt(s string, def int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		i = def
	}

	return i
}
