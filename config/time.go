package config

import (
	"errors"
	"strings"
	"time"
)

type Time time.Time

func (t Time) MarshalYAML() (any, error) {
	return time.Time(t).Format(time.DateTime), nil
}

func (t *Time) UnmarshalYAML(unmarshaller func(any) error) error {
	var v string
	if err := unmarshaller(&v); err != nil {
		return err
	}

	parsed, err := ParseTime(v)
	if err != nil {
		return err
	}

	*t = Time(parsed)
	return nil
}

var ErrParseTime = errors.New("could not parse time")

func ParseTime(v string) (time.Time, error) {
	v = strings.TrimSpace(v)
	if strings.ToLower(v) == "now" {
		return time.Now(), nil
	}
	if duration, err := time.ParseDuration(v); err == nil {
		return time.Now().Add(duration), nil
	}

	layouts := []string{
		"2006-01-02 15:04:05 MST",
		time.Layout,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
		time.DateTime,
		time.DateOnly,
		time.TimeOnly,
	}

	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, v); err == nil {
			return parsed.UTC(), nil
		}
	}

	return time.Time{}, ErrParseTime
}
