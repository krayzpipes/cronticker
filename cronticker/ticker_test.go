package cronticker

import (
	"fmt"
	"testing"
	"log"
	"time"
)


func TestGuaranteeTimeZone_MissingTZReturnsUTC(t *testing.T) {
	schedule := "0 0 0 * * ?"
	expected := fmt.Sprintf("TZ=UTC %s", schedule)

	s, loc, err := guaranteeTimeZone(schedule)
	if expected != s {
		log.Fatalf("expected: %q, got: %q", expected, s)
	}
	expectedLoc, _ := time.LoadLocation("UTC")
	if loc != expectedLoc {
		log.Fatalf("expected: %q, got: %q", expectedLoc, loc)
	}
	if err != nil {
		log.Fatalf("err: %v", err)
	}
}

func TestGuaranteeTimeZone_InvalidTZ(t *testing.T) {
	schedule := "TZ=BadZone 0 0 0 * * ?"

	_, _, err := guaranteeTimeZone(schedule)

	expected := "unknown time zone BadZone"

	if expected != fmt.Sprintf("%s", err) {
		log.Fatalf("expected err: %s, got: %s", expected, err)
	}
}

func TestGuaranteeTimeZone_ValidTZ(t *testing.T) {
	tz := "America/New_York"
	schedule := fmt.Sprintf("TZ=%s 0 0 0 * * ?", tz)
	expectedTZ, _ := time.LoadLocation(tz)

	s, loc, err := guaranteeTimeZone(schedule)

	if s != schedule {
		log.Fatalf("expected %s, got %s", schedule, s)
	}
	locString := fmt.Sprintf("%s", loc)
	expectedString := fmt.Sprintf("%s", expectedTZ)
	if locString != expectedString {
		log.Fatalf("expected time.Location of %s, got %s", expectedTZ, loc)
	}
	if err != nil {
		log.Fatalf("err: %s", err)
	}
}

func TestExtractTZ(t *testing.T) {
	expected := "America/New_York"
	input  := fmt.Sprintf("TZ=%s 0 0 0 * * ?", expected)

	tz := extractTZ(input)

	if expected != tz {
		log.Fatalf("expected %q, got: %q", expected, tz)
	}
}

func TestCronTicker_Stop(t *testing.T) {
    ticker, _ := NewTicker("@daily")

    timeoutTimer := time.NewTimer(2 * time.Second)

    kCopy := ticker.k
    ticker.Stop()
Outer:
    for {
        select {
        case <-kCopy:
            break Outer
        case <-timeoutTimer.C:
            log.Fatal("Expected message on ticker 'k' channel within 2 seconds, but did not receive one")
        }
    }
}

func TestCronTicker_Reset_Error(t *testing.T) {
    ticker, _ := NewTicker("@daily")
    defer ticker.Stop()
    err := ticker.Reset("NOT_VALID_SCHEDULE")
    if err == nil {
        log.Fatal("should have gotten error, but received 'nil'")
    }
}

func TestCronTicker_Reset(t *testing.T) {
    ticker, _ := NewTicker("@daily")
    defer ticker.Stop()
    err := ticker.Reset("@monthly")
    if err != nil {
        log.Fatalf("expected 'nil', got: %q", err)
    }
}

func TestNewTicker_Error(t *testing.T) {
    _, err := NewTicker("NOT_VALID_SCHEDULE")
    if err == nil {
        log.Fatal("expected error, received 'nil'")
    }
}

//func TestnewTicker_ErrorFromGuarantee(t *testing.T) {
//    _, err := NewTicker("TZ=NOT_VALID @daily")
//    if err == nil {
//        log.Fatal("expected error due to TZ parsing, got 'nil'")
//    }
//}

// Examples for documentation

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

	ticker, err := NewTicker("@daily")
	if err != nil {
		log.Fatal(err)
	}
	defer ticker.Stop()

	tick := <-ticker.C
	log.Print(tick)
}

// If you want to change the cron schedule of a ticker
// instead of creating a new one you can reset it.
func ExampleCronTicker_Reset() {
	ticker, err := NewTicker("TZ=UTC 0 0 0 ? * SUN")
	if err != nil {
		log.Fatal(err)
	}
	defer ticker.Stop()

	<-ticker.C
	log.Print("It's Sunday!")

	err = ticker.Reset("TZ=UTC 0 0 0 ? * WED")
	if err != nil {
		log.Fatal(err)
	}

	<-ticker.C
	log.Print("It's Wednesday!")
}

func ExampleCronTicker_Stop() {
	ticker, err := NewTicker("TZ=UTC 0 0 0 ? * SUN")
	if err != nil {
		log.Fatal(err)
	}
	defer ticker.Stop()

	<-ticker.C
}