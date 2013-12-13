package main

import (
	"github.com/maxlinc/gophercloud"
)

func main() {
	provider, username, password, _ := getCredentials()

	_, err := gophercloud.Authenticate(
		provider,
		gophercloud.AuthOptions{
			Username: username,
			Password: password,
		},
	)
	if err != nil {
		panic(err)
	}
}
