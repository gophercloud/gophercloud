/*
Package floatingips provides the ability to manage floating ips through the
Nova API.

This API has been deprecated and will be removed from a future release of the
Nova API service.

For environements that support this extension, this package can be used
regardless of if either Neutron or nova-network is used as the cloud's network
service.

Example to List Floating IPs

	allPages, err := floatingips.List(computeClient).AllPages()
	if err != nil {
		panic("Unable to retrieve Floating IPs: %s", err)
	}

	allFloatingIPs, err := floatingips.ExtractFloatingIPs(allPages)
	if err != nil {
		panic("Unable to extract Floating IPs: %s", err)
	}

	for _, fip := range allFloatingIPs {
		fmt.Println("%+v\n", fip)
	}

Example to Create a Floating IP

	createOpts := floatingips.CreateOpts{
		Pool: "nova",
	}

	fip, err := floatingips.Create(computeClient, createOpts).Extract()
	if err != nil {
		panic("Unable to create Floating IP: %s", err)
	}

Example to Delete a Floating IP

	err := floatingips.Delete(computeClient, "floatingip-id").ExtractErr()
	if err != nil {
		panic("Unable to delete Floating IP: %s", err)
	}

Example to Associate a Floating IP With a Server

	associateOpts := floatingips.AssociateOpts{
		FloatingIP: "10.10.10.2",
	}

	err := floatingips.AssociateInstance(computeClient, "server-id", associateOpts).ExtractErr()
	if err != nil {
		panic("Unable to associate Floating IP: %s", err)
	}

Example to Disassociate a Floating IP From a Server

	disassociateOpts := floatingips.DisassociateOpts{
		FloatingIP: "10.10.10.2",
	}

	err := floatingips.DisassociateInstance(computeClient, "server-id", disassociateOpts).ExtractErr()
	if err != nil {
		panic("Unable to disassocaite Floating IP: %s", err)
	}
*/
package floatingips
