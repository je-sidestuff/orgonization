package templates

import (
	"strings"
	"testing"
	"time"
)

func TestPrintWeekdays_Monday(t *testing.T) {
	loc, err := time.LoadLocation("Atlantic/Reykjavik")
	if err != nil {
		t.Fatal(err)
	}

	// Define a monday date
	date := time.Date(2024, 3, 18, 0, 0, 0, 0, loc)

	// Call PrintWeekdays and get the returned string
	result, err := PrintWeekdays(date)
	if err != nil {
		t.Fatal(err)
	}

	// Validate the result string
	expectedOutputStart := "Monday, Mar 18th:"
	expectedOutputEnd := "Sunday, Mar 24th:"
	if !strings.HasPrefix(result, expectedOutputStart) || !strings.HasSuffix(result, expectedOutputEnd) {
		t.Errorf("PrintWeekdays(%v) output did not match expected format.\nExpected to start with: %v\nExpected to end with: %v\nReceived: %v", date, expectedOutputStart, expectedOutputEnd, result)
	}
}
