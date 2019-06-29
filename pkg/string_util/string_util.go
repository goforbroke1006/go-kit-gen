package string_util

import "strings"

func FirstLetterToLowerCase(str string) string {
	if 0 == len(str) {
		return ""
	}
	if 1 == len(str) {
		return strings.ToLower(str)
	}

	return strings.ToLower(str[0:1]) + str[1:]
}

func FirstLetterToUpperCase(str string) string {
	if 0 == len(str) {
		return ""
	}
	if 1 == len(str) {
		return strings.ToUpper(str)
	}

	return strings.ToUpper(str[0:1]) + str[1:]
}
