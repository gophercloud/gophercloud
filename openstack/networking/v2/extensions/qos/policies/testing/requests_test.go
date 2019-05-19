package testing

import (
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos/policies"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports/65c0ee9f-d634-4522-8954-51021b570b0d", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprintf(w, GetPortResponse)
		th.AssertNoErr(t, err)
	})

	var p struct {
		ports.Port
		policies.QoSPolicyExt
	}
	err := ports.Get(fake.ServiceClient(), "65c0ee9f-d634-4522-8954-51021b570b0d").ExtractInto(&p)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, p.ID, "65c0ee9f-d634-4522-8954-51021b570b0d")
	th.AssertEquals(t, p.QoSPolicyID, "591e0597-39a6-4665-8149-2111d8de9a08")
}

func TestCreatePort(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreatePortRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_, err := fmt.Fprintf(w, CreatePortResponse)
		th.AssertNoErr(t, err)
	})

	var p struct {
		ports.Port
		policies.QoSPolicyExt
	}
	portCreateOpts := ports.CreateOpts{
		NetworkID:    "a87cc70a-3e15-4acf-8205-9b711a3531b7",
	}
	createOpts := policies.PortCreateOptsExt{
		CreateOptsBuilder: portCreateOpts,
		QoSPolicyID: "591e0597-39a6-4665-8149-2111d8de9a08",
	}
	err := ports.Create(fake.ServiceClient(), createOpts).ExtractInto(&p)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, p.NetworkID, "a87cc70a-3e15-4acf-8205-9b711a3531b7")
	th.AssertEquals(t, p.TenantID, "d6700c0c9ffa4f1cb322cd4a1f3906fa")
	th.AssertEquals(t, p.ID, "65c0ee9f-d634-4522-8954-51021b570b0d")
	th.AssertEquals(t, p.QoSPolicyID, "591e0597-39a6-4665-8149-2111d8de9a08")
}

func TestUpdatePortWithPolicy(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports/65c0ee9f-d634-4522-8954-51021b570b0d", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdatePortWithPolicyRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprintf(w, UpdatePortWithPolicyResponse)
		th.AssertNoErr(t, err)
	})

	policyID := "591e0597-39a6-4665-8149-2111d8de9a08"

	var p struct {
		ports.Port
		policies.QoSPolicyExt
	}
	portUpdateOpts := ports.UpdateOpts{}
	updateOpts := policies.PortUpdateOptsExt{
		UpdateOptsBuilder: portUpdateOpts,
		QoSPolicyID: &policyID,
	}
	err := ports.Update(fake.ServiceClient(), "65c0ee9f-d634-4522-8954-51021b570b0d", updateOpts).ExtractInto(&p)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, p.NetworkID, "a87cc70a-3e15-4acf-8205-9b711a3531b7")
	th.AssertEquals(t, p.TenantID, "d6700c0c9ffa4f1cb322cd4a1f3906fa")
	th.AssertEquals(t, p.ID, "65c0ee9f-d634-4522-8954-51021b570b0d")
	th.AssertEquals(t, p.QoSPolicyID, "591e0597-39a6-4665-8149-2111d8de9a08")
}

func TestUpdatePortWithoutPolicy(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports/65c0ee9f-d634-4522-8954-51021b570b0d", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdatePortWithoutPolicyRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprintf(w, UpdatePortWithoutPolicyResponse)
		th.AssertNoErr(t, err)
	})

	policyID := ""

	var p struct {
		ports.Port
		policies.QoSPolicyExt
	}
	portUpdateOpts := ports.UpdateOpts{}
	updateOpts := policies.PortUpdateOptsExt{
		UpdateOptsBuilder: portUpdateOpts,
		QoSPolicyID: &policyID,
	}
	err := ports.Update(fake.ServiceClient(), "65c0ee9f-d634-4522-8954-51021b570b0d", updateOpts).ExtractInto(&p)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, p.NetworkID, "a87cc70a-3e15-4acf-8205-9b711a3531b7")
	th.AssertEquals(t, p.TenantID, "d6700c0c9ffa4f1cb322cd4a1f3906fa")
	th.AssertEquals(t, p.ID, "65c0ee9f-d634-4522-8954-51021b570b0d")
	th.AssertEquals(t, p.QoSPolicyID, "")
}