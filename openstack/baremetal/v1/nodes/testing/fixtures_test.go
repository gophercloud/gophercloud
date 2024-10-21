package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	inventorytest "github.com/gophercloud/gophercloud/v2/openstack/baremetal/inventory/testing"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/nodes"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// NodeListBody contains the canned body of a nodes.List response, without detail.
const NodeListBody = `
 {
  "nodes": [
    {
      "instance_uuid": null,
      "links": [
        {
          "href": "http://ironic.example.com:6385/v1/nodes/d2630783-6ec8-4836-b556-ab427c4b581e",
          "rel": "self"
        },
        {
          "href": "http://ironic.example.com:6385/nodes/d2630783-6ec8-4836-b556-ab427c4b581e",
          "rel": "bookmark"
        }
      ],
      "maintenance": false,
      "name": "foo",
      "power_state": null,
      "provision_state": "enroll"
    },
    {
      "instance_uuid": null,
      "links": [
        {
          "href": "http://ironic.example.com:6385/v1/nodes/08c84581-58f5-4ea2-a0c6-dd2e5d2b3662",
          "rel": "self"
        },
        {
          "href": "http://ironic.example.com:6385/nodes/08c84581-58f5-4ea2-a0c6-dd2e5d2b3662",
          "rel": "bookmark"
        }
      ],
      "maintenance": false,
      "name": "bar",
      "power_state": null,
      "provision_state": "enroll"
    },
    {
      "instance_uuid": null,
      "links": [
        {
          "href": "http://ironic.example.com:6385/v1/nodes/c9afd385-5d89-4ecb-9e1c-68194da6b474",
          "rel": "self"
        },
        {
          "href": "http://ironic.example.com:6385/nodes/c9afd385-5d89-4ecb-9e1c-68194da6b474",
          "rel": "bookmark"
        }
      ],
      "maintenance": false,
      "name": "baz",
      "power_state": null,
      "provision_state": "enroll"
    }
  ]
}
`

// NodeListDetailBody contains the canned body of a nodes.ListDetail response.
const NodeListDetailBody = `
 {
  "nodes": [
    {
      "automated_clean": null,
      "bios_interface": "no-bios",
      "boot_interface": "pxe",
      "chassis_uuid": null,
      "clean_step": {},
      "conductor_group": "",
      "console_enabled": false,
      "console_interface": "no-console",
      "created_at": "2019-01-31T19:59:28+00:00",
      "deploy_interface": "iscsi",
      "deploy_step": {},
      "driver": "ipmi",
      "driver_info": {
        "ipmi_port": "6230",
        "ipmi_username": "admin",
        "deploy_kernel": "http://172.22.0.1/images/tinyipa-stable-rocky.vmlinuz",
        "ipmi_address": "192.168.122.1",
        "deploy_ramdisk": "http://172.22.0.1/images/tinyipa-stable-rocky.gz",
        "ipmi_password": "admin"

      },
      "driver_internal_info": {},
      "extra": {},
      "fault": null,
      "firmware_interface": "no-firmware",
      "inspect_interface": "no-inspect",
      "inspection_finished_at": null,
      "inspection_started_at": null,
      "instance_info": {},
      "instance_uuid": null,
      "last_error": null,
      "links": [
        {
          "href": "http://ironic.example.com:6385/v1/nodes/d2630783-6ec8-4836-b556-ab427c4b581e",
          "rel": "self"
        },
        {
          "href": "http://ironic.example.com:6385/nodes/d2630783-6ec8-4836-b556-ab427c4b581e",
          "rel": "bookmark"
        }
      ],
      "maintenance": false,
      "maintenance_reason": null,
      "management_interface": "ipmitool",
      "name": "foo",
      "network_interface": "flat",
      "portgroups": [
        {
          "href": "http://ironic.example.com:6385/v1/nodes/d2630783-6ec8-4836-b556-ab427c4b581e/portgroups",
          "rel": "self"
        },
        {
          "href": "http://ironic.example.com:6385/nodes/d2630783-6ec8-4836-b556-ab427c4b581e/portgroups",
          "rel": "bookmark"
        }
      ],
      "ports": [
        {
          "href": "http://ironic.example.com:6385/v1/nodes/d2630783-6ec8-4836-b556-ab427c4b581e/ports",
          "rel": "self"
        },
        {
          "href": "http://ironic.example.com:6385/nodes/d2630783-6ec8-4836-b556-ab427c4b581e/ports",
          "rel": "bookmark"
        }
      ],
      "power_interface": "ipmitool",
      "power_state": null,
      "properties": {},
      "provision_state": "enroll",
      "provision_updated_at": "2019-02-15T17:21:29+00:00",
      "raid_config": {},
      "raid_interface": "no-raid",
      "retired": false,
      "retired_reason": "No longer needed",
      "rescue_interface": "no-rescue",
      "reservation": null,
      "resource_class": null,
      "states": [
        {
          "href": "http://ironic.example.com:6385/v1/nodes/d2630783-6ec8-4836-b556-ab427c4b581e/states",
          "rel": "self"
        },
        {
          "href": "http://ironic.example.com:6385/nodes/d2630783-6ec8-4836-b556-ab427c4b581e/states",
          "rel": "bookmark"
        }
      ],
      "storage_interface": "noop",
      "target_power_state": null,
      "target_provision_state": null,
      "target_raid_config": {},
      "traits": [],
      "updated_at": "2019-02-15T19:59:29+00:00",
      "uuid": "d2630783-6ec8-4836-b556-ab427c4b581e",
      "vendor_interface": "ipmitool",
      "volume": [
        {
          "href": "http://ironic.example.com:6385/v1/nodes/d2630783-6ec8-4836-b556-ab427c4b581e/volume",
          "rel": "self"
        },
        {
          "href": "http://ironic.example.com:6385/nodes/d2630783-6ec8-4836-b556-ab427c4b581e/volume",
          "rel": "bookmark"
        }
      ]
    },
    {
      "automated_clean": null,
      "bios_interface": "no-bios",
      "boot_interface": "pxe",
      "chassis_uuid": null,
      "clean_step": {},
      "conductor_group": "",
      "console_enabled": false,
      "console_interface": "no-console",
      "created_at": "2019-01-31T19:59:29+00:00",
      "deploy_interface": "iscsi",
      "deploy_step": {},
      "driver": "ipmi",
      "driver_info": {},
      "driver_internal_info": {},
      "extra": {},
      "fault": null,
      "firmware_interface": "no-firmware",
      "inspect_interface": "no-inspect",
      "inspection_finished_at": "2023-02-02T14:45:59.705249Z",
      "inspection_started_at": "2023-02-02T14:35:59.682403Z",
      "instance_info": {},
      "instance_uuid": null,
      "last_error": null,
      "links": [
        {
          "href": "http://ironic.example.com:6385/v1/nodes/08c84581-58f5-4ea2-a0c6-dd2e5d2b3662",
          "rel": "self"
        },
        {
          "href": "http://ironic.example.com:6385/nodes/08c84581-58f5-4ea2-a0c6-dd2e5d2b3662",
          "rel": "bookmark"
        }
      ],
      "maintenance": false,
      "maintenance_reason": null,
      "management_interface": "ipmitool",
      "name": "bar",
      "network_interface": "flat",
      "portgroups": [
        {
          "href": "http://ironic.example.com:6385/v1/nodes/08c84581-58f5-4ea2-a0c6-dd2e5d2b3662/portgroups",
          "rel": "self"
        },
        {
          "href": "http://ironic.example.com:6385/nodes/08c84581-58f5-4ea2-a0c6-dd2e5d2b3662/portgroups",
          "rel": "bookmark"
        }
      ],
      "ports": [
        {
          "href": "http://ironic.example.com:6385/v1/nodes/08c84581-58f5-4ea2-a0c6-dd2e5d2b3662/ports",
          "rel": "self"
        },
        {
          "href": "http://ironic.example.com:6385/nodes/08c84581-58f5-4ea2-a0c6-dd2e5d2b3662/ports",
          "rel": "bookmark"
        }
      ],
      "power_interface": "ipmitool",
      "power_state": null,
      "properties": {},
      "provision_state": "available",
      "provision_updated_at": null,
      "raid_config": {},
      "raid_interface": "no-raid",
      "retired": false,
      "retired_reason": "No longer needed",
      "rescue_interface": "no-rescue",
      "reservation": null,
      "resource_class": null,
      "states": [
        {
          "href": "http://ironic.example.com:6385/v1/nodes/08c84581-58f5-4ea2-a0c6-dd2e5d2b3662/states",
          "rel": "self"
        },
        {
          "href": "http://ironic.example.com:6385/nodes/08c84581-58f5-4ea2-a0c6-dd2e5d2b3662/states",
          "rel": "bookmark"
        }
      ],
      "storage_interface": "noop",
      "target_power_state": null,
      "target_provision_state": null,
      "target_raid_config": {},
      "traits": [],
      "updated_at": "2019-02-15T19:59:29+00:00",
      "uuid": "08c84581-58f5-4ea2-a0c6-dd2e5d2b3662",
      "vendor_interface": "ipmitool",
      "volume": [
        {
          "href": "http://ironic.example.com:6385/v1/nodes/08c84581-58f5-4ea2-a0c6-dd2e5d2b3662/volume",
          "rel": "self"
        },
        {
          "href": "http://ironic.example.com:6385/nodes/08c84581-58f5-4ea2-a0c6-dd2e5d2b3662/volume",
          "rel": "bookmark"
        }
      ]
    },
    {
      "automated_clean": null,
      "bios_interface": "no-bios",
      "boot_interface": "pxe",
      "chassis_uuid": null,
      "clean_step": {},
      "conductor_group": "",
      "console_enabled": false,
      "console_interface": "no-console",
      "created_at": "2019-01-31T19:59:30+00:00",
      "deploy_interface": "iscsi",
      "deploy_step": {},
      "driver": "ipmi",
      "driver_info": {},
      "driver_internal_info": {},
      "extra": {},
      "fault": null,
      "firmware_interface": "no-firmware",
      "inspect_interface": "no-inspect",
      "inspection_finished_at": null,
      "inspection_started_at": null,
      "instance_info": {},
      "instance_uuid": null,
      "last_error": null,
      "links": [
        {
          "href": "http://ironic.example.com:6385/v1/nodes/c9afd385-5d89-4ecb-9e1c-68194da6b474",
          "rel": "self"
        },
        {
          "href": "http://ironic.example.com:6385/nodes/c9afd385-5d89-4ecb-9e1c-68194da6b474",
          "rel": "bookmark"
        }
      ],
      "maintenance": false,
      "maintenance_reason": null,
      "management_interface": "ipmitool",
      "name": "baz",
      "network_interface": "flat",
      "portgroups": [
        {
          "href": "http://ironic.example.com:6385/v1/nodes/c9afd385-5d89-4ecb-9e1c-68194da6b474/portgroups",
          "rel": "self"
        },
        {
          "href": "http://ironic.example.com:6385/nodes/c9afd385-5d89-4ecb-9e1c-68194da6b474/portgroups",
          "rel": "bookmark"
        }
      ],
      "ports": [
        {
          "href": "http://ironic.example.com:6385/v1/nodes/c9afd385-5d89-4ecb-9e1c-68194da6b474/ports",
          "rel": "self"
        },
        {
          "href": "http://ironic.example.com:6385/nodes/c9afd385-5d89-4ecb-9e1c-68194da6b474/ports",
          "rel": "bookmark"
        }
      ],
      "power_interface": "ipmitool",
      "power_state": null,
      "properties": {},
      "provision_state": "enroll",
      "provision_updated_at": null,
      "raid_config": {},
      "raid_interface": "no-raid",
      "retired": false,
      "retired_reason": "No longer needed",
      "rescue_interface": "no-rescue",
      "reservation": null,
      "resource_class": null,
      "states": [
        {
          "href": "http://ironic.example.com:6385/v1/nodes/c9afd385-5d89-4ecb-9e1c-68194da6b474/states",
          "rel": "self"
        },
        {
          "href": "http://ironic.example.com:6385/nodes/c9afd385-5d89-4ecb-9e1c-68194da6b474/states",
          "rel": "bookmark"
        }
      ],
      "storage_interface": "noop",
      "target_power_state": null,
      "target_provision_state": null,
      "target_raid_config": {},
      "traits": [],
      "updated_at": "2019-02-15T19:59:29+00:00",
      "uuid": "c9afd385-5d89-4ecb-9e1c-68194da6b474",
      "vendor_interface": "ipmitool",
      "volume": [
        {
          "href": "http://ironic.example.com:6385/v1/nodes/c9afd385-5d89-4ecb-9e1c-68194da6b474/volume",
          "rel": "self"
        },
        {
          "href": "http://ironic.example.com:6385/nodes/c9afd385-5d89-4ecb-9e1c-68194da6b474/volume",
          "rel": "bookmark"
        }
      ]
    }
  ]
}
`

// SingleNodeBody is the canned body of a Get request on an existing node.
const SingleNodeBody = `
{
  "automated_clean": null,
  "bios_interface": "no-bios",
  "boot_interface": "pxe",
  "chassis_uuid": null,
  "clean_step": {},
  "conductor_group": "",
  "console_enabled": false,
  "console_interface": "no-console",
  "created_at": "2019-01-31T19:59:28+00:00",
  "deploy_interface": "iscsi",
  "deploy_step": {},
  "driver": "ipmi",
  "driver_info": {
    "ipmi_port": "6230",
    "ipmi_username": "admin",
    "deploy_kernel": "http://172.22.0.1/images/tinyipa-stable-rocky.vmlinuz",
    "ipmi_address": "192.168.122.1",
    "deploy_ramdisk": "http://172.22.0.1/images/tinyipa-stable-rocky.gz",
    "ipmi_password": "admin"
  },
  "driver_internal_info": {},
  "extra": {},
  "fault": null,
  "firmware_interface": "no-firmware",
  "inspect_interface": "no-inspect",
  "inspection_finished_at": null,
  "inspection_started_at": null,
  "instance_info": {},
  "instance_uuid": null,
  "last_error": null,
  "links": [
    {
      "href": "http://ironic.example.com:6385/v1/nodes/d2630783-6ec8-4836-b556-ab427c4b581e",
      "rel": "self"
    },
    {
      "href": "http://ironic.example.com:6385/nodes/d2630783-6ec8-4836-b556-ab427c4b581e",
      "rel": "bookmark"
    }
  ],
  "maintenance": false,
  "maintenance_reason": null,
  "management_interface": "ipmitool",
  "name": "foo",
  "network_interface": "flat",
  "portgroups": [
    {
      "href": "http://ironic.example.com:6385/v1/nodes/d2630783-6ec8-4836-b556-ab427c4b581e/portgroups",
      "rel": "self"
    },
    {
      "href": "http://ironic.example.com:6385/nodes/d2630783-6ec8-4836-b556-ab427c4b581e/portgroups",
      "rel": "bookmark"
    }
  ],
  "ports": [
    {
      "href": "http://ironic.example.com:6385/v1/nodes/d2630783-6ec8-4836-b556-ab427c4b581e/ports",
      "rel": "self"
    },
    {
      "href": "http://ironic.example.com:6385/nodes/d2630783-6ec8-4836-b556-ab427c4b581e/ports",
      "rel": "bookmark"
    }
  ],
  "power_interface": "ipmitool",
  "power_state": null,
  "properties": {},
  "provision_state": "enroll",
  "provision_updated_at": "2019-02-15T17:21:29+00:00",
  "raid_config": {},
  "raid_interface": "no-raid",
  "retired": false,
  "retired_reason": "No longer needed",
  "rescue_interface": "no-rescue",
  "reservation": null,
  "resource_class": null,
  "states": [
    {
      "href": "http://ironic.example.com:6385/v1/nodes/d2630783-6ec8-4836-b556-ab427c4b581e/states",
      "rel": "self"
    },
    {
      "href": "http://ironic.example.com:6385/nodes/d2630783-6ec8-4836-b556-ab427c4b581e/states",
      "rel": "bookmark"
    }
  ],
  "storage_interface": "noop",
  "target_power_state": null,
  "target_provision_state": null,
  "target_raid_config": {},
  "traits": [],
  "updated_at": "2019-02-15T19:59:29+00:00",
  "uuid": "d2630783-6ec8-4836-b556-ab427c4b581e",
  "vendor_interface": "ipmitool",
  "volume": [
    {
      "href": "http://ironic.example.com:6385/v1/nodes/d2630783-6ec8-4836-b556-ab427c4b581e/volume",
      "rel": "self"
    },
    {
      "href": "http://ironic.example.com:6385/nodes/d2630783-6ec8-4836-b556-ab427c4b581e/volume",
      "rel": "bookmark"
    }
  ]
}
`

const NodeValidationBody = `
{
  "bios": {
    "reason": "Driver ipmi does not support bios (disabled or not implemented).",
    "result": false
  },
  "boot": {
    "reason": "Cannot validate image information for node a62b8495-52e2-407b-b3cb-62775d04c2b8 because one or more parameters are missing from its instance_info and insufficent information is present to boot from a remote volume. Missing are: ['ramdisk', 'kernel', 'image_source']",
    "result": false
  },
  "console": {
    "reason": "Driver ipmi does not support console (disabled or not implemented).",
    "result": false
  },
  "deploy": {
    "reason": "Cannot validate image information for node a62b8495-52e2-407b-b3cb-62775d04c2b8 because one or more parameters are missing from its instance_info and insufficent information is present to boot from a remote volume. Missing are: ['ramdisk', 'kernel', 'image_source']",
    "result": false
  },
  "firmware": {
    "reason": "Driver ipmi does not support firmware (disabled or not implemented).",
    "result": false
  },
  "inspect": {
    "reason": "Driver ipmi does not support inspect (disabled or not implemented).",
    "result": false
  },
  "management": {
    "result": true
  },
  "network": {
    "result": true
  },
  "power": {
    "result": true
  },
  "raid": {
    "reason": "Driver ipmi does not support raid (disabled or not implemented).",
    "result": false
  },
  "rescue": {
    "reason": "Driver ipmi does not support rescue (disabled or not implemented).",
    "result": false
  },
  "storage": {
    "result": true
  }
}
`

const NodeBootDeviceBody = `
{
  "boot_device":"pxe",
  "persistent":false
}
`

const NodeSupportedBootDeviceBody = `
{
  "supported_boot_devices": [
    "pxe",
    "disk"
  ]
}
`

const NodeProvisionStateActiveBody = `
{
    "target": "active",
    "configdrive": "http://127.0.0.1/images/test-node-config-drive.iso.gz"
}
`

const NodeProvisionStateActiveBodyWithSteps = `
{
    "target": "active",
    "deploy_steps": [
	{
	    "interface": "deploy",
	    "step": "inject_files",
	    "priority": 50,
	    "args": {
		"files": []
	    }
	}
    ]
}
`

const NodeProvisionStateCleanBody = `
{
    "target": "clean",
    "clean_steps": [
        {
            "interface": "deploy",
            "step": "upgrade_firmware",
            "args": {
                "force": "True"
            }
        }
    ]
}
`

const NodeProvisionStateConfigDriveBody = `
{
	"target": "active",
	"configdrive": {
		"user_data": {
			"ignition": {
			"version": "2.2.0"
			},
			"systemd": {
				"units": [
					{
						"enabled": true,
						"name": "example.service"
					}
				]
			}
		}
	}
}
`

const NodeProvisionStateServiceBody = `
{
    "target": "service",
    "service_steps": [
        {
            "interface": "bios",
            "step": "apply_configuration",
            "args": {
		"settings": []
            }
        }
    ]
}
`

const NodeBIOSSettingsBody = `
{
  "bios": [
   {
      "name": "Proc1L2Cache",
      "value": "10x256 KB"
   },
   {
      "name": "Proc1NumCores",
      "value": "10"
   },
   {
      "name": "ProcVirtualization",
      "value": "Enabled"
   }
   ]
}
`

const NodeDetailBIOSSettingsBody = `
{
  "bios": [
   {
      "created_at": "2021-05-11T21:33:44+00:00",
      "updated_at": null,
      "name": "Proc1L2Cache",
      "value": "10x256 KB",
      "attribute_type": "String",
      "allowable_values": [],
      "lower_bound": null,
      "max_length": 16,
      "min_length": 0,
      "read_only": true,
      "reset_required": null,
      "unique": null,
      "upper_bound": null,
      "links": [
        {
          "href": "http://ironic.example.com:6385/v1/nodes/d26115bf-1296-4ca8-8c86-6f310d8ec375/bios/Proc1L2Cache",
          "rel": "self"
        },
        {
          "href": "http://ironic.example.com:6385/nodes/d26115bf-1296-4ca8-8c86-6f310d8ec375/bios/Proc1L2Cache",
          "rel": "bookmark"
        }
      ]
   },
   {
      "created_at": "2021-05-11T21:33:44+00:00",
      "updated_at": null,
      "name": "Proc1NumCores",
      "value": "10",
      "attribute_type": "Integer",
      "allowable_values": [],
      "lower_bound": 0,
      "max_length": null,
      "min_length": null,
      "read_only": true,
      "reset_required": null,
      "unique": null,
      "upper_bound": 20,
      "links": [
        {
          "href": "http://ironic.example.com:6385/v1/nodes/d26115bf-1296-4ca8-8c86-6f310d8ec375/bios/Proc1NumCores",
          "rel": "self"
        },
        {
          "href": "http://ironic.example.com:6385/nodes/d26115bf-1296-4ca8-8c86-6f310d8ec375/bios/Proc1NumCores",
          "rel": "bookmark"
        }
      ]
   },
   {
      "created_at": "2021-05-11T21:33:44+00:00",
      "updated_at": null,
      "name": "ProcVirtualization",
      "value": "Enabled",
      "attribute_type": "Enumeration",
      "allowable_values": [
        "Enabled",
        "Disabled"
      ],
      "lower_bound": null,
      "max_length": null,
      "min_length": null,
      "read_only": false,
      "reset_required": null,
      "unique": null,
      "upper_bound": null,
      "links": [
        {
          "href": "http://ironic.example.com:6385/v1/nodes/d26115bf-1296-4ca8-8c86-6f310d8ec375/bios/ProcVirtualization",
          "rel": "self"
        },
        {
          "href": "http://ironic.example.com:6385/nodes/d26115bf-1296-4ca8-8c86-6f310d8ec375/bios/ProcVirtualization",
          "rel": "bookmark"
        }
      ]
   }
   ]
}
`

const NodeSingleBIOSSettingBody = `
{
  "Setting": {
      "name": "ProcVirtualization",
      "value": "Enabled"
   }
}
`

const NodeVendorPassthruMethodsBody = `
{
  "create_subscription": {
    "http_methods": [
      "POST"
    ],
    "async": false,
    "description": "",
    "attach": false,
    "require_exclusive_lock": true
  },
  "delete_subscription": {
    "http_methods": [
      "DELETE"
    ],
    "async": false,
    "description": "",
    "attach": false,
    "require_exclusive_lock": true
  },
  "get_subscription": {
    "http_methods": [
      "GET"
    ],
    "async": false,
    "description": "",
    "attach": false,
    "require_exclusive_lock": true
  },
  "get_all_subscriptions": {
    "http_methods": [
      "GET"
    ],
    "async": false,
    "description": "",
    "attach": false,
    "require_exclusive_lock": true
  }
}
`

const NodeGetAllSubscriptionsVnedorPassthruBody = `
{
  "@odata.context": "/redfish/v1/$metadata#EventDestinationCollection.EventDestinationCollection",
  "@odata.id": "/redfish/v1/EventService/Subscriptions",
  "@odata.type": "#EventDestinationCollection.EventDestinationCollection",
  "Description": "List of Event subscriptions",
  "Members": [
    {
      "@odata.id": "/redfish/v1/EventService/Subscriptions/62dbd1b6-f637-11eb-b551-4cd98f20754c"
    }
  ],
  "Members@odata.count": 1,
  "Name": "Event Subscriptions Collection"
}

`

const NodeGetSubscriptionVendorPassthruBody = `
{
  "Context": "Ironic",
  "Destination": "https://192.168.0.1/EventReceiver.php",
  "EventTypes": ["Alert"],
  "Id": "62dbd1b6-f637-11eb-b551-4cd98f20754c",
  "Protocol": "Redfish"
}
`

const NodeCreateSubscriptionVendorPassthruAllParametersBody = `
{
  "Context": "gophercloud",
  "Destination": "https://someurl",
  "EventTypes": ["Alert"],
  "HttpHeaders": [{"Context-Type":"application/json"}],
  "Id": "eaa43e2-018a-424e-990a-cbf47c62ef80",
  "Protocol": "Redfish"
}
`

const NodeCreateSubscriptionVendorPassthruRequiredParametersBody = `
{
  "Context": "",
  "Destination": "https://somedestinationurl",
  "EventTypes": ["Alert"],
  "Id": "344a3e2-978a-444e-990a-cbf47c62ef88",
  "Protocol": "Redfish"
}
`

const NodeSetMaintenanceBody = `
{
  "reason": "I'm tired"
}
`

var NodeInventoryBody = fmt.Sprintf(`
{
  "inventory": %s,
  "plugin_data":{
    "macs":[
      "52:54:00:90:35:d6"
    ],
    "local_gb":10,
    "cpu_arch":"x86_64",
    "memory_mb":2048
  }
}
`, inventorytest.InventorySample)

const NodeFirmwareListBody = `
{
  "firmware": [
   {
      "created_at": "2023-10-03T18:30:00+00:00",
      "updated_at": null,
      "component": "bios",
      "initial_version": "U30 v2.36 (07/16/2020)",
      "current_version": "U30 v2.36 (07/16/2020)",
      "last_version_flashed": null
   },
   {
      "created_at": "2023-10-03T18:30:00+00:00",
      "updated_at": "2023-10-03T18:45:54+00:00",
      "component": "bmc",
      "initial_version": "iLO 5 v2.78",
      "current_version": "iLO 5 v2.81",
      "last_version_flashed": "iLO 5 v2.81"
   }
   ]
}
`

const NodeVirtualMediaAttachBody = `
{
    "image_url": "https://example.com/image",
    "device_type": "cdrom"
}
`

const NodeVirtualMediaAttachBodyWithSource = `
{
    "image_url": "https://example.com/image",
    "device_type": "cdrom",
    "image_download_source": "http"
}
`

var (
	createdAtFoo, _      = time.Parse(time.RFC3339, "2019-01-31T19:59:28+00:00")
	createdAtBar, _      = time.Parse(time.RFC3339, "2019-01-31T19:59:29+00:00")
	createdAtBaz, _      = time.Parse(time.RFC3339, "2019-01-31T19:59:30+00:00")
	updatedAt, _         = time.Parse(time.RFC3339, "2019-02-15T19:59:29+00:00")
	provisonUpdatedAt, _ = time.Parse(time.RFC3339, "2019-02-15T17:21:29+00:00")

	NodeFoo = nodes.Node{
		UUID:                 "d2630783-6ec8-4836-b556-ab427c4b581e",
		Name:                 "foo",
		PowerState:           "",
		TargetPowerState:     "",
		ProvisionState:       "enroll",
		TargetProvisionState: "",
		Maintenance:          false,
		MaintenanceReason:    "",
		Fault:                "",
		LastError:            "",
		Reservation:          "",
		Driver:               "ipmi",
		DriverInfo: map[string]any{
			"ipmi_port":      "6230",
			"ipmi_username":  "admin",
			"deploy_kernel":  "http://172.22.0.1/images/tinyipa-stable-rocky.vmlinuz",
			"ipmi_address":   "192.168.122.1",
			"deploy_ramdisk": "http://172.22.0.1/images/tinyipa-stable-rocky.gz",
			"ipmi_password":  "admin",
		},
		DriverInternalInfo:  map[string]any{},
		Properties:          map[string]any{},
		InstanceInfo:        map[string]any{},
		InstanceUUID:        "",
		ChassisUUID:         "",
		Extra:               map[string]any{},
		ConsoleEnabled:      false,
		RAIDConfig:          map[string]any{},
		TargetRAIDConfig:    map[string]any{},
		CleanStep:           map[string]any{},
		DeployStep:          map[string]any{},
		ResourceClass:       "",
		BIOSInterface:       "no-bios",
		BootInterface:       "pxe",
		ConsoleInterface:    "no-console",
		DeployInterface:     "iscsi",
		FirmwareInterface:   "no-firmware",
		InspectInterface:    "no-inspect",
		ManagementInterface: "ipmitool",
		NetworkInterface:    "flat",
		PowerInterface:      "ipmitool",
		RAIDInterface:       "no-raid",
		RescueInterface:     "no-rescue",
		StorageInterface:    "noop",
		Traits:              []string{},
		VendorInterface:     "ipmitool",
		ConductorGroup:      "",
		Protected:           false,
		ProtectedReason:     "",
		CreatedAt:           createdAtFoo,
		UpdatedAt:           updatedAt,
		ProvisionUpdatedAt:  provisonUpdatedAt,
		Retired:             false,
		RetiredReason:       "No longer needed",
	}

	NodeFooValidation = nodes.NodeValidation{
		BIOS: nodes.DriverValidation{
			Result: false,
			Reason: "Driver ipmi does not support bios (disabled or not implemented).",
		},
		Boot: nodes.DriverValidation{
			Result: false,
			Reason: "Cannot validate image information for node a62b8495-52e2-407b-b3cb-62775d04c2b8 because one or more parameters are missing from its instance_info and insufficent information is present to boot from a remote volume. Missing are: ['ramdisk', 'kernel', 'image_source']",
		},
		Console: nodes.DriverValidation{
			Result: false,
			Reason: "Driver ipmi does not support console (disabled or not implemented).",
		},
		Deploy: nodes.DriverValidation{
			Result: false,
			Reason: "Cannot validate image information for node a62b8495-52e2-407b-b3cb-62775d04c2b8 because one or more parameters are missing from its instance_info and insufficent information is present to boot from a remote volume. Missing are: ['ramdisk', 'kernel', 'image_source']",
		},
		Firmware: nodes.DriverValidation{
			Result: false,
			Reason: "Driver ipmi does not support firmware (disabled or not implemented).",
		},
		Inspect: nodes.DriverValidation{
			Result: false,
			Reason: "Driver ipmi does not support inspect (disabled or not implemented).",
		},
		Management: nodes.DriverValidation{
			Result: true,
		},
		Network: nodes.DriverValidation{
			Result: true,
		},
		Power: nodes.DriverValidation{
			Result: true,
		},
		RAID: nodes.DriverValidation{
			Result: false,
			Reason: "Driver ipmi does not support raid (disabled or not implemented).",
		},
		Rescue: nodes.DriverValidation{
			Result: false,
			Reason: "Driver ipmi does not support rescue (disabled or not implemented).",
		},
		Storage: nodes.DriverValidation{
			Result: true,
		},
	}

	NodeBootDevice = nodes.BootDeviceOpts{
		BootDevice: "pxe",
		Persistent: false,
	}

	NodeSupportedBootDevice = []string{
		"pxe",
		"disk",
	}

	InspectionStartedAt  = time.Date(2023, time.February, 2, 14, 35, 59, 682403000, time.UTC)
	InspectionFinishedAt = time.Date(2023, time.February, 2, 14, 45, 59, 705249000, time.UTC)

	NodeBar = nodes.Node{
		UUID:                 "08c84581-58f5-4ea2-a0c6-dd2e5d2b3662",
		Name:                 "bar",
		PowerState:           "",
		TargetPowerState:     "",
		ProvisionState:       "available",
		TargetProvisionState: "",
		Maintenance:          false,
		MaintenanceReason:    "",
		Fault:                "",
		LastError:            "",
		Reservation:          "",
		Driver:               "ipmi",
		DriverInfo:           map[string]any{},
		DriverInternalInfo:   map[string]any{},
		Properties:           map[string]any{},
		InstanceInfo:         map[string]any{},
		InstanceUUID:         "",
		ChassisUUID:          "",
		Extra:                map[string]any{},
		ConsoleEnabled:       false,
		RAIDConfig:           map[string]any{},
		TargetRAIDConfig:     map[string]any{},
		CleanStep:            map[string]any{},
		DeployStep:           map[string]any{},
		ResourceClass:        "",
		BIOSInterface:        "no-bios",
		BootInterface:        "pxe",
		ConsoleInterface:     "no-console",
		DeployInterface:      "iscsi",
		FirmwareInterface:    "no-firmware",
		InspectInterface:     "no-inspect",
		ManagementInterface:  "ipmitool",
		NetworkInterface:     "flat",
		PowerInterface:       "ipmitool",
		RAIDInterface:        "no-raid",
		RescueInterface:      "no-rescue",
		StorageInterface:     "noop",
		Traits:               []string{},
		VendorInterface:      "ipmitool",
		ConductorGroup:       "",
		Protected:            false,
		ProtectedReason:      "",
		CreatedAt:            createdAtBar,
		UpdatedAt:            updatedAt,
		InspectionStartedAt:  &InspectionStartedAt,
		InspectionFinishedAt: &InspectionFinishedAt,
		Retired:              false,
		RetiredReason:        "No longer needed",
	}

	NodeBaz = nodes.Node{
		UUID:                 "c9afd385-5d89-4ecb-9e1c-68194da6b474",
		Name:                 "baz",
		PowerState:           "",
		TargetPowerState:     "",
		ProvisionState:       "enroll",
		TargetProvisionState: "",
		Maintenance:          false,
		MaintenanceReason:    "",
		Fault:                "",
		LastError:            "",
		Reservation:          "",
		Driver:               "ipmi",
		DriverInfo:           map[string]any{},
		DriverInternalInfo:   map[string]any{},
		Properties:           map[string]any{},
		InstanceInfo:         map[string]any{},
		InstanceUUID:         "",
		ChassisUUID:          "",
		Extra:                map[string]any{},
		ConsoleEnabled:       false,
		RAIDConfig:           map[string]any{},
		TargetRAIDConfig:     map[string]any{},
		CleanStep:            map[string]any{},
		DeployStep:           map[string]any{},
		ResourceClass:        "",
		BIOSInterface:        "no-bios",
		BootInterface:        "pxe",
		ConsoleInterface:     "no-console",
		DeployInterface:      "iscsi",
		FirmwareInterface:    "no-firmware",
		InspectInterface:     "no-inspect",
		ManagementInterface:  "ipmitool",
		NetworkInterface:     "flat",
		PowerInterface:       "ipmitool",
		RAIDInterface:        "no-raid",
		RescueInterface:      "no-rescue",
		StorageInterface:     "noop",
		Traits:               []string{},
		VendorInterface:      "ipmitool",
		ConductorGroup:       "",
		Protected:            false,
		ProtectedReason:      "",
		CreatedAt:            createdAtBaz,
		UpdatedAt:            updatedAt,
		Retired:              false,
		RetiredReason:        "No longer needed",
	}

	ConfigDriveMap = nodes.ConfigDrive{
		UserData: map[string]any{
			"ignition": map[string]string{
				"version": "2.2.0",
			},
			"systemd": map[string]any{
				"units": []map[string]any{{
					"name":    "example.service",
					"enabled": true,
				},
				},
			},
		},
	}

	NodeBIOSSettings = []nodes.BIOSSetting{
		{
			Name:  "Proc1L2Cache",
			Value: "10x256 KB",
		},
		{
			Name:  "Proc1NumCores",
			Value: "10",
		},
		{
			Name:  "ProcVirtualization",
			Value: "Enabled",
		},
	}

	iTrue      = true
	iFalse     = false
	minLength  = 0
	maxLength  = 16
	lowerBound = 0
	upperBound = 20

	NodeDetailBIOSSettings = []nodes.BIOSSetting{
		{
			Name:            "Proc1L2Cache",
			Value:           "10x256 KB",
			AttributeType:   "String",
			AllowableValues: []string{},
			LowerBound:      nil,
			UpperBound:      nil,
			MinLength:       &minLength,
			MaxLength:       &maxLength,
			ReadOnly:        &iTrue,
			ResetRequired:   nil,
			Unique:          nil,
		},
		{
			Name:            "Proc1NumCores",
			Value:           "10",
			AttributeType:   "Integer",
			AllowableValues: []string{},
			LowerBound:      &lowerBound,
			UpperBound:      &upperBound,
			MinLength:       nil,
			MaxLength:       nil,
			ReadOnly:        &iTrue,
			ResetRequired:   nil,
			Unique:          nil,
		},
		{
			Name:            "ProcVirtualization",
			Value:           "Enabled",
			AttributeType:   "Enumeration",
			AllowableValues: []string{"Enabled", "Disabled"},
			LowerBound:      nil,
			UpperBound:      nil,
			MinLength:       nil,
			MaxLength:       nil,
			ReadOnly:        &iFalse,
			ResetRequired:   nil,
			Unique:          nil,
		},
	}

	NodeSingleBIOSSetting = nodes.BIOSSetting{
		Name:  "ProcVirtualization",
		Value: "Enabled",
	}

	NodeVendorPassthruMethods = nodes.VendorPassthruMethods{
		CreateSubscription: nodes.CreateSubscriptionMethod{
			HTTPMethods:          []string{"POST"},
			Async:                false,
			Description:          "",
			Attach:               false,
			RequireExclusiveLock: true,
		},
		DeleteSubscription: nodes.DeleteSubscriptionMethod{
			HTTPMethods:          []string{"DELETE"},
			Async:                false,
			Description:          "",
			Attach:               false,
			RequireExclusiveLock: true,
		},
		GetSubscription: nodes.GetSubscriptionMethod{
			HTTPMethods:          []string{"GET"},
			Async:                false,
			Description:          "",
			Attach:               false,
			RequireExclusiveLock: true,
		},
		GetAllSubscriptions: nodes.GetAllSubscriptionsMethod{
			HTTPMethods:          []string{"GET"},
			Async:                false,
			Description:          "",
			Attach:               false,
			RequireExclusiveLock: true,
		},
	}

	NodeGetAllSubscriptions = nodes.GetAllSubscriptionsVendorPassthru{
		Context:      "/redfish/v1/$metadata#EventDestinationCollection.EventDestinationCollection",
		Etag:         "",
		Id:           "/redfish/v1/EventService/Subscriptions",
		Type:         "#EventDestinationCollection.EventDestinationCollection",
		Description:  "List of Event subscriptions",
		Name:         "Event Subscriptions Collection",
		Members:      []map[string]string{{"@odata.id": "/redfish/v1/EventService/Subscriptions/62dbd1b6-f637-11eb-b551-4cd98f20754c"}},
		MembersCount: 1,
	}

	NodeGetSubscription = nodes.SubscriptionVendorPassthru{
		Id:          "62dbd1b6-f637-11eb-b551-4cd98f20754c",
		Context:     "Ironic",
		Destination: "https://192.168.0.1/EventReceiver.php",
		EventTypes:  []string{"Alert"},
		Protocol:    "Redfish",
	}

	NodeCreateSubscriptionRequiredParameters = nodes.SubscriptionVendorPassthru{
		Id:          "344a3e2-978a-444e-990a-cbf47c62ef88",
		Context:     "",
		Destination: "https://somedestinationurl",
		EventTypes:  []string{"Alert"},
		Protocol:    "Redfish",
	}

	NodeCreateSubscriptionAllParameters = nodes.SubscriptionVendorPassthru{
		Id:          "eaa43e2-018a-424e-990a-cbf47c62ef80",
		Context:     "gophercloud",
		Destination: "https://someurl",
		EventTypes:  []string{"Alert"},
		Protocol:    "Redfish",
	}

	NodeInventoryData = nodes.InventoryData{
		Inventory: inventorytest.Inventory,
	}

	createdAtFirmware, _ = time.Parse(time.RFC3339, "2023-10-03T18:30:00+00:00")
	updatedAtFirmware, _ = time.Parse(time.RFC3339, "2023-10-03T18:45:54+00:00")
	lastVersion          = "iLO 5 v2.81"
	NodeFirmwareList     = []nodes.FirmwareComponent{
		{
			CreatedAt:          createdAtFirmware,
			UpdatedAt:          nil,
			Component:          "bios",
			InitialVersion:     "U30 v2.36 (07/16/2020)",
			CurrentVersion:     "U30 v2.36 (07/16/2020)",
			LastVersionFlashed: "",
		},
		{
			CreatedAt:          createdAtFirmware,
			UpdatedAt:          &updatedAtFirmware,
			Component:          "bmc",
			InitialVersion:     "iLO 5 v2.78",
			CurrentVersion:     "iLO 5 v2.81",
			LastVersionFlashed: lastVersion,
		},
	}
)

// HandleNodeListSuccessfully sets up the test server to respond to a server List request.
func HandleNodeListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/nodes", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Add("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}

		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, NodeListBody)

		case "9e5476bd-a4ec-4653-93d6-72c93aa682ba":
			fmt.Fprintf(w, `{ "servers": [] }`)
		default:
			t.Fatalf("/nodes invoked with unexpected marker=[%s]", marker)
		}
	})
}

// HandleNodeListSuccessfully sets up the test server to respond to a server List request.
func HandleNodeListDetailSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/nodes/detail", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Add("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}

		fmt.Fprintf(w, NodeListDetailBody)
	})
}

// HandleServerCreationSuccessfully sets up the test server to respond to a server creation request
// with a given response.
func HandleNodeCreationSuccessfully(t *testing.T, response string) {
	th.Mux.HandleFunc("/nodes", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `{
          "boot_interface": "pxe",
          "driver": "ipmi",
          "driver_info": {
            "deploy_kernel": "http://172.22.0.1/images/tinyipa-stable-rocky.vmlinuz",
            "deploy_ramdisk": "http://172.22.0.1/images/tinyipa-stable-rocky.gz",
            "ipmi_address": "192.168.122.1",
            "ipmi_password": "admin",
            "ipmi_port": "6230",
            "ipmi_username": "admin"
          },
          "firmware_interface": "no-firmware",
          "name": "foo"
        }`)

		w.WriteHeader(http.StatusAccepted)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, response)
	})
}

// HandleNodeDeletionSuccessfully sets up the test server to respond to a server deletion request.
func HandleNodeDeletionSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/nodes/asdfasdfasdf", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

func HandleNodeGetSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		fmt.Fprintf(w, SingleNodeBody)
	})
}

func HandleNodeUpdateSuccessfully(t *testing.T, response string) {
	th.Mux.HandleFunc("/nodes/1234asdf", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestJSONRequest(t, r, `[{"op": "replace", "path": "/properties", "value": {"root_gb": 25}}]`)

		fmt.Fprintf(w, response)
	})
}

func HandleNodeValidateSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/validate", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		fmt.Fprintf(w, NodeValidationBody)
	})
}

// HandleInjectNMISuccessfully sets up the test server to respond to a node InjectNMI request
func HandleInjectNMISuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/management/inject_nmi", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, "{}")

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleSetBootDeviceSuccessfully sets up the test server to respond to a set boot device request for a node
func HandleSetBootDeviceSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/management/boot_device", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, NodeBootDeviceBody)

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleGetBootDeviceSuccessfully sets up the test server to respond to a get boot device request for a node
func HandleGetBootDeviceSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/management/boot_device", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, NodeBootDeviceBody)
	})
}

// HandleGetBootDeviceSuccessfully sets up the test server to respond to a get boot device request for a node
func HandleGetSupportedBootDeviceSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/management/boot_device/supported", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, NodeSupportedBootDeviceBody)
	})
}

func HandleNodeChangeProvisionStateActive(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/states/provision", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, NodeProvisionStateActiveBody)
		w.WriteHeader(http.StatusAccepted)
	})
}

func HandleNodeChangeProvisionStateActiveWithSteps(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/states/provision", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, NodeProvisionStateActiveBodyWithSteps)
		w.WriteHeader(http.StatusAccepted)
	})
}

func HandleNodeChangeProvisionStateClean(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/states/provision", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, NodeProvisionStateCleanBody)
		w.WriteHeader(http.StatusAccepted)
	})
}

func HandleNodeChangeProvisionStateCleanWithConflict(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/states/provision", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, NodeProvisionStateCleanBody)
		w.WriteHeader(http.StatusConflict)
	})
}

func HandleNodeChangeProvisionStateConfigDrive(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/states/provision", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, NodeProvisionStateConfigDriveBody)
		w.WriteHeader(http.StatusAccepted)
	})
}

func HandleNodeChangeProvisionStateService(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/states/provision", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, NodeProvisionStateServiceBody)
		w.WriteHeader(http.StatusAccepted)
	})
}

// HandleChangePowerStateSuccessfully sets up the test server to respond to a change power state request for a node
func HandleChangePowerStateSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/states/power", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `{
			"target": "power on",
			"timeout": 100
		}`)

		w.WriteHeader(http.StatusAccepted)
	})
}

// HandleChangePowerStateWithConflict sets up the test server to respond to a change power state request for a node with a 409 error
func HandleChangePowerStateWithConflict(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/states/power", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `{
			"target": "power on",
			"timeout": 100
		}`)

		w.WriteHeader(http.StatusConflict)
	})
}

func HandleSetRAIDConfig(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/states/raid", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `
			{
			   "logical_disks" : [
				  {
					 "size_gb" : 100,
					 "is_root_volume" : true,
					 "raid_level" : "1"
				  }
			   ]
			}
		`)

		w.WriteHeader(http.StatusNoContent)
	})
}

func HandleSetRAIDConfigMaxSize(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/states/raid", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `
			{
			   "logical_disks" : [
				  {
					 "size_gb" : "MAX",
					 "is_root_volume" : true,
					 "raid_level" : "1"
				  }
			   ]
			}
		`)

		w.WriteHeader(http.StatusNoContent)
	})
}

func HandleListBIOSSettingsSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/bios", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		fmt.Fprintf(w, NodeBIOSSettingsBody)
	})
}

func HandleListDetailBIOSSettingsSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/bios", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		fmt.Fprintf(w, NodeDetailBIOSSettingsBody)
	})
}

func HandleGetBIOSSettingSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/bios/ProcVirtualization", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		fmt.Fprintf(w, NodeSingleBIOSSettingBody)
	})
}

func HandleGetVendorPassthruMethodsSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/vendor_passthru/methods", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		fmt.Fprintf(w, NodeVendorPassthruMethodsBody)
	})
}

func HandleGetAllSubscriptionsVendorPassthruSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/vendor_passthru", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestFormValues(t, r, map[string]string{"method": "get_all_subscriptions"})

		fmt.Fprintf(w, NodeGetAllSubscriptionsVnedorPassthruBody)
	})
}

func HandleGetSubscriptionVendorPassthruSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/vendor_passthru", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestFormValues(t, r, map[string]string{"method": "get_subscription"})
		th.TestJSONRequest(t, r, `
			{
			   "id" : "62dbd1b6-f637-11eb-b551-4cd98f20754c"
			}
		`)

		fmt.Fprintf(w, NodeGetSubscriptionVendorPassthruBody)
	})
}

func HandleCreateSubscriptionVendorPassthruAllParametersSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/vendor_passthru", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestFormValues(t, r, map[string]string{"method": "create_subscription"})
		th.TestJSONRequest(t, r, `
			{
         "Context":      "gophercloud",
         "EventTypes":   ["Alert"],
         "HttpHeaders":  [{"Content-Type":"application/json"}],
         "Protocol":     "Redfish",
			   "Destination" : "https://someurl"
			}
		`)

		fmt.Fprintf(w, NodeCreateSubscriptionVendorPassthruAllParametersBody)
	})
}

func HandleCreateSubscriptionVendorPassthruRequiredParametersSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/vendor_passthru", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestFormValues(t, r, map[string]string{"method": "create_subscription"})
		th.TestJSONRequest(t, r, `
			{
			   "Destination" : "https://somedestinationurl"
			}
		`)

		fmt.Fprintf(w, NodeCreateSubscriptionVendorPassthruRequiredParametersBody)
	})
}

func HandleDeleteSubscriptionVendorPassthruSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/vendor_passthru", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestFormValues(t, r, map[string]string{"method": "delete_subscription"})
		th.TestJSONRequest(t, r, `
			{
			   "id" : "344a3e2-978a-444e-990a-cbf47c62ef88"
			}
		`)

		w.WriteHeader(http.StatusNoContent)
	})
}

func HandleSetNodeMaintenanceSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/maintenance", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, NodeSetMaintenanceBody)

		w.WriteHeader(http.StatusAccepted)
	})
}

func HandleUnsetNodeMaintenanceSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/maintenance", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusAccepted)
	})
}

// HandleGetInventorySuccessfully sets up the test server to respond to a get inventory request for a node
func HandleGetInventorySuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/inventory", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, NodeInventoryBody)
	})
}

// HandleListFirmware
func HandleListFirmwareSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/nodes/1234asdf/firmware", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, NodeFirmwareListBody)
	})
}

func HandleAttachVirtualMediaSuccessfully(t *testing.T, withSource bool) {
	th.Mux.HandleFunc("/nodes/1234asdf/vmedia", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		if withSource {
			th.TestJSONRequest(t, r, NodeVirtualMediaAttachBodyWithSource)
		} else {
			th.TestJSONRequest(t, r, NodeVirtualMediaAttachBody)
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

func HandleDetachVirtualMediaSuccessfully(t *testing.T, withType bool) {
	th.Mux.HandleFunc("/nodes/1234asdf/vmedia", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		if withType {
			th.TestFormValues(t, r, map[string]string{"device_types": "cdrom"})
		} else {
			th.TestFormValues(t, r, map[string]string{})
		}
		w.WriteHeader(http.StatusNoContent)
	})
}
