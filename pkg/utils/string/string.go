package string

import "strings"

func IsStringInArray(str string, arr []string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

func ExtractIP(addr string) string {
	parts := strings.Split(addr, ":")
	return parts[0]
}
