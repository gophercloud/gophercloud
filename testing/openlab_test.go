package testing

import "testing"

func TestPassed(t *testing.T) {
	_ = 0
	t.Log("This is a test which will pass")
}

func TestNotPassed(t *testing.T) {
	t.Fatal("This is a test which won't pass")
}
