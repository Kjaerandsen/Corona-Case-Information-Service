package function

import "time"

var uptimeStart time.Time // Uptime calculation start value

// Returns the uptime of the service based on
// code found at https://stackoverflow.com/questions/37992660/golang-retrieve-application-uptime
func Uptime() time.Duration {
	return time.Since(uptimeStart)
}

// Starts the timer for the uptime in exchange/v1/diag/
func UptimeInit() {
	uptimeStart = time.Now()
}

