package mockingbird

import (
	"github.com/lazywei/mockingbird/scanner"
	"github.com/taylorchu/toki"
)

const (
	NUMBER toki.Token = iota + 1
	PLUS
	STRING
	LS
)

const (
	NUMBER_PTRN = `(0x)?\d(\d|\.)*`
	PLUS_PTRN   = `\+`
	STRING_PTRN = "[a-z]+"
	LS_PTRN     = `\"(\\.|[^"])*\"`
)

var (
	// Start state on token, ignore anything till the next newline
	singleLineComments = []string{
		`//`, // C
		`--`, // Ada, Haskell, AppleScript
		`#`,  // Ruby
		`%`,  // Tex
		`"`,  // Vim
	}

	// Start state on opening token, ignore anything until the closing
	// token is reached.
	multiLineComments = [][]string{
		[]string{`/*`, `*/`},    // C
		[]string{`<!--`, `-->`}, // XML
		[]string{`{-`, `-}`},    // Haskell
		[]string{`(*`, `*)`},    // Coq
		[]string{`"""`, `"""`},  // Python
		[]string{`'''`, `'''`},  // Python
	}
)

func ExtractTokens(data string) []string {

	s := scanner.NewScanner(data)
	tokens := []string{}

	for s.IsEof() != true {

		if false /* SHEBANG, COMMENTS */ {

		} else if _, ok := s.Scan(`"`); ok {

			if s.Peek(1) == `"` {
				s.Getch()
			} else {
				s.ScanUntil(`[^\\]"`)
			}

		} else if _, ok := s.Scan(`'`); ok {

			if s.Peek(1) == `'` {
				s.Getch()
			} else {
				s.ScanUntil(`[^\\]'`)
			}

		} else if _, ok := s.Scan(`(0x)?\d(\d|\.)*`); ok {
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
