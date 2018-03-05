package testing

// GetWithDHCPOptsResponse represents a raw port response with extra DHCP options.
const GetWithDHCPOptsResponse = `
{
    "port": {
        "status": "ACTIVE",
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "d6700c0c9ffa4f1cb322cd4a1f3906fa",
        "extra_dhcp_opts": [
            {
                "opt_name": "option1",
                "opt_value": "value1",
                "ip_version": 4
            },
            {
                "opt_name": "option2",
                "opt_value": "value2",
                "ip_version": 4
            }
        ],
        "admin_state_up": true,
        "name": "port-with-extra-dhcp-opts",
        "device_owner": "",
        "mac_address": "fa:16:3e:c9:cb:f0",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.4"
            }
        ],
        "id": "65c0ee9f-d634-4522-8954-51021b570b0d",
        "device_id": ""
    }
}
`
