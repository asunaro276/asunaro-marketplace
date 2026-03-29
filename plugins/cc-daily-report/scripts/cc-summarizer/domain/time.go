package domain

import (
	"fmt"
	"time"
)

func ParseTime(s string) (time.Time, error) {
	for _, layout := range []string{time.RFC3339Nano, "2006-01-02T15:04:05.000Z", time.RFC3339} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unparseable time: %s", s)
}

func SessionDuration(startTime, endTime string) float64 {
	start, err := ParseTime(startTime)
	if err != nil {
		return 0
	}
	end, err := ParseTime(endTime)
	if err != nil {
		return 0
	}
	if d := end.Sub(start).Minutes(); d > 0 {
		return d
	}
	return 0
}
