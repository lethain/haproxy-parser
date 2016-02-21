package main

import (
	"bufio"
	"bytes"
	"io"
)

type Token int
const (
	ILLEGAL Token = iota
	EOF
	COMMENT
	SPACE
	NEWLINE
	TAB
	STRING
)

var humanTokens = map[Token]string {
	ILLEGAL: "ILLEGAL",
	EOF: "EOF",
	COMMENT: "COMMENT",
	SPACE: "SPACE",
	NEWLINE: "NEWLINE",
	TAB: "TAB",
	STRING: "STRING",
}



var eof = rune(0)


type tokenTypeChecker func(rune) bool

func isComment(ch rune) bool {
	return ch == '#'
}

func isTab(ch rune) bool {
	return ch == '\t'
}

func isSpace(ch rune) bool {
	return ch == ' '
}

func isNewline(ch rune) bool {
	return ch == '\n'
}

func isNotNewline(ch rune) bool {
	return !isNewline(ch)
}

func isQuote(ch rune) bool {
	return ch == '"'
}


func isLetter(ch rune) bool {
	return !(isQuote(ch) || isNewline(ch) || isSpace(ch) || isTab(ch) || ch == eof)
	// return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '.' || ch == '/' || ch == ':' || (ch >= '0' && ch <= '9') || ch == '-' || ch == '_' || ch == '\\' || ch == '(' || ch == ')'
}

// Scanner represents a lexical scanner.
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (Token, string) {
	// Read the next rune.
	ch := s.read()

	var buf bytes.Buffer
	buf.WriteRune(ch)

	if isComment(ch) {
		return s.ScanContiguous(COMMENT, isNotNewline, buf)
	} else if isNewline(ch) {
		return s.ScanContiguous(NEWLINE, isNewline, buf)
	} else if isTab(ch) {
		return TAB, ""
	} else if isQuote(ch) {
		return s.ScanQuoted(STRING, buf)
	} else if isLetter(ch) {
		return s.ScanContiguous(STRING, isLetter, buf)
	} else if isSpace(ch) {
		return s.ScanContiguous(SPACE, isSpace, buf)
	} else if ch == eof {
		return EOF, ""
	}

	return ILLEGAL, string(ch)
}

func (s *Scanner) ScanContiguous(tok Token, checker tokenTypeChecker, acc bytes.Buffer) (Token, string) {
	ch := s.read()
	if checker(ch) {
		acc.WriteRune(ch)
		return s.ScanContiguous(tok, checker, acc)
	}
	s.unread()
	return tok, acc.String()
}

func (s *Scanner) ScanQuoted(tok Token, acc bytes.Buffer) (Token, string) {
	ch := s.read()
	if isNewline(ch) {
		return ILLEGAL, acc.String()
	}
	acc.WriteRune(ch)
	if isQuote(ch) {
		return tok, acc.String()
	}
	return s.ScanQuoted(tok, acc)
}
