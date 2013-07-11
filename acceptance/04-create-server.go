package main

import (
	"flag"
	"fmt"
	"log"
	"github.com/rackspace/gophercloud"
)

var region, serverName, imageRef, flavorRef *string
var adminPass = flag.String("a", "", "Administrator password (auto-assigned if none)")

func main() {
	provider, username, password := getCredentials()
	region = flag.String("r", "DFW", "Rackspace region in which to create the server")
	serverName = flag.String("n", randomString(16), "Server name (what you see in the control panel)")
	imageRef = flag.String("i", "", "ID of image to deploy onto the server")  // TODO(sfalvo): Make this work in -quiet mode.
	flavorRef = flag.String("f", "", "Flavor of server to deploy image upon") // TODO(sfalvo): Make this work in -quiet mode.

	flag.Parse()

	validations := map[string]string{
		"an image reference (-i flag)": *imageRef,
		"a server flavor (-f flag)":    *flavorRef,
	}
	for flag, value := range validations {
		if value == "" {
			log.Fatal(fmt.Sprintf("You must provide %s", flag))
		}
	}

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

	_, err = servers.CreateServer(gophercloud.NewServer{
		Name:      *serverName,
		ImageRef:  *imageRef,
		FlavorRef: *flavorRef,
		AdminPass: *adminPass,
	})
	if err != nil {
		panic(err)
	}

	allServers, err := servers.ListServers()
	if err != nil {
		panic(err)
	}

	fmt.Printf("ID,Name,Status,Progress\n")
	for _, i := range allServers {
		fmt.Printf("%s,\"%s\",%s,%d\n", i.Id, i.Name, i.Status, i.Progress)
	}
}
