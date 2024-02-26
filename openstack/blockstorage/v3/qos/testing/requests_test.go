package testing

import (
	"context"
	"reflect"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/qos"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockCreateResponse(t)

	options := qos.CreateOpts{
		Name:     "qos-001",
		Consumer: qos.ConsumerFront,
		Specs: map[string]string{
			"read_iops_sec": "20000",
		},
	}
	actual, err := qos.Create(context.TODO(), client.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &createQoSExpected, actual)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockDeleteResponse(t)

	res := qos.Delete(context.TODO(), client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22", qos.DeleteOpts{})
	th.AssertNoErr(t, res.Err)
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListResponse(t)

	pages := 0
	err := qos.List(client.ServiceClient(), nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++
		actual, err := qos.ExtractQoS(page)
		if err != nil {
			return false, err
		}

		expected := []qos.QoS{
			{ID: "1", Consumer: "back-end", Name: "foo", Specs: map[string]string{}},
			{ID: "2", Consumer: "front-end", Name: "bar", Specs: map[string]string{
				"read_iops_sec": "20000",
			},
			},
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Expected %#v, but was %#v", expected, actual)
		}

		return true, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if pages != 1 {
		t.Errorf("Expected one page, got %d", pages)
	}
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockGetResponse(t)

	actual, err := qos.Get(context.TODO(), client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &getQoSExpected, actual)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	MockUpdateResponse(t)

	updateOpts := qos.UpdateOpts{
		Consumer: qos.ConsumerBack,
		Specs: map[string]string{
			"read_iops_sec":  "40000",
			"write_iops_sec": "40000",
		},
	}

	expected := UpdateQos
	actual, err := qos.Update(context.TODO(), client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22", updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, expected, actual)
}

func TestDeleteKeys(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockDeleteKeysResponse(t)

	res := qos.DeleteKeys(context.TODO(), client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22", qos.DeleteKeysOpts{"read_iops_sec"})
	th.AssertNoErr(t, res.Err)
}

func TestAssociate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockAssociateResponse(t)

	associateOpts := qos.AssociateOpts{
		VolumeTypeID: "b596be6a-0ce9-43fa-804a-5c5e181ede76",
	}

	res := qos.Associate(context.TODO(), client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22", associateOpts)
	th.AssertNoErr(t, res.Err)
}

func TestDisssociate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockDisassociateResponse(t)

	disassociateOpts := qos.DisassociateOpts{
		VolumeTypeID: "b596be6a-0ce9-43fa-804a-5c5e181ede76",
	}

	res := qos.Disassociate(context.TODO(), client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22", disassociateOpts)
	th.AssertNoErr(t, res.Err)
}

func TestDissasociateAll(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockDisassociateAllResponse(t)

	res := qos.DisassociateAll(context.TODO(), client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22")
	th.AssertNoErr(t, res.Err)
}

func TestQosAssociationsList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListAssociationsResponse(t)

	expected := []qos.QosAssociation{
		{
			Name:            "foo",
			ID:              "2f954bcf047c4ee9b09a37d49ae6db54",
			AssociationType: "volume_type",
		},
	}

	allPages, err := qos.ListAssociations(client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22").AllPages(context.TODO())
	th.AssertNoErr(t, err)

	actual, err := qos.ExtractAssociations(allPages)
	th.AssertNoErr(t, err)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, but was %#v", expected, actual)
	}
}
