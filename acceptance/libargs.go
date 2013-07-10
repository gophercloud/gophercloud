package main

import (
	"fmt"
	"os"
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
