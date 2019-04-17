package crawler

import (
	"time"
)

func getLocalTime() time.Time {
	t := time.Now()
	local, err := time.LoadLocation("Local")
	handleError(err, "[getTime.go] Get local time err")
	return t.In(local)
}

func getLocalTimeVersion2() string {
	t := time.Now().UTC().Format("2006-01-02T15:04:05-0700")
	return t
}
