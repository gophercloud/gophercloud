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
	// Identity root API
	fakeServer.Mux.HandleFunc("/identity/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			{
				"versions": {
					"values": [
						{
							"id": "v3.14",
							"status": "stable",
							"updated": "2020-04-07T00:00:00Z",
							"links": [
								{
									"rel": "self",
									"href": "%s/identity/v3/"
								}
							],
							"media-types": [
								{
									"base": "application/json",
									"type": "application/vnd.openstack.identity-v3+json"
								}
							]
						}
					]
				}
			}
		`, fakeServer.Server.URL)
	})
	// Identity v3 API
	fakeServer.Mux.HandleFunc("/identity/v3/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			{
				"version": {
					"id": "v3.14",
					"status": "stable",
					"updated": "2020-04-07T00:00:00Z",
					"links": [
						{
							"rel": "self",
							"href": "%s/identity/v3/"
						}
					],
					"media-types": [
						{
							"base": "application/json",
							"type": "application/vnd.openstack.identity-v3+json"
						}
					]
				}
			}
		`, fakeServer.Server.URL)
	})
	// Compute root API
	fakeServer.Mux.HandleFunc("/compute/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			{
				"versions": [
					{
						"id": "v2.0",
						"status": "SUPPORTED",
						"version": "",
						"min_version": "",
						"updated": "2011-01-21T11:33:21Z",
						"links": [
							{
								"rel": "self",
								"href": "%s/compute/v2/"
							}
						]
					},
					{
						"id": "v2.1",
						"status": "CURRENT",
						"version": "2.90",
						"min_version": "2.1",
						"updated": "2013-07-23T11:33:21Z",
						"links": [
							{
								"rel": "self",
								"href": "%s/compute/v2.1/"
							}
						]
					}
				]
			}
		`, fakeServer.Server.URL, fakeServer.Server.URL)
	})
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
	// Container Infra root API
	fakeServer.Mux.HandleFunc("/container-infra/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
		{
			"name": "OpenStack Magnum API",
			"description": "Magnum is an OpenStack project which aims to provide container cluster management.",
			"versions": [
				{
					"id": "v1",
					"links": [
						{
							"href": "%s/v1/",
							"rel": "self"
						}
					],
					"status": "CURRENT",
					"max_version": "1.11",
					"min_version": "1.1"
				}
			]
		}
		`, fakeServer.Server.URL)
	})
	// Container Infra v1 API
	//
	// NOTE(stephenfin): In reality, this returns absolute URLs, but those URLs are wrong since
	// they don't respect the Host header. We're using relative URLs because (a) it's probably
	// what magnum should be doing and (b) it avoids needing 17 odd arguments to Fprintf
	fakeServer.Mux.HandleFunc("/container-infra/v1/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `
			{
				"id": "v1",
				"media_types": [
					{
						"base": "application/json",
						"type": "application/vnd.openstack.magnum.v1+json"
					}
				],
				"links": [
					{
						"href": "/v1/",
						"rel": "self"
					},
					{
						"href": "http://docs.openstack.org/developer/magnum/dev/api-spec-v1.html",
						"rel": "describedby",
						"type": "text/html"
					}
				],
				"clustertemplates": [
					{
						"href": "/v1/clustertemplates/",
						"rel": "self"
					},
					{
						"href": "/clustertemplates/",
						"rel": "bookmark"
					}
				],
				"clusters": [
					{
						"href": "/v1/clusters/",
						"rel": "self"
					},
					{
						"href": "/clusters/",
						"rel": "bookmark"
					}
				],
				"quotas": [
					{
						"href": "/v1/quotas/",
						"rel": "self"
					},
					{
						"href": "/quotas/",
						"rel": "bookmark"
					}
				],
				"certificates": [
					{
						"href": "/v1/certificates/",
						"rel": "self"
					},
					{
						"href": "/certificates/",
						"rel": "bookmark"
					}
				],
				"mservices": [
					{
						"href": "/v1/mservices/",
						"rel": "self"
					},
					{
						"href": "/mservices/",
						"rel": "bookmark"
					}
				],
				"stats": [
					{
						"href": "/v1/stats/",
						"rel": "self"
					},
					{
						"href": "/stats/",
						"rel": "bookmark"
					}
				],
				"federations": [
					{
						"href": "/v1/federations/",
						"rel": "self"
					},
					{
						"href": "/federations/",
						"rel": "bookmark"
					}
				],
				"nodegroups": [
					{
						"href": "/v1/clusters/{cluster_id}/nodegroups",
						"rel": "self"
					},
					{
						"href": "/clusters/{cluster_id}/nodegroups",
						"rel": "bookmark"
					}
				]
			}
		`, fakeServer.Server.URL)
	})
	// Orchestration root API
	fakeServer.Mux.HandleFunc("/heat-api/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			{
				"versions": [
					{
						"id": "v1.0",
						"status": "CURRENT",
						"links": [
							{
								"rel": "self",
								"href": "%s/heat-api/v1/"
							}
						]
					}
				]
			}
		`, fakeServer.Server.URL)
	})
	// Orchestration v1 API (non-existent)
	fakeServer.Mux.HandleFunc("/heat-api/v1/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
	})
	// Workflow root API
	//
	// In reality, this deploys under a port rather than a path (as of Epoxy) but we don't want to
	// have to run multiple fake servers so this will do.
	fakeServer.Mux.HandleFunc("/workflow/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			{
				"versions": [
					{
						"id": "v2.0",
						"status": "CURRENT",
						"links": [
							{
								"href": "%s/workflow/v2",
								"target": "v2",
								"rel": "self"
							}
						]
					}
				]
			}
		`, fakeServer.Server.URL)
	})
	// Workflow v1 API (invalid version document)
	fakeServer.Mux.HandleFunc("/workflow/v2/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			{
				"uri": "%s/v2"
			}
		`, fakeServer.Server.URL)
	})
	// Baremetal root API
	fakeServer.Mux.HandleFunc("/baremetal/", func(w http.ResponseWriter, r *http.Request) {
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
	// Baremetal v1 API
	//
	// NOTE(stephenfin): In reality, this returns absolute URLs and unlike Magnum those URLs are
	// correctly formatted. We're using relative URLs for many of these because, once again, it
	// avoids needing loads of arguments to Fprintf
	fakeServer.Mux.HandleFunc("/baremetal/v1/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
		{
			"id": "v1",
			"links": [
				{
					"href": "/baremetal/v1/",
					"rel": "self"
				},
				{
					"href": "https://docs.openstack.org//ironic/latest/contributor//webapi.html",
					"rel": "describedby",
					"type": "text/html"
				}
			],
			"media_types": {
				"base": "application/json",
				"type": "application/vnd.openstack.ironic.v1+json"
			},
			"chassis": [
				{
					"href": "/baremetal/v1/chassis/",
					"rel": "self"
				},
				{
					"href": "/baremetal/chassis/",
					"rel": "bookmark"
				}
			],
			"nodes": [
				{
					"href": "/baremetal/v1/nodes/",
					"rel": "self"
				},
				{
					"href": "/baremetal/nodes/",
					"rel": "bookmark"
				}
			],
			"ports": [
				{
					"href": "/baremetal/v1/ports/",
					"rel": "self"
				},
				{
					"href": "/baremetal/ports/",
					"rel": "bookmark"
				}
			],
			"drivers": [
				{
					"href": "/baremetal/v1/drivers/",
					"rel": "self"
				},
				{
					"href": "/baremetal/drivers/",
					"rel": "bookmark"
				}
			],
			"version": {
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
		}
		`, fakeServer.Server.URL)
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
