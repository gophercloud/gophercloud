package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/l7policies"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/listeners"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/loadbalancers"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/monitors"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/pools"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// LoadbalancersListBody contains the canned body of a loadbalancer list response.
const LoadbalancersListBody = `
{
	"loadbalancers":[
		{
			"id": "c331058c-6a40-4144-948e-b9fb1df9db4b",
			"project_id": "54030507-44f7-473c-9342-b4d14a95f692",
			"created_at": "2019-06-30T04:15:37",
			"updated_at": "2019-06-30T05:18:49",
			"name": "web_lb",
			"description": "lb config for the web tier",
			"vip_subnet_id": "8a49c438-848f-467b-9655-ea1548708154",
			"vip_address": "10.30.176.47",
			"vip_port_id": "2a22e552-a347-44fd-b530-1f2b1b2a6735",
			"flavor_id": "60df399a-ee85-11e9-81b4-2a2ae2dbcce4",
			"provider": "haproxy",
			"admin_state_up": true,
			"provisioning_status": "ACTIVE",
			"operating_status": "ONLINE",
			"tags": ["test", "stage"]
		},
		{
			"id": "36e08a3e-a78f-4b40-a229-1e7e23eee1ab",
			"project_id": "54030507-44f7-473c-9342-b4d14a95f692",
			"created_at": "2019-06-30T04:15:37",
			"updated_at": "2019-06-30T05:18:49",
			"name": "db_lb",
			"description": "lb config for the db tier",
			"vip_subnet_id": "9cedb85d-0759-4898-8a4b-fa5a5ea10086",
			"vip_address": "10.30.176.48",
			"vip_port_id": "2bf413c8-41a9-4477-b505-333d5cbe8b55",
			"flavor_id": "bba40eb2-ee8c-11e9-81b4-2a2ae2dbcce4",
			"availability_zone": "db_az",
			"provider": "haproxy",
			"admin_state_up": true,
			"provisioning_status": "PENDING_CREATE",
			"operating_status": "OFFLINE",
			"tags": ["test", "stage"],
			"additional_vips": [{"subnet_id": "0d4f6a08-60b7-44ab-8903-f7d76ec54095", "ip_address" : "192.168.10.10"}]
		}
	]
}
`

// SingleLoadbalancerBody is the canned body of a Get request on an existing loadbalancer.
const SingleLoadbalancerBody = `
{
	"loadbalancer": {
		"id": "36e08a3e-a78f-4b40-a229-1e7e23eee1ab",
		"project_id": "54030507-44f7-473c-9342-b4d14a95f692",
		"created_at": "2019-06-30T04:15:37",
		"updated_at": "2019-06-30T05:18:49",
		"name": "db_lb",
		"description": "lb config for the db tier",
		"vip_subnet_id": "9cedb85d-0759-4898-8a4b-fa5a5ea10086",
		"vip_address": "10.30.176.48",
		"vip_port_id": "2bf413c8-41a9-4477-b505-333d5cbe8b55",
		"flavor_id": "bba40eb2-ee8c-11e9-81b4-2a2ae2dbcce4",
		"availability_zone": "db_az",
		"provider": "haproxy",
		"admin_state_up": true,
		"provisioning_status": "PENDING_CREATE",
		"operating_status": "OFFLINE",
		"tags": ["test", "stage"],
		"additional_vips": [{"subnet_id": "0d4f6a08-60b7-44ab-8903-f7d76ec54095", "ip_address" : "192.168.10.10"}]
	}
}
`

// PostUpdateLoadbalancerBody is the canned response body of a Update request on an existing loadbalancer.
const PostUpdateLoadbalancerBody = `
{
	"loadbalancer": {
		"id": "36e08a3e-a78f-4b40-a229-1e7e23eee1ab",
		"project_id": "54030507-44f7-473c-9342-b4d14a95f692",
		"created_at": "2019-06-30T04:15:37",
		"updated_at": "2019-06-30T05:18:49",
		"name": "NewLoadbalancerName",
		"description": "lb config for the db tier",
		"vip_subnet_id": "9cedb85d-0759-4898-8a4b-fa5a5ea10086",
		"vip_address": "10.30.176.48",
		"vip_port_id": "2bf413c8-41a9-4477-b505-333d5cbe8b55",
		"flavor_id": "bba40eb2-ee8c-11e9-81b4-2a2ae2dbcce4",
		"provider": "haproxy",
		"admin_state_up": true,
		"provisioning_status": "PENDING_CREATE",
		"operating_status": "OFFLINE",
		"tags": ["test"]
	}
}
`

// PostFullyPopulatedLoadbalancerBody is the canned response body of a Create request of an fully populated loadbalancer.
const PostFullyPopulatedLoadbalancerBody = `
{
	"loadbalancer": {
		"description": "My favorite load balancer",
		"admin_state_up": true,
		"project_id": "e3cd678b11784734bc366148aa37580e",
		"provisioning_status": "ACTIVE",
		"flavor_id": "",
		"created_at": "2019-06-30T04:15:37",
		"updated_at": "2019-06-30T05:18:49",
		"listeners": [{
			"l7policies": [{
				"description": "",
				"admin_state_up": true,
				"rules": [],
				"project_id": "e3cd678b11784734bc366148aa37580e",
				"listener_id": "95de30ec-67f4-437b-b3f3-22c5d9ef9828",
				"redirect_url": "https://www.example.com/",
				"action": "REDIRECT_TO_URL",
				"position": 1,
				"id": "d0553837-f890-4981-b99a-f7cbd6a76577",
				"name": "redirect_policy"
			}],
			"protocol": "HTTP",
			"description": "",
			"default_tls_container_ref": null,
			"admin_state_up": true,
			"default_pool_id": "c8cec227-410a-4a5b-af13-ecf38c2b0abb",
			"project_id": "e3cd678b11784734bc366148aa37580e",
			"default_tls_container_id": null,
			"connection_limit": -1,
			"sni_container_refs": [],
			"protocol_port": 8080,
			"id": "95de30ec-67f4-437b-b3f3-22c5d9ef9828",
			"name": "redirect_listener"
		}],
		"vip_address": "203.0.113.50",
		"vip_network_id": "d0d217df-3958-4fbf-a3c2-8dad2908c709",
		"vip_subnet_id": "d4af86e1-0051-488c-b7a0-527f97490c9a",
		"vip_port_id": "b4ca07d1-a31e-43e2-891a-7d14f419f342",
		"provider": "octavia",
		"pools": [{
			"lb_algorithm": "ROUND_ROBIN",
			"protocol": "HTTP",
			"description": "",
			"admin_state_up": true,
			"project_id": "e3cd678b11784734bc366148aa37580e",
			"session_persistence": null,
			"healthmonitor": {
				"name": "",
				"admin_state_up": true,
				"project_id": "e3cd678b11784734bc366148aa37580e",
				"delay": 3,
				"expected_codes": "200,201,202",
				"max_retries": 2,
				"http_method": "GET",
				"timeout": 1,
				"max_retries_down": 3,
				"url_path": "/index.html",
				"type": "HTTP",
				"id": "a8a2aa3f-d099-4752-8265-e6472f8147f9"
			},
			"members": [{
				"name": "",
				"weight": 1,
				"admin_state_up": true,
				"subnet_id": "bbb35f84-35cc-4b2f-84c2-a6a29bba68aa",
				"project_id": "e3cd678b11784734bc366148aa37580e",
				"address": "192.0.2.16",
				"protocol_port": 80,
				"id": "7d19ad6c-d549-453e-a5cd-05382c6be96a"
			},{
				"name": "",
				"weight": 1,
				"admin_state_up": true,
				"subnet_id": "bbb35f84-35cc-4b2f-84c2-a6a29bba68aa",
				"project_id": "e3cd678b11784734bc366148aa37580e",
				"address": "192.0.2.19",
				"protocol_port": 80,
				"id": "a167402b-caa6-41d5-b4d4-bde7f2cbfa5e"
			}],
			"id": "c8cec227-410a-4a5b-af13-ecf38c2b0abb",
			"name": "rr_pool"
		}],
		"id": "607226db-27ef-4d41-ae89-f2a800e9c2db",
		"operating_status": "ONLINE",
		"name": "best_load_balancer",
		"availability_zone": "my_az",
		"tags": ["test_tag"]
	}
}
`

// GetLoadbalancerStatusesBody is the canned request body of a Get request on loadbalancer's status.
const GetLoadbalancerStatusesBody = `
{
	"statuses" : {
		"loadbalancer": {
			"id": "36e08a3e-a78f-4b40-a229-1e7e23eee1ab",
			"name": "db_lb",
			"provisioning_status": "PENDING_UPDATE",
			"operating_status": "ACTIVE",
			"tags": ["test", "stage"],
			"listeners": [{
				"id": "db902c0c-d5ff-4753-b465-668ad9656918",
				"name": "db",
				"provisioning_status": "ACTIVE",
				"pools": [{
					"id": "fad389a3-9a4a-4762-a365-8c7038508b5d",
					"name": "db",
					"provisioning_status": "ACTIVE",
					"healthmonitor": {
						"id": "67306cda-815d-4354-9fe4-59e09da9c3c5",
						"type":"PING",
						"provisioning_status": "ACTIVE"
					},
					"members":[{
						"id": "2a280670-c202-4b0b-a562-34077415aabf",
						"name": "db",
						"address": "10.0.2.11",
						"protocol_port": 80,
						"provisioning_status": "ACTIVE"
					}]
				}]
			}]
		}
	}
}
`

// LoadbalancerStatsTree is the canned request body of a Get request on loadbalancer's statistics.
const GetLoadbalancerStatsBody = `
{
    "stats": {
        "active_connections": 0,
        "bytes_in": 9532,
        "bytes_out": 22033,
        "request_errors": 46,
        "total_connections": 112
    }
}
`

var (
	createdTime, _ = time.Parse(time.RFC3339, "2019-06-30T04:15:37Z")
	updatedTime, _ = time.Parse(time.RFC3339, "2019-06-30T05:18:49Z")
)

var (
	LoadbalancerWeb = loadbalancers.LoadBalancer{
		ID:                 "c331058c-6a40-4144-948e-b9fb1df9db4b",
		ProjectID:          "54030507-44f7-473c-9342-b4d14a95f692",
		CreatedAt:          createdTime,
		UpdatedAt:          updatedTime,
		Name:               "web_lb",
		Description:        "lb config for the web tier",
		VipSubnetID:        "8a49c438-848f-467b-9655-ea1548708154",
		VipAddress:         "10.30.176.47",
		VipPortID:          "2a22e552-a347-44fd-b530-1f2b1b2a6735",
		FlavorID:           "60df399a-ee85-11e9-81b4-2a2ae2dbcce4",
		Provider:           "haproxy",
		AdminStateUp:       true,
		ProvisioningStatus: "ACTIVE",
		OperatingStatus:    "ONLINE",
		Tags:               []string{"test", "stage"},
	}
	LoadbalancerDb = loadbalancers.LoadBalancer{
		ID:                 "36e08a3e-a78f-4b40-a229-1e7e23eee1ab",
		ProjectID:          "54030507-44f7-473c-9342-b4d14a95f692",
		CreatedAt:          createdTime,
		UpdatedAt:          updatedTime,
		Name:               "db_lb",
		Description:        "lb config for the db tier",
		VipSubnetID:        "9cedb85d-0759-4898-8a4b-fa5a5ea10086",
		VipAddress:         "10.30.176.48",
		VipPortID:          "2bf413c8-41a9-4477-b505-333d5cbe8b55",
		FlavorID:           "bba40eb2-ee8c-11e9-81b4-2a2ae2dbcce4",
		AvailabilityZone:   "db_az",
		Provider:           "haproxy",
		AdminStateUp:       true,
		ProvisioningStatus: "PENDING_CREATE",
		OperatingStatus:    "OFFLINE",
		Tags:               []string{"test", "stage"},
		AdditionalVips: []loadbalancers.AdditionalVip{
			{
				SubnetID:  "0d4f6a08-60b7-44ab-8903-f7d76ec54095",
				IPAddress: "192.168.10.10",
			},
		},
	}
	LoadbalancerUpdated = loadbalancers.LoadBalancer{
		ID:                 "36e08a3e-a78f-4b40-a229-1e7e23eee1ab",
		ProjectID:          "54030507-44f7-473c-9342-b4d14a95f692",
		CreatedAt:          createdTime,
		UpdatedAt:          updatedTime,
		Name:               "NewLoadbalancerName",
		Description:        "lb config for the db tier",
		VipSubnetID:        "9cedb85d-0759-4898-8a4b-fa5a5ea10086",
		VipAddress:         "10.30.176.48",
		VipPortID:          "2bf413c8-41a9-4477-b505-333d5cbe8b55",
		FlavorID:           "bba40eb2-ee8c-11e9-81b4-2a2ae2dbcce4",
		Provider:           "haproxy",
		AdminStateUp:       true,
		ProvisioningStatus: "PENDING_CREATE",
		OperatingStatus:    "OFFLINE",
		Tags:               []string{"test"},
	}
	FullyPopulatedLoadBalancerDb = loadbalancers.LoadBalancer{
		Description:        "My favorite load balancer",
		AdminStateUp:       true,
		ProjectID:          "e3cd678b11784734bc366148aa37580e",
		UpdatedAt:          updatedTime,
		CreatedAt:          createdTime,
		ProvisioningStatus: "ACTIVE",
		VipSubnetID:        "d4af86e1-0051-488c-b7a0-527f97490c9a",
		VipNetworkID:       "d0d217df-3958-4fbf-a3c2-8dad2908c709",
		VipAddress:         "203.0.113.50",
		VipPortID:          "b4ca07d1-a31e-43e2-891a-7d14f419f342",
		AvailabilityZone:   "my_az",
		ID:                 "607226db-27ef-4d41-ae89-f2a800e9c2db",
		OperatingStatus:    "ONLINE",
		Name:               "best_load_balancer",
		FlavorID:           "",
		Provider:           "octavia",
		Tags:               []string{"test_tag"},
		Listeners: []listeners.Listener{{
			ID:               "95de30ec-67f4-437b-b3f3-22c5d9ef9828",
			ProjectID:        "e3cd678b11784734bc366148aa37580e",
			Name:             "redirect_listener",
			Description:      "",
			Protocol:         "HTTP",
			ProtocolPort:     8080,
			DefaultPoolID:    "c8cec227-410a-4a5b-af13-ecf38c2b0abb",
			AdminStateUp:     true,
			ConnLimit:        -1,
			SniContainerRefs: []string{},
			L7Policies: []l7policies.L7Policy{{
				ID:           "d0553837-f890-4981-b99a-f7cbd6a76577",
				Name:         "redirect_policy",
				ListenerID:   "95de30ec-67f4-437b-b3f3-22c5d9ef9828",
				ProjectID:    "e3cd678b11784734bc366148aa37580e",
				Description:  "",
				Action:       "REDIRECT_TO_URL",
				Position:     1,
				RedirectURL:  "https://www.example.com/",
				AdminStateUp: true,
				Rules:        []l7policies.Rule{},
			}},
		}},
		Pools: []pools.Pool{{
			LBMethod:     "ROUND_ROBIN",
			Protocol:     "HTTP",
			Description:  "",
			AdminStateUp: true,
			Name:         "rr_pool",
			ID:           "c8cec227-410a-4a5b-af13-ecf38c2b0abb",
			ProjectID:    "e3cd678b11784734bc366148aa37580e",
			Members: []pools.Member{{
				Name:         "",
				Address:      "192.0.2.16",
				SubnetID:     "bbb35f84-35cc-4b2f-84c2-a6a29bba68aa",
				AdminStateUp: true,
				ProtocolPort: 80,
				ID:           "7d19ad6c-d549-453e-a5cd-05382c6be96a",
				ProjectID:    "e3cd678b11784734bc366148aa37580e",
				Weight:       1,
			}, {
				Name:         "",
				Address:      "192.0.2.19",
				SubnetID:     "bbb35f84-35cc-4b2f-84c2-a6a29bba68aa",
				AdminStateUp: true,
				ProtocolPort: 80,
				ID:           "a167402b-caa6-41d5-b4d4-bde7f2cbfa5e",
				ProjectID:    "e3cd678b11784734bc366148aa37580e",
				Weight:       1,
			}},
			Monitor: monitors.Monitor{
				ID:             "a8a2aa3f-d099-4752-8265-e6472f8147f9",
				ProjectID:      "e3cd678b11784734bc366148aa37580e",
				Name:           "",
				Type:           "HTTP",
				Timeout:        1,
				MaxRetries:     2,
				Delay:          3,
				MaxRetriesDown: 3,
				HTTPMethod:     "GET",
				URLPath:        "/index.html",
				ExpectedCodes:  "200,201,202",
				AdminStateUp:   true,
			},
		}},
	}
	LoadbalancerStatusesTree = loadbalancers.StatusTree{
		Loadbalancer: &loadbalancers.LoadBalancer{
			ID:                 "36e08a3e-a78f-4b40-a229-1e7e23eee1ab",
			Name:               "db_lb",
			ProvisioningStatus: "PENDING_UPDATE",
			OperatingStatus:    "ACTIVE",
			Tags:               []string{"test", "stage"},
			Listeners: []listeners.Listener{{
				ID:                 "db902c0c-d5ff-4753-b465-668ad9656918",
				Name:               "db",
				ProvisioningStatus: "ACTIVE",
				Pools: []pools.Pool{{
					ID:                 "fad389a3-9a4a-4762-a365-8c7038508b5d",
					Name:               "db",
					ProvisioningStatus: "ACTIVE",
					Monitor: monitors.Monitor{
						ID:                 "67306cda-815d-4354-9fe4-59e09da9c3c5",
						Type:               "PING",
						ProvisioningStatus: "ACTIVE",
					},
					Members: []pools.Member{{
						ID:                 "2a280670-c202-4b0b-a562-34077415aabf",
						Name:               "db",
						Address:            "10.0.2.11",
						ProtocolPort:       80,
						ProvisioningStatus: "ACTIVE",
					}},
				}},
			}},
		},
	}
	LoadbalancerStatsTree = loadbalancers.Stats{
		ActiveConnections: 0,
		BytesIn:           9532,
		BytesOut:          22033,
		RequestErrors:     46,
		TotalConnections:  112,
	}
)

// HandleLoadbalancerListSuccessfully sets up the test server to respond to a loadbalancer List request.
func HandleLoadbalancerListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v2.0/lbaas/loadbalancers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, LoadbalancersListBody)
		case "45e08a3e-a78f-4b40-a229-1e7e23eee1ab":
			fmt.Fprintf(w, `{ "loadbalancers": [] }`)
		default:
			t.Fatalf("/v2.0/lbaas/loadbalancers invoked with unexpected marker=[%s]", marker)
		}
	})
}

// HandleFullyPopulatedLoadbalancerCreationSuccessfully sets up the test server to respond to a
// fully populated loadbalancer creation request with a given response.
func HandleFullyPopulatedLoadbalancerCreationSuccessfully(t *testing.T, response string) {
	th.Mux.HandleFunc("/v2.0/lbaas/loadbalancers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `{
			"loadbalancer": {
				"admin_state_up": true,
				"flavor_id": "bba40eb2-ee8c-11e9-81b4-2a2ae2dbcce4",
				"listeners": [
					{
						"default_pool": {
							"healthmonitor": {
								"delay": 3,
								"expected_codes": "200",
								"http_method": "GET",
								"max_retries": 2,
								"max_retries_down": 3,
								"name": "db",
								"timeout": 1,
								"type": "HTTP",
								"url_path": "/index.html"
							},
							"lb_algorithm": "ROUND_ROBIN",
							"members": [
								{
									"address": "192.0.2.51",
									"protocol_port": 80
								},
								{
									"address": "192.0.2.52",
									"protocol_port": 80
								}
							],
							"name": "Example pool",
							"protocol": "HTTP"
						},
						"l7policies": [
							{
								"action": "REDIRECT_TO_URL",
								"name": "redirect-example.com",
								"redirect_url": "http://www.example.com",
								"rules": [
									{
										"compare_type": "REGEX",
										"type": "PATH",
										"value": "/images*"
									}
								]
							}
						],
						"name": "redirect_listener",
						"protocol": "HTTP",
						"protocol_port": 8080
					}
				],
				"name": "db_lb",
				"provider": "octavia",
				"tags": [
					"test",
					"stage"
				],
				"vip_address": "10.30.176.48",
				"vip_port_id": "2bf413c8-41a9-4477-b505-333d5cbe8b55",
				"vip_subnet_id": "9cedb85d-0759-4898-8a4b-fa5a5ea10086"
			}
		}`)

		w.WriteHeader(http.StatusAccepted)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, response)
	})
}

// HandleLoadbalancerCreationSuccessfully sets up the test server to respond to a loadbalancer creation request
// with a given response.
func HandleLoadbalancerCreationSuccessfully(t *testing.T, response string) {
	th.Mux.HandleFunc("/v2.0/lbaas/loadbalancers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `{
			"loadbalancer": {
				"name": "db_lb",
				"vip_port_id": "2bf413c8-41a9-4477-b505-333d5cbe8b55",
				"vip_subnet_id": "9cedb85d-0759-4898-8a4b-fa5a5ea10086",
				"vip_address": "10.30.176.48",
				"flavor_id": "bba40eb2-ee8c-11e9-81b4-2a2ae2dbcce4",
				"provider": "haproxy",
				"admin_state_up": true,
				"tags": ["test", "stage"],
				"additional_vips": [{"subnet_id": "0d4f6a08-60b7-44ab-8903-f7d76ec54095", "ip_address" : "192.168.10.10"}]
			}
		}`)

		w.WriteHeader(http.StatusAccepted)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, response)
	})
}

// HandleLoadbalancerGetSuccessfully sets up the test server to respond to a loadbalancer Get request.
func HandleLoadbalancerGetSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v2.0/lbaas/loadbalancers/36e08a3e-a78f-4b40-a229-1e7e23eee1ab", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		fmt.Fprintf(w, SingleLoadbalancerBody)
	})
}

// HandleLoadbalancerGetStatusesTree sets up the test server to respond to a loadbalancer Get statuses tree request.
func HandleLoadbalancerGetStatusesTree(t *testing.T) {
	th.Mux.HandleFunc("/v2.0/lbaas/loadbalancers/36e08a3e-a78f-4b40-a229-1e7e23eee1ab/status", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		fmt.Fprintf(w, GetLoadbalancerStatusesBody)
	})
}

// HandleLoadbalancerDeletionSuccessfully sets up the test server to respond to a loadbalancer deletion request.
func HandleLoadbalancerDeletionSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v2.0/lbaas/loadbalancers/36e08a3e-a78f-4b40-a229-1e7e23eee1ab", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleLoadbalancerUpdateSuccessfully sets up the test server to respond to a loadbalancer Update request.
func HandleLoadbalancerUpdateSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v2.0/lbaas/loadbalancers/36e08a3e-a78f-4b40-a229-1e7e23eee1ab", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestJSONRequest(t, r, `{
			"loadbalancer": {
				"name": "NewLoadbalancerName",
				"tags": ["test"]
			}
		}`)

		fmt.Fprintf(w, PostUpdateLoadbalancerBody)
	})
}

// HandleLoadbalancerGetStatsTree sets up the test server to respond to a loadbalancer Get stats tree request.
func HandleLoadbalancerGetStatsTree(t *testing.T) {
	th.Mux.HandleFunc("/v2.0/lbaas/loadbalancers/36e08a3e-a78f-4b40-a229-1e7e23eee1ab/stats", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		fmt.Fprintf(w, GetLoadbalancerStatsBody)
	})
}

// HandleLoadbalancerFailoverSuccessfully sets up the test server to respond to a loadbalancer failover request.
func HandleLoadbalancerFailoverSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v2.0/lbaas/loadbalancers/36e08a3e-a78f-4b40-a229-1e7e23eee1ab/failover", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusAccepted)
	})
}
