package utils

import (
	"strconv"
	"strings"
)

func IntPointertoaOrDefault(value *int, defaultValue string) string {
	if value != nil {
		return strconv.Itoa(*value)
	}
	return defaultValue
}

func StringToIntArray(s string) ([]int, error) {
	if len(s) < 2 {
		return []int{}, nil
	}
	s = s[1 : len(s)-1]
	parts := strings.Split(s, ",")
	var res []int
	for _, part := range parts {
		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		res = append(res, val)
	}
	return res, nil
}
