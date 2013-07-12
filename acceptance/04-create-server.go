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
	serverName = flag.String("n", randomString(16), "Server name (what you see in the control panel)")
	imageRef = flag.String("i", "", "ID of image to deploy onto the server")
	flavorRef = flag.String("f", "", "Flavor of server to deploy image upon")

	flag.Parse()
}

func aSuitableImage(api gophercloud.CloudServersProvider) string {
	images, err := api.ListImages()
	if err != nil {
		panic(err)
	}

	// TODO(sfalvo):
	// Works for Rackspace, might not work for your provider!
	// Need to figure out why ListImages() provides 0 values for
	// Ram and Disk fields.
	//
	// Until then, just return Ubuntu 12.04 LTS.
	for i := 0; i < len(images); i++ {
		if images[i].Id == "6a668bb8-fb5d-407a-9a89-6f957bced767" {
			return images[i].Id
		}
	}
	panic("Image 6a668bb8-fb5d-407a-9a89-6f957bced767 (Ubuntu 12.04 LTS) not found.")
}

func aSuitableFlavor(api gophercloud.CloudServersProvider) string {
	flavors, err := api.ListFlavors()
	if err != nil {
		panic(err)
	}

	// TODO(sfalvo):
	// Works for Rackspace, might not work for your provider!
	// Need to figure out why ListFlavors() provides 0 values for
	// Ram and Disk fields.
	//
	// Until then, just return Ubuntu 12.04 LTS.
	for i := 0; i < len(flavors); i++ {
		if flavors[i].Id == "2" {
			return flavors[i].Id
		}
	}
	panic("Flavor 2 (512MB 1-core 20GB machine) not found.")
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

	if *imageRef == "" {
		*imageRef = aSuitableImage(servers)
	}

	if *flavorRef == "" {
		*flavorRef = aSuitableFlavor(servers)
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

	if !*quiet {
		fmt.Printf("ID,Name,Status,Progress\n")
		for _, i := range allServers {
			fmt.Printf("%s,\"%s\",%s,%d\n", i.Id, i.Name, i.Status, i.Progress)
		}
	}
}
