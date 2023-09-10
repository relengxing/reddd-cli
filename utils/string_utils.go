package utils

import "strings"

func ContainsString(arr []string, target string) bool {
	for _, str := range arr {
		if strings.Contains(str, target) {
			return true
		}
	}
	return false
}
