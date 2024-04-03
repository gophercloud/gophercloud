//go:build acceptance || workflow || executions

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/workflow/v2/executions"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestExecutionsCreate(t *testing.T) {
	client, err := clients.NewWorkflowV2Client()
	th.AssertNoErr(t, err)

	workflow, err := CreateWorkflow(t, client)
	th.AssertNoErr(t, err)
	defer DeleteWorkflow(t, client, workflow)

	execution, err := CreateExecution(t, client, workflow)
	th.AssertNoErr(t, err)
	defer DeleteExecution(t, client, execution)

	tools.PrintResource(t, execution)
}

func TestExecutionsList(t *testing.T) {
	client, err := clients.NewWorkflowV2Client()
	th.AssertNoErr(t, err)

	workflow, err := CreateWorkflow(t, client)
	th.AssertNoErr(t, err)
	defer DeleteWorkflow(t, client, workflow)

	execution, err := CreateExecution(t, client, workflow)
	th.AssertNoErr(t, err)
	defer DeleteExecution(t, client, execution)

	list, err := ListExecutions(t, client, &executions.ListOpts{
		Description: &executions.ListFilter{
			Value: execution.Description,
		},
		CreatedAt: &executions.ListDateFilter{
			Filter: executions.FilterGTE,
			Value:  execution.CreatedAt,
		},
		Input: execution.Input,
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, len(list))

	tools.PrintResource(t, list)
}
