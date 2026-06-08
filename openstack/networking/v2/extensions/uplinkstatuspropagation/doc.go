/*
Package uplinkstatuspropagation provides the ability to retrieve and manage ports
with the uplink status propagation extension through the Neutron API.

Example of Getting a Port with the uplink status propagation extension

	var port struct {
		ports.Port
		uplinkstatuspropagation.PortPropagateUplinkStatusExt
	}

	err := ports.Get(context.TODO(), networkClient, "46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2").ExtractInto(&port)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", port)

Example of Creating Port with the uplink status propagation extension

	iTrue := true
	asu := true
	createOpts := uplinkstatuspropagation.PortPropagateUplinkStatusCreateOptsExt{
		CreateOptsBuilder: ports.CreateOpts{
			Name:         "private-port",
			AdminStateUp: &asu,
			NetworkID:    "a87cc70a-3e15-4acf-8205-9b711a3531b7",
			FixedIPs: []ports.IP{
				{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.2"},
			},
			SecurityGroups: &[]string{"foo"},
		},
		PropagateUplinkStatus: &iTrue,
	}

	var port struct {
		ports.Port
		uplinkstatuspropagation.PortPropagateUplinkStatusExt
	}

	err := ports.Create(context.TODO(), networkClient, createOpts).ExtractInto(&port)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", port)

Example of Updating Port with the uplink status propagation extension

	iFalse := false
	name := "new_port_name"
	updateOpts := uplinkstatuspropagation.PortPropagateUplinkStatusUpdateOptsExt{
		UpdateOptsBuilder: ports.UpdateOpts{
			Name: &name,
			FixedIPs: []ports.IP{
				{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.3"},
			},
			SecurityGroups: &[]string{"f0ac4394-7e4a-4409-9701-ba8be283dbc3"},
		},
		PropagateUplinkStatus: &iFalse,
	}

	var port struct {
		ports.Port
		uplinkstatuspropagation.PortPropagateUplinkStatusExt
	}

	err := ports.Update(context.TODO(), networkClient, "65c0ee9f-d634-4522-8954-51021b570b0d", updateOpts).ExtractInto(&port)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", port)
*/
package uplinkstatuspropagation
