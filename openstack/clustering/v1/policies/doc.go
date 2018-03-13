/*
Package policies provides information and interaction with the policies through
the OpenStack Compute service.

Lists all policies and creates, shows information for, updates, and deletes a policy.

Example to List Policies

	listOpts := policies.ListOpts{
		Limit: 2,
	}

	allPages, err := policies.ListDetail(computeClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allPolicies, err := policies.ExtractImages(allPages)
	if err != nil {
		panic(err)
	}

	for _, policy := range allPolicies {
		fmt.Printf("%+v\n", policy)
	}

Example to Create a policy

	createOpts := receivers.CreateOpts{
    Name: "scaling_policy",
    Spec: map[string]interface{}{
      "type":       "senlin.policy.scaling",
      "version":    "1.0",
      "properties": propsPolicy,
    },
	}
  propsPolicy := &map[string]interface{}{
		"adjustment": adjustsPolicy,
		"event":      "CLUSTER_SCALE_IN",
	}
  adjustsPolicy := &map[string]interface{}{
		"min_step": 1,
		"number":   1,
		"type":     "CHANGE_IN_CAPACITY",
	}

	policy, err := policies.Create(computeClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}
*/
package policies
