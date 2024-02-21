//go:build acceptance || clustering || profiletypes
// +build acceptance clustering profiletypes

package v1

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/clustering/v1/profiletypes"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestProfileTypesList(t *testing.T) {
	client, err := clients.NewClusteringV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.5"

	allPages, err := profiletypes.List(client).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allProfileTypes, err := profiletypes.ExtractProfileTypes(allPages)
	th.AssertNoErr(t, err)

	for _, profileType := range allProfileTypes {
		tools.PrintResource(t, profileType)
	}
}
func TestProfileTypesOpsList(t *testing.T) {
	client, err := clients.NewClusteringV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.5"

	profileTypeName := "os.nova.server-1.0"
	allPages, err := profiletypes.ListOps(client, profileTypeName).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	ops, err := profiletypes.ExtractOps(allPages)
	th.AssertNoErr(t, err)

	for k, v := range ops {
		tools.PrintResource(t, k)
		tools.PrintResource(t, v)
	}
}
