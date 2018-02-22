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
*/
package ipsecpolicies
