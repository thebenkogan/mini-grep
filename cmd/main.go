package main

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

// Usage: echo <input_text> | ./build/main -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Println("usage: mygrep -E <pattern>")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	ok, err := executePattern(os.Stdin, pattern)
	if err != nil {
		fmt.Println("failed to execute pattern on input:", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}
}

func executePattern(reader io.Reader, pattern string) (bool, error) {
	assertStartOfLine := false
	if pattern[0] == '^' {
		assertStartOfLine = true
		pattern = pattern[1:]
	}

	assertEndOfLine := false
	if pattern[len(pattern)-1] == '$' {
		assertEndOfLine = true
		pattern = pattern[:len(pattern)-1]
	}

	matchers, err := parsePattern(pattern)
	if err != nil {
		return false, err
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return false, err
	}
	line := string(data) // TODO: read char by char

	start := 0
	if assertEndOfLine {
		start = len(line) - len(matchers)
	}

	end := len(line) - len(matchers) + 1
	if assertStartOfLine {
		end = 1
	}

	for i := start; i < end; i++ {
		matchesAll := true
		for j := 0; j < len(matchers); j++ {
			if !matchers[j].matches(rune(line[i+j])) {
				matchesAll = false
				break
			}
		}
		if matchesAll {
			return true, nil
		}
		if assertStartOfLine {
			break
		}
	}
	return false, nil
}

func parsePattern(pattern string) ([]Matcher, error) {
	matchers := make([]Matcher, 0)
	i := 0
	for i < len(pattern) {
		switch pattern[i] {
		case '\\':
			matcher, err := NewMetaCharacter(pattern[i : i+2])
			if err != nil {
				return nil, err
			}
			matchers = append(matchers, matcher)
			i += 2
		case '[':
			closing := strings.Index(pattern[i:], "]")
			if closing == -1 {
				return nil, fmt.Errorf("unclosed character group: %q", pattern[i:])
			}
			matchers = append(matchers, NewCharacterGroup(pattern[i+1:closing]))
			i = closing + 1
		default:
			matchers = append(matchers, Char(pattern[i]))
			i += 1
		}
	}
	return matchers, nil
}

type Matcher interface {
	matches(char rune) bool
}

type Char rune

func (c Char) matches(char rune) bool {
	return rune(c) == char
}

func NewMetaCharacter(pattern string) (Matcher, error) {
	switch pattern {
	case `\d`:
		return Digit{}, nil
	case `\w`:
		return Word{}, nil
	default:
		return nil, fmt.Errorf("unsupported meta character: %q", pattern)
	}
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

type CharGroup struct {
	group    string
	negative bool
}

func NewCharacterGroup(group string) CharGroup {
	negative := false
	if group[0] == '^' {
		negative = true
		group = group[1:]
	}
	return CharGroup{group, negative}
}

func (cg CharGroup) matches(char rune) bool {
	inGroup := slices.Contains([]rune(cg.group), char)
	if cg.negative {
		return !inGroup
	}
	return inGroup
}
