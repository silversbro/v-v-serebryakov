package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {

	// Преобразуем строку в срез рун
	runes := []rune(str)

	readyString, err := WriteString(runes, str)

	if err != nil {
		fmt.Println(err, str)
		return str, err
	}

	return readyString, nil
}

func WriteString(runes []rune, text string) (string, error) {
	var builder strings.Builder

	// Распечатаем каждую руну отдельно
	for i, r := range runes {
		num, err := strconv.Atoi(string(r))
		if err != nil {

			builder.WriteString(string(r))
		} else {
			if i == 0 {
				return "", ErrInvalidString
			} else if CheckDecimal(string(r), text) {
				return "", ErrInvalidString
			} else if num == 0 {
				builder.WriteString(string(r))
				index := strings.Index(builder.String(), "0")
				revertVal := builder.String()[:index-1]

				builder.Reset()
				builder.WriteString(revertVal)

				continue
			}

			repeated := strings.Repeat(string(runes[i-1]), num-1)
			builder.WriteString(repeated)
		}
	}

	return builder.String(), nil
}

func CheckDecimal(num string, text string) bool {
	var checkBuilder strings.Builder
	checkBuilder.WriteString(num)
	checkBuilder.WriteString("0")
	index := strings.Index(text, checkBuilder.String())
	println(checkBuilder.String())

	if index == -1 {
		return false
	}

	return true
}
