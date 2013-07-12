package main

import (
	"fmt"
	"os"
	"crypto/rand"
)

// getCredentials will verify existence of needed credential information
// provided through environment variables.  This function will not return
// if at least one piece of required information is missing.
func getCredentials() (provider, username, password string) {
	provider = os.Getenv("SDK_PROVIDER")
	username = os.Getenv("SDK_USERNAME")
	password = os.Getenv("SDK_PASSWORD")

	if (provider == "") || (username == "") || (password == "") {
		fmt.Fprintf(os.Stderr, "One or more of the following environment variables aren't set:\n")
		fmt.Fprintf(os.Stderr, "  SDK_PROVIDER=\"%s\"\n", provider)
		fmt.Fprintf(os.Stderr, "  SDK_USERNAME=\"%s\"\n", username)
		fmt.Fprintf(os.Stderr, "  SDK_PASSWORD=\"%s\"\n", password)
		os.Exit(1)
	}

	return
}

// randomString generates a string of given length, but random content.
// All content will be within the ASCII graphic character set.
// (Implementation from Even Shaw's contribution on
// http://stackoverflow.com/questions/12771930/what-is-the-fastest-way-to-generate-a-long-random-string-in-go).
func randomString(n int) string {
    const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    var bytes = make([]byte, n)
    rand.Read(bytes)
    for i, b := range bytes {
        bytes[i] = alphanum[b % byte(len(alphanum))]
    }
    return string(bytes)
}