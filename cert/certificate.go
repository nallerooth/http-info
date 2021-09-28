package cert

import (
	"time"
)

func CalcRemainingDays(certTimestamp time.Time) int {
	t := time.Now()
	return int(certTimestamp.Sub(t).Hours() / 24)
}
