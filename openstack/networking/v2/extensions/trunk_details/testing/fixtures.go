package testing

// PortWithTrunkDetailsResult represents a raw server response from the
// Neutron API with trunk_details enabled.
// Some fields have been deleted from the response.
const PortWithTrunkDetailsResult = `
{
  "port": {
    "id": "dc3e8758-ee96-402d-94b0-4be5e9396c82",
    "name": "test-port-with-subports",
    "network_id": "42e996cb-6c9e-4cb1-8665-c62aa1610249",
    "tenant_id": "d4aa8944-e8be-4f46-bf93-74331af9c49e",
    "mac_address": "fa:16:3e:1f:de:6d",
    "admin_state_up": true,
    "status": "ACTIVE",
    "device_id": "935f1d9c-1888-457e-98d7-cb57405086cf",
    "device_owner": "compute:nova",
    "fixed_ips": [
      {
        "subnet_id": "f7aea11b-a649-4d23-995f-dcd4f2513f7e",
        "ip_address": "172.16.0.225"
      }
    ],
    "allowed_address_pairs": [],
    "extra_dhcp_opts": [],
    "security_groups": [
      "614f6c36-50b8-4dde-ab59-a46783befeec"
    ],
    "description": "",
    "binding:vnic_type": "normal",
    "qos_policy_id": null,
    "port_security_enabled": true,
    "trunk_details": {
      "trunk_id": "f170c831-8c55-4ceb-ad13-75eab4a121e5",
      "sub_ports": [
        {
          "segmentation_id": 100,
          "segmentation_type": "vlan",
          "port_id": "20c673d8-7f9d-4570-b662-148d9ddcc5bd",
          "mac_address": "fa:16:3e:88:29:a0"
        }
      ]
    },
    "ip_allocation": "immediate",
    "tags": [],
    "created_at": "2023-05-05T10:54:51Z",
    "updated_at": "2023-05-05T16:26:01Z",
    "revision_number": 4,
    "project_id": "d4aa8944-e8be-4f46-bf93-74331af9c49e"
  }
}
`
