package servers

import (
	"encoding/json"
	"testing"
)

// Taken from: http://docs.openstack.org/api/openstack-compute/2/content/List_Servers-d1e2078.html
const goodListServersResult = `
{
	"servers": [
		{
			"id": "52415800-8b69-11e0-9b19-734f6af67565",
			"tenant_id": "1234",
			"user_id": "5678",
			"name": "sample-server",
			"updated": "2010-10-10T12:00:00Z",
			"created": "2010-08-10T12:00:00Z",
			"hostId": "e4d909c290d0fb1ca068ffaddf22cbd0",
			"status": "BUILD",
			"progress": 60,
			"accessIPv4": "67.23.10.132",
			"accessIPv6": "::babe:67.23.10.132",
			"image": {
				"id": "52415800-8b69-11e0-9b19-734f6f006e54",
				"links": [
					{
						"rel": "self",
						"href": "http://servers.api.openstack.org/v2/1234/images/52415800-8b69-11e0-9b19-734f6f006e54"
					},
					{
						"rel": "bookmark",
						"href": "http://servers.api.openstack.org/1234/images/52415800-8b69-11e0-9b19-734f6f006e54"
					}
				]
			},
			"flavor": {
				"id": "52415800-8b69-11e0-9b19-734f216543fd",
				"links": [
					{
						"rel": "self",
						"href": "http://servers.api.openstack.org/v2/1234/flavors/52415800-8b69-11e0-9b19-734f216543fd"
					},
					{
						"rel": "bookmark",
						"href": "http://servers.api.openstack.org/1234/flavors/52415800-8b69-11e0-9b19-734f216543fd"
					}
				]
			},
			"addresses": {
				"public": [
					{
						"version": 4,
						"addr": "67.23.10.132"
					},
					{
						"version": 6,
						"addr": "::babe:67.23.10.132"
					},
					{
						"version": 4,
						"addr": "67.23.10.131"
					},
					{
						"version": 6,
						"addr": "::babe:4317:0A83"
					}
				],
				"private": [
					{
						"version": 4,
						"addr": "10.176.42.16"
					},
					{
						"version": 6,
						"addr": "::babe:10.176.42.16"
					}
				]
			},
			"metadata": {
				"Server Label": "Web Head 1",
				"Image Version": "2.1"
			},
			"links": [
				{
					"rel": "self",
					"href": "http://servers.api.openstack.org/v2/1234/servers/52415800-8b69-11e0-9b19-734f6af67565"
				},
				{
					"rel": "bookmark",
					"href": "http://servers.api.openstack.org/1234/servers/52415800-8b69-11e0-9b19-734f6af67565"
				}
			]
		},
		{
			"id": "52415800-8b69-11e0-9b19-734f1f1350e5",
			"user_id": "5678",
			"name": "sample-server2",
			"tenant_id": "1234",
			"updated": "2010-10-10T12:00:00Z",
			"created": "2010-08-10T12:00:00Z",
			"hostId": "9e107d9d372bb6826bd81d3542a419d6",
			"status": "ACTIVE",
			"accessIPv4": "67.23.10.133",
			"accessIPv6": "::babe:67.23.10.133",
			"image": {
				"id": "52415800-8b69-11e0-9b19-734f5736d2a2",
				"links": [
					{
						"rel": "self",
						"href": "http://servers.api.openstack.org/v2/1234/images/52415800-8b69-11e0-9b19-734f5736d2a2"
					},
					{
						"rel": "bookmark",
						"href": "http://servers.api.openstack.org/1234/images/52415800-8b69-11e0-9b19-734f5736d2a2"
					}
				]
			},
			"flavor": {
				"id": "52415800-8b69-11e0-9b19-734f216543fd",
				"links": [
					{
						"rel": "self",
						"href": "http://servers.api.openstack.org/v2/1234/flavors/52415800-8b69-11e0-9b19-734f216543fd"
					},
					{
						"rel": "bookmark",
						"href": "http://servers.api.openstack.org/1234/flavors/52415800-8b69-11e0-9b19-734f216543fd"
					}
				]
			},
			"addresses": {
				"public": [
					{
						"version": 4,
						"addr": "67.23.10.133"
					},
					{
						"version": 6,
						"addr": "::babe:67.23.10.133"
					}
				],
				"private": [
					{
						"version": 4,
						"addr": "10.176.42.17"
					},
					{
						"version": 6,
						"addr": "::babe:10.176.42.17"
					}
				]
			},
			"metadata": {
				"Server Label": "DB 1"
			},
			"links": [
				{
					"rel": "self",
					"href": "http://servers.api.openstack.org/v2/1234/servers/52415800-8b69-11e0-9b19-734f1f1350e5"
				},
				{
					"rel": "bookmark",
					"href": "http://servers.api.openstack.org/1234/servers/52415800-8b69-11e0-9b19-734f1f1350e5"
				}
			]
		}
	]
}
`

func TestExtractServers(t *testing.T) {
	var listPage ListPage
	err := json.Unmarshal([]byte(goodListServersResult), &listPage.MarkerPageBase.LastHTTPResponse.Body)
	if err != nil {
		t.Fatalf("Error decoding JSON fixture: %v", err)
	}

	svrs, err := ExtractServers(listPage)
	if err != nil {
		t.Fatalf("Error extracting servers: %v", err)
	}

	if len(svrs) != 2 {
		t.Fatalf("Expected 2 servers; got %d", len(svrs))
	}
}
