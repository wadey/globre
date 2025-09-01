package globre

import (
	"fmt"
	"regexp"
	"strings"
)

func MustCompile(glob string) *regexp.Regexp {
	return MustCompileSeparators(glob, "")
}

func MustCompileSeparators(glob string, separators string) *regexp.Regexp {
	r, err := CompileSeparators(glob, separators)
	if err != nil {
		panic(err)
	}
	return r
}

func Compile(glob string) (*regexp.Regexp, error) {
	return CompileSeparators(glob, "")
}

func CompileSeparators(glob string, separators string) (*regexp.Regexp, error) {
	s, err := ConvertSeparators(glob, separators)
	if err != nil {
		return nil, err
	}
	return regexp.Compile(s)
}

func Convert(glob string) (string, error) {
	return ConvertSeparators(glob, "")
}

func ConvertSeparators(glob string, separators string) (string, error) {
	var sb strings.Builder

	var inGroup int
	var inSet bool

	var last rune
	var wildStarted bool
	var backslash bool

	sb.WriteString("^")

	for _, c := range glob {
		if wildStarted && c != '*' {
			sb.WriteString("[^")
			sb.WriteString(regexp.QuoteMeta(separators))
			sb.WriteString("]*")
			wildStarted = false
		}
		if backslash {
			sb.WriteRune(c)
			backslash = false
			// This char was escaped, so don't consider it for lookback
			last = 0
			continue
		}
		if inSet {
			switch c {
			case '!':
				sb.WriteRune('^')
			case ']':
				sb.WriteRune(c)
				inSet = false
			case '\\':
				backslash = true
				sb.WriteRune(c)
			default:
				sb.WriteString(regexp.QuoteMeta(string(c)))
			}

			last = c
			continue
		}

		switch c {
		case '*':
			if wildStarted {
				sb.WriteString(".*")
				wildStarted = false
			} else if len(separators) > 0 {
				wildStarted = true
			} else {
				sb.WriteString(".*")
			}

		case '?':
			if len(separators) > 0 {
				sb.WriteString("[^")
				sb.WriteString(regexp.QuoteMeta(separators))
				sb.WriteRune(']')
			} else {
				sb.WriteRune('.')
			}

		case '{':
			inGroup++
			sb.WriteRune('(')

		case '}':
			inGroup--
			if inGroup < 0 {
				return "", fmt.Errorf("unexpected: }")
			}
			sb.WriteRune(')')

		case '[':
			inSet = true
			sb.WriteRune(c)

		case ']':
			sb.WriteRune(c)

		case '!':
			if last == '[' {
				sb.WriteRune('^')
			} else {
				sb.WriteRune(c)
			}

		case ',':
			if inGroup > 0 {
				sb.WriteRune('|')
			} else {
				sb.WriteRune(c)
			}

		case '\\':
			backslash = true
			sb.WriteRune(c)

		default:
			sb.WriteString(regexp.QuoteMeta(string(c)))
		}

		last = c
	}
	if wildStarted {
		sb.WriteString("[^")
		sb.WriteString(regexp.QuoteMeta(separators))
		sb.WriteString("]*")
	}

	sb.WriteString("$")

	return sb.String(), nil
}
