//go:build acceptance || image || tasks

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/tasks"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestTasksListEachPage(t *testing.T) {
	client, err := clients.NewImageV2Client()
	th.AssertNoErr(t, err)

	listOpts := tasks.ListOpts{
		Limit: 1,
	}

	pager := tasks.List(client, listOpts)
	err = pager.EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		tasks, err := tasks.ExtractTasks(page)
		th.AssertNoErr(t, err)

		for _, task := range tasks {
			tools.PrintResource(t, task)
		}

		return true, nil
	})
	th.AssertNoErr(t, err)
}

func TestTasksListAllPages(t *testing.T) {
	client, err := clients.NewImageV2Client()
	th.AssertNoErr(t, err)

	listOpts := tasks.ListOpts{}

	allPages, err := tasks.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allTasks, err := tasks.ExtractTasks(allPages)
	th.AssertNoErr(t, err)

	for _, i := range allTasks {
		tools.PrintResource(t, i)
	}
}

func TestTaskCreate(t *testing.T) {
	client, err := clients.NewImageV2Client()
	th.AssertNoErr(t, err)

	task, err := CreateTask(t, client, ImportImageURL)
	if err != nil {
		t.Fatalf("Unable to create an Image service task: %v", err)
	}

	tools.PrintResource(t, task)
}
