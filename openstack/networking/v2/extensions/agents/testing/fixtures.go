package testing

import (
	"time"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/agents"
	// "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/bgp/speaker"
)

// AgentsListResult represents raw response for the List request.
const AgentsListResult = `
{
    "agents": [
        {
            "admin_state_up": true,
            "agent_type": "Open vSwitch agent",
            "alive": true,
            "availability_zone": null,
            "binary": "neutron-openvswitch-agent",
            "configurations": {
                "datapath_type": "system",
                "extensions": [
                    "qos"
                ]
            },
            "created_at": "2017-07-26 23:15:44",
            "description": null,
            "heartbeat_timestamp": "2019-01-09 10:28:53",
            "host": "compute1",
            "id": "59186d7b-b512-4fdf-bbaf-5804ffde8811",
            "started_at": "2018-06-26 21:46:19",
            "topic": "N/A"
        },
        {
            "admin_state_up": true,
            "agent_type": "Open vSwitch agent",
            "alive": true,
            "availability_zone": null,
            "binary": "neutron-openvswitch-agent",
            "configurations": {
                "datapath_type": "system",
                "extensions": [
                    "qos"
                ]
            },
            "created_at": "2017-01-22 14:00:50",
            "description": null,
            "heartbeat_timestamp": "2019-01-09 10:28:50",
            "host": "compute2",
            "id": "76af7b1f-d61b-4526-94f7-d2e14e2698df",
            "started_at": "2018-11-06 12:09:17",
            "topic": "N/A"
        }
    ]
}
`

// AgentUpdateRequest represents raw request to update an Agent.
const AgentUpdateRequest = `
{
    "agent": {
        "description": "My OVS agent for OpenStack",
        "admin_state_up": true
    }
}
`

// Agent represents a sample Agent struct.
var Agent = agents.Agent{
	ID:              "43583cf5-472e-4dc8-af5b-6aed4c94ee3a",
	AdminStateUp:    true,
	AgentType:       "Open vSwitch agent",
	Description:     "My OVS agent for OpenStack",
	Alive:           true,
	ResourcesSynced: true,
	Binary:          "neutron-openvswitch-agent",
	Configurations: map[string]interface{}{
		"ovs_hybrid_plug":            false,
		"datapath_type":              "system",
		"vhostuser_socket_dir":       "/var/run/openvswitch",
		"log_agent_heartbeats":       false,
		"l2_population":              true,
		"enable_distributed_routing": false,
	},
	CreatedAt:          time.Date(2017, 7, 26, 23, 2, 5, 0, time.UTC),
	StartedAt:          time.Date(2018, 6, 26, 21, 46, 20, 0, time.UTC),
	HeartbeatTimestamp: time.Date(2019, 1, 9, 11, 43, 01, 0, time.UTC),
	Host:               "compute3",
	Topic:              "N/A",
}

// Agent1 represents first unmarshalled address scope from the
// AgentsListResult.
var Agent1 = agents.Agent{
	ID:           "59186d7b-b512-4fdf-bbaf-5804ffde8811",
	AdminStateUp: true,
	AgentType:    "Open vSwitch agent",
	Alive:        true,
	Binary:       "neutron-openvswitch-agent",
	Configurations: map[string]interface{}{
		"datapath_type": "system",
		"extensions": []interface{}{
			"qos",
		},
	},
	CreatedAt:          time.Date(2017, 7, 26, 23, 15, 44, 0, time.UTC),
	StartedAt:          time.Date(2018, 6, 26, 21, 46, 19, 0, time.UTC),
	HeartbeatTimestamp: time.Date(2019, 1, 9, 10, 28, 53, 0, time.UTC),
	Host:               "compute1",
	Topic:              "N/A",
}

// Agent2 represents second unmarshalled address scope from the
// AgentsListResult.
var Agent2 = agents.Agent{
	ID:           "76af7b1f-d61b-4526-94f7-d2e14e2698df",
	AdminStateUp: true,
	AgentType:    "Open vSwitch agent",
	Alive:        true,
	Binary:       "neutron-openvswitch-agent",
	Configurations: map[string]interface{}{
		"datapath_type": "system",
		"extensions": []interface{}{
			"qos",
		},
	},
	CreatedAt:          time.Date(2017, 1, 22, 14, 00, 50, 0, time.UTC),
	StartedAt:          time.Date(2018, 11, 6, 12, 9, 17, 0, time.UTC),
	HeartbeatTimestamp: time.Date(2019, 1, 9, 10, 28, 50, 0, time.UTC),
	Host:               "compute2",
	Topic:              "N/A",
}

// AgentsGetResult represents raw response for the Get request.
const AgentsGetResult = `
{
    "agent": {
        "binary": "neutron-openvswitch-agent",
        "description": null,
        "availability_zone": null,
        "heartbeat_timestamp": "2019-01-09 11:43:01",
        "admin_state_up": true,
        "alive": true,
        "id": "43583cf5-472e-4dc8-af5b-6aed4c94ee3a",
        "topic": "N/A",
        "host": "compute3",
        "agent_type": "Open vSwitch agent",
        "started_at": "2018-06-26 21:46:20",
        "created_at": "2017-07-26 23:02:05",
        "configurations": {
            "ovs_hybrid_plug": false,
            "datapath_type": "system",
            "vhostuser_socket_dir": "/var/run/openvswitch",
            "log_agent_heartbeats": false,
            "l2_population": true,
            "enable_distributed_routing": false
        }
    }
}
`

// AgentsUpdateResult represents raw response for the Update request.
const AgentsUpdateResult = `
{
    "agent": {
        "binary": "neutron-openvswitch-agent",
        "description": "My OVS agent for OpenStack",
        "availability_zone": null,
        "heartbeat_timestamp": "2019-01-09 11:43:01",
        "admin_state_up": true,
        "alive": true,
        "id": "43583cf5-472e-4dc8-af5b-6aed4c94ee3a",
        "topic": "N/A",
        "host": "compute3",
        "agent_type": "Open vSwitch agent",
        "started_at": "2018-06-26 21:46:20",
        "created_at": "2017-07-26 23:02:05",
	"resources_synced": true,
        "configurations": {
            "ovs_hybrid_plug": false,
            "datapath_type": "system",
            "vhostuser_socket_dir": "/var/run/openvswitch",
            "log_agent_heartbeats": false,
            "l2_population": true,
            "enable_distributed_routing": false
        }
    }
}
`

// AgentDHCPNetworksListResult represents raw response for the ListDHCPNetworks request.
const AgentDHCPNetworksListResult = `
{
    "networks": [
        {
            "admin_state_up": true,
            "availability_zone_hints": [],
            "availability_zones": [
                "nova"
            ],
            "created_at": "2016-03-08T20:19:41",
            "dns_domain": "my-domain.org.",
            "id": "d32019d3-bc6e-4319-9c1d-6722fc136a22",
            "ipv4_address_scope": null,
            "ipv6_address_scope": null,
            "l2_adjacency": false,
            "mtu": 1500,
            "name": "net1",
            "port_security_enabled": true,
            "project_id": "4fd44f30292945e481c7b8a0c8908869",
            "qos_policy_id": "6a8454ade84346f59e8d40665f878b2e",
            "revision_number": 1,
            "router:external": false,
            "shared": false,
            "status": "ACTIVE",
            "subnets": [
                "54d6f61d-db07-451c-9ab3-b9609b6b6f0b"
            ],
            "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
            "updated_at": "2016-03-08T20:19:41",
            "vlan_transparent": true,
            "description": "",
            "is_default": false
        }
    ]
}
`

// ScheduleDHCPNetworkRequest represents raw request for the ScheduleDHCPNetwork request.
const ScheduleDHCPNetworkRequest = `
{
    "network_id": "1ae075ca-708b-4e66-b4a7-b7698632f05f"
}
`

const ListBGPSpeakersResult = `
{
  "bgp_speakers": [
    {
      "peers": [
        "cc4e1b15-e8b1-415e-b39a-3b087ed567b4",
        "4022d79f-835e-4271-b5d1-d90dce5662df"
      ],
      "project_id": "89f56d77-fee7-4b2f-8b1e-583717a93690",
      "name": "gophercloud-testing-speaker",
      "tenant_id": "5c372f0b-051e-485c-a82c-9dd732e7df83",
      "local_as": 12345,
      "advertise_tenant_networks": true,
      "networks": [
        "932d70b1-db21-4542-b520-d5e73ddee407"
      ],
      "ip_version": 4,
      "advertise_floating_ip_host_routes": true,
      "id": "cab00464-284d-4251-9798-2b27db7b1668"
    }
  ]
}
`

const ListDRAgentHostingBGPSpeakersResult = `
{
  "agents": [
    {
      "binary": "neutron-bgp-dragent",
      "description": null,
      "availability_zone": null,
      "heartbeat_timestamp": "2021-09-13 19:55:01",
      "admin_state_up": true,
      "resources_synced": null,
      "alive": true,
      "topic": "bgp_dragent",
      "host": "agent1.example.com",
      "agent_type": "BGP dynamic routing agent",
      "resource_versions": {},
      "created_at": "2020-09-17 20:08:58",
      "started_at": "2021-05-04 11:13:12",
      "id": "60d78b78-b56b-4d91-a174-2c03159f6bb6",
      "configurations": {
        "advertise_routes": 2,
        "bgp_peers": 2,
        "bgp_speakers": 1
      }
    },
    {
      "binary": "neutron-bgp-dragent",
      "description": null,
      "availability_zone": null,
      "heartbeat_timestamp": "2021-09-13 19:54:47",
      "admin_state_up": true,
      "resources_synced": null,
      "alive": true,
      "topic": "bgp_dragent",
      "host": "agent2.example.com",
      "agent_type": "BGP dynamic routing agent",
      "resource_versions": {},
      "created_at": "2020-09-17 20:08:15",
      "started_at": "2021-05-04 11:13:13",
      "id": "d0bdcea2-1d02-4c1d-9e79-b827e77acc22",
      "configurations": {
        "advertise_routes": 2,
        "bgp_peers": 2,
        "bgp_speakers": 1
      }
    }
  ]
}
`

var BGPAgent1 = agents.Agent{
	ID:           "60d78b78-b56b-4d91-a174-2c03159f6bb6",
	AdminStateUp: true,
	AgentType:    "BGP dynamic routing agent",
	Alive:        true,
	Binary:       "neutron-bgp-dragent",
	Configurations: map[string]interface{}{
		"advertise_routes": float64(2),
		"bgp_peers":        float64(2),
		"bgp_speakers":     float64(1),
	},
	CreatedAt:          time.Date(2020, 9, 17, 20, 8, 58, 0, time.UTC),
	StartedAt:          time.Date(2021, 5, 4, 11, 13, 12, 0, time.UTC),
	HeartbeatTimestamp: time.Date(2021, 9, 13, 19, 55, 1, 0, time.UTC),
	Host:               "agent1.example.com",
	Topic:              "bgp_dragent",
}

var BGPAgent2 = agents.Agent{
	ID:           "d0bdcea2-1d02-4c1d-9e79-b827e77acc22",
	AdminStateUp: true,
	AgentType:    "BGP dynamic routing agent",
	Alive:        true,
	Binary:       "neutron-bgp-dragent",
	Configurations: map[string]interface{}{
		"advertise_routes": float64(2),
		"bgp_peers":        float64(2),
		"bgp_speakers":     float64(1),
	},
	CreatedAt:          time.Date(2020, 9, 17, 20, 8, 15, 0, time.UTC),
	StartedAt:          time.Date(2021, 5, 4, 11, 13, 13, 0, time.UTC),
	HeartbeatTimestamp: time.Date(2021, 9, 13, 19, 54, 47, 0, time.UTC),
	Host:               "agent2.example.com",
	Topic:              "bgp_dragent",
}

const ScheduleBGPSpeakerRequest = `
{
    "bgp_speaker_id": "8edb2c68-0654-49a9-b3fe-030f92e3ddf6"
}
`
