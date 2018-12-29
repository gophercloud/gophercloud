/*
Package addressscopes provides the ability to retrieve and manage Address scopes through the Neutron API.

Example of Listing Address scopes

	listOpts := addressscopes.ListOpts{
		IPVersion: 6,
	}

	allPages, err := addressscopes.List(networkClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allAddressScopes, err := addressscopes.ExtractAddressScopes(allPages)
	if err != nil {
		panic(err)
	}

	for _, addressScope := range allAddressScopes {
		fmt.Printf("%+v\n", addressScope)
	}

Example to Get an Address scope

    addressScopeID = "9cc35860-522a-4d35-974d-51d4b011801e"
    addressScope, err := addressscopes.Get(networkClient, addressScopeID).Extract()
    if err != nil {
    	panic(err)
    }
*/
package addressscopes
