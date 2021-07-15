/*
Package qos provides information and interaction with the QoS specifications
for the Openstack Blockstorage service.

Example to create a QoS specification

	createOpts := qos.CreateOpts{
		Name:     "test",
		Consumer: qos.ConsumerFront,
		Specs: map[string]string{
			"read_iops_sec": "20000",
		},
	}

	test, err := qos.Create(client, createOpts).Extract()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("QoS: %+v\n", test)

Example to delete a QoS specification

	qosID := "d6ae28ce-fcb5-4180-aa62-d260a27e09ae"

	deleteOpts := qos.DeleteOpts{
		Force: false,
	}

	err = qos.Delete(client, qosID, deleteOpts).ExtractErr()
	if err != nil {
		log.Fatal(err)
	}

Example to list QoS specifications

	listOpts := qos.ListOpts{}

	allPages, err := qos.List(client, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allQoS, err := qos.ExtractQoS(allPages)
	if err != nil {
		panic(err)
	}

	for _, qos := range allQoS {
		fmt.Printf("List: %+v\n", qos)
	}

*/
package qos
