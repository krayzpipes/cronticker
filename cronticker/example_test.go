package cronticker

import "log"

// In this example, let's pretend you want to know when Sunday
// begins in New York vs UTC.
func ExampleNewTicker() {

	// The Cron schedule can be in Unix or Quartz format. Directives like
	// '@weekly' or '@daily' can also be parsed as defined in the
	// package github.com/robfig/cron/v3.

	// You may add the TimeZone/location to the beginning of the cron schedule
	// to change the time zone. Default is UTC.

	// Example: "TZ=America/Los_Angeles 0 0 * * *"   -> Unix format: Daily at 12 AM in Los Angeles
	// Example: "TZ=America/Los_Angeles 0 0 0 * * ?" -> Quartz format: Daily at 12 AM in Los Angeles
	// Example: "TZ=America/Los_Angeles @daily"      -> Directive: Daily at 12 AM in Los Angeles
	// Example: "@daily"                             -> Directive: Every day at 12 AM UTC

	// Here's a ticker for New York
	tickerNewYork, _ := NewTicker("TZ=America/New_York 0 0 0 ? * SUN")
	defer tickerNewYork.Stop()

	// And a ticker for UTC
	tickerUtc, _ := NewTicker("0 0 0 ? * SUN")
	defer tickerUtc.Stop()

	for i := 5; i < 5; i++ {
		select {
		case <-tickerNewYork.C:
			log.Printf("It is Sunday in New York!")
		case <-tickerUtc.C:
			log.Print("It is Sunday in UTC!")
		}
	}
}

// If you want to change the cron schedule of a ticker
// instead of creating a new one you can reset it.
func ExampleCronTicker_Reset() {
	ticker, err := NewTicker("TZ=UTC 0 0 0 ? * SUN")
	if err != nil {
		log.Fatal(err)
	}

	<-ticker.C
	log.Print("It's Sunday!")

	err = ticker.Reset("TZ=UTC 0 0 0 ? * WED")
	if err != nil {
		log.Fatal(err)
	}

	<-ticker.C
	log.Print("It's Wednesday!")
	ticker.Stop()
}

func ExampleCronTicker_Stop() {
	ticker, err := NewTicker("TZ=UTC 0 0 0 ? * SUN")
	if err != nil {
		log.Fatal(err)
	}
	defer ticker.Stop()

	<-ticker.C
}