package phone

import (
	"errors"
	"regexp"
)

var (
	nonDigitsPattern   = regexp.MustCompile(`\D`)
	phoneNumberPattern = regexp.MustCompile(`^\d{10}$`)
)

// Normalize returns a normalized phone number parsed from a string.
// If the number if invalid, an error is returned along with the emtpy string.
func Normalize(s string) (string, error) {
	cleaned := nonDigitsPattern.ReplaceAllString(s, "")

	if !phoneNumberPattern.MatchString(cleaned) {
		return "", errors.New("invalid phone number: " + s)
	}

	return cleaned, nil
}
