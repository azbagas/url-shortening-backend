package helper

import "time"

func FormatToUTCString(t time.Time) string {
	return t.Format(time.RFC3339)
}
