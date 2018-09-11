package v2

import (
	"testing"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/workflow/v2/crontriggers"
	"github.com/gophercloud/gophercloud/openstack/workflow/v2/workflows"
	th "github.com/gophercloud/gophercloud/testhelper"
)

// CreateCronTrigger creates a cron trigger for the given workflow.
func CreateCronTrigger(t *testing.T, client *gophercloud.ServiceClient, workflow *workflows.Workflow) (*crontriggers.CronTrigger, error) {
	crontriggerName := tools.RandomString("crontrigger_", 5)
	t.Logf("Attempting to create cron trigger: %s", crontriggerName)

	firstExecution := time.Now().AddDate(1, 0, 0)
	createOpts := crontriggers.CreateOpts{
		WorkflowID: workflow.ID,
		Name:       crontriggerName,
		Pattern:    "0 0 1 1 *",
		WorkflowInput: map[string]interface{}{
			"msg": "Hello World!",
		},
		FirstExecutionTime: &firstExecution,
	}
	crontrigger, err := crontriggers.Create(client, createOpts).Extract()
	if err != nil {
		return crontrigger, err
	}
	t.Logf("Cron trigger created: %s", crontriggerName)
	th.AssertEquals(t, crontrigger.Name, crontriggerName)
	return crontrigger, nil
}
