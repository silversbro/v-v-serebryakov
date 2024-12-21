package main

import (
	"errors"
	"unpack"
)

func main() (string, error) {
	// Place your code here.
	// Исходная строка
	str := "d\n5abc"

	readyStr, err := unpack.Unpack(str)

	if err != nil {
		err = errors.New("Error в строке:" + str)

		return "", err
	}

	return readyStr, nil
}
