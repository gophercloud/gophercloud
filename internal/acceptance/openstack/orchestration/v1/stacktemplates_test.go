//go:build acceptance || orchestration || stacktemplates

package v1

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/orchestration/v1/stacktemplates"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestStackTemplatesCRUD(t *testing.T) {
	client, err := clients.NewOrchestrationV1Client()
	th.AssertNoErr(t, err)

	stack, err := CreateStack(t, client)
	th.AssertNoErr(t, err)
	defer DeleteStack(t, client, stack.Name, stack.ID)

	tmpl, err := stacktemplates.Get(context.TODO(), client, stack.Name, stack.ID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, tmpl)
}

func TestStackTemplatesValidate(t *testing.T) {
	client, err := clients.NewOrchestrationV1Client()
	th.AssertNoErr(t, err)

	validateOpts := stacktemplates.ValidateOpts{
		Template: validateTemplate,
	}

	validatedTemplate, err := stacktemplates.Validate(context.TODO(), client, validateOpts).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, validatedTemplate)
}

func TestStackTemplateWithFile(t *testing.T) {
	client, err := clients.NewOrchestrationV1Client()
	th.AssertNoErr(t, err)

	stack, err := CreateStackWithFile(t, client)
	th.AssertNoErr(t, err)
	defer DeleteStack(t, client, stack.Name, stack.ID)

	tmpl, err := stacktemplates.Get(context.TODO(), client, stack.Name, stack.ID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, tmpl)
}
