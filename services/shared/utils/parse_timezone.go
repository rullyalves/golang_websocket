package utils

import (
	"time"
)

func ParseToBRTimezone(value string) (*time.Time, error) {
	return parseDateWithTimezone(value, "America/Sao_Paulo")
}

func WithTimeZone(date time.Time) (*time.Time, error) {
	return toBrTimeZone(date, "America/Sao_Paulo")
}

func toBrTimeZone(date time.Time, timeZone string) (*time.Time, error) {
	loc, err := time.LoadLocation(timeZone)

	if err != nil {
		return nil, err
	}

	result := date.In(loc)

	return &result, nil
}

func parseDateWithTimezone(value string, timezone string) (*time.Time, error) {

	layout := "2006-01-02T15:04:05.999999999"

	date, err := time.Parse(layout, value)

	if err != nil {
		return nil, err
	}

	return toBrTimeZone(date, timezone)
}
