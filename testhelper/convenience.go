package testhelper

import (
	"reflect"
	"testing"
)

// AssertEquals compares two arbitrary values and performs a comparison. If the
// comparison fails, a fatal error is raised that will fail the test
func AssertEquals(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Fatalf("Expected [%#v] but got [%#v]", expected, actual)
	}
}

// CheckEquals is similar to AssertEquals, except with a non-fatal error
func CheckEquals(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Errorf("Expected [%#v] but got [%#v]", expected, actual)
	}
}

// AssertDeepEquals - like Equals - performs a comparison - but on more complex
// structures that requires deeper inspection
func AssertDeepEquals(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %#v but got %#v", expected, actual)
	}
}

// CheckDeepEquals is similar to AssertDeepEquals, except with a non-fatal error
func CheckDeepEquals(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %#v but got %#v", expected, actual)
	}
}

// AssertNoErr is a convenience function for checking whether an error value is
// an actual error
func AssertNoErr(t *testing.T, e error) {
	if e != nil {
		t.Fatalf("Unexpected error: %#v", e)
	}
}

// CheckNoErr is similar to AssertNoErr, except with a non-fatal error
func CheckNoErr(t *testing.T, e error) {
	if e != nil {
		t.Errorf("Unexpected error: %#v", e)
	}
}
