package main

import (
	"errors"
	"github.com/silversbro/hw02_unpack_string"
)

func main() (string, error) {
	// Place your code here.
	// Исходная строка
	str := "d\n5abc"

	readyStr, err := hw02_unpack_string.Unpack(str)

	if err != nil {
		err = errors.New("Error в строке:" + str)

		return "", err
	}

	return readyStr, nil
}
