package testhelper

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func prefix() string {
	_, file, line, _ := runtime.Caller(3)
	return fmt.Sprintf("Failure in %s, line %d:", filepath.Base(file), line)
}

func green(str interface{}) string {
	return fmt.Sprintf("\033[0m\033[1;32m%#v\033[0m\033[1;31m", str)
}

func yellow(str interface{}) string {
	return fmt.Sprintf("\033[0m\033[1;33m%#v\033[0m\033[1;31m", str)
}

func logFatal(t *testing.T, str string) {
	t.Fatalf("\033[1;31m%s %s\033[0m", prefix(), str)
}

func logError(t *testing.T, str string) {
	t.Errorf("\033[1;31m%s %s\033[0m", prefix(), str)
}

// AssertEquals compares two arbitrary values and performs a comparison. If the
// comparison fails, a fatal error is raised that will fail the test
func AssertEquals(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		logFatal(t, fmt.Sprintf("expected %s but got %s", green(expected), yellow(actual)))
	}
}

// CheckEquals is similar to AssertEquals, except with a non-fatal error
func CheckEquals(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		logError(t, fmt.Sprintf("expected %s but got %s", green(expected), yellow(actual)))
	}
}

// AssertDeepEquals - like Equals - performs a comparison - but on more complex
// structures that requires deeper inspection
func AssertDeepEquals(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		logFatal(t, fmt.Sprintf("expected %s but got %s", green(expected), yellow(actual)))
	}
}

// CheckDeepEquals is similar to AssertDeepEquals, except with a non-fatal error
func CheckDeepEquals(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		logError(t, fmt.Sprintf("expected %s but got %s", green(expected), yellow(actual)))
	}
}

// AssertNoErr is a convenience function for checking whether an error value is
// an actual error
func AssertNoErr(t *testing.T, e error) {
	if e != nil {
		logFatal(t, fmt.Sprintf("unexpected error %s", yellow(e.Error())))
	}
}

// CheckNoErr is similar to AssertNoErr, except with a non-fatal error
func CheckNoErr(t *testing.T, e error) {
	if e != nil {
		logError(t, fmt.Sprintf("unexpected error %s", yellow(e.Error())))
	}
}
