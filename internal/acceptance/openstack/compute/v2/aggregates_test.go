//go:build acceptance || compute || aggregates

package v2

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/aggregates"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/hypervisors"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestAggregatesList(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	allPages, err := aggregates.List(client).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allAggregates, err := aggregates.ExtractAggregates(allPages)
	th.AssertNoErr(t, err)

	for _, v := range allAggregates {
		tools.PrintResource(t, v)
	}
}

func TestAggregatesCRUD(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	aggregate, err := CreateAggregate(t, client)
	th.AssertNoErr(t, err)

	defer DeleteAggregate(t, client, aggregate)

	tools.PrintResource(t, aggregate)

	updateOpts := aggregates.UpdateOpts{
		Name:             "new_aggregate_name",
		AvailabilityZone: "new_azone",
	}

	updatedAggregate, err := aggregates.Update(context.TODO(), client, aggregate.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, aggregate)

	th.AssertEquals(t, updatedAggregate.Name, "new_aggregate_name")
	th.AssertEquals(t, updatedAggregate.AvailabilityZone, "new_azone")
}

func TestAggregatesAddRemoveHost(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	hostToAdd, err := getHypervisor(t, client)
	th.AssertNoErr(t, err)

	aggregate, err := CreateAggregate(t, client)
	th.AssertNoErr(t, err)
	defer DeleteAggregate(t, client, aggregate)

	addHostOpts := aggregates.AddHostOpts{
		Host: hostToAdd,
	}

	aggregateWithNewHost, err := aggregates.AddHost(context.TODO(), client, aggregate.ID, addHostOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, aggregateWithNewHost)

	th.AssertEquals(t, aggregateWithNewHost.Hosts[0], hostToAdd)

	removeHostOpts := aggregates.RemoveHostOpts{
		Host: hostToAdd,
	}

	aggregateWithRemovedHost, err := aggregates.RemoveHost(context.TODO(), client, aggregate.ID, removeHostOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, aggregateWithRemovedHost)

	th.AssertEquals(t, len(aggregateWithRemovedHost.Hosts), 0)
}

func TestAggregatesSetRemoveMetadata(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	aggregate, err := CreateAggregate(t, client)
	th.AssertNoErr(t, err)
	defer DeleteAggregate(t, client, aggregate)

	opts := aggregates.SetMetadataOpts{
		Metadata: map[string]any{"key": "value"},
	}

	aggregateWithMetadata, err := aggregates.SetMetadata(context.TODO(), client, aggregate.ID, opts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, aggregateWithMetadata)

	if _, ok := aggregateWithMetadata.Metadata["key"]; !ok {
		t.Fatalf("aggregate %s did not contain metadata", aggregateWithMetadata.Name)
	}

	optsToRemove := aggregates.SetMetadataOpts{
		Metadata: map[string]any{"key": nil},
	}

	aggregateWithRemovedKey, err := aggregates.SetMetadata(context.TODO(), client, aggregate.ID, optsToRemove).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, aggregateWithRemovedKey)

	if _, ok := aggregateWithRemovedKey.Metadata["key"]; ok {
		t.Fatalf("aggregate %s still contains metadata", aggregateWithRemovedKey.Name)
	}
}

func getHypervisor(t *testing.T, client *gophercloud.ServiceClient) (string, error) {
	allPages, err := hypervisors.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allHypervisors, err := hypervisors.ExtractHypervisors(allPages)
	th.AssertNoErr(t, err)

	for _, h := range allHypervisors {
		// Nova API takes Hostnames, not FQDNs, so we need to strip the domain.
		host := strings.Split(h.HypervisorHostname, ".")[0]
		return host, nil
	}

	return "", fmt.Errorf("Unable to get hypervisor")
}
