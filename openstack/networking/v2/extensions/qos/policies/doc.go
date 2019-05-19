/*
Package policies provides information and interaction with the QoS policy extension
for the OpenStack Networking service.

Example to Get a Port with a QoS policy

	var portWithQoS struct {
        ports.Port
        policies.QoSPolicyExt
    }

    portID := "46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2"

    err = ports.Get(client, portID).ExtractInto(&portWithQoS)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Port: %+v\n", portWithQoS)

Example to Create a Port with a QoS policy

	var portWithQoS struct {
		ports.Port
		policies.QoSPolicyExt
	}

	policyID := "d6ae28ce-fcb5-4180-aa62-d260a27e09ae"
	networkID := "7069db8d-e817-4b39-a654-d2dd76e73d36"

	portCreateOpts := ports.CreateOpts{
		NetworkID: networkID,
	}

	createOpts := policies.PortCreateOptsExt{
		CreateOptsBuilder: portCreateOpts,
		QoSPolicyID:       policyID,
	}

	err = ports.Create(client, createOpts).ExtractInto(&portWithQoS)
	if err != nil {
		panic(err)
	}

	fmt.Printf("port: %+v\n", portWithQoS)

Example to add a QoS policy to an existing Port

    var portWithQoS struct {
        ports.Port
        policies.QoSPolicyExt
    }

    portUpdateOpts := ports.UpdateOpts{}

    policyID := "d6ae28ce-fcb5-4180-aa62-d260a27e09ae"

    updateOpts := policies.PortUpdateOptsExt{
        UpdateOptsBuilder: portUpdateOpts,
        QoSPolicyID: &policyID,
    }

    err := ports.Update(client, "65c0ee9f-d634-4522-8954-51021b570b0d", updateOpts).ExtractInto(&portWithQoS)
    if err != nil {
        panic(err)
    }

    fmt.Printf("port: %+v\n", portWithQoS)

Example to delete a QoS policy from the existing Port

    var portWithQoS struct {
        ports.Port
        policies.QoSPolicyExt
    }

    portUpdateOpts := ports.UpdateOpts{}

    policyID := ""

    updateOpts := policies.PortUpdateOptsExt{
        UpdateOptsBuilder: portUpdateOpts,
        QoSPolicyID: &policyID,
    }

    err := ports.Update(client, "65c0ee9f-d634-4522-8954-51021b570b0d", updateOpts).ExtractInto(&portWithQoS)
    if err != nil {
        panic(err)
    }

    fmt.Printf("port: %+v\n", portWithQoS)

 */
package policies
