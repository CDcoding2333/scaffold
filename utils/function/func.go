package function

import (
	"regexp"
)

// CheckTel ...
func CheckTel(mobileNum string) bool {
	reg := regexp.MustCompile("^\\d{5,13}$")
	return reg.MatchString(mobileNum)
}
