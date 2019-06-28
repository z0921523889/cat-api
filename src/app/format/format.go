package format

import "time"

var timeFormatter = "2006-01-02T15:04:05"

func ParseTime(timeString string) (time.Time, error) {
	return time.Parse(timeFormatter, timeString)
}
