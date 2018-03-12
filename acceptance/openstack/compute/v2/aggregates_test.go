// +build acceptance compute aggregates

package v2

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/aggregates"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/hypervisors"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestAggregatesList(t *testing.T) {
	clients.RequireAdmin(t)

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

	for _, v := range allAggregates {
		tools.PrintResource(t, v)
	}
}

func TestAggregatesCRUD(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewComputeV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	aggregate, err := CreateAggregate(t, client)
	if err != nil {
		t.Fatalf("Unable to create an aggregate: %v", err)
	}
	defer DeleteAggregate(t, client, aggregate)

	tools.PrintResource(t, aggregate)

	updateOpts := aggregates.UpdateOpts{
		Name:             "new_aggregate_name",
		AvailabilityZone: "new_azone",
	}

	updatedAggregate, err := aggregates.Update(client, aggregate.ID, updateOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to update an aggregate: %v", err)
	}

	tools.PrintResource(t, aggregate)

	th.AssertEquals(t, updatedAggregate.Name, "new_aggregate_name")
	th.AssertEquals(t, updatedAggregate.AvailabilityZone, "new_azone")
}

func TestAggregatesAddRemoveHost(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewComputeV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	hostToAdd, err := getHypervisor(t, client)
	if err != nil {
		t.Fatal(err)
	}

	aggregate, err := CreateAggregate(t, client)
	if err != nil {
		t.Fatalf("Unable to create an aggregate: %v", err)
	}
	defer DeleteAggregate(t, client, aggregate)

	addHostOpts := aggregates.AddHostOpts{
		Host: hostToAdd.HypervisorHostname,
	}

	aggregateWithNewHost, err := aggregates.AddHost(client, aggregate.ID, addHostOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to add host to aggregate: %v", err)
	}

	tools.PrintResource(t, aggregateWithNewHost)

	th.AssertEquals(t, aggregateWithNewHost.Hosts[0], hostToAdd.HypervisorHostname)

	removeHostOpts := aggregates.RemoveHostOpts{
		Host: hostToAdd.HypervisorHostname,
	}

	aggregateWithRemovedHost, err := aggregates.RemoveHost(client, aggregate.ID, removeHostOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to remove host from aggregate: %v", err)
	}

	tools.PrintResource(t, aggregateWithRemovedHost)

	th.AssertEquals(t, len(aggregateWithRemovedHost.Hosts), 0)
}

func TestAggregatesSetRemoveMetadata(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewComputeV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	aggregate, err := CreateAggregate(t, client)
	if err != nil {
		t.Fatalf("Unable to create an aggregate: %v", err)
	}
	defer DeleteAggregate(t, client, aggregate)

	opts := aggregates.SetMetadataOpts{
		Metadata: map[string]interface{}{"key": "value"},
	}

	aggregateWithMetadata, err := aggregates.SetMetadata(client, aggregate.ID, opts).Extract()
	if err != nil {
		t.Fatalf("Unable to set metadata to aggregate: %v", err)
	}

	tools.PrintResource(t, aggregateWithMetadata)

	if _, ok := aggregateWithMetadata.Metadata["key"]; !ok {
		t.Fatalf("aggregate %s did not contain metadata", aggregateWithMetadata.Name)
	}

	optsToRemove := aggregates.SetMetadataOpts{
		Metadata: map[string]interface{}{"key": nil},
	}

	aggregateWithRemovedKey, err := aggregates.SetMetadata(client, aggregate.ID, optsToRemove).Extract()
	if err != nil {
		t.Fatalf("Unable to set metadata to aggregate: %v", err)
	}

	tools.PrintResource(t, aggregateWithRemovedKey)

	if _, ok := aggregateWithRemovedKey.Metadata["key"]; ok {
		t.Fatalf("aggregate %s still contains metadata", aggregateWithRemovedKey.Name)
	}
}

func getHypervisor(t *testing.T, client *gophercloud.ServiceClient) (*hypervisors.Hypervisor, error) {
	allPages, err := hypervisors.List(client).AllPages()
	if err != nil {
		t.Fatalf("Unable to list hypervisors: %v", err)
	}

	allHypervisors, err := hypervisors.ExtractHypervisors(allPages)
	if err != nil {
		t.Fatal("Unable to extract hypervisors")
	}

	for _, h := range allHypervisors {
		return &h, nil
	}

	return nil, fmt.Errorf("Unable to get hypervisor")
}
