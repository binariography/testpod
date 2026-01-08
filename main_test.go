package main

import (
	"testing"
)

func TestPass(t *testing.T) {
	a := "a"
	if a != "a" {
		t.Errorf("Expected a, got %s", a)
	}
}
