package templates

import (
	"strings"
	"time"
)

// PrintWeekdays is a function that is expected to take a date and return the full date of each weekday of that week,
// separated by dividers
func PrintWeekdays(date time.Time) (string, error) {
	var result strings.Builder

	// Need to find the previous Monday
	weekday := date.Weekday()
	offsetDays := int(weekday-time.Monday) % 7

	// Move back to the previous monday
	date = date.AddDate(0, 0, -offsetDays)

	for i := 0; i < 7; i++ {
		day := date.Day()
		var suffix string
		switch day {
		case 1, 21, 31:
			suffix = "st"
		case 2, 22:
			suffix = "nd"
		case 3, 23:
			suffix = "rd"
		default:
			suffix = "th"
		}
		result.WriteString(date.Format("Monday, Jan 2") + suffix + ":")
		date = date.AddDate(0, 0, 1)

		// Add dividers if this is not the last day
		if i < 6 {
			result.WriteString("\n")      // Add newlines after each date
			result.WriteString("---\n") // Add newlines after each date
			result.WriteString("\n")      // Add newlines after each date
		}
	}
	return result.String(), nil
}
