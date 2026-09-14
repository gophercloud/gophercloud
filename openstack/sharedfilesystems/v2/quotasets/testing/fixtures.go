package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	client "github.com/gophercloud/gophercloud/testhelper/client"
)

var (
	ShareType = "default"
	tenantID = "7b8b4e4e93774e07ab9a05867129fd17"
	userID = "d8a9bde6cb724fad9e6ec83cbc7f43e9"
)

const ExpectedInitialQuotaSet = `
{
	quota_set: {
		"giga_bytes": 10,
		"snapshots": 10,
		"shares": 10,
		"snapchot_gigabytes": 10,
		"share_networks": 10,
		"share_groups": 10,
		"share_group_snapshots": 10,
		"share_replicas": 10,
		"share_replica_gigabytes": 10,
		"per_share_gigabytes": 10,
	}
}
`

const ExpectedUpdatedQuotaSet = `
{
	quota_set: {
		"giga_bytes": 100,
		"snapshots": 100,
		"shares": 100,
		"snapchot_gigabytes": 100,
		"share_networks": 100,
		"share_groups": 100,
		"share_group_snapshots": 100,
		"share_replicas": 100,
		"share_replica_gigabytes": 100,
		"per_share_gigabytes": 100,
	}
}
`

func HandleGetQuotaSetSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/quota-sets/"+tenantID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ExpectedInitialQuotaSet)
	})
}

func HandleUpdateQuotaSetSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/quota-sets/"+tenantID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ExpectedUpdatedQuotaSet)
	})
}

func HandleGetByShareTypeSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/quota-sets/"+tenantID+"?share_type="+ShareType, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ExpectedInitialQuotaSet)
	})
}

func HandleUpdateByShareTypeSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/quota-sets/"+tenantID+"?share_type="+ShareType, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ExpectedUpdatedQuotaSet)
	})
}

func HandleGetByUserSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/quota-sets/"+tenantID+"?user_id="+userID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ExpectedInitialQuotaSet)
	})
}

func HandleUpdateByUserSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/quota-sets/"+tenantID+"?user_id="+userID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ExpectedUpdatedQuotaSet)
	})
}
