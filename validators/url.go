package validators

import "regexp"

var (
	rePattern = `^(?:http(s)?:\/\/)?[\w.-]+(?:\.[\w\.-]+)+`
	reURL     = regexp.MustCompile(rePattern)
)

// ValidateURL :
func ValidateURL(URL string) bool {
	ok, _ := regexp.MatchString(rePattern, URL)

	return ok
}
