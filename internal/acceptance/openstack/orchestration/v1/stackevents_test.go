//go:build acceptance || orchestration || stackevents

package v1

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/openstack/orchestration/v1/stackevents"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestStackEvents(t *testing.T) {
	client, err := clients.NewOrchestrationV1Client()
	th.AssertNoErr(t, err)

	stack, err := CreateStack(t, client)
	th.AssertNoErr(t, err)
	defer DeleteStack(t, client, stack.Name, stack.ID)

	allPages, err := stackevents.List(client, stack.Name, stack.ID, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allEvents, err := stackevents.ExtractEvents(allPages)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, len(allEvents), 4)

	/*
			allPages is currently broke
		allPages, err = stackevents.ListResourceEvents(client, stack.Name, stack.ID, basicTemplateResourceName, nil).AllPages(context.TODO())
		th.AssertNoErr(t, err)
		allEvents, err = stackevents.ExtractEvents(allPages)
		th.AssertNoErr(t, err)

		for _, v := range allEvents {
			tools.PrintResource(t, v)
		}
	*/
}
