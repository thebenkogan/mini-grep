package regex

import (
	"fmt"
	"math"
	"slices"
	"strings"
)

type matcher struct {
	symbol     symbol
	quantifier *quantifier
}

type quantifier struct {
	low  int
	high int
}

var runeToQuantifier = map[rune]quantifier{
	'+': {1, math.MaxInt},
	'*': {0, math.MaxInt},
	'?': {0, 1},
}

func parseMatchers(pattern string) ([]matcher, error) {
	matchers := make([]matcher, 0)

	i := 0
	for i < len(pattern) {
		var symbol symbol
		var err error
		switch pattern[i] {
		case '\\':
			symbol, err = newMetaCharacter(pattern[i : i+2])
			if err != nil {
				return nil, err
			}
			i += 2
		case '[':
			closing := strings.Index(pattern[i:], "]")
			if closing == -1 {
				return nil, fmt.Errorf("unclosed character group: %q", pattern[i:])
			}
			symbol = newCharacterGroup(pattern[i+1 : closing])
			i = closing + 1
		case '.':
			symbol = wildcard{}
			i += 1
		default:
			symbol = char(pattern[i])
			i += 1
		}

		matcher := matcher{symbol, nil}
		if i < len(pattern) {
			if r, ok := runeToQuantifier[rune(pattern[i])]; ok {
				matcher.quantifier = &r
				i += 1
			}
		}
		matchers = append(matchers, matcher)
	}

	return matchers, nil
}

type symbol interface {
	matches(char rune) bool
}

type char rune

func (c char) matches(char rune) bool {
	return rune(c) == char
}

func newMetaCharacter(pattern string) (symbol, error) {
	switch pattern {
	case `\d`:
		return digit{}, nil
	case `\w`:
		return word{}, nil
	default:
		return nil, fmt.Errorf("unsupported meta character: %q", pattern)
	}
}

type digit struct{}

func (_ digit) matches(char rune) bool {
	code := int(char)
	return code >= 48 && code <= 57
}

type word struct{}

func (_ word) matches(char rune) bool {
	code := int(char)
	isDigit := code >= 48 && code <= 57
	isCapitalLetter := code >= 65 && code <= 90
	isLowerLetter := code >= 97 && code <= 122
	isUnderscore := code == 95
	return isDigit || isCapitalLetter || isLowerLetter || isUnderscore
}

type wildcard struct{}

func (_ wildcard) matches(char rune) bool {
	return true
}

type charGroup struct {
	group    string
	negative bool
}

func newCharacterGroup(group string) charGroup {
	negative := false
	if group[0] == '^' {
		negative = true
		group = group[1:]
	}
	return charGroup{group, negative}
}

func (cg charGroup) matches(char rune) bool {
	inGroup := slices.Contains([]rune(cg.group), char)
	if cg.negative {
		return !inGroup
	}
	return inGroup
}
