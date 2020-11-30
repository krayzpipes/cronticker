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

// CronTicker is the struct/object returned to the user as a proxy
// to the ticker. The user can check the ticker channel for the next
// 'tick' via CronTicker.C (kind of what you can do with 'time.Time.C').
type CronTicker struct {
	C chan time.Time
	k chan bool
}

// Stop sends the appropriate message on the control channel to
// kill the CronTicker goroutines.
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

// NewTicker returns a CronTicker object.
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

	scheduleWithTZ, tz := guaranteeTimeZone(schedule)
	parser := getScheduleParser()

	cronSchedule, err := parser.Parse(scheduleWithTZ)
	if err != nil {
		return err
	}

	go cronRunner(cronSchedule, tz, c, k)

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
func guaranteeTimeZone(schedule string) (string, string) {
	tz := "UTC"
	if !strings.HasPrefix(schedule, "TZ=") {
		return fmt.Sprintf("TZ=%v %v", tz, schedule), tz
	}
	end := strings.Index(schedule, " ")
	eq := strings.Index(schedule, "=")  // TODO - finish this...
	return schedule, tz
}

// cronRunner handles calculating the next 'tick'. It communicates to
// the CronTicker via a channel and will stop/return whenever it recieves
// a bool on the `k` channel.
func cronRunner(schedule cron.Schedule, tz string, c chan time.Time, k <-chan bool) {
	// TODO - Should pull location from the schedule and add that to time.now()?
	// TODO - Otherwise, time.Now() might not be before the next scheduled tick.
	nextTick := schedule.Next(time.Now())
	timer := time.NewTimer(time.Until(nextTick))
	for {
		select {
		case <-k:
			timer.Stop()
			return
		case tickTime := <-timer.C:
			c <- tickTime
			nextTick = schedule.Next(tickTime)
			timer.Reset(time.Until(nextTick))
		}
	}
}
