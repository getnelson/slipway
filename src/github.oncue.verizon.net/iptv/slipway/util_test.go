package main

import (
	"testing"
)

func TestUnitNameExtractor(t *testing.T) {
	e1 := "aloha"
	e2 := "1.0.9"
	provided := "docker.oncue.verizon.net/units/aloha-1.0:1.0.9"
	last, tag := getUnitNameFromDockerContainer(provided)

	if tag != e2 {
		t.Error(tag, e2)
	}
	if last != e1 {
		t.Error(last, e1)
	}
}
