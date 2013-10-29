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
		withServerApi(acc, func(api gophercloud.CloudServersProvider) {
			log("Creating server")
			id, err := createServer(api, "", "", "", "")
			if err != nil {
				panic(err)
			}
			waitForServerState(api, id, "ACTIVE")
			defer api.DeleteServerById(id)
			
			tryAllAddresses(id, api)
			tryAddressesByNetwork("public", id, api)
			tryAddressesByNetwork("private", id, api)

			log("Done")
		})
	})
}

func tryAllAddresses(id string, api gophercloud.CloudServersProvider) {
	log("Getting list of all addresses...")
	addresses, err := api.ListAddresses(id)
        if (err != nil) && (err != gophercloud.WarnUnauthoritative) {
        	panic(err)
        }
        if err == gophercloud.WarnUnauthoritative {
	        log("Uh oh -- got a response back, but it's not authoritative for some reason.")
        }
        for _, addr := range addresses.Public {
        	log("Address:", addr.Addr, "  IPv", addr.Version)
        }
}

func tryAddressesByNetwork(networkLabel string, id string, api gophercloud.CloudServersProvider) {
	log("Getting list of addresses on", networkLabel, "network...")
        network, err := api.ListAddressesByNetwork(id, networkLabel)
        if (err != nil) && (err != gophercloud.WarnUnauthoritative) {
	        panic(err) 
        }       
        if err == gophercloud.WarnUnauthoritative {
                log("Uh oh -- got a response back, but it's not authoritative for some reason.")
        }
	for _, addr := range network[networkLabel]{
		log("Address:", addr.Addr, "  IPv", addr.Version)
	}
}

func log(s... interface{}) {
	if !*quiet {
		fmt.Println(s...)
	}
}
