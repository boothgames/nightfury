package nightfury

import (
	"regexp"
	"strings"
)

var slugRegex = regexp.MustCompile("[ _]")

// Slug generate slug for the string
func Slug(value string) string {
	return strings.ToLower(slugRegex.ReplaceAllString(value, "-"))
}
