package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestCronTriggersCreateDelete(t *testing.T) {
	client, err := clients.NewWorkflowV2Client()
	th.AssertNoErr(t, err)

	workflow, err := CreateWorkflow(t, client)
	th.AssertNoErr(t, err)

	trigger, err := CreateCronTrigger(t, client, workflow)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, trigger)
}
