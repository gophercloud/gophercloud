//go:build acceptance || baremetal || conductors

package v1

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/conductors"
	"github.com/gophercloud/gophercloud/v2/pagination"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestConductorsListAndGet(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = "1.49"

	err = conductors.List(client, conductors.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		conductorList, err := conductors.ExtractConductors(page)
		if err != nil {
			return false, err
		}

		tools.PrintResource(t, conductorList)

		if len(conductorList) > 0 {
			conductor, err := conductors.Get(context.TODO(), client, conductorList[0].Hostname).Extract()
			th.AssertNoErr(t, err)

			tools.PrintResource(t, conductor)
		}

		return true, nil
	})
	th.AssertNoErr(t, err)
}
