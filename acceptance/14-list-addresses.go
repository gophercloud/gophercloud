package main

import (
	"fmt"
	"flag"
	"github.com/rackspace/gophercloud"
)

var quiet = flag.Bool("quiet", false, "Quiet mode, for acceptance testing.  $? still indicates errors though.")

func main() {
	flag.Parse()
	withIdentity(false, func(acc gophercloud.AccessProvider) {
		withServerApi(acc, func(servers gophercloud.CloudServersProvider) {
			log("Creating server")
			id, err := createServer(servers, "", "", "", "")
			if err != nil {
				panic(err)
			}
			waitForServerState(servers, id, "ACTIVE")
			defer servers.DeleteServerById(id)

			log("Getting list of addresses...")
			addresses, err := servers.ListAddresses(id)
			if (err != nil) && (err != gophercloud.WarnUnauthoritative) {
				panic(err)
			}
			if err == gophercloud.WarnUnauthoritative {
				log("Uh oh -- got a response back, but it's not authoritative for some reason.")
			}
			for _, addr := range addresses.Public {
				log("Address:", addr.Addr, "  IPv", addr.Version)
			}

			log("Done")
		})
	})
}

func log(s... interface{}) {
	if !*quiet {
		fmt.Println(s...)
	}
}
