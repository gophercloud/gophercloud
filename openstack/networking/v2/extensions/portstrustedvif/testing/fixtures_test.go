package testing

const ListResponse = `
{
    "ports": [
        {
            "status": "ACTIVE",
            "binding:host_id": "devstack",
            "name": "",
            "admin_state_up": true,
            "network_id": "70c1db1f-b701-45bd-96e0-a313ee3430b3",
            "tenant_id": "",
            "device_owner": "network:router_gateway",
            "mac_address": "fa:16:3e:58:42:ed",
            "binding:vnic_type": "normal",
            "fixed_ips": [
                {
                    "subnet_id": "008ba151-0b8c-4a67-98b5-0d2b87666062",
                    "ip_address": "172.24.4.2"
                }
            ],
            "id": "d80b1a3b-4fc1-49f3-952e-1e2ab7081d8b",
            "security_groups": [],
            "dns_name": "test-port",
            "dns_assignment": [
              {
                "hostname": "test-port",
                "ip_address": "172.24.4.2",
                "fqdn": "test-port.openstack.local."
              }
            ],
            "device_id": "9ae135f4-b6e0-4dad-9e91-3c223e385824",
            "port_security_enabled": false,
			"trusted": true,
            "created_at": "2019-06-30T04:15:37",
            "updated_at": "2019-06-30T05:18:49"
        }
    ]
}
`

const GetPortTrustedVIFStatusResponse = `
{
    "port": {
        "status": "DOWN",
        "name": "private-port",
        "admin_state_up": true,
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "d6700c0c9ffa4f1cb322cd4a1f3906fa",
        "device_owner": "",
        "mac_address": "fa:16:3e:c9:cb:f0",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.1"
            }
        ],
        "id": "46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2",
        "trusted": false,
        "device_id": ""
    }
}
`
const GetPortTrustedVIFUnsetResponse = `
{
    "port": {
        "status": "DOWN",
        "name": "private-port",
        "admin_state_up": true,
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "d6700c0c9ffa4f1cb322cd4a1f3906fa",
        "device_owner": "",
        "mac_address": "fa:16:3e:c9:cb:f0",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.1"
            }
        ],
        "id": "46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2",
        "device_id": ""
    }
}
`

const CreatePortTrustedVIFRequest = `
{
    "port": {
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "name": "private-port",
        "admin_state_up": true,
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.2"
            }
        ],
        "trusted": true
    }
}`

const CreatePortTrustedVIFResponse = `
{
    "port": {
        "status": "DOWN",
        "name": "private-port",
        "admin_state_up": true,
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "d6700c0c9ffa4f1cb322cd4a1f3906fa",
        "device_owner": "",
        "mac_address": "fa:16:3e:c9:cb:f0",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.2"
            }
        ],
        "id": "65c0ee9f-d634-4522-8954-51021b570b0d",
        "trusted": true,
        "device_id": ""
    }
}`

const UpdatePortTrustedVIFRequest = `
{
    "port": {
        "trusted": false
    }
}`

const UpdatePortTrustedVIFResponse = `
{
    "port": {
        "status": "DOWN",
        "name": "private-port",
        "admin_state_up": true,
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "d6700c0c9ffa4f1cb322cd4a1f3906fa",
        "device_owner": "",
        "mac_address": "fa:16:3e:c9:cb:f0",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.3"
            }
        ],
        "allowed_address_pairs": [
          {
            "ip_address": "10.0.0.4",
            "mac_address": "fa:16:3e:c9:cb:f0"
          }
        ],
        "id": "65c0ee9f-d634-4522-8954-51021b570b0d",
        "security_groups": [
            "f0ac4394-7e4a-4409-9701-ba8be283dbc3"
        ],
        "trusted": false,
        "device_id": ""
    }
}
`
