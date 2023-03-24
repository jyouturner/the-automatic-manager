package main

import "testing"

func TestWatchCalendar(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "test monitor google calendar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WatchCalendar()
		})
	}
}
