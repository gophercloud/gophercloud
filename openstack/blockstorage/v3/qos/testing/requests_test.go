package testing

import (
	"reflect"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/qos"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
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
	actual, err := qos.Create(client.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &createQoSExpected, actual)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockDeleteResponse(t)

	res := qos.Delete(client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22", qos.DeleteOpts{})
	th.AssertNoErr(t, res.Err)
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListResponse(t)

	pages := 0
	err := qos.List(client.ServiceClient(), nil).EachPage(func(page pagination.Page) (bool, error) {
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
