package util

import (
	"fmt"
	"runtime/debug"
	"time"
)

// TimeHelper converters a string to a time.Time struct.
func TimeHelper(v string) time.Time {
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, v)

	if err != nil {
		debug.PrintStack()
		fmt.Printf("Warning, Time field error: %v\n", err)
	}

	return t
}
