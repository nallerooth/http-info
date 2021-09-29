package cert

import (
	"fmt"
	"time"

	"github.com/nallerooth/http-info/colors"
)

// CalcRemainingDays returns the number of days until a time.Time occurs
func CalcRemainingDays(certTimestamp time.Time) int {
	t := time.Now()
	return int(certTimestamp.Sub(t).Hours() / 24)
}

// CalcRemainingDaysColor returns a colored string of of the number of days
// remaining until a time.Time occurs
func CalcRemainingDaysColor(certTimestamp time.Time) string {
	days := CalcRemainingDays(certTimestamp)
	color := colors.Green

	if days <= 30 {
		color = colors.Yellow
	}
	if days <= 15 {
		color = colors.Red
	}
	return color(fmt.Sprintf("%d days remaining", days))
}
