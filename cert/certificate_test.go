package cert

import (
	"strings"
	"testing"
	"time"
)

func TestCertificate(t *testing.T) {
	t.Run("CalcRemainingDays", func(t *testing.T) {
		testCases := []struct {
			hours int
			days  int
		}{
			// Hours needs a 1 hour padding, as the calculation takes some time
			{25, 1},
			{241, 10},
			{36, 1},
			{78, 3},
			{24*365 + 1, 365},
			{24*768 + 1, 768},
		}

		for _, tc := range testCases {
			testTime := time.Now().Add(time.Duration(tc.hours) * time.Hour)
			result := CalcRemainingDays(testTime)
			if tc.days != result {
				t.Errorf("Expected %d, got %d", tc.days, result)
			}
		}
	})

	t.Run("CalcRemainingDaysColor", func(t *testing.T) {
		testCases := []struct {
			hours int
			color string
		}{
			{24 * 50, "[32m"}, // about 50 days -> green
			{24 * 29, "[33m"}, // about 29 days -> yellow
			{24 * 10, "[31m"}, // about 10 days -> red
		}

		for _, tc := range testCases {
			testTime := time.Now().Add(time.Duration(tc.hours) * time.Hour)
			result := CalcRemainingDaysColor(testTime)
			if !strings.Contains(result, tc.color) {
				t.Errorf("Color '%s' not found in '%s'", tc.color, result)
			}
			if !strings.Contains(result, "[0m") {
				t.Errorf("Reset color not found in %s", result)
			}
		}

	})
}
