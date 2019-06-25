package string_util

import "strings"

func FirstLetterToLowerCase(str string) string {
	if len(str) == 0 {
		return ""
	}
	if len(str) == 1 {
		return strings.ToLower(str)
	}

	return strings.ToLower(str[0:1]) + str[1:]
}
