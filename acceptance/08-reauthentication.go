package main

import (
	"fmt"
	"flag"
	"github.com/rackspace/gophercloud"
)

var quiet = flag.Bool("quiet", false, "Quiet mode for acceptance testing.  $? non-zero on error though.")
var rgn = flag.String("r", "DFW", "Datacenter region to interrogate.")

func main() {
	provider, username, password := getCredentials()
	flag.Parse()

	// Authenticate initially against the service.
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

	// Cache our initial authentication token.
	token1 := auth.AuthToken()

	// Acquire access to the cloud servers API.
	servers, err := gophercloud.ServersApi(auth, gophercloud.ApiCriteria{
		Name:      "cloudServersOpenStack",
		Region:    *rgn,
		VersionId: "2",
		UrlChoice: gophercloud.PublicURL,
	})
	if err != nil {
		panic(err)
	}

	// Just to confirm everything works, we should be able to list images without error.
	_, err = servers.ListImages()
	if err != nil {
		panic(err)
	}

	// Revoke our current authentication token.
	auth.Revoke(auth.AuthToken())

	// Attempt to list images again.  This should _succeed_, because we enabled re-authentication.
	_, err = servers.ListImages()
	if err != nil {
		panic(err)
	}

	// However, our new authentication token should differ.
	token2 := auth.AuthToken()

	if !*quiet {
		fmt.Println("Old authentication token: ", token1)
		fmt.Println("New authentication token: ", token2)
	}
}
