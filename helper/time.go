package helper

import "time"

func FormatToUTCString(t time.Time) string {
	return t.Format(time.RFC3339)
}

func GetCurrentTime(timezone string) time.Time {
	loc, err := time.LoadLocation(timezone)
	PanicIfError(err)

	return time.Now().In(loc)
}