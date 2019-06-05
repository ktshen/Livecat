package crawler

import "time"

func getLocalTime() time.Time {
	t := time.Now()
	local, err := time.LoadLocation("Local")
	handleError(err, "Get local time err")
	return t.In(local)
}
