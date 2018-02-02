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

	createdAggregate, err := CreateAggregate(t, client)
	if err != nil {
		t.Fatalf("Unable to create an aggregate: %v", err)
	}
	defer DeleteAggregate(t, client, createdAggregate)

	tools.PrintResource(t, createdAggregate)
}

func TestAggregatesGet(t *testing.T) {
	client, err := clients.NewComputeV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	createdAggregate, err := CreateAggregate(t, client)
	if err != nil {
		t.Fatalf("Unable to create an aggregate: %v", err)
	}
	defer DeleteAggregate(t, client, createdAggregate)

	aggregate, err := aggregates.Get(client, createdAggregate.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get an aggregate: %v", err)
	}

	tools.PrintResource(t, aggregate)
}

func TestAggregatesUpdate(t *testing.T) {
	client, err := clients.NewComputeV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	createdAggregate, err := CreateAggregate(t, client)
	if err != nil {
		t.Fatalf("Unable to create an aggregate: %v", err)
	}
	defer DeleteAggregate(t, client, createdAggregate)

	updateOpts := aggregates.UpdateOpts{
		Name:             "new_aggregate_name",
		AvailabilityZone: "new_azone",
	}

	updatedAggregate, err := aggregates.Update(client, createdAggregate.ID, updateOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to update an aggregate: %v", err)
	}

	tools.PrintResource(t, updatedAggregate)
}
