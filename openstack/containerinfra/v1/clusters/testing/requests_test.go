package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/clusters"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestCreateCluster(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleCreateClusterSuccessfully(t)

	masterCount := 1
	nodeCount := 1
	createTimeout := 30
	masterLBEnabled := true
	opts := clusters.CreateOpts{
		ClusterTemplateID: "0562d357-8641-4759-8fed-8173f02c9633",
		CreateTimeout:     &createTimeout,
		DiscoveryURL:      "",
		FlavorID:          "m1.small",
		Keypair:           "my_keypair",
		Labels:            map[string]string{},
		MasterCount:       &masterCount,
		MasterFlavorID:    "m1.small",
		MasterLBEnabled:   &masterLBEnabled,
		Name:              "k8s",
		NodeCount:         &nodeCount,
		FloatingIPEnabled: gophercloud.Enabled,
		FixedNetwork:      "private_network",
		FixedSubnet:       "private_subnet",
		MergeLabels:       gophercloud.Enabled,
	}

	sc := fake.ServiceClient()
	sc.Endpoint = sc.Endpoint + "v1/"
	res := clusters.Create(context.TODO(), sc, opts)
	th.AssertNoErr(t, res.Err)

	requestID := res.Header.Get("X-OpenStack-Request-Id")
	th.AssertEquals(t, requestUUID, requestID)

	actual, err := res.Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, clusterUUID, actual)
}

func TestGetCluster(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleGetClusterSuccessfully(t)

	sc := fake.ServiceClient()
	sc.Endpoint = sc.Endpoint + "v1/"
	actual, err := clusters.Get(context.TODO(), sc, "746e779a-751a-456b-a3e9-c883d734946f").Extract()
	th.AssertNoErr(t, err)
	actual.CreatedAt = actual.CreatedAt.UTC()
	actual.UpdatedAt = actual.UpdatedAt.UTC()
	th.AssertDeepEquals(t, ExpectedCluster, *actual)
}

func TestListClusters(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleListClusterSuccessfully(t)

	count := 0
	sc := fake.ServiceClient()
	sc.Endpoint = sc.Endpoint + "v1/"
	err := clusters.List(sc, clusters.ListOpts{Limit: 2}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := clusters.ExtractClusters(page)
		th.AssertNoErr(t, err)
		for idx := range actual {
			actual[idx].CreatedAt = actual[idx].CreatedAt.UTC()
			actual[idx].UpdatedAt = actual[idx].UpdatedAt.UTC()
		}
		th.AssertDeepEquals(t, ExpectedClusters, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestListDetailClusters(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleListDetailClusterSuccessfully(t)

	count := 0
	sc := fake.ServiceClient()
	sc.Endpoint = sc.Endpoint + "v1/"
	err := clusters.ListDetail(sc, clusters.ListOpts{Limit: 2}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := clusters.ExtractClusters(page)
		th.AssertNoErr(t, err)
		for idx := range actual {
			actual[idx].CreatedAt = actual[idx].CreatedAt.UTC()
			actual[idx].UpdatedAt = actual[idx].UpdatedAt.UTC()
		}
		th.AssertDeepEquals(t, ExpectedClusters, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestUpdateCluster(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleUpdateClusterSuccessfully(t)

	updateOpts := []clusters.UpdateOptsBuilder{
		clusters.UpdateOpts{
			Op:    clusters.ReplaceOp,
			Path:  "/master_lb_enabled",
			Value: "True",
		},
		clusters.UpdateOpts{
			Op:    clusters.ReplaceOp,
			Path:  "/registry_enabled",
			Value: "True",
		},
	}

	sc := fake.ServiceClient()
	sc.Endpoint = sc.Endpoint + "v1/"
	res := clusters.Update(context.TODO(), sc, clusterUUID, updateOpts)
	th.AssertNoErr(t, res.Err)

	requestID := res.Header.Get("X-OpenStack-Request-Id")
	th.AssertEquals(t, requestUUID, requestID)

	actual, err := res.Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, clusterUUID, actual)
}

func TestUpgradeCluster(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleUpgradeClusterSuccessfully(t)

	opts := clusters.UpgradeOpts{
		ClusterTemplate: "0562d357-8641-4759-8fed-8173f02c9633",
	}

	sc := fake.ServiceClient()
	sc.Endpoint = sc.Endpoint + "v1/"
	res := clusters.Upgrade(context.TODO(), sc, clusterUUID, opts)
	th.AssertNoErr(t, res.Err)

	requestID := res.Header.Get("X-OpenStack-Request-Id")
	th.AssertEquals(t, requestUUID, requestID)

	actual, err := res.Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, clusterUUID, actual)
}

func TestDeleteCluster(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleDeleteClusterSuccessfully(t)

	sc := fake.ServiceClient()
	sc.Endpoint = sc.Endpoint + "v1/"
	r := clusters.Delete(context.TODO(), sc, clusterUUID)
	err := r.ExtractErr()
	th.AssertNoErr(t, err)

	uuid := ""
	idKey := "X-Openstack-Request-Id"
	if len(r.Header[idKey]) > 0 {
		uuid = r.Header[idKey][0]
		if uuid == "" {
			t.Errorf("No value for header [%s]", idKey)
		}
	} else {
		t.Errorf("Missing header [%s]", idKey)
	}

	th.AssertEquals(t, requestUUID, uuid)
}

func TestResizeCluster(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleResizeClusterSuccessfully(t)

	nodeCount := 2

	opts := clusters.ResizeOpts{
		NodeCount: &nodeCount,
	}

	sc := fake.ServiceClient()
	sc.Endpoint = sc.Endpoint + "v1/"
	res := clusters.Resize(context.TODO(), sc, clusterUUID, opts)
	th.AssertNoErr(t, res.Err)

	requestID := res.Header.Get("X-OpenStack-Request-Id")
	th.AssertEquals(t, requestUUID, requestID)

	actual, err := res.Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, clusterUUID, actual)
}
