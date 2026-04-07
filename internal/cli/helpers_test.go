package cli

import (
	"slices"
	"testing"
)

func TestParseRate(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedCents int64
		expectError   bool
	}{
		{"Standard rate", "40.00", 4000, false},
		{"Whole number", "40", 4000, false},
		{"Too many decimals", "40.001", 0, true},
		{"No leading number", ".40", 40, true},
		{"Empty string", "", 0, false},
		{"Non number", "this is a test", 0, true},
		{"Negative number in front", "-40.50", 0, true},
		{"Negative number in back", "40.-50", 0, true},
		{"Non number in front", "test.44", 0, true},
		{"Non number in back", "44.test", 0, true},
		{"One digit cent", "44.4", 4440, false},
		{"No zeroes", "505.55", 50555, false},
		{"Zero", "0", 0, false},
		{"Zero with decimals", "0.00", 0, false},
		{"Trailing decimal", "40.", 0, true},
		{"Leading zeroes", "040.00", 4000, false},
		{"Spaces", " 40.00 ", 0, true},
		{"Only spaces", "   ", 0, true},
		{"Negative zero dollars", "-0.50", 0, true},
		{"Negative zero cents", "0.-50", 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cents, err := parseRate(tc.input)
			if (err != nil) != tc.expectError {
				t.Errorf("Expected err: %t, got err: %v", tc.expectError, err)
			}
			if !tc.expectError && tc.expectedCents != cents {
				t.Errorf("expected cents: %v, got cents %v", tc.expectedCents, cents)
			}
		})
	}
}

func TestParseTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"Empty string", "", nil},
		{"Single tag", "tag1", []string{"tag1"}},
		{"Multiple tags", "tag1,tag2", []string{"tag1", "tag2"}},
		{"Tags with spaces", " tag1 , tag2 ", []string{"tag1", "tag2"}},
		{"Leading and trailing commas", ",tag1,", []string{"tag1"}},
		{"Consecutive commas", "tag1,,tag2", []string{"tag1", "tag2"}},
		{"Only commas and spaces", " , , ", nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := parseTags(tc.input)
			if !slices.Equal(result, tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}
