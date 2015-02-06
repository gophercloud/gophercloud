package stackresources

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/rackspace/gophercloud"
	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

var FindExpected = []Resource{
	Resource{
		Name: "hello_world",
		Links: []gophercloud.Link{
			gophercloud.Link{
				Href: "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b/resources/hello_world",
				Rel:  "self",
			},
			gophercloud.Link{
				Href: "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b",
				Rel:  "stack",
			},
		},
		LogicalID:    "hello_world",
		StatusReason: "state changed",
		UpdatedTime:  time.Date(2015, 2, 5, 21, 33, 11, 0, time.UTC),
		RequiredBy:   []interface{}{},
		Status:       "CREATE_IN_PROGRESS",
		PhysicalID:   "49181cd6-169a-4130-9455-31185bbfc5bf",
		Type:         "OS::Nova::Server",
	},
}

const FindOutput = `
{
  "resources": [
  {
    "resource_name": "hello_world",
    "links": [
      {
      "href": "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b/resources/hello_world",
      "rel": "self"
      },
      {
        "href": "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b",
        "rel": "stack"
      }
    ],
    "logical_resource_id": "hello_world",
    "resource_status_reason": "state changed",
    "updated_time": "2015-02-05T21:33:11Z",
    "required_by": [],
    "resource_status": "CREATE_IN_PROGRESS",
    "physical_resource_id": "49181cd6-169a-4130-9455-31185bbfc5bf",
    "resource_type": "OS::Nova::Server"
  }
  ]
}`

// HandleFindSuccessfully creates an HTTP handler at `/stacks/hello_world/resources`
// on the test handler mux that responds with a `Find` response.
func HandleFindSuccessfully(t *testing.T, output string) {
	th.Mux.HandleFunc("/stacks/hello_world/resources", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, output)
	})
}

var ListExpected = []Resource{
	Resource{
		Name: "hello_world",
		Links: []gophercloud.Link{
			gophercloud.Link{
				Href: "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b/resources/hello_world",
				Rel:  "self",
			},
			gophercloud.Link{
				Href: "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b",
				Rel:  "stack",
			},
		},
		LogicalID:    "hello_world",
		StatusReason: "state changed",
		UpdatedTime:  time.Date(2015, 2, 5, 21, 33, 11, 0, time.UTC),
		RequiredBy:   []interface{}{},
		Status:       "CREATE_IN_PROGRESS",
		PhysicalID:   "49181cd6-169a-4130-9455-31185bbfc5bf",
		Type:         "OS::Nova::Server",
	},
}

const ListOutput = `{
  "resources": [
  {
    "resource_name": "hello_world",
    "links": [
    {
      "href": "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b/resources/hello_world",
      "rel": "self"
    },
    {
      "href": "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b",
      "rel": "stack"
    }
    ],
    "logical_resource_id": "hello_world",
    "resource_status_reason": "state changed",
    "updated_time": "2015-02-05T21:33:11Z",
    "required_by": [],
    "resource_status": "CREATE_IN_PROGRESS",
    "physical_resource_id": "49181cd6-169a-4130-9455-31185bbfc5bf",
    "resource_type": "OS::Nova::Server"
  }
]
}`

// HandleListSuccessfully creates an HTTP handler at `/stacks/hello_world/49181cd6-169a-4130-9455-31185bbfc5bf/resources`
// on the test handler mux that responds with a `List` response.
func HandleListSuccessfully(t *testing.T, output string) {
	th.Mux.HandleFunc("/stacks/hello_world/49181cd6-169a-4130-9455-31185bbfc5bf/resources", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Set("Content-Type", "application/json")
		r.ParseForm()
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, output)
		default:
			t.Fatalf("Unexpected marker: [%s]", marker)
		}
	})
}

var GetExpected = &Resource{
	Name: "wordpress_instance",
	Links: []gophercloud.Link{
		gophercloud.Link{
			Href: "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/teststack/0b1771bd-9336-4f2b-ae86-a80f971faf1e/resources/wordpress_instance",
			Rel:  "self",
		},
		gophercloud.Link{
			Href: "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/teststack/0b1771bd-9336-4f2b-ae86-a80f971faf1e",
			Rel:  "stack",
		},
	},
	LogicalID:    "wordpress_instance",
	StatusReason: "state changed",
	UpdatedTime:  time.Date(2014, 12, 10, 18, 34, 35, 0, time.UTC),
	RequiredBy:   []interface{}{},
	Status:       "CREATE_COMPLETE",
	PhysicalID:   "00e3a2fe-c65d-403c-9483-4db9930dd194",
	Type:         "OS::Nova::Server",
}

const GetOutput = `
{
  "resource": {
    "resource_name": "wordpress_instance",
    "description": "",
    "links": [
    {
      "href": "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/teststack/0b1771bd-9336-4f2b-ae86-a80f971faf1e/resources/wordpress_instance",
      "rel": "self"
    },
    {
      "href": "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/teststack/0b1771bd-9336-4f2b-ae86-a80f971faf1e",
      "rel": "stack"
    }
    ],
    "logical_resource_id": "wordpress_instance",
    "resource_status": "CREATE_COMPLETE",
    "updated_time": "2014-12-10T18:34:35Z",
    "required_by": [],
    "resource_status_reason": "state changed",
    "physical_resource_id": "00e3a2fe-c65d-403c-9483-4db9930dd194",
    "resource_type": "OS::Nova::Server"
  }
}`

// HandleGetSuccessfully creates an HTTP handler at `/stacks/teststack/0b1771bd-9336-4f2b-ae86-a80f971faf1e/resources/wordpress_instance`
// on the test handler mux that responds with a `Get` response.
func HandleGetSuccessfully(t *testing.T, output string) {
	th.Mux.HandleFunc("/stacks/teststack/0b1771bd-9336-4f2b-ae86-a80f971faf1e/resources/wordpress_instance", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, output)
	})
}

var MetadataExpected = map[string]string{
	"number": "7",
	"animal": "auk",
}

const MetadataOutput = `
{
    "metadata": {
      "number": "7",
      "animal": "auk"
    }
}`

// HandleMetadataSuccessfully creates an HTTP handler at `/stacks/teststack/0b1771bd-9336-4f2b-ae86-a80f971faf1e/resources/wordpress_instance/metadata`
// on the test handler mux that responds with a `Metadata` response.
func HandleMetadataSuccessfully(t *testing.T, output string) {
	th.Mux.HandleFunc("/stacks/teststack/0b1771bd-9336-4f2b-ae86-a80f971faf1e/resources/wordpress_instance/metadata", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, output)
	})
}

var ListTypesExpected = []string{
	"OS::Nova::Server",
	"OS::Heat::RandomString",
	"OS::Swift::Container",
	"OS::Trove::Instance",
	"OS::Nova::FloatingIPAssociation",
	"OS::Cinder::VolumeAttachment",
	"OS::Nova::FloatingIP",
	"OS::Nova::KeyPair",
}

const ListTypesOutput = `
{
  "resource_types": [
    "OS::Nova::Server",
    "OS::Heat::RandomString",
    "OS::Swift::Container",
    "OS::Trove::Instance",
    "OS::Nova::FloatingIPAssociation",
    "OS::Cinder::VolumeAttachment",
    "OS::Nova::FloatingIP",
    "OS::Nova::KeyPair"
  ]
}`

// HandleListTypesSuccessfully creates an HTTP handler at `/resource_types`
// on the test handler mux that responds with a `ListTypes` response.
func HandleListTypesSuccessfully(t *testing.T, output string) {
	th.Mux.HandleFunc("/resource_types", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, output)
	})
}

var GetSchemaExpected = &TypeSchema{
	Attributes: map[string]interface{}{
		"an_attribute": map[string]interface{}{
			"description": "An attribute description .",
		},
	},
	Properties: map[string]interface{}{
		"a_property": map[string]interface{}{
			"update_allowed": false,
			"required":       true,
			"type":           "string",
			"description":    "A resource description.",
		},
	},
	ResourceType: "OS::Heat::AResourceName",
}

const GetSchemaOutput = `
{
  "attributes": {
    "an_attribute": {
      "description": "An attribute description ."
    }
  },
  "properties": {
    "a_property": {
      "update_allowed": false,
      "required": true,
      "type": "string",
      "description": "A resource description."
    }
  },
  "resource_type": "OS::Heat::AResourceName"
}`

// HandleGetSchemaSuccessfully creates an HTTP handler at `/resource_types/OS::Heat::AResourceName`
// on the test handler mux that responds with a `Schema` response.
func HandleGetSchemaSuccessfully(t *testing.T, output string) {
	th.Mux.HandleFunc("/resource_types/OS::Heat::AResourceName", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, output)
	})
}

const GetTemplateOutput = `
{
  "Outputs": {
    "addresses": {
      "Description": "A dict of all network addresses with correspondingport_id.",
      "Value": "{\"Fn::GetAtt\": [\"Server\", \"addresses\"]}"
    },
    "first_address": {
      "Description": "Convenience attribute to fetch the first assigned network address, or an empty string if nothing has been assigned at this time. Result may not be predictable if the server has addresses from more than one network.",
      "Value": "{\"Fn::GetAtt\": [\"Server\", \"first_address\"]}"
    },
    "show": {
      "Description": "A dict of all server details as returned by the API.",
      "Value": "{\"Fn::GetAtt\": [\"Server\", \"show\"]}"
    },
    "instance_name": {
      "Description": "AWS compatible instance name.",
      "Value": "{\"Fn::GetAtt\": [\"Server\", \"instance_name\"]}"
    },
    "accessIPv4": {
      "Description": "The manually assigned alternative public IPv4 address of the server.",
      "Value": "{\"Fn::GetAtt\": [\"Server\", \"accessIPv4\"]}"
    },
    "accessIPv6": {
      "Description": "The manually assigned alternative public IPv6 address of the server.",
      "Value": "{\"Fn::GetAtt\": [\"Server\", \"accessIPv6\"]}"
    },
    "networks": {
      "Description": "A dict of assigned network addresses of the form: {\"public\": [ip1, ip2...], \"private\": [ip3, ip4]}.",
      "Value": "{\"Fn::GetAtt\": [\"Server\", \"networks\"]}"
    }
  },
  "HeatTemplateFormatVersion": "2012-12-12",
  "Parameters": {
    "scheduler_hints": {
      "Type": "Json",
      "Description": "Arbitrary key-value pairs specified by the client to help boot a server."
    },
    "admin_pass": {
      "Type": "String",
      "Description": "The administrator password for the server."
    },
    "user_data_format": {
      "Default": "HEAT_CFNTOOLS",
      "Type": "String",
      "Description": "How the user_data should be formatted for the server. For HEAT_CFNTOOLS, the user_data is bundled as part of the heat-cfntools cloud-init boot configuration data. For RAW the user_data is passed to Nova unmodified. For SOFTWARE_CONFIG user_data is bundled as part of the software config data, and metadata is derived from any associated SoftwareDeployment resources.",
      "AllowedValues": [
      "HEAT_CFNTOOLS",
      "RAW",
      "SOFTWARE_CONFIG"
      ]
    },
    "admin_user": {
      "Type": "String",
      "Description": "Name of the administrative user to use on the server. This property will be removed from Juno in favor of the default cloud-init user set up for each image (e.g. \"ubuntu\" for Ubuntu 12.04+, \"fedora\" for Fedora 19+ and \"cloud-user\" for CentOS/RHEL 6.5)."
    },
    "name": {
      "Type": "String",
      "Description": "Server name."
    },
    "block_device_mapping": {
      "Type": "CommaDelimitedList",
      "Description": "Block device mappings for this server."
    },
    "key_name": {
      "Type": "String",
      "Description": "Name of keypair to inject into the server."
    },
    "image": {
      "Type": "String",
      "Description": "The ID or name of the image to boot with."
    },
    "availability_zone": {
      "Type": "String",
      "Description": "Name of the availability zone for server placement."
    },
    "image_update_policy": {
      "Default": "REPLACE",
      "Type": "String",
      "Description": "Policy on how to apply an image-id update; either by requesting a server rebuild or by replacing the entire server",
      "AllowedValues": [
      "REBUILD",
      "REPLACE",
      "REBUILD_PRESERVE_EPHEMERAL"
      ]
    },
    "software_config_transport": {
      "Default": "POLL_SERVER_CFN",
      "Type": "String",
      "Description": "How the server should receive the metadata required for software configuration. POLL_SERVER_CFN will allow calls to the cfn API action DescribeStackResource authenticated with the provided keypair. POLL_SERVER_HEAT will allow calls to the Heat API resource-show using the provided keystone credentials.",
      "AllowedValues": [
      "POLL_SERVER_CFN",
      "POLL_SERVER_HEAT"
      ]
    },
    "metadata": {
      "Type": "Json",
      "Description": "Arbitrary key/value metadata to store for this server. Both keys and values must be 255 characters or less.  Non-string values will be serialized to JSON (and the serialized string must be 255 characters or less)."
    },
    "personality": {
      "Default": {},
        "Type": "Json",
        "Description": "A map of files to create/overwrite on the server upon boot. Keys are file names and values are the file contents."
    },
    "user_data": {
      "Default": "",
      "Type": "String",
      "Description": "User data script to be executed by cloud-init."
    },
    "flavor_update_policy": {
      "Default": "RESIZE",
      "Type": "String",
      "Description": "Policy on how to apply a flavor update; either by requesting a server resize or by replacing the entire server.",
      "AllowedValues": [
      "RESIZE",
      "REPLACE"
      ]
    },
    "flavor": {
      "Type": "String",
      "Description": "The ID or name of the flavor to boot onto."
    },
    "diskConfig": {
      "Type": "String",
      "Description": "Control how the disk is partitioned when the server is created.",
      "AllowedValues": [
      "AUTO",
      "MANUAL"
      ]
    },
    "reservation_id": {
      "Type": "String",
      "Description": "A UUID for the set of servers being requested."
    },
    "networks": {
      "Type": "CommaDelimitedList",
      "Description": "An ordered list of nics to be added to this server, with information about connected networks, fixed ips, port etc."
    },
    "security_groups": {
      "Default": [],
      "Type": "CommaDelimitedList",
      "Description": "List of security group names or IDs. Cannot be used if neutron ports are associated with this server; assign security groups to the ports instead."
    },
    "config_drive": {
      "Type": "String",
      "Description": "value for config drive either boolean, or volume-id."
    }
  },
  "Resources": {
    "Server": {
      "Type": "OS::Nova::Server",
      "Properties": {
        "scheduler_hints": {
          "Ref": "scheduler_hints"
        },
        "admin_pass": {
          "Ref": "admin_pass"
        },
        "user_data_format": {
          "Ref": "user_data_format"
        },
        "admin_user": {
          "Ref": "admin_user"
        },
        "name": {
          "Ref": "name"
        },
        "block_device_mapping": {
          "Fn::Split": [
          ",",
          {
            "Ref": "block_device_mapping"
          }
          ]
        },
        "key_name": {
          "Ref": "key_name"
        },
        "image": {
          "Ref": "image"
        },
        "availability_zone": {
          "Ref": "availability_zone"
        },
        "image_update_policy": {
          "Ref": "image_update_policy"
        },
        "software_config_transport": {
          "Ref": "software_config_transport"
        },
        "metadata": {
          "Ref": "metadata"
        },
        "personality": {
          "Ref": "personality"
        },
        "user_data": {
          "Ref": "user_data"
        },
        "flavor_update_policy": {
          "Ref": "flavor_update_policy"
        },
        "flavor": {
          "Ref": "flavor"
        },
        "diskConfig": {
          "Ref": "diskConfig"
        },
        "reservation_id": {
          "Ref": "reservation_id"
        },
        "networks": {
          "Fn::Split": [
          ",",
          {
            "Ref": "networks"
          }
          ]
        },
        "security_groups": {
          "Fn::Split": [
          ",",
          {
            "Ref": "security_groups"
          }
          ]
        },
        "config_drive": {
          "Ref": "config_drive"
        }
      }
    }
  }
}`
