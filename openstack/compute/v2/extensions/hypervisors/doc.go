/*
Package hypervisors returns details about the hypervisors in the OpenStack
cloud.

Example of Retrieving Details of All Hypervisors

	allPages, err := hypervisors.List(computeClient).AllPages()
	if err != nil {
		panic("Unable to retrieve hypervisors: %s", err)
	}

	allHypervisors, err := hypervisors.ExtractHypervisors(allPages)
	if err != nil {
		panic("Unable to extract hypervisors: %s", err)
	}

	for _, hypervisor := range allHypervisors {
		fmt.Println("%+v\n", hypervisor)
	}
*/
package hypervisors
