package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/clustering/v1/clusters"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestCreateCluster(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/clusters", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("X-OpenStack-Request-ID", "req-781e9bdc-4163-46eb-91c9-786c53188bbb")
		w.Header().Add("Location", "http://senlin.cloud.blizzard.net:8778/v1/actions/625628cd-f877-44be-bde0-fec79f84e13d")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
			"cluster": {
				"config": {},
				"created_at": "2015-02-10T14:26:14Z",
				"data": {},
				"dependents": {},
				"desired_capacity": 3,
				"domain": null,
				"id": "7d85f602-a948-4a30-afd4-e84f47471c15",
				"init_at": "2015-02-10T15:26:14Z",
				"max_size": 20,
				"metadata": {},
				"min_size": 1,
				"name": "cluster1",
				"nodes": [
					"b07c57c8-7ab2-47bf-bdf8-e894c0c601b9",
					"ecc23d3e-bb68-48f8-8260-c9cf6bcb6e61",
					"da1e9c87-e584-4626-a120-022da5062dac"
				],
				"policies": [],
				"profile_id": "edc63d0a-2ca4-48fa-9854-27926da76a4a",
				"profile_name": "mystack",
				"project": "6e18cc2bdbeb48a5b3cad2dc499f6804",
				"status": "ACTIVE",
				"status_reason": "Cluster scale-in succeeded",
				"timeout": 3600,
				"updated_at": "2015-02-10T16:26:14Z",
				"user": "5e5bf8027826429c96af157f68dc9072"
			}
		}`)
	})

	opts := clusters.CreateOpts{
		Name:            "cluster1",
		DesiredCapacity: 3,
		ProfileID:       "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
		MinSize:         1,
		MaxSize:         20,
		Timeout:         3600,
		Metadata:        map[string]interface{}{},
		Config:          map[string]interface{}{},
	}

	createdAt, _ := time.Parse(time.RFC3339, "2015-02-10T14:26:14Z")
	initAt, _ := time.Parse(time.RFC3339, "2015-02-10T15:26:14Z")
	updatedAt, _ := time.Parse(time.RFC3339, "2015-02-10T16:26:14Z")

	createResult := clusters.Create(fake.ServiceClient(), opts)
	if createResult.Err != nil {
		t.Error("Error creating cluster. error=", createResult.Err)
	}

	location := createResult.Header.Get("Location")
	th.AssertEquals(t, "http://senlin.cloud.blizzard.net:8778/v1/actions/625628cd-f877-44be-bde0-fec79f84e13d", location)

	actionID, err := clusters.ExtractActionFromLocation(location)
	if err != nil {
		t.Error("Error Extracting Action ID. error=", err)
	}
	th.AssertEquals(t, "625628cd-f877-44be-bde0-fec79f84e13d", actionID)

	actual, err := createResult.Extract()
	if err != nil {
		t.Error("Error creating cluster. error=", err)
	} else {
		expected := clusters.Cluster{
			Config:          map[string]interface{}{},
			CreatedAt:       createdAt,
			Data:            map[string]interface{}{},
			Dependents:      map[string]interface{}{},
			DesiredCapacity: 3,
			Domain:          "",
			ID:              "7d85f602-a948-4a30-afd4-e84f47471c15",
			InitAt:          initAt,
			MaxSize:         20,
			Metadata:        map[string]interface{}{},
			MinSize:         1,
			Name:            "cluster1",
			Nodes: []string{
				"b07c57c8-7ab2-47bf-bdf8-e894c0c601b9",
				"ecc23d3e-bb68-48f8-8260-c9cf6bcb6e61",
				"da1e9c87-e584-4626-a120-022da5062dac",
			},
			Policies:     []string{},
			ProfileID:    "edc63d0a-2ca4-48fa-9854-27926da76a4a",
			ProfileName:  "mystack",
			Project:      "6e18cc2bdbeb48a5b3cad2dc499f6804",
			Status:       "ACTIVE",
			StatusReason: "Cluster scale-in succeeded",
			Timeout:      3600,
			UpdatedAt:    updatedAt,
			User:         "5e5bf8027826429c96af157f68dc9072",
		}
		th.AssertDeepEquals(t, expected, *actual)
	}
}

func TestCreateClusterEmptyTime(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/clusters", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
			"cluster": {
				"config": {},
				"created_at": null,
				"data": {},
				"dependents": {},
				"desired_capacity": 3,
				"domain": null,
				"id": "7d85f602-a948-4a30-afd4-e84f47471c15",
				"init_at": null,
				"max_size": 20,
				"metadata": {},
				"min_size": 1,
				"name": "cluster1",
				"nodes": [
					"b07c57c8-7ab2-47bf-bdf8-e894c0c601b9",
					"ecc23d3e-bb68-48f8-8260-c9cf6bcb6e61",
					"da1e9c87-e584-4626-a120-022da5062dac"
				],
				"policies": [],
				"profile_id": "edc63d0a-2ca4-48fa-9854-27926da76a4a",
				"profile_name": "mystack",
				"project": "6e18cc2bdbeb48a5b3cad2dc499f6804",
				"status": "ACTIVE",
				"status_reason": "Cluster scale-in succeeded",
				"timeout": 3600,
				"updated_at": null,
				"user": "5e5bf8027826429c96af157f68dc9072"
			}
		}`)
	})

	opts := clusters.CreateOpts{
		Name:            "cluster1",
		DesiredCapacity: 3,
		ProfileID:       "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
		MinSize:         1,
		MaxSize:         20,
		Timeout:         3600,
		Metadata:        map[string]interface{}{},
		Config:          map[string]interface{}{},
	}

	actual, err := clusters.Create(fake.ServiceClient(), opts).Extract()
	if err != nil {
		t.Error("Error creating cluster. error=", err)
	} else {
		expected := clusters.Cluster{
			Config:          map[string]interface{}{},
			CreatedAt:       time.Time{},
			Data:            map[string]interface{}{},
			Dependents:      map[string]interface{}{},
			DesiredCapacity: 3,
			Domain:          "",
			ID:              "7d85f602-a948-4a30-afd4-e84f47471c15",
			InitAt:          time.Time{},
			MaxSize:         20,
			Metadata:        map[string]interface{}{},
			MinSize:         1,
			Name:            "cluster1",
			Nodes: []string{
				"b07c57c8-7ab2-47bf-bdf8-e894c0c601b9",
				"ecc23d3e-bb68-48f8-8260-c9cf6bcb6e61",
				"da1e9c87-e584-4626-a120-022da5062dac",
			},
			Policies:     []string{},
			ProfileID:    "edc63d0a-2ca4-48fa-9854-27926da76a4a",
			ProfileName:  "mystack",
			Project:      "6e18cc2bdbeb48a5b3cad2dc499f6804",
			Status:       "ACTIVE",
			StatusReason: "Cluster scale-in succeeded",
			Timeout:      3600,
			UpdatedAt:    time.Time{},
			User:         "5e5bf8027826429c96af157f68dc9072",
		}
		th.AssertDeepEquals(t, expected, *actual)
	}
}

func TestCreateClusterInvalidTimeFloat(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/clusters", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
			"cluster": {
				"config": {},
				"created_at": 123456789.0,
				"data": {},
				"dependents": {},
				"desired_capacity": 3,
				"domain": null,
				"id": "7d85f602-a948-4a30-afd4-e84f47471c15",
				"init_at": 123456789.0,
				"max_size": 20,
				"metadata": {},
				"min_size": 1,
				"name": "cluster1",
				"nodes": [
					"b07c57c8-7ab2-47bf-bdf8-e894c0c601b9",
					"ecc23d3e-bb68-48f8-8260-c9cf6bcb6e61",
					"da1e9c87-e584-4626-a120-022da5062dac"
				],
				"policies": [],
				"profile_id": "edc63d0a-2ca4-48fa-9854-27926da76a4a",
				"profile_name": "mystack",
				"project": "6e18cc2bdbeb48a5b3cad2dc499f6804",
				"status": "ACTIVE",
				"status_reason": "Cluster scale-in succeeded",
				"timeout": 3600,
				"updated_at": 123456789.0,
				"user": "5e5bf8027826429c96af157f68dc9072"
			}
		}`)
	})

	opts := clusters.CreateOpts{
		Name:            "cluster1",
		DesiredCapacity: 3,
		ProfileID:       "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
		MinSize:         1,
		MaxSize:         20,
		Timeout:         3600,
		Metadata:        map[string]interface{}{},
		Config:          map[string]interface{}{},
	}

	_, err := clusters.Create(fake.ServiceClient(), opts).Extract()
	th.AssertEquals(t, false, err == nil)
}

func TestCreateClusterInvalidTimeString(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/clusters", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
			"cluster": {
				"config": {},
				"created_at": "invalid",
				"data": {},
				"dependents": {},
				"desired_capacity": 3,
				"domain": null,
				"id": "7d85f602-a948-4a30-afd4-e84f47471c15",
				"init_at": "invalid",
				"max_size": 20,
				"metadata": {},
				"min_size": 1,
				"name": "cluster1",
				"nodes": [
					"b07c57c8-7ab2-47bf-bdf8-e894c0c601b9",
					"ecc23d3e-bb68-48f8-8260-c9cf6bcb6e61",
					"da1e9c87-e584-4626-a120-022da5062dac"
				],
				"policies": [],
				"profile_id": "edc63d0a-2ca4-48fa-9854-27926da76a4a",
				"profile_name": "mystack",
				"project": "6e18cc2bdbeb48a5b3cad2dc499f6804",
				"status": "ACTIVE",
				"status_reason": "Cluster scale-in succeeded",
				"timeout": 3600,
				"updated_at": "invalid",
				"user": "5e5bf8027826429c96af157f68dc9072"
			}
		}`)
	})

	opts := clusters.CreateOpts{
		Name:            "cluster1",
		DesiredCapacity: 3,
		ProfileID:       "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
		MinSize:         1,
		MaxSize:         20,
		Timeout:         3600,
		Metadata:        map[string]interface{}{},
		Config:          map[string]interface{}{},
	}

	_, err := clusters.Create(fake.ServiceClient(), opts).Extract()
	th.AssertEquals(t, false, err == nil)
}
