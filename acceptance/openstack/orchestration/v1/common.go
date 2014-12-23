// +build acceptance

package v1

import (
	"os"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	th "github.com/rackspace/gophercloud/testhelper"
)

const template = `
{
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
`

func newClient(t *testing.T) *gophercloud.ServiceClient {
	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	client, err := openstack.AuthenticatedClient(ao)
	th.AssertNoErr(t, err)

	c, err := openstack.NewOrchestrationV1(client, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
	th.AssertNoErr(t, err)
	return c
}
