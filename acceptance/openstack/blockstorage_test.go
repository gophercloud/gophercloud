// +build acceptance

package openstack

import (
	"fmt"
	blockstorage "github.com/rackspace/gophercloud/openstack/blockstorage/v1"
	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/snapshots"
	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/volumes"
	identity "github.com/rackspace/gophercloud/openstack/identity/v2"
	"github.com/rackspace/gophercloud/openstack/utils"
	"os"
	"strconv"
	"testing"
	"time"
)

var numVols = 2

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

	var cv volumes.Volume
	for i := 0; i < numVols; i++ {
		cv, err := volumes.Create(client, volumes.CreateOpts{
			"size":         1,
			"display_name": "test-volume" + strconv.Itoa(i),
		})
		if err != nil {
			t.Error(err)
			return
		}
		defer func() {
			time.Sleep(10000 * time.Millisecond)
			err = volumes.Delete(client, volumes.DeleteOpts{
				"id": cv.Id,
			})
			if err != nil {
				t.Error(err)
				return
			}
		}()
	}

	vols, err := volumes.List(client, volumes.ListOpts{
		"full": true,
	})
	if err != nil {
		t.Error(err)
		return
	}
	if len(vols) != numVols {
		t.Errorf("Expected %d volumes, got %d", numVols, len(vols))
		return
	}

	vols, err = volumes.List(client, volumes.ListOpts{
		"full": false,
	})
	if err != nil {
		t.Error(err)
		return
	}
	if len(vols) != numVols {
		t.Errorf("Expected %d volumes, got %d", numVols, len(vols))
		return
	}

	_, err = volumes.Get(client, volumes.GetOpts{
		"id": cv.Id,
	})
	if err != nil {
		t.Error(err)
		return
	}

}

func TestSnapshots(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Error(err)
		return
	}

	var css snapshots.Snapshot

	cv, err := volumes.Create(client, volumes.CreateOpts{
		"size":         1,
		"display_name": "test-volume",
	})
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		for i := 0; i < 60; i++ {
			gss, _ := snapshots.Get(client, snapshots.GetOpts{
				"id": css.Id,
			})
			if gss.Status == "" {
				err = volumes.Delete(client, volumes.DeleteOpts{
					"id": cv.Id,
				})
				if err != nil {
					t.Error(err)
					return
				}
				break
			}
			time.Sleep(5000 * time.Millisecond)
		}
	}()

	for i := 0; i < 60; i++ {
		gv, err := volumes.Get(client, volumes.GetOpts{
			"id": cv.Id,
		})
		if err != nil {
			t.Error(err)
			return
		}
		if gv.Status == "available" {
			break
		}
		time.Sleep(2000 * time.Millisecond)
	}

	css, err = snapshots.Create(client, snapshots.CreateOpts{
		"volume_id":    cv.Id,
		"display_name": "test-snapshot",
	})
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		for i := 0; i < 60; i++ {
			gss, err := snapshots.Get(client, snapshots.GetOpts{
				"id": css.Id,
			})
			if err != nil {
				t.Error(err)
				return
			}
			if gss.Status == "available" {
				err = snapshots.Delete(client, snapshots.DeleteOpts{
					"id": css.Id,
				})
				if err != nil {
					t.Error(err)
					return
				}
				break
			}
			time.Sleep(2000 * time.Millisecond)
		}
	}()

	lss, err := snapshots.List(client, snapshots.ListOpts{
		Full: true,
	})
	if err != nil {
		t.Error(err)
		return
	}
}
