/*
Package extradhcpopts allow to work with extra DHCP functionality of Neutron ports.

Example to Get a Port with DHCP opts

	portID := "46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2"
	var s struct {
		ports.Port
		extradhcpopts.ExtraDHCPOptsExt
	}

	err := ports.Get(networkClient, portID).ExtractInto(&s)
	if err != nil {
		panic(err)
	}
*/
package extradhcpopts
