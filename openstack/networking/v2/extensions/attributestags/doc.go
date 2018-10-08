/*
Package attributestags manages Tags on Resources created by the OpenStack Neutron Service.

This enables tagging via a standard interface for resources types which support it.

See https://developer.openstack.org/api-ref/network/v2/#standard-attributes-tag-extension for more information on the underlying API.

Example to ReplaceAll Resource Tags

    network, err := networks.Create(conn, createOpts).Extract()

    tagReplaceAllOpts := attributestags.ReplaceAllOpts{
        Tags:         []string{"abc", "123"},
    }
    attributestags.ReplaceAll(conn, "networks", network.ID, tagReplaceAllOpts)

Example to List all Resource Tags

	tags, err = attributestags.List(conn, "networks", network.ID).Extract()

Example to Delete all Resource Tags

	err = attributestags.DeleteAll(conn, "networks", network.ID).ExtractErr()
*/
package attributestags
