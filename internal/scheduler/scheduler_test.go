package scheduler

import "testing"

func TestParseScheduleTime(t *testing.T) {
	tests := []struct {
		name          string
		scheduleTime  string
		expectedHour  int
		expectedMin   int
		expectedError bool
	}{
		{
			name:         "valid time",
			scheduleTime: "14:30",
			expectedHour: 14,
			expectedMin:  30,
		},
		{
			name:         "midnight",
			scheduleTime: "00:00",
			expectedHour: 0,
			expectedMin:  0,
		},
		{
			name:         "end of day",
			scheduleTime: "23:59",
			expectedHour: 23,
			expectedMin:  59,
		},
		{
			name:          "invalid format - no colon",
			scheduleTime:  "1430",
			expectedError: true,
		},
		{
			name:          "invalid format - too many parts",
			scheduleTime:  "14:30:00",
			expectedError: true,
		},
		{
			name:          "invalid hour - too high",
			scheduleTime:  "25:30",
			expectedError: true,
		},
		{
			name:          "invalid hour - negative",
			scheduleTime:  "-1:30",
			expectedError: true,
		},
		{
			name:          "invalid minute - too high",
			scheduleTime:  "14:60",
			expectedError: true,
		},
		{
			name:          "invalid minute - negative",
			scheduleTime:  "14:-1",
			expectedError: true,
		},
		{
			name:          "non-numeric hour",
			scheduleTime:  "ab:30",
			expectedError: true,
		},
		{
			name:          "non-numeric minute",
			scheduleTime:  "14:cd",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hour, minute, err := parseScheduleTime(tt.scheduleTime)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error for input %q, but got none", tt.scheduleTime)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for input %q: %v", tt.scheduleTime, err)
				return
			}

			if hour != tt.expectedHour {
				t.Errorf("Hour mismatch for %q: expected %d, got %d", tt.scheduleTime, tt.expectedHour, hour)
			}

			if minute != tt.expectedMin {
				t.Errorf("Minute mismatch for %q: expected %d, got %d", tt.scheduleTime, tt.expectedMin, minute)
			}
		})
	}
}
