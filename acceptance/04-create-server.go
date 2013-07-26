package main

import (
	"flag"
	"fmt"
	"github.com/rackspace/gophercloud"
)

var provider, username, password string

var region, serverName, imageRef, flavorRef *string
var adminPass = flag.String("a", "", "Administrator password (auto-assigned if none)")
var quiet = flag.Bool("quiet", false, "Quiet mode for acceptance tests.  $? non-zero if error.")

func configure() {
	provider, username, password = getCredentials()
	region = flag.String("r", "DFW", "Rackspace region in which to create the server")
	serverName = flag.String("n", randomString("ACPTTEST--", 16), "Server name (what you see in the control panel)")
	imageRef = flag.String("i", "", "ID of image to deploy onto the server")
	flavorRef = flag.String("f", "", "Flavor of server to deploy image upon")

	flag.Parse()
}

func main() {
	configure()

	auth, err := gophercloud.Authenticate(
		provider,
		gophercloud.AuthOptions{
			Username: username,
			Password: password,
		},
	)
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

	_, err = createServer(servers, *imageRef, *flavorRef, *serverName, *adminPass)
	if err != nil {
		panic(err)
	}

	allServers, err := servers.ListServers()
	if err != nil {
		panic(err)
	}

	if !*quiet {
		fmt.Printf("ID,Name,Status,Progress\n")
		for _, i := range allServers {
			fmt.Printf("%s,\"%s\",%s,%d\n", i.Id, i.Name, i.Status, i.Progress)
		}
	}
}
