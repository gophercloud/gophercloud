package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clustertemplates"
	th "github.com/gophercloud/gophercloud/testhelper"
)

// CreateClusterTemplate will create a COE cluster template
func CreateClusterTemplate(t *testing.T, client *gophercloud.ServiceClient) (*clustertemplates.ClusterTemplate, error) {
	choices, err := clients.AcceptanceTestChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	// Create a new cluster template
	opts := clustertemplates.CreateOpts{Name: "k8s",
		COE:               "kubernetes",
		ImageID:           choices.COEImageID,
		KeypairID:         "default",
		DockerVolumeSize:  5,
		FlavorID:          choices.FlavorID,
		MasterFlavorID:    choices.FlavorID,
		ExternalNetworkID: choices.ExternalNetworkID,
	}
	ct, err := clustertemplates.Create(client, opts).Extract()
	th.AssertNoErr(t, err)
	// defer clustertemplates.Delete(client, ct.ID)
	th.AssertEquals(t, "calico", ct.NetworkDriver)
	clusterTemplateID := ct.ID
	t.Logf(clusterTemplateID)

	return ct, err
}
