package mockingbird

import "github.com/lazywei/mockingbird/scanner"

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

// This function is silly ... silly Go ...
func scanOrNot(s *scanner.Scanner, ptrn string) bool {
	_, ok := s.Scan(ptrn)
	return ok
}

func ExtractTokens(data string) []string {

	s := scanner.NewScanner(data)
	tokens := []string{}

	for s.IsEos() != true {

		if false /* SHEBANG, COMMENTS */ {
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
