package validation

import "regexp"

const EMAIL_REGEX = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

func IsEmailFormatOK(_string string) bool {
	re := regexp.MustCompile(EMAIL_REGEX)
	return re.MatchString(_string)
}
