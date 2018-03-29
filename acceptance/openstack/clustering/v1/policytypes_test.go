// +build acceptance clustering policytypes

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/clustering/v1/policytypes"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestPolicyTypeList(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewClusteringV1Client()
	th.AssertNoErr(t, err)

	allPages, err := policytypes.List(client).AllPages()
	th.AssertNoErr(t, err)

	allPolicyTypes, err := policytypes.ExtractPolicyTypes(allPages)
	th.AssertNoErr(t, err)

	for _, v := range allPolicyTypes {
		tools.PrintResource(t, v)
	}
}
