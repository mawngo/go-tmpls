package internal

import (
	"strconv"
	"time"
)

// now returns the current time.
func now() time.Time {
	return time.Now()
}

// datefmt format time.
// The default format is time.DateOnly.
//
// Example usage: datefmt (now).
func datefmt(v time.Time, format ...string) string {
	if len(format) == 0 {
		return v.Format(time.DateOnly)
	}
	return v.Format(format[0])
}

// datetimefmt format time.
// The default format is time.DateTime.
//
// Example usage: datetimefmt (now).
func datetimefmt(v time.Time, format ...string) string {
	if len(format) == 0 {
		return v.Format(time.DateTime)
	}
	return v.Format(format[0])
}

// date formats a date.
//
// Example usage: now | date "2006-01-02".
func date(fmt string, date any) string {
	return dateInZone(fmt, date, "Local")
}

// duration Formats a given number of seconds as a time.Duration.
//
// Example usage: duration 95 => 1m35s.
func duration(sec any) string {
	var n int64
	switch value := sec.(type) {
	default:
		n = 0
	case string:
		n, _ = strconv.ParseInt(value, 10, 64)
	case int:
		n = int64(value)
	case int8:
		n = int64(value)
	case int16:
		n = int64(value)
	case int32:
		n = int64(value)
	case int64:
		n = value
	case uint:
		n = int64(value)
	case uint8:
		n = int64(value)
	case uint16:
		n = int64(value)
	case uint32:
		n = int64(value)
	case uint64:
		n = int64(value)
	case float32:
		n = int64(value)
	case float64:
		n = int64(value)
	}
	return (time.Duration(n) * time.Second).String()
}

// durationRound Rounds a given duration to the most significant unit.
// Strings and time.Duration gets parsed as a duration, while a time.Time is calculated as the duration since.
//
// Example usage: durationRound "2400h10m5s" => "3mo".
func durationRound(duration any) string {
	var d time.Duration
	switch duration := duration.(type) {
	default:
		d = 0
	case string:
		d, _ = time.ParseDuration(duration)
	case int:
		d = time.Duration(int64(duration))
	case int8:
		d = time.Duration(duration)
	case int16:
		d = time.Duration(duration)
	case int32:
		d = time.Duration(duration)
	case int64:
		d = time.Duration(duration)
	case uint:
		d = time.Duration(int64(duration))
	case uint8:
		d = time.Duration(int64(duration))
	case uint16:
		d = time.Duration(int64(duration))
	case uint32:
		d = time.Duration(int64(duration))
	case uint64:
		d = time.Duration(int64(duration))
	case float32:
		d = time.Duration(int64(duration))
	case float64:
		d = time.Duration(int64(duration))
	case time.Time:
		d = time.Since(duration)
	case time.Duration:
		d = duration
	}

	u := uint64(d)
	neg := d < 0
	if neg {
		u = -u
	}

	var (
		year   = uint64(time.Hour) * 24 * 365
		month  = uint64(time.Hour) * 24 * 30
		day    = uint64(time.Hour) * 24
		hour   = uint64(time.Hour)
		minute = uint64(time.Minute)
		second = uint64(time.Second)
	)
	switch {
	case u > year:
		return strconv.FormatUint(u/year, 10) + "y"
	case u > month:
		return strconv.FormatUint(u/month, 10) + "mo"
	case u > day:
		return strconv.FormatUint(u/day, 10) + "d"
	case u > hour:
		return strconv.FormatUint(u/hour, 10) + "h"
	case u > minute:
		return strconv.FormatUint(u/minute, 10) + "m"
	case u > second:
		return strconv.FormatUint(u/second, 10) + "s"
	}
	return "0s"
}

// dateInZone format a date in a specific zone.
//
// Example usage: dateInZone "2006-01-02" (now) "UTC".
func dateInZone(fmt string, date any, zone string) string {
	var t time.Time
	switch date := date.(type) {
	default:
		t = time.Now()
	case time.Time:
		t = date
	case *time.Time:
		t = *date
	case int64:
		t = time.Unix(date, 0)
	case int:
		t = time.Unix(int64(date), 0)
	case int32:
		t = time.Unix(int64(date), 0)
	}

	loc, err := time.LoadLocation(zone)
	if err != nil {
		loc, _ = time.LoadLocation("UTC")
	}

	return t.In(loc).Format(fmt)
}

// toDate parse a date string.
// The first argument is the date layout and the second the date string.
// If the string can’t be converted it returns the zero value.
//
// This is useful when you want to convert a string date to another format (using pipe). The example below converts
// Example usage: toDate "2006-01-02" "2017-12-31" | date "02/01/2006".
func toDate(fmt, str string) time.Time {
	t, _ := time.ParseInLocation(fmt, str, time.Local)
	return t
}
