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
	HEADING
	KEY
	VAL
	QUOTED
)

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


func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
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
	}
	
	return ILLEGAL, string(ch)
}

func (s *Scanner) ScanContiguous(tok Token, checker tokenTypeChecker, acc bytes.Buffer) (Token, string) {
	next := s.Read()
	if checker(next) {
		acc.WriteRune(next)
		return s.ScanContiguous(tok, checker, acc)
	}
	s.unread()
	return tok, acc.String()
}



	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	/*
	if isSpace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		s.unread()
		return s.scanIdent()
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return EOF, ""
	case '*':
		return ASTERISK, string(ch)
	case ',':
		return COMMA, string(ch)

	}
*/
