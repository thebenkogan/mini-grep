package main

import (
	"fmt"
	"io"
	"math"
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
	regex, err := NewRegex(pattern)
	if err != nil {
		return false, err
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return false, err
	}
	line := string(data)

	return regex.match(line), nil
}

type Regex struct {
	matchers    []Matcher
	assertStart bool
	assertEnd   bool
}

func NewRegex(regex string) (*Regex, error) {
	assertStart := false
	if regex[0] == '^' {
		assertStart = true
		regex = regex[1:]
	}

	assertEnd := false
	if regex[len(regex)-1] == '$' {
		assertEnd = true
		regex = regex[:len(regex)-1]
	}

	matchers, err := parsePattern(regex)
	if err != nil {
		return nil, err
	}

	return &Regex{
		matchers,
		assertStart,
		assertEnd,
	}, nil
}

func (r *Regex) match(text string) bool {
	if r.assertStart {
		return matchHere(r.matchers, text, r.assertEnd)
	}

	for i := 0; i < len(text); i++ {
		if matchHere(r.matchers, text[i:], r.assertEnd) {
			return true
		}
	}

	return false
}

func matchHere(matchers []Matcher, text string, assertEnd bool) bool {
	if len(matchers) == 0 {
		return !assertEnd || text == ""
	}
	matcher := matchers[0]
	if len(text) == 0 && matcher.quantifier == nil {
		return false
	}

	if matcher.quantifier != nil {
		for i := 0; i < matcher.quantifier.low; i++ {
			if !matcher.symbol.matches(rune(text[i])) {
				return false
			}
		}

		for i := matcher.quantifier.low; i < len(text)+1 && i <= matcher.quantifier.high; i++ {
			if matchHere(matchers[1:], text[i:], assertEnd) {
				return true
			}
			if !matcher.symbol.matches(rune(text[i])) {
				break
			}
		}
		return false
	}

	if !matcher.symbol.matches(rune(text[0])) {
		return false
	}

	if len(text) > 0 {
		return matchHere(matchers[1:], text[1:], assertEnd)
	}

	return false
}

type Range struct {
	low  int
	high int
}

var quantifierToRange = map[rune]Range{
	'+': {1, math.MaxInt},
	'?': {0, 1},
}

type Matcher struct {
	symbol     Symbol
	quantifier *Range
}

func parsePattern(pattern string) ([]Matcher, error) {
	matchers := make([]Matcher, 0)

	i := 0
	for i < len(pattern) {
		var symbol Symbol
		var err error
		switch pattern[i] {
		case '\\':
			symbol, err = NewMetaCharacter(pattern[i : i+2])
			if err != nil {
				return nil, err
			}
			i += 2
		case '[':
			closing := strings.Index(pattern[i:], "]")
			if closing == -1 {
				return nil, fmt.Errorf("unclosed character group: %q", pattern[i:])
			}
			symbol = NewCharacterGroup(pattern[i+1 : closing])
			i = closing + 1
		default:
			symbol = Char(pattern[i])
			i += 1
		}

		matcher := Matcher{symbol, nil}
		if i < len(pattern) {
			if r, ok := quantifierToRange[rune(pattern[i])]; ok {
				matcher.quantifier = &r
				i += 1
			}
		}
		matchers = append(matchers, matcher)
	}

	return matchers, nil
}

type Symbol interface {
	matches(char rune) bool
}

type Char rune

func (c Char) matches(char rune) bool {
	return rune(c) == char
}

func NewMetaCharacter(pattern string) (Symbol, error) {
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
