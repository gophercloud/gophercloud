// +build acceptance blockstorage

package v1

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/volumes"
	"github.com/rackspace/gophercloud/openstack/utils"
	"github.com/rackspace/gophercloud/pagination"
)

var numVols = 1

func newClient() (*gophercloud.ServiceClient, error) {
	ao, err := utils.AuthOptions()
	if err != nil {
		return nil, err
	}

	client, err := openstack.AuthenticatedClient(ao)
	if err != nil {
		return nil, err
	}

	return openstack.NewBlockStorageV1(client, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
}

func TestVolumes(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Failed to create Block Storage v1 client: %v", err)
	}

	var cv *volumes.Volume
	for i := 0; i < numVols; i++ {
		cv, err = volumes.Create(client, volumes.CreateOpts{
			Size: 1,
			Name: "gophercloud-test-volume-" + strconv.Itoa(i),
		})
		if err != nil {
			t.Error(err)
			return
		}
		defer func() {
			time.Sleep(10000 * time.Millisecond)
			err = volumes.Delete(client, cv.ID)
			if err != nil {
				t.Error(err)
				return
			}
		}()

	}

	_, err = volumes.Update(client, cv.ID, volumes.UpdateOpts{
		Name: "gophercloud-updated-volume",
	})
	if err != nil {
		t.Error(err)
		return
	}

	gr, err := volumes.Get(client, cv.ID)
	if err != nil {
		t.Error(err)
		return
	}
	v, err := volumes.ExtractVolume(gr)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("Got volume: %+v\n", v)

	if v.Name != "gophercloud-updated-volume" {
		t.Errorf("Unable to update volume: Expected name: gophercloud-updated-volume\nActual name: %s", v.Name)
	}

	pager := volumes.List(client, volumes.ListOpts{})
	if err != nil {
		t.Error(err)
		return
	}
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		vols, err := volumes.ExtractVolumes(page)
		if len(vols) != numVols {
			t.Errorf("Expected %d volumes, got %d", numVols, len(vols))
		}
		return true, err
	})
}
