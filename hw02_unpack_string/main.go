package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	// Place your code here.
	// Исходная строка
	str := "d\n5abc"

	// Преобразуем строку в срез рун
	runes := []rune(str)

	readyString, err := writeString(runes, str)

	if err != nil {
		fmt.Println(err, str)
		return
	}

	fmt.Printf("Преобразованное : %s", readyString)
}

func writeString(runes []rune, text string) (string, error) {
	var builder strings.Builder

	// Распечатаем каждую руну отдельно
	for i, r := range runes {
		num, err := strconv.Atoi(string(r))
		if err != nil {

			builder.WriteString(string(r))
		} else {
			if i == 0 {
				return "", errors.New("Некорректная строка:")
			} else if checkDecimal(string(r), text) {
				return "", errors.New("Некорректная строка:")
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

func checkDecimal(num string, text string) bool {
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
