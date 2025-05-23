package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/quotas"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestCreateQuota(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleCreateQuotaSuccessfully(t, fakeServer)

	opts := quotas.CreateOpts{
		ProjectID: "aa5436ab58144c768ca4e9d2e9f5c3b2",
		Resource:  "Cluster",
		HardLimit: 10,
	}

	sc := client.ServiceClient(fakeServer)
	sc.Endpoint = sc.Endpoint + "v1/"

	res := quotas.Create(context.TODO(), sc, opts)
	th.AssertNoErr(t, res.Err)

	requestID := res.Header.Get("X-OpenStack-Request-Id")
	th.AssertEquals(t, requestUUID, requestID)

	quota, err := res.Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, projectID, quota.ProjectID)
}
