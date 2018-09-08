package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestWorkflowsCreateGetDelete(t *testing.T) {
	client, err := clients.NewWorkflowV2Client()
	th.AssertNoErr(t, err)

	workflow, err := CreateWorkflow(t, client)
	th.AssertNoErr(t, err)
	defer DeleteWorkflow(t, client, workflow)

	workflowget, err := GetWorkflow(t, client, workflow.ID)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, workflowget)
}
