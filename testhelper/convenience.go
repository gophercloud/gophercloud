package testhelper

import (
	"reflect"
	"testing"
)

// This function compares two arbitrary values and performs a comparison. If the
// comparison fails, a fatal error is raised that will fail the test
func Equals(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Fatalf("Expected [%#v] but got [%#v]", expected, actual)
	}
}

// This function, like Equals, performs a comparison - but on more complex
// structures that requires deeper inspection
func DeepEquals(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %#v but got %#v", expected, actual)
	}
}

// A convenience function for checking whether an error value is an actual error
func CheckErr(t *testing.T, e error) {
	if e != nil {
		t.Fatalf("Unexpected error: %#v", e)
	}
}
