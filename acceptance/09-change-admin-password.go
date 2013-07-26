package main

import (
	"flag"
	"fmt"
	"github.com/rackspace/gophercloud"
	"time"
)

var quiet = flag.Bool("quiet", false, "Quiet mode, for acceptance testing.  $? still indicates errors though.")
var serverId = flag.String("i", "", "ID of server whose admin password is to be changed.")
var newPass = flag.String("p", "", "New password for the server.")

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

	// If user doesn't explicitly provide a server ID, create one dynamically.
	if *serverId == "" {
		var err error
		*serverId, err = createServer(api, "", "", "", "")
		if err != nil {
			panic(err)
		}

		// Wait for server to finish provisioning.
		for {
			s, err := api.ServerById(*serverId)
			if err != nil {
				panic(err)
			}
			if s.Status == "ACTIVE" {
				break
			}
			time.Sleep(10 * time.Second)
		}
	}

	// If no password is provided, create one dynamically.
	if *newPass == "" {
		*newPass = randomString("", 16)
	}

	err = api.SetAdminPassword(*serverId, *newPass)
	if err != nil {
		panic(err)
	}

	if !*quiet {
		fmt.Println("Password change request submitted.")
	}
}
