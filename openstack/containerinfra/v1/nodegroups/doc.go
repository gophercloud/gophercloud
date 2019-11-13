/*
Package nodegroups provides methods for interacting with the Magnum node group API.

All node group actions must be performed on a specific cluster,
so the cluster UUID/name is required as a parameter in each method.


Create a client to use:

    opts, err := openstack.AuthOptionsFromEnv()
    if err != nil {
        panic(err)
    }

    provider, err := openstack.AuthenticatedClient(opts)
    if err != nil {
        panic(err)
    }

    client, err := openstack.NewContainerInfraV1(provider, gophercloud.EndpointOpts{Region: os.Getenv("OS_REGION_NAME")})
    if err != nil {
        panic(err)
    }

    client.Microversion = "1.9"


Example of Getting a node group:

    ng, err := nodegroups.Get(client, clusterUUID, nodeGroupUUID).Extract()
    if err != nil {
    	panic(err)
    }
    fmt.Printf("%#v\n", ng)


Example of Listing node groups:

    listOpts := nodegroup.ListOpts{
    	Role: "worker",
    }

    allPages, err := nodegroups.List(client, clusterUUID, listOpts).AllPages()
    if err != nil {
    	panic(err)
    }

    ngs, err := nodegroups.ExtractNodeGroups(allPages)
    if err != nil {
    	panic(err)
    }

    for _, ng := range ngs {
    	fmt.Printf("%#v\n", ng)
    }
*/
package nodegroups
