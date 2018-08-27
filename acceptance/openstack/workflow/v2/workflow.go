package v2

import (
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/workflow/v2/workflows"
	th "github.com/gophercloud/gophercloud/testhelper"
)

// CreateWorkflow creates a workflow on Mistral API.
// The created workflow is a dummy workflow that performs a simple echo.
func CreateWorkflow(t *testing.T, client *gophercloud.ServiceClient) (*workflows.Workflow, error) {
	workflowName := tools.RandomString("workflow_create_vm_", 5)

	definition := `---
version: '2.0'

` + workflowName + `:
  description: Simple workflow example
  type: direct

  tasks:
    test:
      action: std.echo output="Hello World!"`

	t.Logf("Attempting to create workflow: %s", workflowName)

	opts := &workflows.CreateOpts{
		Namespace:  "some-namespace",
		Scope:      "private",
		Definition: strings.NewReader(definition),
	}
	workflowList, err := workflows.Create(client, opts).Extract()
	if err != nil {
		return nil, err
	}
	th.AssertEquals(t, 1, len(workflowList))

	workflow := workflowList[0]

	t.Logf("Workflow created: %s", workflowName)

	th.AssertEquals(t, workflowName, workflow.Name)

	return &workflow, nil
}
