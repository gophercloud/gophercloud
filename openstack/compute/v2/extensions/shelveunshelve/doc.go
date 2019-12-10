/*
Package startstop provides functionality to start and stop servers that have
been provisioned by the OpenStack Compute service.

Example to Shelve, Shelve-offload and Unshelve a Server

	serverID := "47b6b7b7-568d-40e4-868c-d5c41735532e"

	err := startstop.Shelve(computeClient, serverID).ExtractErr()
	if err != nil {
		panic(err)
	}

	err := startstop.ShelveOffload(computeClient, serverID).ExtractErr()
	if err != nil {
		panic(err)
	}

	err := startstop.Unshelve(computeClient, serverID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package shelveunshelve
