package testhelper

import "testing"

func Compare(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Fatalf("Expected [%#v] but got [%#v]", expected, actual)
	}
}

func CheckErr(t *testing.T, e error) {
	if e != nil {
		t.Fatalf("Unexpected error: %#v", e)
	}
}
