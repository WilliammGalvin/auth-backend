package validators

import (
	"fmt"
	"regexp"
	"unicode/utf8"
)

// ValidateEmail ensures email is a valid format
func ValidateEmail(email string) error {
	if utf8.RuneCountInString(email) > 100 {
		return fmt.Errorf("email too long (max 100 characters)")
	}

	re, err := regexp.Compile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z0-9-]+$`)

	if err != nil {
		return fmt.Errorf("error compiling regex: %v", err)
	}

	if !re.MatchString(email) {
		return fmt.Errorf("email not valid")
	}

	return nil
}
