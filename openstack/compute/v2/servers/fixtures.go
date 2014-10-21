// +build fixtures
package servers

// ServerListBody contains the canned body of a servers.List response.
const ServerListBody = `
{
	"servers": [
		{
			"status": "ACTIVE",
			"updated": "2014-09-25T13:10:10Z",
			"hostId": "29d3c8c896a45aa4c34e52247875d7fefc3d94bbcc9f622b5d204362",
			"OS-EXT-SRV-ATTR:host": "devstack",
			"addresses": {
				"private": [
					{
						"OS-EXT-IPS-MAC:mac_addr": "fa:16:3e:7c:1b:2b",
						"version": 4,
						"addr": "10.0.0.32",
						"OS-EXT-IPS:type": "fixed"
					}
				]
			},
			"links": [
				{
					"href": "http://104.130.131.164:8774/v2/fcad67a6189847c4aecfa3c81a05783b/servers/ef079b0c-e610-4dfb-b1aa-b49f07ac48e5",
					"rel": "self"
				},
				{
					"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/servers/ef079b0c-e610-4dfb-b1aa-b49f07ac48e5",
					"rel": "bookmark"
				}
			],
			"key_name": null,
			"image": {
				"id": "f90f6034-2570-4974-8351-6b49732ef2eb",
				"links": [
					{
						"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/images/f90f6034-2570-4974-8351-6b49732ef2eb",
						"rel": "bookmark"
					}
				]
			},
			"OS-EXT-STS:task_state": null,
			"OS-EXT-STS:vm_state": "active",
			"OS-EXT-SRV-ATTR:instance_name": "instance-0000001e",
			"OS-SRV-USG:launched_at": "2014-09-25T13:10:10.000000",
			"OS-EXT-SRV-ATTR:hypervisor_hostname": "devstack",
			"flavor": {
				"id": "1",
				"links": [
					{
						"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/flavors/1",
						"rel": "bookmark"
					}
				]
			},
			"id": "ef079b0c-e610-4dfb-b1aa-b49f07ac48e5",
			"security_groups": [
				{
					"name": "default"
				}
			],
			"OS-SRV-USG:terminated_at": null,
			"OS-EXT-AZ:availability_zone": "nova",
			"user_id": "9349aff8be7545ac9d2f1d00999a23cd",
			"name": "herp",
			"created": "2014-09-25T13:10:02Z",
			"tenant_id": "fcad67a6189847c4aecfa3c81a05783b",
			"OS-DCF:diskConfig": "MANUAL",
			"os-extended-volumes:volumes_attached": [],
			"accessIPv4": "",
			"accessIPv6": "",
			"progress": 0,
			"OS-EXT-STS:power_state": 1,
			"config_drive": "",
			"metadata": {}
		},
		{
			"status": "ACTIVE",
			"updated": "2014-09-25T13:04:49Z",
			"hostId": "29d3c8c896a45aa4c34e52247875d7fefc3d94bbcc9f622b5d204362",
			"OS-EXT-SRV-ATTR:host": "devstack",
			"addresses": {
				"private": [
					{
						"OS-EXT-IPS-MAC:mac_addr": "fa:16:3e:9e:89:be",
						"version": 4,
						"addr": "10.0.0.31",
						"OS-EXT-IPS:type": "fixed"
					}
				]
			},
			"links": [
				{
					"href": "http://104.130.131.164:8774/v2/fcad67a6189847c4aecfa3c81a05783b/servers/9e5476bd-a4ec-4653-93d6-72c93aa682ba",
					"rel": "self"
				},
				{
					"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/servers/9e5476bd-a4ec-4653-93d6-72c93aa682ba",
					"rel": "bookmark"
				}
			],
			"key_name": null,
			"image": {
				"id": "f90f6034-2570-4974-8351-6b49732ef2eb",
				"links": [
					{
						"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/images/f90f6034-2570-4974-8351-6b49732ef2eb",
						"rel": "bookmark"
					}
				]
			},
			"OS-EXT-STS:task_state": null,
			"OS-EXT-STS:vm_state": "active",
			"OS-EXT-SRV-ATTR:instance_name": "instance-0000001d",
			"OS-SRV-USG:launched_at": "2014-09-25T13:04:49.000000",
			"OS-EXT-SRV-ATTR:hypervisor_hostname": "devstack",
			"flavor": {
				"id": "1",
				"links": [
					{
						"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/flavors/1",
						"rel": "bookmark"
					}
				]
			},
			"id": "9e5476bd-a4ec-4653-93d6-72c93aa682ba",
			"security_groups": [
				{
					"name": "default"
				}
			],
			"OS-SRV-USG:terminated_at": null,
			"OS-EXT-AZ:availability_zone": "nova",
			"user_id": "9349aff8be7545ac9d2f1d00999a23cd",
			"name": "derp",
			"created": "2014-09-25T13:04:41Z",
			"tenant_id": "fcad67a6189847c4aecfa3c81a05783b",
			"OS-DCF:diskConfig": "MANUAL",
			"os-extended-volumes:volumes_attached": [],
			"accessIPv4": "",
			"accessIPv6": "",
			"progress": 0,
			"OS-EXT-STS:power_state": 1,
			"config_drive": "",
			"metadata": {}
		}
	]
}
`

// SingleServerBody is the canned body of a Get request on an existing server.
const SingleServerBody = `
{
	"server": {
		"status": "ACTIVE",
		"updated": "2014-09-25T13:04:49Z",
		"hostId": "29d3c8c896a45aa4c34e52247875d7fefc3d94bbcc9f622b5d204362",
		"OS-EXT-SRV-ATTR:host": "devstack",
		"addresses": {
			"private": [
				{
					"OS-EXT-IPS-MAC:mac_addr": "fa:16:3e:9e:89:be",
					"version": 4,
					"addr": "10.0.0.31",
					"OS-EXT-IPS:type": "fixed"
				}
			]
		},
		"links": [
			{
				"href": "http://104.130.131.164:8774/v2/fcad67a6189847c4aecfa3c81a05783b/servers/9e5476bd-a4ec-4653-93d6-72c93aa682ba",
				"rel": "self"
			},
			{
				"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/servers/9e5476bd-a4ec-4653-93d6-72c93aa682ba",
				"rel": "bookmark"
			}
		],
		"key_name": null,
		"image": {
			"id": "f90f6034-2570-4974-8351-6b49732ef2eb",
			"links": [
				{
					"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/images/f90f6034-2570-4974-8351-6b49732ef2eb",
					"rel": "bookmark"
				}
			]
		},
		"OS-EXT-STS:task_state": null,
		"OS-EXT-STS:vm_state": "active",
		"OS-EXT-SRV-ATTR:instance_name": "instance-0000001d",
		"OS-SRV-USG:launched_at": "2014-09-25T13:04:49.000000",
		"OS-EXT-SRV-ATTR:hypervisor_hostname": "devstack",
		"flavor": {
			"id": "1",
			"links": [
				{
					"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/flavors/1",
					"rel": "bookmark"
				}
			]
		},
		"id": "9e5476bd-a4ec-4653-93d6-72c93aa682ba",
		"security_groups": [
			{
				"name": "default"
			}
		],
		"OS-SRV-USG:terminated_at": null,
		"OS-EXT-AZ:availability_zone": "nova",
		"user_id": "9349aff8be7545ac9d2f1d00999a23cd",
		"name": "derp",
		"created": "2014-09-25T13:04:41Z",
		"tenant_id": "fcad67a6189847c4aecfa3c81a05783b",
		"OS-DCF:diskConfig": "MANUAL",
		"os-extended-volumes:volumes_attached": [],
		"accessIPv4": "",
		"accessIPv6": "",
		"progress": 0,
		"OS-EXT-STS:power_state": 1,
		"config_drive": "",
		"metadata": {}
	}
}
`

var (
	// ServerHerp is a Server struct that should correspond to the first result in ServerListBody.
	ServerHerp = Server{
		Status:  "ACTIVE",
		Updated: "2014-09-25T13:10:10Z",
		HostID:  "29d3c8c896a45aa4c34e52247875d7fefc3d94bbcc9f622b5d204362",
		Addresses: map[string]interface{}{
			"private": []interface{}{
				map[string]interface{}{
					"OS-EXT-IPS-MAC:mac_addr": "fa:16:3e:7c:1b:2b",
					"version":                 float64(4),
					"addr":                    "10.0.0.32",
					"OS-EXT-IPS:type":         "fixed",
				},
			},
		},
		Links: []interface{}{
			map[string]interface{}{
				"href": "http://104.130.131.164:8774/v2/fcad67a6189847c4aecfa3c81a05783b/servers/ef079b0c-e610-4dfb-b1aa-b49f07ac48e5",
				"rel":  "self",
			},
			map[string]interface{}{
				"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/servers/ef079b0c-e610-4dfb-b1aa-b49f07ac48e5",
				"rel":  "bookmark",
			},
		},
		Image: map[string]interface{}{
			"id": "f90f6034-2570-4974-8351-6b49732ef2eb",
			"links": []interface{}{
				map[string]interface{}{
					"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/images/f90f6034-2570-4974-8351-6b49732ef2eb",
					"rel":  "bookmark",
				},
			},
		},
		Flavor: map[string]interface{}{
			"id": "1",
			"links": []interface{}{
				map[string]interface{}{
					"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/flavors/1",
					"rel":  "bookmark",
				},
			},
		},
		ID:       "ef079b0c-e610-4dfb-b1aa-b49f07ac48e5",
		UserID:   "9349aff8be7545ac9d2f1d00999a23cd",
		Name:     "herp",
		Created:  "2014-09-25T13:10:02Z",
		TenantID: "fcad67a6189847c4aecfa3c81a05783b",
		Metadata: map[string]interface{}{},
	}

	// ServerDerp is a Server struct that should correspond to the second server in ServerListBody.
	ServerDerp = Server{
		Status:  "ACTIVE",
		Updated: "2014-09-25T13:04:49Z",
		HostID:  "29d3c8c896a45aa4c34e52247875d7fefc3d94bbcc9f622b5d204362",
		Addresses: map[string]interface{}{
			"private": []interface{}{
				map[string]interface{}{
					"OS-EXT-IPS-MAC:mac_addr": "fa:16:3e:9e:89:be",
					"version":                 float64(4),
					"addr":                    "10.0.0.31",
					"OS-EXT-IPS:type":         "fixed",
				},
			},
		},
		Links: []interface{}{
			map[string]interface{}{
				"href": "http://104.130.131.164:8774/v2/fcad67a6189847c4aecfa3c81a05783b/servers/9e5476bd-a4ec-4653-93d6-72c93aa682ba",
				"rel":  "self",
			},
			map[string]interface{}{
				"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/servers/9e5476bd-a4ec-4653-93d6-72c93aa682ba",
				"rel":  "bookmark",
			},
		},
		Image: map[string]interface{}{
			"id": "f90f6034-2570-4974-8351-6b49732ef2eb",
			"links": []interface{}{
				map[string]interface{}{
					"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/images/f90f6034-2570-4974-8351-6b49732ef2eb",
					"rel":  "bookmark",
				},
			},
		},
		Flavor: map[string]interface{}{
			"id": "1",
			"links": []interface{}{
				map[string]interface{}{
					"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/flavors/1",
					"rel":  "bookmark",
				},
			},
		},
		ID:       "9e5476bd-a4ec-4653-93d6-72c93aa682ba",
		UserID:   "9349aff8be7545ac9d2f1d00999a23cd",
		Name:     "derp",
		Created:  "2014-09-25T13:04:41Z",
		TenantID: "fcad67a6189847c4aecfa3c81a05783b",
		Metadata: map[string]interface{}{},
	}
)
