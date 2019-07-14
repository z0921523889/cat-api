package format

import "time"

const (
	DateFormatter = "2006-01-02"
	TimeFormatter = "2006-01-02T15:04:05"
)

func ParseTime(timeString string, timeFormatter string) (time.Time, error) {
	var t time.Time
	var err error
	var local *time.Location
	local, err = time.LoadLocation("Local")
	if err == nil {
		t, err = time.ParseInLocation(timeFormatter, timeString, local)
	}
	return t, err
}

func ReplaceDateOfTime(origin time.Time, replace time.Time) time.Time {
	return time.Date(
		replace.Year(),
		replace.Month(),
		replace.Day(),
		origin.Hour(),
		origin.Minute(),
		origin.Second(),
		origin.Nanosecond(),
		origin.Location(),
	)
}
