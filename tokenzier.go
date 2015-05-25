package mockingbird

import (
	"regexp"
	"strings"

	"github.com/lazywei/mockingbird/scanner"
)

var (
	// Start state on token, ignore anything till the next newline

	// Why would we need the trailing space here?
	singleLineCmtPtrn = `\s*\/\/ |\s*\-\- |\s*# |\s*% |\s*" `

	// Start state on opening token, ignore anything until the closing
	// token is reached.
	multiLineCmtPtrn = `\/\*|<!--|{-|\(\*|"""|'''`

	multiLineCmtPairs = map[string]string{
		`/*`:   `\*\/`, // C
		`<!--`: `-->`,  // XML
		`{-`:   `-}`,   // Haskell
		`(*`:   `\*\)`, // Coq
		`"""`:  `"""`,  // Python
		`'''`:  `'''`,  // Python
	}
)

func ExtractTokens(data string) []string {

	s := scanner.NewScanner(data)
	tokens := []string{}

	for s.IsEos() != true {

		if token, ok := s.Scan(`^#!.+`); ok {
			name, ok := extractShebang(token)
			if ok {
				tokens = append(tokens, "SHEBANG#!"+name)
			}
		} else if s.IsBol() && scanOrNot(s, singleLineCmtPtrn) {

			s.SkipUntil(`\n|\z`)

		} else if startToken, ok := s.Scan(multiLineCmtPtrn); ok {

			closeToken := multiLineCmtPairs[startToken]
			s.SkipUntil(closeToken)

		} else if scanOrNot(s, `"`) {

			if s.Peek(1) == `"` {
				s.Getch()
			} else {
				s.ScanUntil(`[^\\]"`)
			}

		} else if scanOrNot(s, `'`) {

			if s.Peek(1) == `'` {
				s.Getch()
			} else {
				s.ScanUntil(`[^\\]'`)
			}

		} else if scanOrNot(s, `(0x)?\d(\d|\.)*`) {
			// Skip number literals

		} else if rtn, ok := s.Scan(`<[^\s<>][^<>]*>`); ok {

			for _, tkn := range extractSgmlTokens(rtn) {
				tokens = append(tokens, tkn)
			}

		} else if rtn, ok := s.Scan(`;|\{|\}|\(|\)|\[|\]`); ok {
			// Common programming punctuation
			tokens = append(tokens, rtn)

		} else if rtn, ok := s.Scan(`[\w\.@#\/\*]+`); ok {
			// Regular token
			tokens = append(tokens, rtn)

		} else if rtn, ok := s.Scan(`<<?|\+|\-|\*|\/|%|&&?|\|\|?`); ok {
			// Common operators
			tokens = append(tokens, rtn)

		} else {
			s.Getch()
		}

	}

	return tokens
}

// This function is silly ... silly Go ...
func scanOrNot(s *scanner.Scanner, ptrn string) bool {
	_, ok := s.Scan(ptrn)
	return ok
}

func extractSgmlTokens(token string) []string {
	s := scanner.NewScanner(token)
	tokens := []string{}

	for s.IsEos() != true {
		if token, ok := s.Scan(`<\/?[^\s>]+`); ok {
			tokens = append(tokens, token+">")
		} else if token, ok := s.Scan(`\w+=`); ok {
			tokens = append(tokens, token)

			// Then skip over attribute value

			if scanOrNot(s, `"`) {
				s.SkipUntil(`[^\\]"`)
			} else if scanOrNot(s, `'`) {
				s.SkipUntil(`[^\\]'`)
			} else {
				s.SkipUntil(`\w+`)
			}
		} else if token, ok := s.Scan(`\w+`); ok {
			tokens = append(tokens, token)
		} else if scanOrNot(s, `\w+`) {
			s.Terminate()
		} else {
			s.Getch()
		}
	}

	return tokens
}

func extractShebang(token string) (string, bool) {
	s := scanner.NewScanner(token)

	path, ok := s.Scan(`^#!\s*\S+`)

	if !ok {
		return "", false
	}

	paths := strings.Split(path, `/`)

	if len(paths) == 0 {
		return "", false
	}

	name := paths[len(paths)-1]

	if name == `env` {
		s.Scan(`\s+`)
		s.Scan(`.*=[^\s]+\s+`)
		name, ok = s.Scan(`\S+`)
	}

	if ok {
		name = regexp.MustCompile(`[^\d]+`).FindString(name)
		return name, true
	} else {
		return "", false
	}
}
