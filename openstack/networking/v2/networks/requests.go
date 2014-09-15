package networks

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
)

// User-defined options sent to the API when creating or updating a network.
type NetworkOpts struct {
	// The administrative state of the network, which is up (true) or down (false).
	AdminStateUp bool `json:"admin_state_up"`
	// The network name (optional)
	Name string `json:"name"`
	// Indicates whether this network is shared across all tenants. By default,
	// only administrative users can change this value.
	Shared bool `json:"shared"`
	// Admin-only. The UUID of the tenant that will own the network. This tenant
	// can be different from the tenant that makes the create network request.
	// However, only administrative users can specify a tenant ID other than their
	// own. You cannot change this value through authorization policies.
	TenantID string `json:"tenant_id"`
}

func APIVersions(c *gophercloud.ServiceClient) (*APIVersionsList, error) {
	var resp APIVersionsList
	_, err := perigee.Request("GET", APIVersionsURL(c), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		Results:     &resp,
		OkCodes:     []int{200},
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func APIInfo(c *gophercloud.ServiceClient, v string) (*APIInfoList, error) {
	var resp APIInfoList
	_, err := perigee.Request("GET", APIInfoURL(c, v), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		Results:     &resp,
		OkCodes:     []int{200},
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func GetExtension(c *gophercloud.ServiceClient, name string) (*Extension, error) {
	var ext Extension
	_, err := perigee.Request("GET", ExtensionURL(c, name), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		Results: &struct {
			Extension *Extension `json:"extension"`
		}{&ext},
		OkCodes: []int{200},
	})

	if err != nil {
		return nil, err
	}
	return &ext, nil
}
