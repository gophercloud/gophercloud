/*
Package injectnetworkinfo provides functionality to inject the network info into
a server that has been provisioned by the OpenStack Compute service. This action
requires admin privileges and Nova configured with a Xen hypervisor driver.

Example to Inject a Network Info into a Server

	serverID := "47b6b7b7-568d-40e4-868c-d5c41735532e"
	err := injectnetworkinfo.InjectNetworkInfo(client, id).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package injectnetworkinfo
