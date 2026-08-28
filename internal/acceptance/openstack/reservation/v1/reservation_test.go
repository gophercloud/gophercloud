//go:build acceptance || reservation

package v1

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestReservationEndpoint(t *testing.T) {
	client, err := clients.NewReservationV1Client()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, client.Endpoint)
	tools.PrintResource(t, client.ResourceBase)

	var result gophercloud.Result
	_, err = client.Get(context.TODO(), client.ServiceURL("leases"), &result.Body, nil)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, result.Body)
}
