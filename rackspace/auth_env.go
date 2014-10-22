package utils

import (
	"fmt"
	"os"

	"github.com/rackspace/gophercloud"
)

var nilOptions = gophercloud.AuthOptions{}

// ErrNoAuthUrl, ErrNoUsername, and ErrNoPassword errors indicate of the
// required RAX_AUTH_URL, RAX_USERNAME, or RAX_PASSWORD environment variables,
// respectively, remain undefined.  See the AuthOptions() function for more details.
var (
	ErrNoAuthURL  = fmt.Errorf("Environment variable RAX_AUTH_URL needs to be set.")
	ErrNoUsername = fmt.Errorf("Environment variable RAX_USERNAME needs to be set.")
	ErrNoPassword = fmt.Errorf("Environment variable RAX_API_KEY or RAX_PASSWORD needs to be set.")
)

// AuthOptionsFromEnv fills out an identity.AuthOptions structure with the
// settings found on the various Rackspace RAX_* environment variables.
func AuthOptionsFromEnv() (gophercloud.AuthOptions, error) {
	authURL := os.Getenv("RAX_AUTH_URL")
	username := os.Getenv("RAX_USERNAME")
	password := os.Getenv("RAX_PASSWORD")
	apiKey := os.Getenv("RAX_API_KEY")

	if authURL == "" {
		return nilOptions, ErrNoAuthURL
	}

	if username == "" {
		return nilOptions, ErrNoUsername
	}

	if password == "" && apiKey == "" {
		return nilOptions, ErrNoPassword
	}

	ao := gophercloud.AuthOptions{
		IdentityEndpoint: authURL,
		Username:         username,
		Password:         password,
		APIKey:           apiKey,
	}

	return ao, nil
}
