package main

import (
	"strings"
	"testing"
)

func TestSimple(t *testing.T) {
	str := "# test\n"
	str += "global\n"
	str += "\tlog /dev/log local0 err\n"
	str += "\t# this is a test\n"
	str += "\tmaxconn 1000 # another comment\n"
	str += "\tdaemon\n"
	fin := strings.NewReader(str)

	sects, err := Parse(fin)
	if err != nil {
		t.Errorf("did not expect error: %v", err)
	}

	if len(sects) != 1 {
		t.Errorf("expected one section, got %v", len(sects))
	}

	sect := sects[0]
	if sect.Heading != "global" {
		t.Errorf("expected heading to be global, was %v", sect.Heading)
	}
}
