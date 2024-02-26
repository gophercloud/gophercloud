//go:build acceptance || clustering || events
// +build acceptance clustering events

package v1

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/clustering/v1/events"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestEventsList(t *testing.T) {
	client, err := clients.NewClusteringV1Client()
	th.AssertNoErr(t, err)

	opts := events.ListOpts{
		Limit: 200,
	}

	allPages, err := events.List(client, opts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allEvents, err := events.ExtractEvents(allPages)
	th.AssertNoErr(t, err)

	for _, event := range allEvents {
		tools.PrintResource(t, event)
	}
}
