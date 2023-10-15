package lib

import (
	"time"
)

// Duration is a mapping on a time.Duration
type Duration time.Duration

// MarshalCSV exports the duration
func (it Duration) MarshalCSV() (string, error) {
	return time.Duration(it).String(), nil
}

// UnmarshalCSV imports the duration
func (it *Duration) UnmarshalCSV(s string) error {
	duration, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*it = Duration(duration)
	return nil
}

// Bool is a mapping on a boolean
type Bool bool

// MarshalCSV exports the boolean
func (it Bool) MarshalCSV() (string, error) {
	if it {
		return "true", nil
	}
	return "false", nil
}

// UnmarshalCSV imports the boolean
func (it *Bool) UnmarshalCSV(s string) error {
	*it = (s == "true")
	return nil
}

// Date is a mapping on a time.Time
type Date time.Time

// MarshalCSV exports the date
func (it Date) MarshalCSV() (string, error) {
	return time.Time(it).Format(time.DateOnly), nil
}

// UnmarshalCSV imports the date
func (it *Date) UnmarshalCSV(s string) error {
	tm, err := time.Parse(time.DateOnly, s)
	*it = Date(tm)
	return err
}
