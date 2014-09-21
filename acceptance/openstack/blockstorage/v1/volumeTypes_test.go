// +build acceptance

package v1

import (
	"strconv"
	"testing"

	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/volumeTypes"
)

var numVolTypes = 1

func TestVolumeTypes(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Failed to create Block Storage v1 client: %v", err)
	}

	var cvt *volumeTypes.VolumeType
	for i := 0; i < numVolTypes; i++ {
		cvt, err = volumeTypes.Create(client, volumeTypes.CreateOpts{
			ExtraSpecs: map[string]string{
				"capabilities": "gpu",
			},
			Name: "gophercloud-test-volumeType-200" + strconv.Itoa(i),
		})
		if err != nil {
			t.Error(err)
			return
		} /*
			defer func() {
				time.Sleep(10000 * time.Millisecond)
				err = volumeTypes.Delete(client, cvt.ID)
				if err != nil {
					t.Error(err)
					return
				}
			}*/
		t.Logf("created volume type: %+v\n", cvt)
	}

}
