/*
Package defaultrules provides information and interaction with Security Group
Default Rules for the OpenStack Networking service. These rules are templates
used by Neutron to populate the rules of newly created security groups; they
require the "security-groups-default-rules" API extension.

Example to List Default Security Group Rules

	listOpts := defaultrules.ListOpts{
		Protocol: "tcp",
	}

	allPages, err := defaultrules.List(networkClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allDefaultRules, err := defaultrules.ExtractDefaultRules(allPages)
	if err != nil {
		panic(err)
	}

	for _, rule := range allDefaultRules {
		fmt.Printf("%+v\n", rule)
	}

Example to Create a Default Security Group Rule

	createOpts := defaultrules.CreateOpts{
		Direction:     rules.DirIngress,
		EtherType:     rules.EtherType4,
		Protocol:      rules.ProtocolTCP,
		PortRangeMin:  gophercloud.IntToPointer(443),
		PortRangeMax:  gophercloud.IntToPointer(443),
		RemoteGroupID: "PARENT",
	}

	rule, err := defaultrules.Create(context.TODO(), networkClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Default Security Group Rule

	ruleID := "37d94f8a-d136-465c-ae46-144f0d8ef141"
	err := defaultrules.Delete(context.TODO(), networkClient, ruleID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package defaultrules
