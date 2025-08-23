package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type User struct {
	Email string `json:"email"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain)
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	targetSuffix := "." + domain

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var user User
		if err := json.Unmarshal(line, &user); err != nil {
			return nil, fmt.Errorf("json unmarshal error: %w", err)
		}

		if strings.HasSuffix(strings.ToLower(user.Email), targetSuffix) {
			parts := strings.SplitN(user.Email, "@", 2)
			if len(parts) == 2 {
				domain := strings.ToLower(parts[1])
				result[domain]++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return result, nil
}
