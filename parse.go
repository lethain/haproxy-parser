/*
Initial approach based largely on https://blog.gopheracademy.com/advent-2014/parsers-lexers/
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"errors"
)

var filepath = flag.String("filepath", "haproxy.cfg", "path to haproxy configuration file")
var EOFError = errors.New("end of file")


type Server struct {
	Name string
	IP string
	FullText string
}

func (s *Server) String() string {
	return fmt.Sprintf("Server(%v, %v)", s.Name, s.IP)
}

type Section struct {
	Heading  string
	FullText string
	Servers []Server
}

func (s *Section) String() string {
	return fmt.Sprintf("Section(%v, %v, %v bytes)", s.Heading, s.Servers, len(s.FullText))
}

type TokenString struct {
	T Token
	S string
}

func CollectUntil(tok Token, sc *Scanner, ts []TokenString) ([]TokenString, Token, string) {
	for {
		nextTok, lit := sc.Scan()
		if tok != nextTok {
			ts = append(ts, TokenString{nextTok, lit})
		}
		if tok == nextTok {
			return ts, nextTok, lit
		}
	}
}

// Variant on above which collects token strings from a slice
func CollectTokenStringsUntil(tok Token, ts []TokenString) ([]TokenString, []TokenString) {
	for i, t := range ts {
		if t.T == tok {
			return ts[:i], ts[i:]
		}
	}
	return ts, ts[len(ts):]
}

// Discard tokens until you reach token of desired type
func DiscardUntil(tok Token, sc *Scanner) (Token, string) {
	for {
		nextTok, lit := sc.Scan()
		if tok == nextTok {
			return nextTok, lit
		}
	}
}

func NewSection(ts []TokenString) (Section, error) {
	sect := new(Section)

	// collect full text
	for _, t := range ts {
		sect.FullText += t.S
	}
	
	// get heading
	headerTs, restTs := CollectTokenStringsUntil(NEWLINE, ts)
	for _, t := range headerTs {
		sect.Heading += t.S
	}
	restTs = restTs[1:]

	var lineTs []TokenString
	// also want to capture bind info
	// bind 127.0.0.1:5526
	for {
		if len(restTs) == 0 {
			break
		}
		lineTs, restTs = CollectTokenStringsUntil(NEWLINE, restTs)
		// skip past newline
		restTs = restTs[1:]
		
		first := lineTs[0]
		if first.T == STRING && first.S == "server" {
			serv := Server{lineTs[2].S, lineTs[4].S, ""}
			for _, t := range lineTs {
				serv.FullText += t.S
			}
			sect.Servers = append(sect.Servers, serv)
		}
	}
	return *sect, nil
}

func Parse(r io.Reader) ([]Section, error) {
	sc := NewScanner(r)
	sections := make([]Section, 0)

	tok, lit := sc.Scan()
	ts := make([]TokenString, 0)
	for {
		if tok == EOF {
			break
		} else if tok == ILLEGAL {
			return sections, fmt.Errorf("illegal token: %v", lit)
		}
		if tok != STRING {
			tok, lit = DiscardUntil(STRING, sc)
			continue
		}
		
		ts = append(ts, TokenString{tok, lit})
		ts, tok, lit = CollectUntil(NEWLINE, sc, ts)
		if tok == EOF || tok == ILLEGAL {
			return sections, fmt.Errorf("Unexpected token: %s, %s", humanTokens[tok], lit)
		}
		ts = append(ts, TokenString{tok, lit})
		for {
			tok, lit = sc.Scan()
			if tok != TAB {
				break
			}
			ts, tok, lit = CollectUntil(NEWLINE, sc, ts)
			ts = append(ts, TokenString{tok, lit})			
		}

		// transform collected tokens into a section
		section, err := NewSection(ts)
		if err != nil {
			panic(err)
		}
		sections = append(sections, section)
		ts = make([]TokenString, 0)
	}

	// translate..

	return sections, nil
}

func main() {
	flag.Parse()

	fin, err := os.Open(*filepath)
	if err != nil {
		panic(err)
	}
	parsed, err := Parse(fin)
	if err != nil {
		panic(err)
	}
	marshalled, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		panic(err)
	}	
	os.Stdout.Write(marshalled)
}
