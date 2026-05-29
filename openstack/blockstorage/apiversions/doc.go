/*
Package apiversions provides information and interaction with the different
API versions for the OpenStack Block Storage service, code-named Cinder.

Example of Retrieving all API Versions

	allPages, err := apiversions.List(client).AllPages(context.TODO())
	if err != nil {
		panic("unable to get API versions: " + err.Error())
	}

	allVersions, err := apiversions.ExtractAPIVersions(allPages)
	if err != nil {
		panic("unable to extract API versions: " + err.Error())
	}

	for _, version := range allVersions {
		fmt.Printf("%+v\n", version)
	}

Example of Retrieving an API Version

	version, err := apiversions.Get(context.TODO(), client, "v3").Extract()
	if err != nil {
		panic("unable to get API version: " + err.Error())
	}

	fmt.Printf("%+v\n", version)
*/
package apiversions
