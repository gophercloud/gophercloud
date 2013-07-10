package main

import (
	"fmt"
	"github.com/rackspace/gophercloud"
	"flag"
)

var quiet = flag.Bool("quiet", false, "Quiet mode, for acceptance testing.  $? still indicates errors though.")

func main() {
	provider, username, password := getCredentials()
	flag.Parse()

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

	api, err := gophercloud.ServersApi(acc, gophercloud.ApiCriteria{
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

	if !*quiet {
		for _, s := range servers {
			fmt.Printf("%s\n", s.Id)
		}
	}
}
