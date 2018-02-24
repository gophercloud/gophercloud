/*
Package ipsecpolicies allows management and retrieval of IPSec Policies in the
OpenStack Networking Service.

Example to Create a Policy

	createOpts := ipsecpolicies.CreateOpts{
		Name:        "IPSecPolicy_1",
	}

	policy, err := policies.Create(networkClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Policy

	err := ipsecpolicies.Delete(client, "5291b189-fd84-46e5-84bd-78f40c05d69c").ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package ipsecpolicies
