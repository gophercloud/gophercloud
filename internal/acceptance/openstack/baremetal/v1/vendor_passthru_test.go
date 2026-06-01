//go:build acceptance || baremetal || nodes || drivers

package v1

import (
	"context"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/drivers"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestVendorPassthruMethods(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = "1.38"

	node, err := CreateFakeNode(t, client)
	th.AssertNoErr(t, err)
	defer DeleteNode(t, client, node)

	nodeMethods, err := nodes.ListVendorPassthruMethods(context.TODO(), client, node.UUID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, nodeMethods)

	err = nodes.CallVendorPassthru(context.TODO(), client, node.UUID, http.MethodPost, nodes.VendorPassthruCallOpts{
		Method: "gophercloud_acceptance_missing_method",
		Body: map[string]any{
			"noop": true,
		},
	}).Err
	if err == nil {
		t.Fatalf("expected missing node vendor passthru method to fail")
	}

	var driverName string
	err = drivers.ListDrivers(client, drivers.ListDriversOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		driverList, err := drivers.ExtractDrivers(page)
		if err != nil {
			return false, err
		}
		if len(driverList) == 0 {
			return true, nil
		}

		driverName = driverList[0].Name
		return false, nil
	})
	th.AssertNoErr(t, err)
	if driverName == "" {
		t.Skip("no baremetal drivers available")
	}

	driverMethods, err := drivers.ListVendorPassthruMethods(context.TODO(), client, driverName).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, driverMethods)

	err = drivers.CallVendorPassthru(context.TODO(), client, driverName, http.MethodPost, drivers.VendorPassthruCallOpts{
		Method: "gophercloud_acceptance_missing_method",
		Body: map[string]any{
			"noop": true,
		},
	}).Err
	if err == nil {
		t.Fatalf("expected missing driver vendor passthru method to fail")
	}
}
