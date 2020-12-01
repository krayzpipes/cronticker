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
