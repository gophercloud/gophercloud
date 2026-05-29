package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/orchestration/v1/stacktemplates"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestGetTemplate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetSuccessfully(t, fakeServer, GetOutput)

	actual, err := stacktemplates.Get(context.TODO(), client.ServiceClient(fakeServer), "postman_stack", "16ef0584-4458-41eb-87c8-0dc8d5f66c87").Extract()
	th.AssertNoErr(t, err)

	expected := GetExpected
	th.AssertDeepEquals(t, expected, string(actual))
}

func TestValidateTemplate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleValidateSuccessfully(t, fakeServer, ValidateOutput)

	opts := stacktemplates.ValidateOpts{
		Template: `{
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
		      "type": "OS::Nova::Server",
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
		}`,
	}
	actual, err := stacktemplates.Validate(context.TODO(), client.ServiceClient(fakeServer), opts).Extract()
	th.AssertNoErr(t, err)

	expected := ValidateExpected
	th.AssertDeepEquals(t, expected, actual)
}
