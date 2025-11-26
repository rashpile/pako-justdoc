package model

import "regexp"

var validName = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,128}$`)

// IsValidName checks if a channel or document name is valid.
// Valid names contain only alphanumeric characters, hyphens, and underscores,
// and are between 1 and 128 characters long.
func IsValidName(name string) bool {
	return validName.MatchString(name)
}