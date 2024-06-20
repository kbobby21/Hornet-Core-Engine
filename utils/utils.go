package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func CommaSeparateArr(s string) []string {
	return strings.Split(s, ",")
}

func IsEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	fmt.Printf("emailRegex.MatchString(e) & email is :%v: %v\n", emailRegex.MatchString(e), e)
	return emailRegex.MatchString(e)
}

func BigLength(a string, b string) ([]string, []string) {

	ContactEmail := CommaSeparateArr(a)
	ContactNumber := CommaSeparateArr(b)
	l := 0
	if len(ContactEmail) > len(ContactNumber) {
		l = len(ContactEmail)
	} else {
		l = len(ContactNumber)
	}
	a1 := make([]string, l)
	b1 := make([]string, l)
	copy(a1, ContactEmail)
	copy(b1, ContactNumber)

	return a1, b1
}
