/*
package portforwarding enables management and retrieval of port forwarding resources for Floating IPs from the
OpenStack Networking service.

Example to Create a Port Forwarding

	createOpts := floatingips.CreateOpts{
		FloatingNetworkID: "a6917946-38ab-4ffd-a55a-26c0980ce5ee",
	}

	fip, err := floatingips.Create(networkingClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}
*/
package portforwarding
