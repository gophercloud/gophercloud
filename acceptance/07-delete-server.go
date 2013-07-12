package main

import (
	"fmt"
	"flag"
	"github.com/rackspace/gophercloud"
)

var quiet = flag.Bool("quiet", false, "Quiet operation for acceptance tests.  $? non-zero if problem.")
var region = flag.String("r", "DFW", "Datacenter region")

func main() {
	provider, username, password := getCredentials()
	flag.Parse()

	auth, err := gophercloud.Authenticate(provider, gophercloud.AuthOptions{
		Username: username,
		Password: password,
	})
	if err != nil {
		panic(err)
	}

	servers, err := gophercloud.ServersApi(auth, gophercloud.ApiCriteria{
		Name:      "cloudServersOpenStack",
		Region:    *region,
		VersionId: "2",
		UrlChoice: gophercloud.PublicURL,
	})
	if err != nil {
		panic(err)
	}

	ss, err := servers.ListServers()
	if err != nil {
		panic(err)
	}

	n := 0
	for _, s := range ss {
		if len(s.Name) < 10 {
			continue
		}
		if s.Name[0:10] == "ACPTTEST--" {
			err := servers.DeleteServerById(s.Id)
			if err != nil {
				panic(err)
			}
			n++
		}
	}

	if !*quiet {
		fmt.Printf("%d servers removed.\n", n)
	}
}
