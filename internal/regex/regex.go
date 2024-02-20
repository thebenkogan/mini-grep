package regex

type Regex struct {
	matchers    []matcher
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

	matchers, err := parseMatchers(regex)
	if err != nil {
		return nil, err
	}

	return &Regex{
		matchers,
		assertStart,
		assertEnd,
	}, nil
}

func (r *Regex) Match(text string) bool {
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

func matchHere(matchers []matcher, text string, assertEnd bool) bool {
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
