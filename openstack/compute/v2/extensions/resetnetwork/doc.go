/*
Package resetnetwork provides functionality to reset the network of a server
that has been provisioned by the OpenStack Compute service. This action
requires admin privileges and Nova configured with a Xen hypervisor driver.

Example to Reset a Network of a Server

	serverID := "47b6b7b7-568d-40e4-868c-d5c41735532e"
	err := resetnetwork.ResetNetwork(client, id).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package resetnetwork
