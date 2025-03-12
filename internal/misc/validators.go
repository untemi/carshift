package misc

import (
	"errors"
	"os"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var usernamePattern = regexp.MustCompile("^[a-zA-Z0-9_]+$")

func ValidateUsername(t string, requireLength bool) bool {
	Length := utf8.RuneCountInString(t)
	validPattern := usernamePattern.MatchString(t)

	max := Length < 26 || !requireLength
	min := Length > 3 || !requireLength

	return max && min && validPattern
}

func ValidateName(t string, IsNullable bool) bool {
	if len(strings.TrimSpace(t)) < 2 {
		if IsNullable {
			return true
		}
		return false
	}
	for _, char := range t {
		if !unicode.IsLetter(char) && !unicode.IsSpace(char) && char != '-' && char != '\'' {
			return false
		}
	}
	return true
}

func ValidatePassword(t string) bool {
	if len(t) < 8 {
		return false
	}

	var hasDigit, hasLower, hasUpper bool
	for _, char := range t {
		if unicode.IsDigit(char) {
			hasDigit = true
		} else if unicode.IsLower(char) {
			hasLower = true
		} else if unicode.IsUpper(char) {
			hasUpper = true
		}

		if hasDigit && hasLower && hasUpper {
			return true
		}
	}

	return hasDigit && hasLower && hasUpper
}

func IsFileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !errors.Is(err, os.ErrNotExist)
}
