package main

import (
	"fmt"
	"github.com/rackspace/gophercloud"
	"os"
)

func main() {
	provider := os.Getenv("SDK_PROVIDER")
	username := os.Getenv("SDK_USERNAME")
	password := os.Getenv("SDK_PASSWORD")

	if (provider == "") || (username == "") || (password == "") {
		fmt.Fprintf(os.Stderr, "One or more of the following environment variables aren't set:\n")
		fmt.Fprintf(os.Stderr, "	SDK_PROVIDER=\"%s\"\n", provider)
		fmt.Fprintf(os.Stderr, "	SDK_USERNAME=\"%s\"\n", username)
		fmt.Fprintf(os.Stderr, "	SDK_PASSWORD=\"%s\"\n", password)
		os.Exit(1)
	}

	acc, err := gophercloud.Authenticate(
		provider,
		gophercloud.AuthOptions{
			Username: username,
			Password: password,
		},
	)
	if err != nil {
		panic(err)
	}

	api, err := gophercloud.ComputeApi(acc, gophercloud.ApiCriteria{
		Name:      "cloudServersOpenStack",
		Region:    "DFW",
		VersionId: "2",
		UrlChoice: gophercloud.PublicURL,
	})
	if err != nil {
		panic(err)
	}

	servers, err := api.ListServers()
	if err != nil {
		panic(err)
	}

	fmt.Println("Server ID")
	fmt.Println("----------------------------------------")
	for _, s := range servers {
		fmt.Printf("	%s\n", s.Id)
	}
}
