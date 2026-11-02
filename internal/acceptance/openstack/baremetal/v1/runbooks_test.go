//go:build acceptance || baremetal || runbooks

package v1

import (
	"context"
	"sort"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/runbooks"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestRunbooksCRUD(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = "1.92"

	runbook, err := CreateRunbook(t, client)
	th.AssertNoErr(t, err)
	defer DeleteRunbook(t, client, runbook)

	found := false
	err = runbooks.List(client, runbooks.ListOpts{Detail: true}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		runbookList, err := runbooks.ExtractRunbooks(page)
		if err != nil {
			return false, err
		}
		for _, r := range runbookList {
			if r.UUID == runbook.UUID {
				found = true
				th.AssertEquals(t, runbook.Name, r.Name)
				th.AssertEquals(t, 1, len(r.Steps))
				return false, nil
			}
		}
		return true, nil
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, true, found)

	byName, err := runbooks.Get(context.TODO(), client, runbook.Name).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, runbook.UUID, byName.UUID)

	updated, err := runbooks.Update(context.TODO(), client, runbook.UUID, runbooks.UpdateOpts{
		runbooks.UpdateOperation{
			Op:    runbooks.ReplaceOp,
			Path:  "/disable_ramdisk",
			Value: true,
		},
		runbooks.UpdateOperation{
			Op:   runbooks.RemoveOp,
			Path: "/extra/purpose",
		},
	}).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, true, updated.DisableRamdisk)
	th.AssertEquals(t, 0, len(updated.Extra))
}

func TestRunbooksTraits(t *testing.T) {
	clients.RequireLong(t)
	clients.SkipReleasesBelow(t, "stable/2026.2")

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = "1.112"

	runbook, err := CreateRunbook(t, client)
	th.AssertNoErr(t, err)
	defer DeleteRunbook(t, client, runbook)

	updated, err := runbooks.Update(context.TODO(), client, runbook.UUID, runbooks.UpdateOpts{
		runbooks.UpdateOperation{
			Op:    runbooks.AddOp,
			Path:  "/description",
			Value: "Gophercloud acceptance test",
		},
	}).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "Gophercloud acceptance test", *updated.Description)

	err = runbooks.SetTraits(context.TODO(), client, runbook.UUID, []string{"CUSTOM_GOPHER_A", "CUSTOM_GOPHER_B"}).ExtractErr()
	th.AssertNoErr(t, err)
	err = runbooks.AddTrait(context.TODO(), client, runbook.UUID, "CUSTOM_GOPHER_C").ExtractErr()
	th.AssertNoErr(t, err)
	err = runbooks.DeleteTrait(context.TODO(), client, runbook.UUID, "CUSTOM_GOPHER_A").ExtractErr()
	th.AssertNoErr(t, err)

	traits, err := runbooks.ListTraits(context.TODO(), client, runbook.UUID).Extract()
	th.AssertNoErr(t, err)
	sort.Strings(traits)
	th.AssertDeepEquals(t, []string{"CUSTOM_GOPHER_B", "CUSTOM_GOPHER_C"}, traits)

	got, err := runbooks.Get(context.TODO(), client, runbook.UUID).Extract()
	th.AssertNoErr(t, err)
	sort.Strings(got.Traits)
	th.AssertDeepEquals(t, traits, got.Traits)

	err = runbooks.DeleteTraits(context.TODO(), client, runbook.UUID).ExtractErr()
	th.AssertNoErr(t, err)
	traits, err = runbooks.ListTraits(context.TODO(), client, runbook.UUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 0, len(traits))
}
