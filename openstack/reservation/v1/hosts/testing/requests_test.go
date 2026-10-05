package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/reservation/v1/hosts"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListHosts(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleListHosts(t, fakeServer)

	allPages, err := hosts.List(client.ServiceClient(fakeServer), hosts.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	actual, err := hosts.ExtractHosts(allPages)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedHostsList, actual)

	// Blazar reports a null status and updated_at for an unmodified host.
	th.AssertEquals(t, "", actual[0].Status)
	th.AssertEquals(t, true, actual[0].UpdatedAt == nil)
}

func TestCreateHost(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleCreateHost(t, fakeServer)

	createOpts := hosts.CreateOpts{
		Name: "compute-1.example.com",
		ExtraCapabilities: map[string]any{
			"gpu":  "a100",
			"rack": "b12",
		},
	}

	actual, err := hosts.Create(context.TODO(), client.ServiceClient(fakeServer), createOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &ExpectedHostWithCapabilities, actual)
}

// Blazar requires a name, so the request should not reach the server without one.
func TestCreateHostWithoutName(t *testing.T) {
	_, err := hosts.CreateOpts{}.ToHostCreateMap()
	th.AssertEquals(t, true, err != nil)
}

// An extra capability must not be able to overwrite a host field.
func TestCreateHostCapabilityCollision(t *testing.T) {
	createOpts := hosts.CreateOpts{
		Name:              "compute-1.example.com",
		ExtraCapabilities: map[string]any{"name": "somethingelse"},
	}

	_, err := createOpts.ToHostCreateMap()
	th.AssertEquals(t, true, err != nil)
}

func TestGetHost(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleGetHost(t, fakeServer)

	actual, err := hosts.Get(context.TODO(), client.ServiceClient(fakeServer), "18").Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &ExpectedHostWithCapabilities, actual)
}

func TestUpdateHost(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleUpdateHost(t, fakeServer)

	updateOpts := hosts.UpdateOpts{
		ExtraCapabilities: map[string]any{"gpu": "h100"},
	}

	actual, err := hosts.Update(context.TODO(), client.ServiceClient(fakeServer), "18", updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "h100", actual.ExtraCapabilities["gpu"])
	th.AssertEquals(t, false, actual.UpdatedAt == nil)
}

// A nil capability must reach Blazar as an explicit null, which removes it.
func TestUpdateHostRemoveCapability(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleRemoveHostCapability(t, fakeServer)

	updateOpts := hosts.UpdateOpts{
		ExtraCapabilities: map[string]any{"gpu": nil},
	}

	actual, err := hosts.Update(context.TODO(), client.ServiceClient(fakeServer), "18", updateOpts).Extract()
	th.AssertNoErr(t, err)

	_, ok := actual.ExtraCapabilities["gpu"]
	th.AssertEquals(t, false, ok)
	th.AssertEquals(t, "b12", actual.ExtraCapabilities["rack"])
}

// Blazar rejects an empty body, so the request should not be sent at all.
func TestUpdateHostWithoutCapabilities(t *testing.T) {
	_, err := hosts.UpdateOpts{}.ToHostUpdateMap()
	th.AssertEquals(t, true, err != nil)
}

func TestDeleteHost(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleDeleteHost(t, fakeServer)

	err := hosts.Delete(context.TODO(), client.ServiceClient(fakeServer), "18").ExtractErr()
	th.AssertNoErr(t, err)
}

// Blazar flattens extra capabilities into the host object.
func TestListHostsWithCapabilities(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleListHostsWithCapabilities(t, fakeServer)

	allPages, err := hosts.List(client.ServiceClient(fakeServer), hosts.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	actual, err := hosts.ExtractHosts(allPages)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedHostsListWithCapabilities, actual)

	// The timestamps are tagged "-"
	_, ok := actual[0].ExtraCapabilities["created_at"]
	th.AssertEquals(t, false, ok)
	_, ok = actual[0].ExtraCapabilities["updated_at"]
	th.AssertEquals(t, false, ok)
}

func TestListAllocations(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleListAllocations(t, fakeServer)

	listOpts := hosts.ListAllocationsOpts{
		LeaseID:       "98c3544d-0afe-4251-8556-700c847a127f",
		ReservationID: "c04d56e0-6b31-40b2-a57e-7283958100bd",
	}

	allPages, err := hosts.ListAllocations(client.ServiceClient(fakeServer), listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	actual, err := hosts.ExtractAllocations(allPages)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedAllocationsList, actual)
}

func TestGetAllocation(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleGetAllocation(t, fakeServer)

	getOpts := hosts.GetAllocationOpts{
		LeaseID:       "98c3544d-0afe-4251-8556-700c847a127f",
		ReservationID: "c04d56e0-6b31-40b2-a57e-7283958100bd",
	}

	actual, err := hosts.GetAllocation(context.TODO(), client.ServiceClient(fakeServer), "18", getOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &ExpectedAllocation, actual)
}

func TestListResourceProperties(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleListResourceProperties(t, fakeServer)

	allPages, err := hosts.ListResourceProperties(client.ServiceClient(fakeServer), nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	actual, err := hosts.ExtractResourceProperties(allPages)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, []hosts.ResourceProperty{{Property: "gophercloud_test"}}, actual)
}

func TestListResourcePropertiesDetail(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleListResourceProperties(t, fakeServer)

	listOpts := hosts.ListResourcePropertiesOpts{
		Detail: true,
		All:    true,
	}

	allPages, err := hosts.ListResourceProperties(client.ServiceClient(fakeServer), listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	actual, err := hosts.ExtractResourceProperties(allPages)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedResourcePropertiesDetail, actual)
}

// Blazar reads detail and all as strings, so even "false" would enable them.
// Unset options must be left out of the query.
func TestListResourcePropertiesOptsOmitsFalse(t *testing.T) {
	query, err := hosts.ListResourcePropertiesOpts{}.ToResourcePropertyListQuery()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "", query)
}

func TestUpdateResourceProperty(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleUpdateResourceProperty(t, fakeServer)

	updateOpts := hosts.UpdateResourcePropertyOpts{
		Private: true,
	}

	actual, err := hosts.UpdateResourceProperty(context.TODO(), client.ServiceClient(fakeServer), "gophercloud_test", updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &ExpectedUpdatedResourceProperty, actual)
}

// Making a property public again must send private as false rather than drop
// it.
func TestUpdateResourcePropertyOptsSendsFalse(t *testing.T) {
	b, err := hosts.UpdateResourcePropertyOpts{Private: false}.ToResourcePropertyUpdateMap()
	th.AssertNoErr(t, err)

	private, ok := b["private"]
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, false, private)
}

// Property names are free-form, so characters such as # and ? must be escaped.
func TestUpdateResourcePropertyEscapesName(t *testing.T) {
	for name, escaped := range map[string]string{
		"gpu#model": "gpu%23model",
		"gpu?model": "gpu%3Fmodel",
	} {
		t.Run(name, func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()

			fakeServer.Mux.HandleFunc("/os-hosts/properties/gpu", func(w http.ResponseWriter, r *http.Request) {
				t.Errorf("request for %q reached the gpu property", name)
			})
			fakeServer.Mux.HandleFunc("/os-hosts/properties/"+name, func(w http.ResponseWriter, r *http.Request) {
				th.TestMethod(t, r, "PATCH")
				th.AssertEquals(t, "/os-hosts/properties/"+escaped, r.URL.EscapedPath())

				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				fmt.Fprint(w, ResourcePropertyUpdateResult)
			})

			_, err := hosts.UpdateResourceProperty(context.TODO(), client.ServiceClient(fakeServer), name, hosts.UpdateResourcePropertyOpts{Private: true}).Extract()
			th.AssertNoErr(t, err)
		})
	}
}
