package testhelper

import "testing"

// This function compares two arbitrary values and performs a comparison. If the
// comparison fails, a fatal error is raised that will fail the test
func Compare(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Fatalf("Expected [%#v] but got [%#v]", expected, actual)
	}
}

// A convenience function for checking whether an error value is an actual error
func CheckErr(t *testing.T, e error) {
	if e != nil {
		t.Fatalf("Unexpected error: %#v", e)
	}
}
