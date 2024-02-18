package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestGrep(t *testing.T) {
	grepTests := []struct {
		input   string
		pattern string
		want    bool
	}{
		{"apple", "a", true},
		{"dog", "a", false},
	}

	for _, tt := range grepTests {
		t.Run(fmt.Sprintf("input: %q, pattern: %q", tt.input, tt.pattern), func(t *testing.T) {
			got, _ := executePattern(bytes.NewBufferString(tt.input), tt.pattern)
			if got != tt.want {
				t.Errorf("got %v want %v", got, tt.want)
			}
		})
	}
}
