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

// 将下划线命名转为驼峰命名
func ToCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := 0; i < len(parts); i++ {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}
