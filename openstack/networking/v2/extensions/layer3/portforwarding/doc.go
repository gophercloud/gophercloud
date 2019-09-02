/*
package portforwarding enables management and retrieval of port forwarding resources for Floating IPs from the
OpenStack Networking service.

Example to Create a Port Forwarding for a floating IP

	createOpts := &portforwarding.CreateOpts{
		Protocol:          "tcp",
		InternalPort:      25,
		ExternalPort:      2230,
		InternalIPAddress: internalIP,
		InternalPortID:    portID,
	}

	pf, err := portforwarding.Create(networkingClient, floatingIPID, createOpts).Extract()
	if err != nil {
		panic(err)
	}
*/
package portforwarding
