package function

import (
	"fmt"
	"strconv"
	"time"
)

/*  Time based functions including:
	Uptime
	Checking if dates are in order
	Checking if dates are in scope
 */

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

// Checks if start-date is before end-date
// Returns a boolean
func StartEnd(dates []string) bool{
	var intDates [6]int
	var err error

	// TODO: document, If the input is not a number it defaults to 0 when converting might be a conversion issue
	// TODO: document, Values over expected values in the time command compound into a later date.
	// Converts the string array to an integer array
	for i := 0; i < 6 ; i++ {
		intDates[i], err = strconv.Atoi(dates[i])
		//fmt.Println(i, " ", intDates[i], " ", dates[i])
		if err != nil {
			//fmt.Println(i, "ERROR ", intDates[i], " ", dates[i])
			return false
		}
	}

	// Sets the start and end dates from the input
	startDate := time.Date(intDates[0], time.Month(intDates[1]), intDates[2], 0, 0, 0, 0, time.UTC)
	endDate   := time.Date(intDates[3], time.Month(intDates[4]), intDates[5], 0, 0, 0, 0, time.UTC)

	// Check if the difference between start and end date more than or equal to a day
	if !DateConsecutive(startDate, endDate) {
		return false
	}

	return true
}

// Check if the first date is before the second date
func DateConsecutive(date1 time.Time, date2 time.Time) bool {
	// Check if the difference between start and end date more than or equal to a day
	if date2.Sub(date1).Hours() < 24 {
		return false
	}
	return true
}

// Returns the date -7, -10 and -13 days as a string array for the stringency endpoint
func DateToday() [3]string {
	var intMonth int
	var output [3]string

	// Set the date to the current date
	var date = time.Now()

	// Set the date to seven days prior the the current date
	date = date.AddDate(0,0,-7)

	// Fill in the dates starting from a seven days ago, then 10 then 13.
	for i:=0; i<3; i++ {
		year, month, day := date.Date()

		intMonth = int(month)

		if int(month) < 10 {
			output[i] = fmt.Sprintf("%v-0%v-", year, intMonth)
		} else {
			output[i] = fmt.Sprintf("%v-%v-", year, intMonth)
		}

		if day < 10 {
			output[i] = fmt.Sprintf("%s0%v", output[i], day)
		} else {
			output[i] = fmt.Sprintf("%s%v", output[i], day)
		}

		// Subtract three days, as in fill in the date three days prior in the next loop
		date = date.AddDate(0,0,-3)
	}

	return output
}