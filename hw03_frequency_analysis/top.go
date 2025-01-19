package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type CountWord struct {
	Word  string
	Count int
}

func Top10(text string) []string {
	words := strings.Fields(text)
	countMap := make(map[string]int)
	countWords := make([]CountWord, 0)
	readySlice := make([]string, 0)

	for _, word := range words {
		countMap[word] += 1
	}

	for i, countWord := range countMap {
		el := CountWord{Word: i, Count: countWord}
		countWords = append(countWords, el)
	}

	sort.SliceStable(countWords, func(i, j int) bool {
		if countWords[i].Count == countWords[j].Count {
			return countWords[i].Word < countWords[j].Word
		} else {
			return countWords[i].Count > countWords[j].Count
		}
	})

	for i, countWord := range countWords {
		readySlice = append(readySlice, countWord.Word)
		if i == 9 {
			break
		}
	}

	return readySlice
}
