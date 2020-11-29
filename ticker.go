package cronticker

// MIN    Minute field   0 to 50
// HOUR   Hour field     0 to 23
// DOM    Day of month	 1 - 31
// MON    Month field	 1 - 12
// DOW	  Day of Week	 0 - 6

import (
	"strings"
	"time"
	"strconv"
)

// Constants for cron location when cron fields are
// split into a slice
const (
	min int = iota
	hour
	dom
	mon
	dow
	ANY string = "*"
	STEPS string = "/"
	RANGE string = "-"
	SEPARATOR string = ","
	MIN string = "MIN"
	HOUR = "HOUR"
	DOM = "DOM"
	MON = "MON"
	DOW = "DOW"
)

type Cron

type CronType string

const (
	Period CronType = "PERIOD"
	Calculated CronType = "CALCULATED"
)

type CronTicker struct {
	Config string
	Ticker time.Ticker
	Timer time.Timer
	Type   CronType
	C <-chan time.Time
}

func (t *CronTicker) Stop() {
	t.Ticker.Stop()
	// Close ticker channel?
}

func (t *CronTicker) Reset(cron string) {
	var d time.Duration

	// Calculate next

	switch t.Type {
	case Period:

	}
	if t.Type == Calculated {
		// lots of junk to figure out 'd'
	} else if t.Type == Period {

	}
}

func getDuration(cron string) (time.Duration, CronType, error) {

}

type Cron struct {
	Raw string
	Parsed map[string]int
}

func (c *Cron) Parse() {
	fields := strings.Split(c.Raw, " ")


}

func getMinutes(field string) ([]int, error) {
	var minuteList []int

	if field == "*" {
		minuteList = append(minuteList, -1)
		return minuteList, nil
	}

	minutes := strings.Split(field, ",")
	for _, minute := range minutes {

		// Handle ranges
		if strings.Contains(minute, "-") {
			minuteRange, err := getRange(minute)
			if err != nil {
				return minuteRange, err
			}
			minuteList = append(minuteList, minuteRange...)

			// Handle periods
		} else if strings.Contains(minute, "/") {

			// Do something
			// Handle normal integers
		} else {
			integer, err := strconv.Atoi(minute)
			if err != nil {
				return minuteList, err
			}
			minuteList = append(minuteList, integer)

		}

		// Handle normal integers

	}
	return minuteList, nil
}

func getRange(field string, min int, max int) ([]int, error){
	var minutes []int
	bounds := strings.Split(field, "-")
	startInt, err := strconv.Atoi(bounds[0])
	if err != nil {
		return minutes, err
	}

	endInt, err := strconv.Atoi(bounds[1])
	if err != nil {
		return minutes, err
	}

	difference := endInt - startInt
	for i := 0; i < difference; i++ {
		minute := startInt + i + 1
		minutes = append(minutes, minute)
	}
	return minutes, nil
}

func getPeriod(field string, min int, max int) ([]int, error) {
	var minutes []int
	start := int
	var err error = nil
	bounds := strings.Split(field, "/")
	if bounds[0] != "*" {
		start, err = strconv.Atoi(bounds[0])
		if err != nil {
			return minutes, err
		}
	}

}

// Calculate next
	// Calculate timer based on exact second (time.Second * time.Duration(60-date.Second())