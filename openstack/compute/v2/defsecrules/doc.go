/*
Package defsecrules enables management of default security group rules.

Default security group rules are rules that are managed in the "default"
security group.

This is only applicable in environments running nova-network. This package will
not work if the OpenStack environment is running the OpenStack Networking
(Neutron) service.

Example of Listing Default Security Group Rules

	allPages, err := defsecrules.List(computeClient).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allDefaultRules, err := defsecrules.ExtractDefaultRules(allPages)
	if err != nil {
		panic(err)
	}

	for _, df := range allDefaultRules {
		fmt.Printf("%+v\n", df)
	}

Example of Retrieving a Default Security Group Rule

	rule, err := defsecrules.Get(context.TODO(), computeClient, "rule-id").Extract()
	if err != nil {
		panic(err)
	}

Example of Creating a Default Security Group Rule

	createOpts := defsecrules.CreateOpts{
		IPProtocol: "TCP",
		FromPort:   80,
		ToPort:     80,
		CIDR:       "10.10.12.0/24",
	}

	rule, err := defsecrules.Create(context.TODO(), computeClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example of Deleting a Default Security Group Rule

	err := defsecrules.Delete(context.TODO(), computeClient, "rule-id").ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package defsecrules
