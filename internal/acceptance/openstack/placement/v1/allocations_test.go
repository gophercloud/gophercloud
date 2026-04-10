//go:build acceptance || placement || allocations

package v1

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/allocations"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestGetAllocationsSuccess(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	consumerUUID := fmt.Sprintf("%08x-0000-0000-0000-000000000000", rand.Int31())

	// Assert: We don't have any allocations for this random UUID.
	// We get an empty allocations map, not 404.
	allocs, err := allocations.Get(context.TODO(), client, consumerUUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 0, len(allocs.Allocations))
}
