// +build acceptance

package openstack

import (
	"fmt"
	blockstorage "github.com/rackspace/gophercloud/openstack/blockstorage/v1"
	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/volumes"
	"github.com/rackspace/gophercloud/openstack/identity"
	"github.com/rackspace/gophercloud/openstack/utils"
	"os"
	"testing"
)

func getClient() (*blockstorage.Client, error) {
	ao, err := utils.AuthOptions()
	if err != nil {
		return nil, err
	}

	r, err := identity.Authenticate(ao)
	if err != nil {
		return nil, err
	}

	sc, err := identity.GetServiceCatalog(r)
	if err != nil {
		return nil, err
	}

	ces, err := sc.CatalogEntries()
	if err != nil {
		return nil, err
	}

	var eps []identity.Endpoint
	for _, ce := range ces {
		if ce.Type == "volume" {
			eps = ce.Endpoints
		}
	}

	region := os.Getenv("OS_REGION_NAME")
	rep := ""
	for _, ep := range eps {
		if ep.Region == region {
			rep = ep.PublicURL
		}
	}

	client := blockstorage.NewClient(rep, r, ao)

	return client, nil

}

func TestVolumes(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Error(err)
		return
	}
	v, err := volumes.Create(client, volumes.CreateOpts{
		Size: 1,
	})
	if err != nil {
		t.Error(err)
		return
	}
}
