package utils

func Truncate(s string, max_len int) string {
	l := len(s)
	if l < max_len {
		return s
	}
	return "*" + s[(l-max_len+2):]
}
