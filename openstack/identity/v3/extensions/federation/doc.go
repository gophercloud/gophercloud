/*
Package federation provides information and interaction with OS-FEDERATION API for the
Openstack Identity service.

Example to List Mappings

	allPages, err := federation.ListMappings(identityClient).AllPages()
	if err != nil {
		panic(err)
	}
	allMappings, err := federation.ExtractMappings(allPages)
	if err != nil {
		panic(err)
	}
*/
package federation
