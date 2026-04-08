package helper

import (
	"fmt"
	"regexp"
	"strings"
)

func NormalizePhoneNumber(phone string) (string, error) {
	phone = strings.TrimSpace(phone)

	if phone == "" {
		return "", fmt.Errorf("phone number is required")
	}

	// validasi karakter
	matched, _ := regexp.MatchString(`^\+?[0-9]+$`, phone)
	if !matched {
		return "", fmt.Errorf("invalid phone number format")
	}

	// normalize ke format 62
	if strings.HasPrefix(phone, "08") {
		phone = "62" + phone[1:]
	} else if strings.HasPrefix(phone, "+62") {
		phone = phone[1:]
	}

	// validasi prefix Indonesia
	if !strings.HasPrefix(phone, "62") {
		return "", fmt.Errorf("phone number must start with 62")
	}

	// panjang
	if len(phone) < 11 || len(phone) > 15 {
		return "", fmt.Errorf("invalid phone number length")
	}

	return phone, nil
}
