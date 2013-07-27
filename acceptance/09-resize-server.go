package main

import (
	"fmt"
	"flag"
	"github.com/rackspace/gophercloud"
	"time"
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

	// These tests are going to take some time to complete.
	// So, we'll do two tests at the same time to help amortize test time.
	done := make(chan bool)
	go resizeRejectTest(api, done)
	go resizeAcceptTest(api, done)
	_ = <- done
	_ = <- done

	if !*quiet {
		fmt.Println("Done.")
	}
}

func resizeRejectTest(api gophercloud.CloudServersProvider, done chan bool) {
	withServer(api, func(id string) {
		newFlavorId := findAlternativeFlavor()
		err := api.ResizeServer(id, randomString("ACPTTEST", 24), newFlavorId, "")
		if err != nil {
			panic(err)
		}

		waitForVerifyResize(api, id)

		err = api.RevertResize(id)
		if err != nil {
			panic(err)
		}
	})
	done <- true
}

func resizeAcceptTest(api gophercloud.CloudServersProvider, done chan bool) {
	withServer(api, func(id string) {
		newFlavorId := findAlternativeFlavor()
		err := api.ResizeServer(id, randomString("ACPTTEST", 24), newFlavorId, "")
		if err != nil {
			panic(err)
		}

		waitForVerifyResize(api, id)

		err = api.ConfirmResize(id)
		if err != nil {
			panic(err)
		}
	})
	done <- true
}

func waitForVerifyResize(api gophercloud.CloudServersProvider, id string) {
	for {
		s, err := api.ServerById(id)
		if err != nil {
			panic(err)
		}
		if s.Status == "VERIFY_RESIZE" {
			break
		}
		time.Sleep(10 * time.Second)
	}
}

func withServer(api gophercloud.CloudServersProvider, f func(string)) {
	id, err := createServer(api, "", "", "", "")
	if err != nil {
		panic(err)
	}

	for {
		s, err := api.ServerById(id)
		if err != nil {
			panic(err)
		}
		if s.Status == "ACTIVE" {
			break
		}
		time.Sleep(10 * time.Second)
	}

	f(id)

	err = api.DeleteServerById(id)
	if err != nil {
		panic(err)
	}
}
