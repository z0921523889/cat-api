package format

import "time"

const (
	DateFormatter = "2006-01-02"
	TimeFormatter = "2006-01-02T15:04:05"
)


func ParseTime(timeString string,timeFormatter string) (time.Time, error) {
	return time.Parse(timeFormatter, timeString)
}
