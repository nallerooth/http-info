package colors

import (
	"strings"
	"testing"
)

func TestColors(t *testing.T) {
	testCases := []struct {
		text  string
		color string
		fn    func(string) string
	}{
		{"red", "\033[31m", Red},
		{"green", "\033[32m", Green},
		{"yellow", "\033[33m", Yellow},
		{"blue", "\033[34m", Blue},
		{"purple", "\033[35m", Purple},
		{"cyan", "\033[36m", Cyan},
		{"gray", "\033[37m", Gray},
		{"white", "\033[97m", White},
	}

	for _, tc := range testCases {
		t.Run(tc.text, func(t *testing.T) {
			result := tc.fn(tc.text)
			if !strings.Contains(result, tc.text) {
				t.Errorf("Text: %s not found in %s", tc.text, result)
			}
			if !strings.Contains(result, tc.color) {
				t.Errorf("Color: %s not found in %s", tc.color, result)
			}
			if !strings.Contains(result, "\033[0m") {
				t.Errorf("Color not reset at end of string")
			}
		})

	}
}
