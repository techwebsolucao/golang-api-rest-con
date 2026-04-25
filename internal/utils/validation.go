package utils

import (
	"fmt"
	"unicode"
)

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("a senha deve ter no mínimo 8 caracteres")
	}

	hasNumber := false
	for _, char := range password {
		if unicode.IsDigit(char) {
			hasNumber = true
			break
		}
	}

	if !hasNumber {
		return fmt.Errorf("a senha deve conter pelo menos um número")
	}

	return nil
}