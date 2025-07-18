package v1

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	idv3 "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/identity/v3"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/clusters"
	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/clustertemplates"
	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/quotas"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// CreateClusterTemplateCOE will create a random cluster template for the specified orchestration engine.
// An error will be returned if the cluster template could not be created.
func CreateClusterTemplateCOE(t *testing.T, client *gophercloud.ServiceClient, coe string) (*clustertemplates.ClusterTemplate, error) {
	choices, err := clients.AcceptanceTestChoicesFromEnv()
	if err != nil {
		return nil, err
	}

	name := tools.RandomString("TESTACC-", 8)
	t.Logf("Attempting to create %s cluster template: %s", coe, name)

	boolFalse := false
	labels := map[string]string{
		"test": "test",
	}
	createOpts := clustertemplates.CreateOpts{
		COE:                 coe,
		DNSNameServer:       "8.8.8.8",
		DockerStorageDriver: "overlay2",
		ExternalNetworkID:   choices.ExternalNetworkID,
		FlavorID:            choices.FlavorID,
		FloatingIPEnabled:   &boolFalse,
		ImageID:             choices.MagnumImageID,
		MasterFlavorID:      choices.FlavorID,
		MasterLBEnabled:     &boolFalse,
		Name:                name,
		Public:              &boolFalse,
		RegistryEnabled:     &boolFalse,
		ServerType:          "vm",
		// workaround for https://bugs.launchpad.net/magnum/+bug/2109685
		Labels: labels,
	}

	res := clustertemplates.Create(context.TODO(), client, createOpts)
	if res.Err != nil {
		return nil, res.Err
	}

	requestID := res.Header.Get("X-OpenStack-Request-Id")
	th.AssertEquals(t, true, requestID != "")

	t.Logf("Cluster Template %s request ID: %s", name, requestID)

	clusterTemplate, err := res.Extract()
	if err != nil {
		return nil, err
	}

	t.Logf("Successfully created cluster template: %s", clusterTemplate.Name)

	tools.PrintResource(t, clusterTemplate)
	tools.PrintResource(t, clusterTemplate.CreatedAt)

	th.AssertEquals(t, name, clusterTemplate.Name)
	th.AssertDeepEquals(t, labels, clusterTemplate.Labels)
	th.AssertEquals(t, choices.ExternalNetworkID, clusterTemplate.ExternalNetworkID)
	th.AssertEquals(t, choices.MagnumImageID, clusterTemplate.ImageID)

	return clusterTemplate, nil
}

// CreateClusterTemplate will create a random swarm cluster template.
// An error will be returned if the cluster template could not be created.
func CreateClusterTemplate(t *testing.T, client *gophercloud.ServiceClient) (*clustertemplates.ClusterTemplate, error) {
	return CreateClusterTemplateCOE(t, client, "swarm")
}

// CreateKubernetesClusterTemplate will create a random kubernetes cluster template.
// An error will be returned if the cluster template could not be created.
func CreateKubernetesClusterTemplate(t *testing.T, client *gophercloud.ServiceClient) (*clustertemplates.ClusterTemplate, error) {
	return CreateClusterTemplateCOE(t, client, "kubernetes")
}

// DeleteClusterTemplate will delete a given cluster-template. A fatal error will occur if the
// cluster-template could not be deleted. This works best as a deferred function.
func DeleteClusterTemplate(t *testing.T, client *gophercloud.ServiceClient, id string) {
	t.Logf("Attempting to delete cluster-template: %s", id)

	err := clustertemplates.Delete(context.TODO(), client, id).ExtractErr()
	if err != nil {
		t.Fatalf("Error deleting cluster-template %s: %s:", id, err)
	}

	t.Logf("Successfully deleted cluster-template: %s", id)
}

// CreateClusterTimeout will create a random cluster and wait for it to reach CREATE_COMPLETE status
// within the given timeout duration. An error will be returned if the cluster could not be created.
func CreateClusterTimeout(t *testing.T, client *gophercloud.ServiceClient, clusterTemplateID string, timeout time.Duration) (string, error) {
	clusterName := tools.RandomString("TESTACC-", 8)
	t.Logf("Attempting to create cluster: %s using template %s", clusterName, clusterTemplateID)

	choices, err := clients.AcceptanceTestChoicesFromEnv()
	if err != nil {
		return "", err
	}

	masterCount := 1
	nodeCount := 1
	// createTimeout is the creation timeout on the magnum side in minutes
	createTimeout := int(math.Ceil(timeout.Minutes()))
	createOpts := clusters.CreateOpts{
		ClusterTemplateID: clusterTemplateID,
		CreateTimeout:     &createTimeout,
		FlavorID:          choices.FlavorID,
		Keypair:           choices.MagnumKeypair,
		Labels:            map[string]string{},
		MasterCount:       &masterCount,
		MasterFlavorID:    choices.FlavorID,
		Name:              clusterName,
		NodeCount:         &nodeCount,
	}

	createResult := clusters.Create(context.TODO(), client, createOpts)
	th.AssertNoErr(t, createResult.Err)
	if len(createResult.Header["X-Openstack-Request-Id"]) > 0 {
		t.Logf("Cluster Create Request ID: %s", createResult.Header["X-Openstack-Request-Id"][0])
	}

	clusterID, err := createResult.Extract()
	if err != nil {
		return "", err
	}

	t.Logf("Cluster created: %+v", clusterID)

	err = WaitForCluster(client, clusterID, "CREATE_COMPLETE", timeout)
	if err != nil {
		return clusterID, err
	}

	t.Logf("Successfully created cluster: %s id: %s", clusterName, clusterID)
	return clusterID, nil
}

// CreateCluster will create a random cluster. An error will be returned if the
// cluster could not be created. Has a timeout of 300 seconds.
func CreateCluster(t *testing.T, client *gophercloud.ServiceClient, clusterTemplateID string) (string, error) {
	return CreateClusterTimeout(t, client, clusterTemplateID, 300*time.Second)
}

// CreateKubernetesCluster is the same as CreateCluster with a longer timeout necessary for creating a kubernetes cluster
func CreateKubernetesCluster(t *testing.T, client *gophercloud.ServiceClient, clusterTemplateID string) (string, error) {
	return CreateClusterTimeout(t, client, clusterTemplateID, 900*time.Second)
}

func DeleteCluster(t *testing.T, client *gophercloud.ServiceClient, id string) {
	t.Logf("Attempting to delete cluster: %s", id)

	r := clusters.Delete(context.TODO(), client, id)
	err := clusters.Delete(context.TODO(), client, id).ExtractErr()
	deleteRequestID := ""
	idKey := "X-Openstack-Request-Id"
	if len(r.Header[idKey]) > 0 {
		deleteRequestID = r.Header[idKey][0]
	}
	if err != nil {
		t.Fatalf("Error deleting cluster. requestID=%s clusterID=%s: err%s:", deleteRequestID, id, err)
	}

	err = WaitForCluster(client, id, "DELETE_COMPLETE", 300*time.Second)
	if err != nil {
		t.Fatalf("Error deleting cluster %s: %s:", id, err)
	}

	t.Logf("Successfully deleted cluster: %s", id)
}

func WaitForCluster(client *gophercloud.ServiceClient, clusterID string, status string, timeout time.Duration) error {
	return tools.WaitForTimeout(func(ctx context.Context) (bool, error) {
		cluster, err := clusters.Get(ctx, client, clusterID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) && status == "DELETE_COMPLETE" {
				return true, nil
			}

			return false, err
		}

		if cluster.Status == status {
			return true, nil
		}

		if strings.Contains(cluster.Status, "FAILED") {
			return false, fmt.Errorf("cluster %s FAILED. Status=%s StatusReason=%s", clusterID, cluster.Status, cluster.StatusReason)
		}

		return false, nil
	}, timeout)
}

// CreateQuota will create a random quota. An error will be returned if the
// quota could not be created.
func CreateQuota(t *testing.T, client *gophercloud.ServiceClient) (*quotas.Quotas, error) {
	name := tools.RandomString("TESTACC-", 8)
	t.Logf("Attempting to create quota: %s", name)

	idClient, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	project, err := idv3.CreateProject(t, idClient, nil)
	th.AssertNoErr(t, err)
	defer idv3.DeleteProject(t, idClient, project.ID)

	createOpts := quotas.CreateOpts{
		Resource:  "Cluster",
		ProjectID: project.ID,
		HardLimit: 10,
	}

	res := quotas.Create(context.TODO(), client, createOpts)
	if res.Err != nil {
		return nil, res.Err
	}

	requestID := res.Header.Get("X-OpenStack-Request-Id")
	th.AssertEquals(t, true, requestID != "")

	t.Logf("Quota %s request ID: %s", name, requestID)

	quota, err := res.Extract()
	if err == nil {
		t.Logf("Successfully created quota: %s", quota.ProjectID)

		tools.PrintResource(t, quota)

		th.AssertEquals(t, project.ID, quota.ProjectID)
		th.AssertEquals(t, "Cluster", quota.Resource)
		th.AssertEquals(t, 10, quota.HardLimit)
	}

	return quota, err
}
