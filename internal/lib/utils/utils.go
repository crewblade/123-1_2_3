package utils

import (
	"strconv"
	"strings"
)

func StringToIntArray(s string) ([]int, error) {
	if len(s) < 2 { // Проверка на пустой массив "{}"
		return []int{}, nil
	}
	s = s[1 : len(s)-1] // Удаление фигурных скобок
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
