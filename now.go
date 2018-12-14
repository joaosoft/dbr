package dbr

import "time"

const timeFormat = "2006-01-02 15:04:05.000000"

func Now() string {
	return time.Now().UTC().Format(timeFormat)
}
