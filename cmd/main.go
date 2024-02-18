package main

import (
	"fmt"
	"io"
	"os"
	"slices"
)

// Usage: echo <input_text> | your_grep.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Println("usage: mygrep -E <pattern>")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	ok, err := executePattern(os.Stdin, pattern)
	if err != nil {
		fmt.Println("failed to execute pattern on input", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}
}

func executePattern(reader io.Reader, pattern string) (bool, error) {
	matcher, err := NewMatcher(pattern)
	if err != nil {
		return false, err
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return false, err
	}
	line := string(data)

	for _, c := range line {
		if matcher.matches(c) {
			return true, nil
		}
	}
	return false, nil
}

type Matcher interface {
	matches(char rune) bool
}

func NewMatcher(pattern string) (Matcher, error) {
	switch pattern {
	case `\d`:
		return Digit{}, nil
	default:
		if len(pattern) != 1 {
			return nil, fmt.Errorf("unsupported pattern: %q", pattern)
		}
		return Char(pattern[0]), nil
	}
}

type Char rune

func (c Char) matches(char rune) bool {
	return rune(c) == char
}

type Digit struct{}

var DIGITS = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

func (_ Digit) matches(char rune) bool {
	return slices.Contains(DIGITS, char)
}
