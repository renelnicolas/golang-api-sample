package helpers

// TrimWhitespaceFn : is a strings.TrimFunc function that removes all white space.
func TrimWhitespaceFn(r rune) bool {
	return r <= 32
}
