//go:build acceptance || orchestration || stacks

package v1

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/orchestration/v1/stacks"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestStacksPagination(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	client, err := clients.NewOrchestrationV1Client()
	th.AssertNoErr(t, err)

	expected := make(map[string]int)
	for range 2 {
		name := tools.RandomString("ACPTTEST-pagination-", 16)
		template := new(stacks.Template)
		template.Bin = []byte(basicTemplate)
		stack, err := stacks.Create(ctx, client, stacks.CreateOpts{
			Name:         name,
			TemplateOpts: template,
			Timeout:      2,
		}).Extract()
		th.AssertNoErr(t, err)
		t.Cleanup(func() {
			cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cleanupCancel()
			if err := stacks.Delete(cleanupCtx, client, name, stack.ID).ExtractErr(); err != nil {
				t.Errorf("Unable to delete stack %s: %v", stack.ID, err)
			}
		})
		th.AssertNoErr(t, gophercloud.WaitFor(ctx, func(ctx context.Context) (bool, error) {
			current, err := stacks.Get(ctx, client, name, stack.ID).Extract()
			if err != nil {
				return false, err
			}
			return current.Status == "CREATE_COMPLETE", nil
		}))
		expected[stack.ID] = 1
		t.Logf("Created stack %s (%s)", name, stack.ID)
	}

	t.Run("EachPage", func(t *testing.T) {
		actual := make(map[string]int)
		pages := 0
		err := stacks.List(client, stacks.ListOpts{Limit: 1}).EachPage(ctx, func(_ context.Context, page pagination.Page) (bool, error) {
			items, err := stacks.ExtractStacks(page)
			if err != nil {
				return false, err
			}
			pages++
			th.CheckEquals(t, 1, len(items))
			for _, item := range items {
				if _, ok := expected[item.ID]; ok {
					actual[item.ID]++
				}
			}
			return true, nil
		})
		th.AssertNoErr(t, err)
		t.Logf("Visited %d non-empty pages", pages)
		th.CheckTrue(t, pages >= 2)
		th.CheckDeepEquals(t, expected, actual)
	})

	t.Run("AllPages", func(t *testing.T) {
		page, err := stacks.List(client, stacks.ListOpts{Limit: 1}).AllPages(ctx)
		th.AssertNoErr(t, err)
		items, err := stacks.ExtractStacks(page)
		th.AssertNoErr(t, err)
		actual := make(map[string]int)
		for _, item := range items {
			if _, ok := expected[item.ID]; ok {
				actual[item.ID]++
			}
		}
		th.CheckDeepEquals(t, expected, actual)
	})
}

func TestStacksCRUD(t *testing.T) {
	client, err := clients.NewOrchestrationV1Client()
	th.AssertNoErr(t, err)

	createdStack, err := CreateStack(t, client)
	th.AssertNoErr(t, err)
	defer DeleteStack(t, client, createdStack.Name, createdStack.ID)

	tools.PrintResource(t, createdStack)
	tools.PrintResource(t, createdStack.CreationTime)

	template := new(stacks.Template)
	template.Bin = []byte(basicTemplate)
	updateOpts := stacks.UpdateOpts{
		TemplateOpts: template,
		Timeout:      20,
	}

	err = stacks.Update(context.TODO(), client, createdStack.Name, createdStack.ID, updateOpts).ExtractErr()
	th.AssertNoErr(t, err)

	err = WaitForStackStatus(client, createdStack.Name, createdStack.ID, "UPDATE_COMPLETE")
	th.AssertNoErr(t, err)

	var found bool
	allPages, err := stacks.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allStacks, err := stacks.ExtractStacks(allPages)
	th.AssertNoErr(t, err)

	for _, v := range allStacks {
		if v.ID == createdStack.ID {
			found = true
		}
	}

	th.AssertTrue(t, found)
}
