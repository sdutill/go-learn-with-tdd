package main

import (
	"fmt"
	"testing"
)

func TestRepeat(t *testing.T) {
	repeated := Repeat("a", 25)
	expected := "aaaaaaaaaaaaaaaaaaaaaaaaa"

	if repeated != expected {
		t.Errorf("expected %q but got %q", expected, repeated)
	}
}

// Laughs repeatedly
func ExampleRepeat() {
	fmt.Println(Repeat("he", 4))
	// Output: hehehehe
}

func BenchmarkRepeat(b *testing.B) {
	for b.Loop() {
		Repeat("a", 25)
	}
}
