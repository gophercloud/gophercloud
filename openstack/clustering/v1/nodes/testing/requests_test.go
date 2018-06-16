package testing

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/clustering/v1/nodes"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestCreateNode(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/nodes", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("X-OpenStack-Request-ID", "req-3791a089-9d46-4671-a3f9-55e95e55d2b4")
		w.Header().Add("Location", "http://senlin.cloud.blizzard.net:8778/v1/actions/ffd94dd8-6266-4887-9a8c-5b78b72136da")

		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
			"node": {
				"cluster_id": "e395be1e-8d8e-43bb-bd6c-943eccf76a6d",
				"created_at": "2016-05-13T07:02:20Z",
				"data": {
					"internal_ports": [
						{
							"network_id": "847e4f65-1ff1-42b1-9e74-74e6a109ad11",
							"security_group_ids": ["8db277ab-1d98-4148-ba72-724721789427"],
							"fixed_ips": [
								{
								  "subnet_id": "863b20c0-c011-4650-85c2-ad531f4570a4",
								  "ip_address": "10.63.177.162"
								}
							],
							"id": "43aa53d7-a70b-4f40-812f-4feecb687018",
							"remove": true
						}
					],
					"placement": {
						"zone": "nova"
					}
				},
				"dependents": {},
				"domain": "1235be1e-8d8e-43bb-bd6c-943eccf76a6d",
				"id": "82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1",
				"index": 2,
				"init_at": "2016-05-13T08:02:04Z",
				"metadata": {
					"test": {
						"nil_interface": null,
						"bool_value": false,
						"string_value": "test_string",
						"float_value": 123.3
					},
					"foo": "bar"
				},
				"name": "node-e395be1e-002",
				"physical_id": "66a81d68-bf48-4af5-897b-a3bfef7279a8",
				"profile_id": "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
				"profile_name": "pcirros",
				"project_id": "eee0b7c083e84501bdd50fb269d2a10e",
				"role": "",
				"status": "ACTIVE",
				"status_reason": "Creation succeeded",
				"updated_at": "2016-05-13T09:02:04Z",
				"user": "ab79b9647d074e46ac223a8fa297b846"				
			}
		}`)
	})

	opts := nodes.CreateOpts{
		ClusterID: "e395be1e-8d8e-43bb-bd6c-943eccf76a6d",
		Metadata: map[string]interface{}{
			"foo": "bar",
			"test": map[string]interface{}{
				"nil_interface": interface{}(nil),
				"float_value":   float64(123.3),
				"string_value":  "test_string",
				"bool_value":    false,
			},
		},
		Name:      "node-e395be1e-002",
		ProfileID: "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
		Role:      "",
	}

	createResult := nodes.Create(fake.ServiceClient(), opts)
	if createResult.Err != nil {
		t.Error("Error creating node. error=", createResult.Err)
	}

	requestID := createResult.Header.Get("X-Openstack-Request-Id")
	th.AssertEquals(t, "req-3791a089-9d46-4671-a3f9-55e95e55d2b4", requestID)

	location := createResult.Header.Get("Location")
	th.AssertEquals(t, "http://senlin.cloud.blizzard.net:8778/v1/actions/ffd94dd8-6266-4887-9a8c-5b78b72136da", location)

	actionID := ""
	locationFields := strings.Split(location, "actions/")
	if len(locationFields) >= 2 {
		actionID = locationFields[1]
	}
	th.AssertEquals(t, "ffd94dd8-6266-4887-9a8c-5b78b72136da", actionID)

	actual, err := createResult.Extract()
	if err != nil {
		t.Error("Error creating nodes. error=", err)
	} else {
		createdAt, _ := time.Parse(time.RFC3339, "2016-05-13T07:02:20Z")
		initAt, _ := time.Parse(time.RFC3339, "2016-05-13T08:02:04Z")
		updatedAt, _ := time.Parse(time.RFC3339, "2016-05-13T09:02:04Z")

		expected := nodes.Node{
			ClusterID: "e395be1e-8d8e-43bb-bd6c-943eccf76a6d",
			CreatedAt: createdAt,
			Data: map[string]interface{}{
				"internal_ports": []map[string]interface{}{
					{
						"network_id": "847e4f65-1ff1-42b1-9e74-74e6a109ad11",
						"security_group_ids": []interface{}{
							"8db277ab-1d98-4148-ba72-724721789427",
						},
						"fixed_ips": []interface{}{
							map[string]interface{}{
								"subnet_id":  "863b20c0-c011-4650-85c2-ad531f4570a4",
								"ip_address": "10.63.177.162",
							},
						},
						"id":     "43aa53d7-a70b-4f40-812f-4feecb687018",
						"remove": true,
					},
				},
				"placement": map[string]interface{}{
					"zone": "nova",
				},
			},
			Dependents: map[string]interface{}{},
			Domain:     "1235be1e-8d8e-43bb-bd6c-943eccf76a6d",
			ID:         "82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1",
			Index:      2,
			InitAt:     initAt,
			Metadata: map[string]interface{}{
				"foo": "bar",
				"test": map[string]interface{}{
					"nil_interface": interface{}(nil),
					"float_value":   float64(123.3),
					"string_value":  "test_string",
					"bool_value":    false,
				},
			},
			Name:         "node-e395be1e-002",
			PhysicalID:   "66a81d68-bf48-4af5-897b-a3bfef7279a8",
			ProfileID:    "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
			ProfileName:  "pcirros",
			ProjectID:    "eee0b7c083e84501bdd50fb269d2a10e",
			Role:         "",
			Status:       "ACTIVE",
			StatusReason: "Creation succeeded",
			UpdatedAt:    updatedAt,
			User:         "ab79b9647d074e46ac223a8fa297b846",
		}
		th.AssertDeepEquals(t, expected, *actual)
	}
}

func TestCreateNodeEmptyTime(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/nodes", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
			"node": {
				"cluster_id": "e395be1e-8d8e-43bb-bd6c-943eccf76a6d",
				"created_at": null,
				"data": {
					"internal_ports": [
						{
							"network_id": "847e4f65-1ff1-42b1-9e74-74e6a109ad11",
							"security_group_ids": ["8db277ab-1d98-4148-ba72-724721789427"],
							"fixed_ips": [
								{
								  "subnet_id": "863b20c0-c011-4650-85c2-ad531f4570a4",
								  "ip_address": "10.63.177.162"
								}
							],
							"floating": {
								"id": "e906af80-ce13-4ec3-9fba-fa20581c2695",
							  	"floating_network_id": "c87774b5-95a4-4efb-8e6b-883e2212d67b",
							  	"floating_ip_address": "10.0.0.20",
								"remove": false
							},
							"id": "43aa53d7-a70b-4f40-812f-4feecb687018",
							"remove": true
						}
					],
					"placement": {
						"zone": "nova"
					}
				},
				"dependents": {},
				"domain": "1235be1e-8d8e-43bb-bd6c-943eccf76a6d",
				"id": "82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1",
				"index": 2,
				"init_at": null,
				"metadata": {},
				"name": "node-e395be1e-002",
				"physical_id": "66a81d68-bf48-4af5-897b-a3bfef7279a8",
				"profile_id": "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
				"profile_name": "pcirros",
				"project_id": "eee0b7c083e84501bdd50fb269d2a10e",
				"role": "",
				"status": "ACTIVE",
				"status_reason": "Creation succeeded",
				"updated_at": null,
				"user": "ab79b9647d074e46ac223a8fa297b846"				
			}
		}`)
	})

	opts := nodes.CreateOpts{
		ClusterID: "e395be1e-8d8e-43bb-bd6c-943eccf76a6d",
		Metadata:  map[string]interface{}{},
		Name:      "node-e395be1e-002",
		ProfileID: "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
		Role:      "",
	}

	actual, err := nodes.Create(fake.ServiceClient(), opts).Extract()
	if err != nil {
		t.Error("Error creating nodes. error=", err)
	} else {
		expected := nodes.Node{
			ClusterID: "e395be1e-8d8e-43bb-bd6c-943eccf76a6d",
			CreatedAt: time.Time{},
			Data: map[string]interface{}{
				"internal_ports": []map[string]interface{}{
					{
						"network_id": "847e4f65-1ff1-42b1-9e74-74e6a109ad11",
						"security_group_ids": []interface{}{
							"8db277ab-1d98-4148-ba72-724721789427",
						},
						"fixed_ips": []interface{}{
							map[string]interface{}{
								"subnet_id":  "863b20c0-c011-4650-85c2-ad531f4570a4",
								"ip_address": "10.63.177.162",
							},
						},
						"floating": map[string]interface{}{
							"floating_network_id": "c87774b5-95a4-4efb-8e6b-883e2212d67b",
							"floating_ip_address": "10.0.0.20",
							"remove":              false,
							"id":                  "e906af80-ce13-4ec3-9fba-fa20581c2695",
						},
						"id":     "43aa53d7-a70b-4f40-812f-4feecb687018",
						"remove": true,
					},
				},
				"placement": map[string]interface{}{
					"zone": "nova",
				},
			},

			Dependents:   map[string]interface{}{},
			Domain:       "1235be1e-8d8e-43bb-bd6c-943eccf76a6d",
			ID:           "82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1",
			Index:        2,
			InitAt:       time.Time{},
			Metadata:     map[string]interface{}{},
			Name:         "node-e395be1e-002",
			PhysicalID:   "66a81d68-bf48-4af5-897b-a3bfef7279a8",
			ProfileID:    "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
			ProfileName:  "pcirros",
			ProjectID:    "eee0b7c083e84501bdd50fb269d2a10e",
			Role:         "",
			Status:       "ACTIVE",
			StatusReason: "Creation succeeded",
			UpdatedAt:    time.Time{},
			User:         "ab79b9647d074e46ac223a8fa297b846",
		}
		th.AssertDeepEquals(t, expected, *actual)
	}
}

func TestCreateNodeInvalidTimeFloat(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/nodes", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
			"node": {
				"cluster_id": "e395be1e-8d8e-43bb-bd6c-943eccf76a6d",
				"created_at": 123456789.0,
				"data": {
					"internal_ports": [
						{
							"network_id": "847e4f65-1ff1-42b1-9e74-74e6a109ad11",
							"security_group_ids": ["8db277ab-1d98-4148-ba72-724721789427"],
							"fixed_ips": [
								{
								  "subnet_id": "863b20c0-c011-4650-85c2-ad531f4570a4",
								  "ip_address": "10.63.177.162"
								}
							],
							"id": "43aa53d7-a70b-4f40-812f-4feecb687018",
							"remove": true
						}
					],
					"placement": {
						"zone": "nova"
					}
				},
				"dependents": {},
				"domain": "1235be1e-8d8e-43bb-bd6c-943eccf76a6d",
				"id": "82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1",
				"index": 2,
				"init_at": 123456789.0,
				"metadata": {},
				"name": "node-e395be1e-002",
				"physical_id": "66a81d68-bf48-4af5-897b-a3bfef7279a8",
				"profile_id": "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
				"profile_name": "pcirros",
				"project_id": "eee0b7c083e84501bdd50fb269d2a10e",
				"role": "",
				"status": "ACTIVE",
				"status_reason": "Creation succeeded",
				"updated_at": 123456789.0,
				"user": "ab79b9647d074e46ac223a8fa297b846"				
			}
		}`)
	})

	opts := nodes.CreateOpts{
		ClusterID: "e395be1e-8d8e-43bb-bd6c-943eccf76a6d",
		Metadata:  map[string]interface{}{},
		Name:      "node-e395be1e-002",
		ProfileID: "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
		Role:      "",
	}

	_, err := nodes.Create(fake.ServiceClient(), opts).Extract()
	th.AssertEquals(t, false, err == nil)
}

func TestCreateNodeInvalidTimeString(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/nodes", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
			"node": {
				"cluster_id": "e395be1e-8d8e-43bb-bd6c-943eccf76a6d",
				"created_at": "invalid",
				"data": {
					"internal_ports": [
						{
							"network_id": "847e4f65-1ff1-42b1-9e74-74e6a109ad11",
							"security_group_ids": ["8db277ab-1d98-4148-ba72-724721789427"],
							"fixed_ips": [
								{
								  "subnet_id": "863b20c0-c011-4650-85c2-ad531f4570a4",
								  "ip_address": "10.63.177.162"
								}
							],
							"id": "43aa53d7-a70b-4f40-812f-4feecb687018",
							"remove": true
						}
					],
					"placement": {
						"zone": "nova"
					}
				},
				"dependents": {},
				"domain": "1235be1e-8d8e-43bb-bd6c-943eccf76a6d",
				"id": "82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1",
				"index": 2,
				"init_at": "invalid",
				"metadata": {},
				"name": "node-e395be1e-002",
				"physical_id": "66a81d68-bf48-4af5-897b-a3bfef7279a8",
				"profile_id": "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
				"profile_name": "pcirros",
				"project_id": "eee0b7c083e84501bdd50fb269d2a10e",
				"role": "",
				"status": "ACTIVE",
				"status_reason": "Creation succeeded",
				"updated_at": "invalid",
				"user": "ab79b9647d074e46ac223a8fa297b846"				
			}
		}`)
	})

	opts := nodes.CreateOpts{
		ClusterID: "e395be1e-8d8e-43bb-bd6c-943eccf76a6d",
		Metadata:  map[string]interface{}{},
		Name:      "node-e395be1e-002",
		ProfileID: "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
		Role:      "",
	}

	_, err := nodes.Create(fake.ServiceClient(), opts).Extract()
	th.AssertEquals(t, false, err == nil)
}
