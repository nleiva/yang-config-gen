package common

import (
	"strconv"
)

func SplitOnFirstNumber(s string) (string, string) {
	for i, r := range s {
		if _, err := strconv.Atoi(string(r)); err == nil {
			return s[:i], s[i:]
		}
	}
	return s, "" // No number found, return original string and empty string
}
