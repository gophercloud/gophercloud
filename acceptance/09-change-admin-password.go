package main

import (
	"flag"
	"fmt"
	"github.com/rackspace/gophercloud"
)

var quiet = flag.Bool("quiet", false, "Quiet mode, for acceptance testing.  $? still indicates errors though.")
var serverId = flag.String("i", "", "ID of server whose admin password is to be changed.")
var newPass = flag.String("p", "", "New password for the server.")

func main() {
	provider, username, password := getCredentials()
	flag.Parse()

	if *serverId == "" {
		panic("Server ID expected [use -i option]")
	}

	if *newPass == "" {
		panic("Password expected [use -p option]")
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

	api, err := gophercloud.ServersApi(acc, gophercloud.ApiCriteria{
		Name:      "cloudServersOpenStack",
		Region:    "DFW",
		VersionId: "2",
		UrlChoice: gophercloud.PublicURL,
	})
	if err != nil {
		panic(err)
	}

	err = api.SetAdminPassword(*serverId, *newPass)
	if err != nil {
		panic(err)
	}

	if !*quiet {
		fmt.Println("Password change request submitted.")
	}
}
