package main

import (
	"testing"
)

func TestWatchJira(t *testing.T) {

	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "test monitor jira",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WatchJira()
		})
	}
}
