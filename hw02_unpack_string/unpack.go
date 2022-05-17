package hw02unpackstring

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(incomingString string) (string, error) {
	var processingChar rune
	var escapeFlag bool
	var resultString strings.Builder

	for _, char := range incomingString {
		switch {
		case char == 92:
			if processingChar != 0 {
				resultString.WriteRune(processingChar)
				processingChar = 0
			}
			if escapeFlag == true {
				processingChar = char
				escapeFlag = false
			} else {
				escapeFlag = true
			}
		case unicode.IsDigit(char):
			if escapeFlag {
				processingChar = char
				escapeFlag = false
			} else if processingChar != 0 {
				letter := fmt.Sprintf("%c", processingChar)
				str := strings.Repeat(letter, int(char-48))
				resultString.WriteString(str)
				processingChar = 0
			} else {
				return "", ErrInvalidString
			}
		default:
			escapeFlag = false
			if processingChar != 0 {
				resultString.WriteRune(processingChar)
			}
			processingChar = char

		}

	}
	if processingChar != 0 {
		resultString.WriteRune(processingChar)
	}
	return resultString.String(), nil
}
