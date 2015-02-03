package stacks

import (
	"testing"

	os "github.com/rackspace/gophercloud/openstack/orchestration/v1/stacks"
	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

func TestCreateStack(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	os.HandleCreateSuccessfully(t, CreateOutput)

	createOpts := os.CreateOpts{
		Name:    "stackcreated",
		Timeout: 60,
		Template: `{
      "outputs": {
        "db_host": {
          "value": {
            "get_attr": [
            "db",
            "hostname"
            ]
          }
        }
      },
      "heat_template_version": "2014-10-16",
      "description": "HEAT template for creating a Cloud Database.\n",
      "parameters": {
        "db_name": {
          "default": "wordpress",
          "type": "string",
          "description": "the name for the database",
          "constraints": [
          {
            "length": {
              "max": 64,
              "min": 1
            },
            "description": "must be between 1 and 64 characters"
          },
          {
            "allowed_pattern": "[a-zA-Z][a-zA-Z0-9]*",
            "description": "must begin with a letter and contain only alphanumeric characters."
          }
          ]
        },
        "db_instance_name": {
          "default": "Cloud_DB",
          "type": "string",
          "description": "the database instance name"
        },
        "db_username": {
          "default": "admin",
          "hidden": true,
          "type": "string",
          "description": "database admin account username",
          "constraints": [
          {
            "length": {
              "max": 16,
              "min": 1
                },
              "description": "must be between 1 and 16 characters"
            },
            {
              "allowed_pattern": "[a-zA-Z][a-zA-Z0-9]*",
              "description": "must begin with a letter and contain only alphanumeric characters."
            }
          ]
          },
          "db_volume_size": {
            "default": 30,
            "type": "number",
            "description": "database volume size (in GB)",
            "constraints": [
            {
              "range": {
                "max": 1024,
                "min": 1
              },
              "description": "must be between 1 and 1024 GB"
            }
            ]
          },
          "db_flavor": {
            "default": "1GB Instance",
            "type": "string",
            "description": "database instance size",
            "constraints": [
            {
              "description": "must be a valid cloud database flavor",
              "allowed_values": [
              "1GB Instance",
              "2GB Instance",
              "4GB Instance",
              "8GB Instance",
              "16GB Instance"
              ]
            }
            ]
          },
        "db_password": {
          "default": "admin",
          "hidden": true,
          "type": "string",
          "description": "database admin account password",
          "constraints": [
          {
            "length": {
              "max": 41,
              "min": 1
            },
            "description": "must be between 1 and 14 characters"
          },
          {
            "allowed_pattern": "[a-zA-Z0-9]*",
            "description": "must contain only alphanumeric characters."
          }
          ]
        }
      },
      "resources": {
        "db": {
          "type": "OS::Trove::Instance",
          "properties": {
            "flavor": {
              "get_param": "db_flavor"
            },
            "size": {
              "get_param": "db_volume_size"
            },
            "users": [
            {
              "password": {
                "get_param": "db_password"
              },
              "name": {
                "get_param": "db_username"
              },
              "databases": [
              {
                "get_param": "db_name"
              }
              ]
            }
            ],
            "name": {
              "get_param": "db_instance_name"
            },
            "databases": [
            {
              "name": {
                "get_param": "db_name"
              }
            }
            ]
          }
        }
      }
    }`,
		DisableRollback: os.Disable,
	}
	actual, err := Create(fake.ServiceClient(), createOpts).Extract()
	th.AssertNoErr(t, err)

	expected := CreateExpected
	th.AssertDeepEquals(t, expected, actual)
}
