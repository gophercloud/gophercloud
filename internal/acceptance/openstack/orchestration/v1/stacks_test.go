//go:build acceptance || orchestration || stacks

package v1

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/orchestration/v1/stacks"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

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
		Tags:         []string{"tag1", "tag2"},
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

	// Check tag filters using multiple tags to exercise comma-separated queries.
	for _, test := range []struct {
		name  string
		opts  stacks.ListOpts
		found bool
	}{
		{
			name:  "tags matching all",
			opts:  stacks.ListOpts{Tags: []string{"tag1", "tag2"}},
			found: true,
		},
		{
			name: "tags missing one",
			opts: stacks.ListOpts{Tags: []string{"tag1", "tag3"}},
		},
		{
			name:  "tags-any matching one",
			opts:  stacks.ListOpts{TagsAny: []string{"tag1", "tag3"}},
			found: true,
		},
		{
			name: "tags-any matching none",
			opts: stacks.ListOpts{TagsAny: []string{"tag3", "tag4"}},
		},
		{
			name:  "not-tags missing one",
			opts:  stacks.ListOpts{NotTags: []string{"tag1", "tag3"}},
			found: true,
		},
		{
			name: "not-tags matching all",
			opts: stacks.ListOpts{NotTags: []string{"tag1", "tag2"}},
		},
		{
			name:  "not-tags-any matching none",
			opts:  stacks.ListOpts{NotTagsAny: []string{"tag3", "tag4"}},
			found: true,
		},
		{
			name: "not-tags-any matching one",
			opts: stacks.ListOpts{NotTagsAny: []string{"tag1", "tag3"}},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			allPages, err := stacks.List(client, test.opts).AllPages(context.TODO())
			th.AssertNoErr(t, err)

			allStacks, err := stacks.ExtractStacks(allPages)
			th.AssertNoErr(t, err)

			found := false
			for _, stack := range allStacks {
				if stack.ID == createdStack.ID {
					found = true
					break
				}
			}
			th.AssertEquals(t, test.found, found)
		})
	}
}
