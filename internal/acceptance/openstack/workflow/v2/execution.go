package v2

import (
	"context"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/workflow/v2/executions"
	"github.com/gophercloud/gophercloud/v2/openstack/workflow/v2/workflows"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// CreateExecution creates an execution for the given workflow.
func CreateExecution(t *testing.T, client *gophercloud.ServiceClient, workflow *workflows.Workflow) (*executions.Execution, error) {
	executionDescription := tools.RandomString("execution_", 5)

	t.Logf("Attempting to create execution: %s", executionDescription)
	createOpts := executions.CreateOpts{
		ID:                executionDescription,
		WorkflowID:        workflow.ID,
		WorkflowNamespace: workflow.Namespace,
		Description:       executionDescription,
		Input: map[string]any{
			"msg": "Hello World!",
		},
	}
	execution, err := executions.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return execution, err
	}

	t.Logf("Execution created: %s", executionDescription)

	th.AssertEquals(t, execution.Description, executionDescription)

	t.Logf("Wait for execution status SUCCESS: %s", executionDescription)
	th.AssertNoErr(t, tools.WaitFor(func(ctx context.Context) (bool, error) {
		latest, err := executions.Get(ctx, client, execution.ID).Extract()
		if err != nil {
			return false, err
		}

		if latest.State == "SUCCESS" {
			execution = latest
			return true, nil
		}

		if latest.State == "ERROR" {
			return false, fmt.Errorf("Execution in ERROR state")
		}

		return false, nil
	}))
	t.Logf("Execution success: %s", executionDescription)

	return execution, nil
}

// DeleteExecution deletes an execution.
func DeleteExecution(t *testing.T, client *gophercloud.ServiceClient, execution *executions.Execution) {
	err := executions.Delete(context.TODO(), client, execution.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete executions %s: %v", execution.Description, err)
	}
	t.Logf("Deleted executions: %s", execution.Description)
}

// ListExecutions lists the executions.
func ListExecutions(t *testing.T, client *gophercloud.ServiceClient, opts executions.ListOptsBuilder) ([]executions.Execution, error) {
	allPages, err := executions.List(client, opts).AllPages(context.TODO())
	if err != nil {
		t.Fatalf("Unable to list executions: %v", err)
	}

	executionsList, err := executions.ExtractExecutions(allPages)
	if err != nil {
		t.Fatalf("Unable to extract executions: %v", err)
	}

	t.Logf("Executions list find, length: %d", len(executionsList))
	return executionsList, err
}
