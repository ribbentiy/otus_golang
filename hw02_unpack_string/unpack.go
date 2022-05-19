package hw02unpackstring

import (
	"errors"
	"strconv"
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
		case string(char) == `\`:
			if processingChar != 0 {
				resultString.WriteRune(processingChar)
				processingChar = 0
			}
			if escapeFlag {
				processingChar = char
				escapeFlag = false
			} else {
				escapeFlag = true
			}
		case unicode.IsDigit(char) && escapeFlag:
			processingChar = char
			escapeFlag = false
		case unicode.IsDigit(char) && processingChar != 0:
			letter := string(processingChar)
			count, err := strconv.Atoi(string(char))
			if err != nil {
				return "", err
			}
			str := strings.Repeat(letter, count)
			resultString.WriteString(str)
			processingChar = 0
		case unicode.IsDigit(char) && (!escapeFlag || processingChar == 0):
			return "", ErrInvalidString
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
