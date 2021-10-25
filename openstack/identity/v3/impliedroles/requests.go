package impliedroles

import (
	"github.com/gophercloud/gophercloud"
)

// Get the implied roles associated with the provided role
func GetImpliesRoles(client *gophercloud.ServiceClient, priorRoleId string) (r GetImpliedRoleResult) {
	url := getURL(client, priorRoleId)
	resp, err := client.Get(url, &r.Body, nil)

	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)

	return
}
