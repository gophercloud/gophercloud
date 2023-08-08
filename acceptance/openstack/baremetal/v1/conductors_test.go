//go:build acceptance || baremetal || conductors
// +build acceptance baremetal conductors

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/conductors"
	"github.com/gophercloud/gophercloud/pagination"

	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestConductorsListAndGet(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = "1.49"

	err = conductors.List(client, conductors.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		conductorList, err := conductors.ExtractConductors(page)
		if err != nil {
			return false, err
		}

		tools.PrintResource(t, conductorList)

		if len(conductorList) > 0 {
			conductor, err := conductors.Get(client, conductorList[0].Hostname).Extract()
			th.AssertNoErr(t, err)

			tools.PrintResource(t, conductor)
		}

		return true, nil
	})
	th.AssertNoErr(t, err)
}
