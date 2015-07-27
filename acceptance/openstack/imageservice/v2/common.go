// +build acceptance imageservice

package v2

import (
	"os"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	th "github.com/rackspace/gophercloud/testhelper"
)

func newClient(t *testing.T) *gophercloud.ServiceClient {

	authURL := os.Getenv("OS_AUTH_URL")
	username := os.Getenv("OS_USERNAME")
	password := os.Getenv("OS_PASSWORD")
	tenantName := os.Getenv("OS_TENANT_NAME")
	domainName := os.Getenv("OS_DOMAIN_NAME")

	t.Logf("Credentials used: OS_AUTH_URL='%s' OS_USERNAME='%s' OS_PASSWORD='*****' OS_TENANT_NAME='%s' OS_TENANT_NAME='%s' \n",
		authURL, username, tenantName, domainName)

	client, err := openstack.NewClient(authURL)
	th.AssertNoErr(t, err)

	ao := gophercloud.AuthOptions{
		Username:   username,
		Password:   password,
		TenantName: tenantName,
		DomainName: domainName,
	}

	err = openstack.AuthenticateV3(client, ao)
	th.AssertNoErr(t, err)
	t.Logf("Token is %v", client.TokenID)

	c, err := openstack.NewImageServiceV2(client, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	th.AssertNoErr(t, err)
	return c
}
