package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/orchestration/v1/stacks"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestCreateStack(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateSuccessfully(t, fakeServer, CreateOutput)
	template := new(stacks.Template)
	template.Bin = []byte(`
		{
			"heat_template_version": "2013-05-23",
			"description": "Simple template to test heat commands",
			"parameters": {
				"flavor": {
					"default": "m1.tiny",
					"type": "string"
				}
			}
		}`)
	createOpts := stacks.CreateOpts{
		Name:            "stackcreated",
		Timeout:         60,
		TemplateOpts:    template,
		DisableRollback: gophercloud.Disabled,
	}
	actual, err := stacks.Create(context.TODO(), client.ServiceClient(fakeServer), createOpts).Extract()
	th.AssertNoErr(t, err)

	expected := CreateExpected
	th.AssertDeepEquals(t, expected, actual)
}

func TestCreateStackMissingRequiredInOpts(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateSuccessfully(t, fakeServer, CreateOutput)
	template := new(stacks.Template)
	template.Bin = []byte(`
		{
			"heat_template_version": "2013-05-23",
			"description": "Simple template to test heat commands",
			"parameters": {
				"flavor": {
					"default": "m1.tiny",
					"type": "string"
				}
			}
		}`)
	createOpts := stacks.CreateOpts{
		DisableRollback: gophercloud.Disabled,
	}
	r := stacks.Create(context.TODO(), client.ServiceClient(fakeServer), createOpts)
	th.AssertEquals(t, "error creating the options map: Missing input for argument [Name]", r.Err.Error())
}

func TestAdoptStack(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateSuccessfully(t, fakeServer, CreateOutput)
	template := new(stacks.Template)
	template.Bin = []byte(`
{
  "stack_name": "postman_stack",
  "template": {
	"heat_template_version": "2013-05-23",
	"description": "Simple template to test heat commands",
	"parameters": {
	  "flavor": {
		"default": "m1.tiny",
		"type": "string"
	  }
	},
	"resources": {
	  "hello_world": {
		"type":"OS::Nova::Server",
		"properties": {
		  "key_name": "heat_key",
		  "flavor": {
			"get_param": "flavor"
		  },
		  "image": "ad091b52-742f-469e-8f3c-fd81cadf0743",
		  "user_data": "#!/bin/bash -xv\necho \"hello world\" &gt; /root/hello-world.txt\n"
		}
	  }
	}
  }
}`)
	adoptOpts := stacks.AdoptOpts{
		AdoptStackData:  `{environment{parameters{}}}`,
		Name:            "stackcreated",
		Timeout:         60,
		TemplateOpts:    template,
		DisableRollback: gophercloud.Disabled,
	}
	actual, err := stacks.Adopt(context.TODO(), client.ServiceClient(fakeServer), adoptOpts).Extract()
	th.AssertNoErr(t, err)

	expected := CreateExpected
	th.AssertDeepEquals(t, expected, actual)
}

func TestListStack(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListSuccessfully(t, fakeServer, FullListOutput)

	count := 0
	err := stacks.List(client.ServiceClient(fakeServer), nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := stacks.ExtractStacks(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ListExpected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestGetStack(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetSuccessfully(t, fakeServer, GetOutput)

	actual, err := stacks.Get(context.TODO(), client.ServiceClient(fakeServer), "postman_stack", "16ef0584-4458-41eb-87c8-0dc8d5f66c87").Extract()
	th.AssertNoErr(t, err)

	expected := GetExpected
	th.AssertDeepEquals(t, expected, actual)
}

func TestFindStack(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleFindSuccessfully(t, fakeServer, GetOutput)

	actual, err := stacks.Find(context.TODO(), client.ServiceClient(fakeServer), "16ef0584-4458-41eb-87c8-0dc8d5f66c87").Extract()
	th.AssertNoErr(t, err)

	expected := GetExpected
	th.AssertDeepEquals(t, expected, actual)
}

func TestUpdateStack(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUpdateSuccessfully(t, fakeServer)

	template := new(stacks.Template)
	template.Bin = []byte(`
		{
			"heat_template_version": "2013-05-23",
			"description": "Simple template to test heat commands",
			"parameters": {
				"flavor": {
					"default": "m1.tiny",
					"type": "string"
				}
			}
		}`)
	updateOpts := &stacks.UpdateOpts{
		TemplateOpts: template,
	}
	err := stacks.Update(context.TODO(), client.ServiceClient(fakeServer), "gophercloud-test-stack-2", "db6977b2-27aa-4775-9ae7-6213212d4ada", updateOpts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestUpdateStackNoTemplate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUpdateSuccessfully(t, fakeServer)

	parameters := make(map[string]any)
	parameters["flavor"] = "m1.tiny"

	updateOpts := &stacks.UpdateOpts{
		Parameters: parameters,
	}
	expected := stacks.ErrTemplateRequired{}

	err := stacks.Update(context.TODO(), client.ServiceClient(fakeServer), "gophercloud-test-stack-2", "db6977b2-27aa-4775-9ae7-6213212d4ada", updateOpts).ExtractErr()
	th.AssertEquals(t, expected, err)
}

func TestUpdatePatchStack(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUpdatePatchSuccessfully(t, fakeServer)

	parameters := make(map[string]any)
	parameters["flavor"] = "m1.tiny"

	updateOpts := &stacks.UpdateOpts{
		Parameters: parameters,
	}
	err := stacks.UpdatePatch(context.TODO(), client.ServiceClient(fakeServer), "gophercloud-test-stack-2", "db6977b2-27aa-4775-9ae7-6213212d4ada", updateOpts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestDeleteStack(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteSuccessfully(t, fakeServer)

	err := stacks.Delete(context.TODO(), client.ServiceClient(fakeServer), "gophercloud-test-stack-2", "db6977b2-27aa-4775-9ae7-6213212d4ada").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestPreviewStack(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePreviewSuccessfully(t, fakeServer, GetOutput)

	template := new(stacks.Template)
	template.Bin = []byte(`
		{
			"heat_template_version": "2013-05-23",
			"description": "Simple template to test heat commands",
			"parameters": {
				"flavor": {
					"default": "m1.tiny",
					"type": "string"
				}
			}
		}`)
	previewOpts := stacks.PreviewOpts{
		Name:            "stackcreated",
		Timeout:         60,
		TemplateOpts:    template,
		DisableRollback: gophercloud.Disabled,
	}
	actual, err := stacks.Preview(context.TODO(), client.ServiceClient(fakeServer), previewOpts).Extract()
	th.AssertNoErr(t, err)

	expected := PreviewExpected
	th.AssertDeepEquals(t, expected, actual)
}

func TestAbandonStack(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAbandonSuccessfully(t, fakeServer, AbandonOutput)

	actual, err := stacks.Abandon(context.TODO(), client.ServiceClient(fakeServer), "postman_stack", "16ef0584-4458-41eb-87c8-0dc8d5f66c8").Extract()
	th.AssertNoErr(t, err)

	expected := AbandonExpected
	th.AssertDeepEquals(t, expected, actual)
}
