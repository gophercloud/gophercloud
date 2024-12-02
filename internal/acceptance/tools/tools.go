package tools

import (
	"context"
	"encoding/json"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
)

// WaitFor uses WaitForTimeout to poll a predicate function once per second to
// wait for a certain state to arrive, with a default timeout of 600 seconds.
func WaitFor(predicate func(context.Context) (bool, error)) error {
	return WaitForTimeout(predicate, 600*time.Second)
}

// WaitForTimeout polls a predicate function once per second to wait for a
// certain state to arrive, or until the given timeout is reached.
func WaitForTimeout(predicate func(context.Context) (bool, error), timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	return gophercloud.WaitFor(ctx, predicate)
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
func PrintResource(t *testing.T, resource any) {
	b, _ := json.MarshalIndent(resource, "", "  ")
	t.Log(string(b))
}
