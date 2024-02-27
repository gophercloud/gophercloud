/*
Package endpoints provides information and interaction with the service
OS-EP-FILTER/endpoints API resource in the OpenStack Identity service.

For more information, see:
https://docs.openstack.org/api-ref/identity/v3-ext/#list-associations-by-project

Example to List Project Endpoints

	projectD := "e629d6e599d9489fb3ae5d9cc12eaea3"

	allPages, err := projectendpoints.List(identityClient, projectID).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allEndpoints, err := projectendpoints.ExtractEndpoints(allPages)
	if err != nil {
		panic(err)
	}

	for _, endpoint := range allEndpoints {
		fmt.Printf("%+v\n", endpoint)
	}
*/
package projectendpoints
