package internal

var FirstNotEmpty = func(s ...string) string {
	for _, v := range s {
		if v != "" {
			return v
		}
	}

	return ""
}
