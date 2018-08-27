package testing

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/workflow/v2/crontriggers"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestCreateCronTrigger(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/cron_triggers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusCreated)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `
			{
				"created_at": "1970-01-01 00:00:00",
				"id": "1",
				"name": "trigger",
				"pattern": "* * * * *",
				"project_id": "p1",
				"remaining_executions": 42,
				"scope": "private",
				"updated_at": "1970-01-01 00:00:00",
				"first_execution_time": "1970-01-01 00:00:00",
				"next_execution_time": "1970-01-01 00:00:00",
				"workflow_id": "w1",
				"workflow_input": "{\"msg\": \"hello\"}",
				"workflow_name": "my_wf",
				"workflow_params": "{\"msg\": \"world\"}"
			}
		`)
	})

	firstExecution := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	opts := &crontriggers.CreateOpts{
		WorkflowID:         "w1",
		Name:               "trigger",
		FirstExecutionTime: &firstExecution,
		WorkflowParams: map[string]interface{}{
			"msg": "world",
		},
		WorkflowInput: map[string]interface{}{
			"msg": "hello",
		},
	}

	actual, err := crontriggers.Create(fake.ServiceClient(), opts).Extract()
	if err != nil {
		t.Fatalf("Unable to create cron trigger: %v", err)
	}

	expected := &crontriggers.CronTrigger{
		ID:                  "1",
		Name:                "trigger",
		Pattern:             "* * * * *",
		ProjectID:           "p1",
		RemainingExecutions: 42,
		Scope:               "private",
		WorkflowID:          "w1",
		WorkflowName:        "my_wf",
		WorkflowParams: map[string]interface{}{
			"msg": "world",
		},
		WorkflowInput: map[string]interface{}{
			"msg": "hello",
		},
		CreatedAt:          time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
		FirstExecutionTime: &firstExecution,
		NextExecutionTime:  &firstExecution,
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, but was %#v", expected, actual)
	}
}

func TestDeleteCronTrigger(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/cron_triggers/1", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.WriteHeader(http.StatusAccepted)
	})

	res := crontriggers.Delete(fake.ServiceClient(), "1")
	th.AssertNoErr(t, res.Err)
}
