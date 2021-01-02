// Package cronticker creates a ticker (similar to time.Ticker) from a
// cron schedule.
//
// The Cron schedule can be in Unix or Quartz format. Directives like
// '@weekly' or '@daily' can also be parsed as defined in the
// package github.com/robfig/cron/v3.
//
// You may add the TimeZone/location to the beginning of the cron schedule
// to change the time zone. Default is UTC.
//
// See the NewTicker section for examples.
//
// Cronticker calculates the duration until the next scheduled 'tick'
// based on the cron schedule, and starts a new timer based on the
// duration calculated.
//
// You can access the channel with `CronTicker.C`.
package cronticker

import (
	"fmt"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

// CronTicker is the struct returned to the user as a proxy
// to the ticker. The user can check the ticker channel for the next
// 'tick' via CronTicker.C (similar to the user of time.Timer).
type CronTicker struct {
	C chan time.Time
	k chan bool
}

// Stop sends the appropriate message on the control channel to
// kill the CronTicker goroutines. It's good practice to use `defer CronTicker.Stop()`.
func (c *CronTicker) Stop() {
	c.k <- true
}

// Reset kills the ticker and starts it up again with
// the new schedule. The channel remains the same.
func (c *CronTicker) Reset(schedule string) error {
	var err error
	c.Stop()
	err = newTicker(schedule, c.C, c.k)
	if err != nil {
		return err
	}
	return nil
}

// NewTicker returns a CronTicker struct.
// You can check the ticker channel for the next tick by
// `CronTicker.C`
func NewTicker(schedule string) (CronTicker, error) {
	var cronTicker CronTicker
	var err error

	cronTicker.C = make(chan time.Time, 1)
	cronTicker.k = make(chan bool, 1)

	err = newTicker(schedule, cronTicker.C, cronTicker.k)
	if err != nil {
		return cronTicker, err
	}
	return cronTicker, nil
}

// newTicker prepares the channels, parses the schedule, and kicks off
// the goroutine that handles scheduling of each 'tick'.
func newTicker(schedule string, c chan time.Time, k <-chan bool) error {
	var err error

	scheduleWithTZ, loc, err := guaranteeTimeZone(schedule)
	if err != nil {
		return err
	}
	parser := getScheduleParser()

	cronSchedule, err := parser.Parse(scheduleWithTZ)
	if err != nil {
		return err
	}

	go cronRunner(cronSchedule, loc, c, k)

	return nil

}

// getScheduleParser returns a new parser that allows the use of the 'seconds' field
// like in the Quarts cron format, as well as descriptors such as '@weekly'.
func getScheduleParser() cron.Parser {
	parser := cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	return parser
}

// guaranteeTimeZone sets the `TZ=` value to `UTC` if there is none
// already in the cron schedule string.
func guaranteeTimeZone(schedule string) (string, *time.Location, error) {
	var loc *time.Location

	// If time zone is not included, set default to UTC
	if !strings.HasPrefix(schedule, "TZ=") {
		schedule = fmt.Sprintf("TZ=%s %s", "UTC", schedule)
	}

	tz := extractTZ(schedule)

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return schedule, loc, err
	}

	return schedule, loc, nil
}

func extractTZ(schedule string) string {
	end := strings.Index(schedule, " ")
	eq := strings.Index(schedule, "=")
	return schedule[eq+1 : end]
}

// cronRunner handles calculating the next 'tick'. It communicates to
// the CronTicker via a channel and will stop/return whenever it receives
// a bool on the `k` channel.
func cronRunner(schedule cron.Schedule, loc *time.Location, c chan time.Time, k <-chan bool) {
	nextTick := schedule.Next(time.Now().In(loc))
	timer := time.NewTimer(time.Until(nextTick))
	for {
		select {
		case <-k:
			timer.Stop()
			return
		case tickTime := <-timer.C:
			c <- tickTime
			nextTick = schedule.Next(tickTime.In(loc))
			timer.Reset(time.Until(nextTick))
		}
	}
}
