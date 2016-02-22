package main

import (
	"testing"
	"strings"
)

func MatchTokens(t *testing.T, expected []TokenString, sc *Scanner) {
	for _, expect := range expected {
		tok, lit := sc.Scan()
		if tok != expect.T {
			t.Errorf("expected token of type %v (%v), got %v (%v)", humanTokens[expect.T], expect.T, humanTokens[tok], tok)
		}
		if expect.S != "" && expect.S != lit {
			t.Errorf("expected lit of %v, got %v", expect.S, lit)
		}
		if tok == EOF || tok == ILLEGAL {
			break
		}

	}
}

func TestScanEOF(t *testing.T) {
	fin := strings.NewReader("")
	sc := NewScanner(fin)
	tokens := []TokenString{
		TokenString{EOF, ""},
	}
	MatchTokens(t, tokens, sc)
}

func TestScanComment(t * testing.T) {
	fin := strings.NewReader("# this is stupid\n# this is still stupid")
	sc := NewScanner(fin)
	tokens := []TokenString{
		TokenString{COMMENT, ""},
		TokenString{NEWLINE, ""},
		TokenString{COMMENT, ""},
		TokenString{EOF, ""},
	}
	MatchTokens(t, tokens, sc)
}

func TestScanSection(t * testing.T) {
	str := "# test\n"
	str += "global\n"
	str += "\tlog /dev/log local0 err\n"
	str += "\t# this is a test\n"
	str += "\tmaxconn 1000 # another comment\n"
	str += "\tdaemon\n"
	fin := strings.NewReader(str)
	sc := NewScanner(fin)
	tokens := []TokenString{
		TokenString{COMMENT, "# test"},
		TokenString{NEWLINE, ""},
		TokenString{STRING, "global"},
		TokenString{NEWLINE, ""},
		TokenString{TAB, ""},
		TokenString{STRING, "log"},
		TokenString{SPACE, ""},
		TokenString{STRING, "/dev/log"},
		TokenString{SPACE, ""},
		TokenString{STRING, "local0"},
		TokenString{SPACE, ""},
		TokenString{STRING, "err"},
		TokenString{NEWLINE, ""},
		TokenString{TAB, ""},
		TokenString{COMMENT, "# this is a test"},
		TokenString{NEWLINE, ""},
		TokenString{TAB, ""},
		TokenString{STRING, "maxconn"},
		TokenString{SPACE, ""},
		TokenString{STRING, "1000"},
		TokenString{SPACE, ""},
		TokenString{COMMENT, "# another comment"},
		TokenString{NEWLINE, ""},
		TokenString{TAB, ""},
		TokenString{STRING, "daemon"},
		TokenString{NEWLINE, ""},
		TokenString{EOF, ""},
	}
	MatchTokens(t, tokens, sc)
}
