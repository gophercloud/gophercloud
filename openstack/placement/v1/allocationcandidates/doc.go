/*
Package allocationcandidates queries allocation candidates from the
OpenStack Placement service.

Allocation candidates API requests are available starting from version 1.10.

The response format changed in version 1.12: the allocations field in
allocation_requests changed from an array to a dictionary keyed by
resource provider UUID. Use the microversions.go types for 1.10-1.11.

Example to list allocation candidates

	listOpts := allocationcandidates.ListOpts{
		Resources: "VCPU:1,MEMORY_MB:1024",
		Required:  []string{"HW_CPU_X86_SSE"},
	}

	page, err := allocationcandidates.List(placementClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}
	allocationCandidates, err := allocationcandidates.ExtractAllocationCandidates(page)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", allocationCandidates)

Example to list allocation candidates with granular resource groups (microversion >= 1.33)

	listOpts := allocationcandidates.ListOpts{
		Resources:   "VCPU:1,MEMORY_MB:1024",
		GroupPolicy: "isolate",
		ResourceGroups: map[string]allocationcandidates.ResourceGroup{
			"1": {
				Resources: "SRIOV_NET_VF:1",
				Required:  []string{"CUSTOM_PHYSNET1"},
			},
			"_NET": {
				Resources: "NET_BW_EGR_KILOBIT_PER_SEC:10",
				MemberOf:  "in:42896e0d-205d-4fe3-bd1e-100924931787,5e08ea53-c4c6-448e-9334-ac4953de3cfa",
			},
		},
	}

	placementClient.Microversion = "1.33"
	page, err := allocationcandidates.List(placementClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}
	allocationCandidates, err := allocationcandidates.ExtractAllocationCandidates(page)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", allocationCandidates)
*/
package allocationcandidates
