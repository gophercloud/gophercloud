package v3

import (
	"testing"

	"github.com/gophercloud/gophercloud/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/qos"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestQoS(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	qos1, err := CreateQoS(t, client)
	th.AssertNoErr(t, err)
	defer DeleteQoS(t, client, qos1)

	qos2, err := CreateQoS(t, client)
	th.AssertNoErr(t, err)
	defer DeleteQoS(t, client, qos2)

	getQoS2, err := qos.Get(client, qos2.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, qos2, getQoS2)

	err = qos.DeleteKeys(client, qos2.ID, qos.DeleteKeysOpts{"read_iops_sec"}).ExtractErr()
	th.AssertNoErr(t, err)

	updateOpts := qos.UpdateOpts{
		Consumer: qos.ConsumerBack,
		Specs: map[string]string{
			"read_iops_sec":  "40000",
			"write_iops_sec": "40000",
		},
	}

	expectedQosSpecs := map[string]string{
		"consumer":       "back-end",
		"read_iops_sec":  "40000",
		"write_iops_sec": "40000",
	}

	updatedQosSpecs, err := qos.Update(client, qos2.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, updatedQosSpecs, expectedQosSpecs)

	listOpts := qos.ListOpts{
		Limit: 1,
	}

	err = qos.List(client, listOpts).EachPage(func(page pagination.Page) (bool, error) {
		actual, err := qos.ExtractQoS(page)
		th.AssertNoErr(t, err)
		th.AssertEquals(t, 1, len(actual))

		var found bool
		for _, q := range actual {
			if q.ID == qos1.ID || q.ID == qos2.ID {
				found = true
			}
		}

		th.AssertEquals(t, found, true)

		return true, nil
	})

	th.AssertNoErr(t, err)
}

func TestQoSAssociations(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	qos1, err := CreateQoS(t, client)
	th.AssertNoErr(t, err)
	defer DeleteQoS(t, client, qos1)

	vt, err := CreateVolumeType(t, client)
	th.AssertNoErr(t, err)
	defer DeleteVolumeType(t, client, vt)

	associateOpts := qos.AssociateOpts{
		VolumeTypeID: vt.ID,
	}

	err = qos.Associate(client, qos1.ID, associateOpts).ExtractErr()
	th.AssertNoErr(t, err)

	allQosAssociations, err := qos.ListAssociations(client, qos1.ID).AllPages()
	th.AssertNoErr(t, err)

	allAssociations, err := qos.ExtractAssociations(allQosAssociations)
	th.AssertNoErr(t, err)
	tools.PrintResource(t, allAssociations)
	th.AssertEquals(t, 1, len(allAssociations))
	th.AssertEquals(t, vt.ID, allAssociations[0].ID)

	disassociateOpts := qos.DisassociateOpts{
		VolumeTypeID: vt.ID,
	}

	err = qos.Disassociate(client, qos1.ID, disassociateOpts).ExtractErr()
	th.AssertNoErr(t, err)

	allQosAssociations, err = qos.ListAssociations(client, qos1.ID).AllPages()
	th.AssertNoErr(t, err)

	allAssociations, err = qos.ExtractAssociations(allQosAssociations)
	th.AssertNoErr(t, err)
	tools.PrintResource(t, allAssociations)
	th.AssertEquals(t, 0, len(allAssociations))

	err = qos.Associate(client, qos1.ID, associateOpts).ExtractErr()
	th.AssertNoErr(t, err)

	err = qos.DisassociateAll(client, qos1.ID).ExtractErr()
	th.AssertNoErr(t, err)
}
