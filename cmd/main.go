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
	case `\w`:
		return Word{}, nil
	default:
		if pattern[0] == '[' && pattern[len(pattern)-1] == ']' {
			return CharGroup(pattern[1 : len(pattern)-1]), nil
		}

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

func (_ Digit) matches(char rune) bool {
	code := int(char)
	return code >= 48 && code <= 57
}

type Word struct{}

func (_ Word) matches(char rune) bool {
	code := int(char)
	isDigit := code >= 48 && code <= 57
	isCapitalLetter := code >= 65 && code <= 90
	isLowerLetter := code >= 97 && code <= 122
	isUnderscore := code == 95
	return isDigit || isCapitalLetter || isLowerLetter || isUnderscore
}

type CharGroup string

func (g CharGroup) matches(char rune) bool {
	return slices.Contains([]rune(g), char)
}
