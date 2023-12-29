package util

import (
	"crypto/rand"
	"fmt"
	"halodeksik-be/app/appconstant"
	"strings"
	"time"
	"unicode"
)

func IsEmptyString(str string) bool {
	return str == ""
}

func GetCurrentDateAndTime() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())
}

func ParseDateTime(timeStr string, timeFormat ...string) (time.Time, error) {
	if len(timeFormat) > 0 {
		return time.Parse(timeFormat[0], timeStr)
	}
	return time.Parse(appconstant.TimeFormatQueryParam, timeStr)
}

func RandomToken(marker string) (string, error) {
	bTime, _ := time.Now().MarshalText()
	bMarker := []byte(marker)

	b := append(bMarker, bTime...)
	_, err := rand.Read(b)

	return fmt.Sprintf("%x", b), err
}

func PascalToSnake(input string) string {
	var builder strings.Builder

	for i, char := range input {
		if unicode.IsUpper(char) {
			if i > 0 {
				builder.WriteRune('_')
			}
			builder.WriteRune(unicode.ToLower(char))
		} else {
			builder.WriteRune(char)
		}
	}

	return builder.String()
}
