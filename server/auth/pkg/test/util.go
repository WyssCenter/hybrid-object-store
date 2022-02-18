package test

import (
	"testing"
)

// AssertEqual is a helper function to check 2 values
func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}
