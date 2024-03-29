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

		{"apple123", `\d`, true},
		{"c", `\d`, false},

		{"foo101", `\w`, true},
		{"$!?", `\w`, false},

		{"apple", `[abc]`, true},
		{"dog", `[abc]`, false},
		{"[]", `[abc]`, false},

		{"dog", `[^abc]`, true},
		{"cab", `[^abc]`, false},

		{"1 apple", `\d apple`, true},
		{"app - BK: 123-456", `\w\w: \d\d\d-\d\d\d`, true},
		{"app - BK: 12-3456", `\w\w: \d\d\d-\d\d\d`, false},

		{"log", `^log`, true},
		{"slog", `^log`, false},

		{"dog", `dog$`, true},
		{"dogs", `dog$`, false},
		{"dog", `^dog$`, true},

		{"apple", `a+`, true},
		{"SaaS", `a+`, true},
		{"dog", `a+`, false},
		{"cake", `cake+`, true},
		{"baaadaaac", `a+e`, false},
		{"baaad", `a+a`, true},

		{"dog", `dogs?`, true},
		{"dogs", `dogs?`, true},
		{"-3", `-?\d`, true},
		{"3", `-?\d`, true},
		{"33", `^-?\d$`, false},
		{"#bag", `#?tag`, false},

		{"dog", `do*g`, true},
		{"dooooog", `do*g`, true},
		{"dg", `do*g`, true},
		{"dooot", `do*g`, false},

		{"dog", `d.g`, true},
		{"cog", `d.g`, false},
	}

	for _, tt := range grepTests {
		t.Run(fmt.Sprintf("input: %q, pattern: %q", tt.input, tt.pattern), func(t *testing.T) {
			got, err := executePattern(bytes.NewBufferString(tt.input), tt.pattern)

			if err != nil {
				t.Fatal(err)
			}

			if got != tt.want {
				t.Errorf("got %v want %v", got, tt.want)
			}
		})
	}
}
