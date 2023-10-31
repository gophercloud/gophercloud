package tools

import (
	"encoding/json"
	"errors"
	"math/rand"
	"strings"
	"testing"
	"time"
)

// ErrTimeout is returned if WaitFor/WaitForTimeout take longer than their timeout duration.
var ErrTimeout = errors.New("Timed out")

// WaitFor uses WaitForTimeout to poll a predicate function once per second to
// wait for a certain state to arrive, with a default timeout of 300 seconds.
func WaitFor(predicate func() (bool, error)) error {
	return WaitForTimeout(predicate, 600*time.Second)
}

// WaitForTimeout polls a predicate function once per second to wait for a
// certain state to arrive, or until the given timeout is reached.
func WaitForTimeout(predicate func() (bool, error), timeout time.Duration) error {
	startTime := time.Now()
	for time.Since(startTime) < timeout {
		time.Sleep(2 * time.Second)

		satisfied, err := predicate()
		if err != nil {
			return err
		}
		if satisfied {
			return nil
		}
	}
	return ErrTimeout
}

// MakeNewPassword generates a new string that's guaranteed to be different than the given one.
func MakeNewPassword(oldPass string) string {
	randomPassword := RandomString("", 16)
	for randomPassword == oldPass {
		randomPassword = RandomString("", 16)
	}
	return randomPassword
}

// RandomString generates a string of given length, but random content.
// All content will be within the ASCII graphic character set.
func RandomString(prefix string, n int) string {
	charset := []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	return prefix + randomString(charset, n)
}

// RandomFunnyString returns a random string of the given length filled with
// funny Unicode code points.
func RandomFunnyString(length int) string {
	charset := []rune("012abc \n\t🤖👾👩🏾‍🚀+.,;:*`~|\"'/\\]êà²×c師☷")
	return randomString(charset, length)
}

// RandomFunnyStringNoSlash returns a random string of the given length filled with
// funny Unicode code points, but no forward slash.
func RandomFunnyStringNoSlash(length int) string {
	charset := []rune("012abc \n\t🤖👾👩🏾‍🚀+.,;:*`~|\"'\\]êà²×c師☷")
	return randomString(charset, length)
}

func randomString(charset []rune, length int) string {
	var s strings.Builder
	for i := 0; i < length; i++ {
		s.WriteRune(charset[rand.Intn(len(charset))])
	}
	return s.String()
}

// RandomInt will return a random integer between a specified range.
func RandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

// Elide returns the first bit of its input string with a suffix of "..." if it's longer than
// a comfortable 40 characters.
func Elide(value string) string {
	if len(value) > 40 {
		return value[0:37] + "..."
	}
	return value
}

// PrintResource returns a resource as a readable structure
func PrintResource(t *testing.T, resource interface{}) {
	b, _ := json.MarshalIndent(resource, "", "  ")
	t.Logf(string(b))
}
