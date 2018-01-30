// +build acceptance compute aggregates

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/aggregates"
)

func TestAggregatesList(t *testing.T) {
	client, err := clients.NewComputeV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	allPages, err := aggregates.List(client).AllPages()
	if err != nil {
		t.Fatalf("Unable to list aggregates: %v", err)
	}

	allAggregates, err := aggregates.ExtractAggregates(allPages)
	if err != nil {
		t.Fatalf("Unable to extract aggregates")
	}

	for _, h := range allAggregates {
		tools.PrintResource(t, h)
	}
}

func TestAggregatesCreateDelete(t *testing.T) {
	client, err := clients.NewComputeV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	opts := aggregates.CreateOpts{
		Name:             "name",
		AvailabilityZone: "london",
	}

	aggregate, err := aggregates.Create(client, opts).Extract()
	if err != nil {
		t.Fatalf("Unable to create an aggregate: %v", err)
	}

	err = aggregates.Delete(client, aggregate.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete an aggregate: %v", err)
	}
}
