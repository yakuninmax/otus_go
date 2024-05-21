package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/example/stringutil"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(packedString string) (string, error) {
	// Check if string is empty
	if packedString == "" {
		return "", nil
	}

	// Get runes from packed string
	runes := []rune(packedString)

	// Check if first rune is digit
	if unicode.IsDigit(runes[0]) {
		return "", ErrInvalidString
	}

	// Unpacking
	// Start from end of string
	var reversedResult strings.Builder

	for i := len(runes) - 1; i >= 0; i-- {
		currentRune := runes[i]

		// Check if rune is digit
		if unicode.IsDigit(currentRune) {
			// Check if next rune is digit
			nextRune := runes[i-1]

			if unicode.IsDigit(nextRune) {
				return "", ErrInvalidString
			}

			// Get number of chars from digit
			number, err := strconv.Atoi(string(currentRune))
			if err != nil {
				return "", err
			}

			// Append chars
			reversedResult.WriteString(strings.Repeat(string(nextRune), number))

			// Decrease index i and continue
			i--
			continue
		}

		// If rune is letter, append it
		reversedResult.WriteString(string(currentRune))
	}

	// Revers string and return
	result := stringutil.Reverse(reversedResult.String())
	return result, nil
}
