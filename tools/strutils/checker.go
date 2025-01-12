package strutils

import (
	"regexp"
)

func CheckNumber(s string) bool {
	re := regexp.MustCompile(`(?U)^-?\d+$`)
	return re.MatchString(s)
}

func CheckMultiIds(s string) bool {
	re := regexp.MustCompile("(?U)^[0-9,]+$")
	return re.MatchString(s)
}

func CheckRegularString(s string) bool {
	re := regexp.MustCompile(`(?U)^\w+$`)
	return re.MatchString(s)
}

// a-zA-Z0-9_- supports underscore, hyphen, and space
func CheckRegularString2(s string) bool {
	re := regexp.MustCompile(`(?U)^[a-zA-Z0-9_\- ]+$`)
	return re.MatchString(s)
}

func CheckUUID(s string) bool {
	re := regexp.MustCompile(`(?U)^[A-Za-z0-9\-]+$`)
	return re.MatchString(s)
}

func CheckAddress(s string) bool {
	re := regexp.MustCompile("(?U)^0x[A-Za-z0-9]{40}$")
	return re.MatchString(s)
}

func CheckUrl(s string) bool {
	re := regexp.MustCompile(`(?U)^[.:_/\-?=&%#\w\s]+$`)
	return re.MatchString(s)
}
