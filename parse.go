/*
Initial approach based largely on https://blog.gopheracademy.com/advent-2014/parsers-lexers/
*/
package main

import (
	"flag"
	"io"
	"fmt"
	"os"
)


var filepath = flag.String("filepath", "haproxy.cfg", "path to haproxy configuration file")

type Section struct {
	Heading string
	Settings map[string]string
	Flags []string
}


func Parse(r io.Reader) ([]Section, error) {
	sc := NewScanner(r)

	for {
		tok, lit := sc.Scan()
		fmt.Printf("%v: %v\n", tok, lit)
		if tok == EOF || tok == ILLEGAL {
			break
		}		
	}

	// translate..
	sections := make([]Section, 0)
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
	fmt.Printf("%v\n", parsed)
}
