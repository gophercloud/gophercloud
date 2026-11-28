/*
Package remoteconsoles provides the ability to create server remote consoles
through the Compute API.
You need to specify at least "2.6" microversion for the ComputeClient to use
that API.

Example of Creating a new RemoteConsole

	computeClient, err := openstack.NewComputeV2(context.TODO(), providerClient, endpointOptions)
	computeClient.Microversion = "2.6"

	createOpts := remoteconsoles.CreateOpts{
	  Protocol: remoteconsoles.ConsoleProtocolVNC,
	  Type:     remoteconsoles.ConsoleTypeNoVNC,
	}
	serverID := "b16ba811-199d-4ffd-8839-ba96c1185a67"

	remtoteConsole, err := remoteconsoles.Create(context.TODO(), computeClient, serverID, createOpts).Extract()
	if err != nil {
	  panic(err)
	}

	fmt.Printf("Console URL: %s\n", remtoteConsole.URL)
*/

/*
Package remoteconsoles provides the ability to get console connection information
associated with a console authentication token for a server.

Since 'spice-direct' consoles were added in microversion 2.99, the client must
set Microversion >= 2.99 to avoid a Bad Request when calling this API.

# Example of getting console connection information

// Compute client
client, err := openstack.NewComputeV2(ctx, providerClient, eo)

	if err != nil {
		panic(err)
	}

// Remote console requires microversion >= 2.99
client.Microversion = "2.99"

console := "c58656f2-6657-4a67-ad5f-edc4aa5e04a2"
auth, err := remoteconsoles.Get(ctx, client, console).Extract()

	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"Console details: instance=%s host=%s port=%d tls_port=%d internal_access_path=%s\n",
		auth.InstanceUUID, auth.Host, auth.Port, auth.TLSPort auth.InternalAccessPath,
	)
*/
package remoteconsoles
