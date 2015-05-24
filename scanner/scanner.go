// This package provides Ruby's StringScanner-like functions.
package scanner

import "regexp"

type Scanner struct {
	input string
	pos   int
}

func NewScanner(input string) *Scanner {
	return &Scanner{
		input: input,
		pos:   0,
	}
}

func (scn *Scanner) Scan(regExpStr string) (rtn string, ok bool) {
	if scn.IsEos() {
		return "", false
	}

	strForScan := scn.input[scn.pos:]
	loc := regexp.MustCompile(regExpStr).FindStringIndex(strForScan)

	if loc == nil {
		rtn = ""
		ok = false
		return
	} else if loc[0] != 0 {
		rtn = ""
		ok = false
		return
	} else {

		rtn = scn.input[scn.pos+loc[0] : scn.pos+loc[1]]
		ok = true

		scn.pos = loc[1] + scn.pos

		return
	}
}

func (scn *Scanner) ScanUntil(regExpStr string) (rtn string, ok bool) {
	if scn.IsEos() {
		return "", false
	}

	strForScan := scn.input[scn.pos:]
	loc := regexp.MustCompile(regExpStr).FindStringIndex(strForScan)

	if loc == nil {
		rtn = ""
		ok = false
		return
	} else {

		rtn = scn.input[scn.pos : scn.pos+loc[1]]
		ok = true

		scn.pos = loc[1] + scn.pos

		return
	}
}

func (scn *Scanner) SkipUntil(regExpStr string) (rtn int, ok bool) {
	rtnStr, ok := scn.ScanUntil(regExpStr)
	if ok {
		rtn = len(rtnStr)
	} else {
		rtn = 0
	}

	return
}

func (scn *Scanner) Getch() (rtn string, ok bool) {
	if scn.IsEos() {
		return "", false
	}

	rtn = scn.input[scn.pos : scn.pos+1]
	scn.pos = scn.pos + 1
	ok = true
	return
}

func (scn *Scanner) IsEos() bool {
	return scn.pos >= len(scn.input)
}
func (scn *Scanner) IsBol() bool {

	if scn.pos == 0 {
		return true
	} else if scn.pos >= 2 && scn.input[scn.pos-1:scn.pos] == "\n" {
		return true
	} else {
		return false
	}
}

func (scn *Scanner) Peek(length int) string {
	if scn.IsEos() {
		return ""
	} else {
		end := scn.pos + length

		if end >= len(scn.input) {
			end = len(scn.input)
		}

		return scn.input[scn.pos:end]
	}
}
