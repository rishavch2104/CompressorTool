package main

import (
	"os"
	"testing"
)

func TestCountCharacters(t *testing.T) {
	file, _ := os.Open("135-0.txt")
	got, err := countCharacters(file)
	if err != nil {
		t.Error("got error")
	}
	if got["X"] != 333 {
		t.Errorf("got %v, want 333", got["X"])
	}

}
