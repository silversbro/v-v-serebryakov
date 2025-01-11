package hw02unpackstring

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrInvalidString = errors.New("invalid string")
	validNumber      = regexp.MustCompile(`\d+\d+`)
	replaceZeroChar  = regexp.MustCompile(`.0`)
)

func Unpack(str string) (string, error) {
	if validNumber.MatchString(str) {
		return "", ErrInvalidString
	}

	if replaceZeroChar.MatchString(str) {
		str = replaceZeroChar.ReplaceAllString(str, "")
	}

	readyString, err := WriteString(str)
	if err != nil {
		return "", err
	}

	return readyString, nil
}

func WriteString(text string) (string, error) {
	var builder strings.Builder
	runes := []rune(text)

	if len(runes) < 2 {
		return text, nil
	}

	// Распечатаем каждую руну отдельно
	for i, r := range runes {
		num, err := strconv.Atoi(string(r))
		if err != nil {
			builder.WriteString(string(r))
		} else {
			if i == 0 {
				return "", ErrInvalidString
			}

			repeated := strings.Repeat(string(runes[i-1]), num-1)
			builder.WriteString(repeated)
		}
	}

	return builder.String(), nil
}
