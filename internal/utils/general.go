package utils

import "regexp"

func Sanitize(str string) string {
	sanitize, _ := regexp.Compile(`['"]`)
	return sanitize.ReplaceAllString(str, "")
}
