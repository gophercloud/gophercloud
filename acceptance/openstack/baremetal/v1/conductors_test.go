//go:build acceptance || baremetal || ports
// +build acceptance baremetal ports

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/conductors"
	"github.com/gophercloud/gophercloud/pagination"

	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestConductorsList(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = "1.53"

	err = conductors.List(client, conductors.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		_, err := conductors.ExtractConductors(page)
		if err != nil {
			return false, err
		}

		return false, nil
	})
	th.AssertNoErr(t, err)
}
