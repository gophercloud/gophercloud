package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/flavors"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListFlavors(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleList(t)

	pages := 0
	err := flavors.List(fake.ServiceClient()).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := flavors.ExtractFlavors(page)
		if err != nil {
			return false, err
		}

		expected := []flavors.Flavor{
			{
				ID:   1,
				Name: "m1.tiny",
				RAM:  512,
				Links: []gophercloud.Link{
					{Href: "https://openstack.example.com/v1.0/1234/flavors/1", Rel: "self"},
					{Href: "https://openstack.example.com/flavors/1", Rel: "bookmark"},
				},
				StrID: "1",
			},
			{
				ID:   2,
				Name: "m1.small",
				RAM:  1024,
				Links: []gophercloud.Link{
					{Href: "https://openstack.example.com/v1.0/1234/flavors/2", Rel: "self"},
					{Href: "https://openstack.example.com/flavors/2", Rel: "bookmark"},
				},
				StrID: "2",
			},
			{
				ID:   3,
				Name: "m1.medium",
				RAM:  2048,
				Links: []gophercloud.Link{
					{Href: "https://openstack.example.com/v1.0/1234/flavors/3", Rel: "self"},
					{Href: "https://openstack.example.com/flavors/3", Rel: "bookmark"},
				},
				StrID: "3",
			},
			{
				ID:   4,
				Name: "m1.large",
				RAM:  4096,
				Links: []gophercloud.Link{
					{Href: "https://openstack.example.com/v1.0/1234/flavors/4", Rel: "self"},
					{Href: "https://openstack.example.com/flavors/4", Rel: "bookmark"},
				},
				StrID: "4",
			},
			{
				ID:   0,
				Name: "ds512M",
				RAM:  512,
				Links: []gophercloud.Link{
					{Href: "https://openstack.example.com/v1.0/1234/flavors/d1", Rel: "self"},
					{Href: "https://openstack.example.com/flavors/d1", Rel: "bookmark"},
				},
				StrID: "d1",
			},
		}

		th.AssertDeepEquals(t, expected, actual)
		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, pages)
}

func TestGetFlavor(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGet(t)

	actual, err := flavors.Get(context.TODO(), fake.ServiceClient(), flavorID).Extract()
	th.AssertNoErr(t, err)

	expected := &flavors.Flavor{
		ID:   1,
		Name: "m1.tiny",
		RAM:  512,
		Links: []gophercloud.Link{
			{Href: "https://openstack.example.com/v1.0/1234/flavors/1", Rel: "self"},
		},
		StrID: "1",
	}

	th.AssertDeepEquals(t, expected, actual)
}
