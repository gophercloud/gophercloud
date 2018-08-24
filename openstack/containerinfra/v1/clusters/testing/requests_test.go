package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clusters"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestCreateCluster(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleCreateClusterSuccessfully(t)

	masterCount := 1
	nodeCount := 1
	createTimeout := 30
	opts := clusters.CreateOpts{
		ClusterTemplateID: "0562d357-8641-4759-8fed-8173f02c9633",
		CreateTimeout:     &createTimeout,
		DiscoveryURL:      "",
		FlavorID:          "m1.small",
		Keypair:           "my_keypair",
		Labels:            map[string]string{},
		MasterCount:       &masterCount,
		MasterFlavorID:    "m1.small",
		Name:              "k8s",
		NodeCount:         &nodeCount,
	}

	sc := fake.ServiceClient()
	sc.Endpoint = sc.Endpoint + "v1/"
	res := clusters.Create(sc, opts)
	th.AssertNoErr(t, res.Err)

	requestID := res.Header.Get("X-OpenStack-Request-Id")
	th.AssertEquals(t, requestUUID, requestID)

	actual, err := res.Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, clusterUUID, actual)
}
