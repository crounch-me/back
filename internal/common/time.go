package common

import "time"

func IsInRange(current, start, end time.Time) bool {
	is_date_on_or_after_start := current.After(start) || current.Equal(start)
	is_date_on_or_before_end := current.Before(end) || current.Equal(end)

	return is_date_on_or_after_start && is_date_on_or_before_end
}
