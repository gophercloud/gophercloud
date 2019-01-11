package testing

import (
	"time"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/agents"
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
