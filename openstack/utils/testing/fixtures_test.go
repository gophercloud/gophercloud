package testing

import (
	"fmt"
	"net/http"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func setupIdentityVersionHandler(fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			{
				"versions": {
					"values": [
						{
							"status": "stable",
							"id": "v3.0",
							"links": [
								{ "href": "%s/v3.0", "rel": "self" }
							]
						},
						{
							"status": "stable",
							"id": "v2.0",
							"links": [
								{ "href": "%s/v2.0", "rel": "self" }
							]
						}
					]
				}
			}
		`, fakeServer.Server.URL, fakeServer.Server.URL)
	})
}

func setupMultiServiceVersionHandler(fakeServer th.FakeServer) {
	// Compute v2 API
	fakeServer.Mux.HandleFunc("/compute/v2/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			{
				"version": {
					"id": "v2.0",
					"status": "SUPPORTED",
					"version": "",
					"min_version": "",
					"updated": "2011-01-21T11:33:21Z",
					"links": [
						{
							"rel": "self",
							"href": "%s/compute/v2/"
						},
						{
							"rel": "describedby",
							"type": "text/html",
							"href": "http://docs.openstack.org/"
						}
					],
					"media-types": [
						{
							"base": "application/json",
							"type": "application/vnd.openstack.compute+json;version=2"
						}
					]
				}
			}
		`, fakeServer.Server.URL)
	})
	// Compute v2.1 API
	fakeServer.Mux.HandleFunc("/compute/v2.1/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			{
				"version": {
					"id": "v2.1",
					"status": "CURRENT",
					"version": "2.90",
					"min_version": "2.1",
					"updated": "2013-07-23T11:33:21Z",
					"links": [
						{
							"rel": "self",
							"href": "%s/compute/v2.1/"
						},
						{
							"rel": "describedby",
							"type": "text/html",
							"href": "http://docs.openstack.org/"
						}
					],
					"media-types": [
						{
							"base": "application/json",
							"type": "application/vnd.openstack.compute+json;version=2.1"
						}
					]
				}
			}
		`, fakeServer.Server.URL)
	})
	// Baremetal API
	fakeServer.Mux.HandleFunc("/baremetal/v1/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
		{
			"name": "OpenStack Ironic API",
			"description": "Ironic is an OpenStack project which enables the provision and management of baremetal machines.",
			"default_version": {
				"id": "v1",
				"links": [
					{
						"href": "%s/baremetal/v1/",
						"rel": "self"
					}
				],
				"status": "CURRENT",
				"min_version": "1.1",
				"version": "1.87"
			},
			"versions": [
				{
					"id": "v1",
					"links": [
						{
							"href": "%s/baremetal/v1/",
							"rel": "self"
						}
					],
					"status": "CURRENT",
					"min_version": "1.1",
					"version": "1.87"
				}
			]
		}
		`, fakeServer.Server.URL, fakeServer.Server.URL)
	})
	// Fictional multi-version API
	fakeServer.Mux.HandleFunc("/multi-version/v1.2/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
		{
			"name": "Multi-version API",
			"description": "A fictional API with multiple microversions.",
			"versions": [
				{
					"id": "v1",
					"links": [
						{
							"href": "%s/multi-version/v1/",
							"rel": "self"
						}
					],
					"status": "CURRENT",
					"min_version": "1.1",
					"version": "1.87"
				},
				{
					"id": "v1.2",
					"links": [
						{
							"href": "%s/multi-version/v1/",
							"rel": "self"
						}
					],
					"status": "CURRENT",
					"min_version": "1.2",
					"version": "1.90"
				}
			]
		}
		`, fakeServer.Server.URL, fakeServer.Server.URL)
	})
}
