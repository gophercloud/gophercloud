package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgpvpns"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	filterProjectID := []string{"b7549121395844bea941bb92feb3fad9"}
	fields := []string{"id", "name"}
	listOpts := bgpvpns.ListOpts{
		Fields:    fields,
		ProjectID: filterProjectID[0],
	}
	th.Mux.HandleFunc("/v2.0/bgpvpn/bgpvpns",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

			r.ParseForm()
			th.AssertDeepEquals(t, r.Form["fields"], fields)
			th.AssertDeepEquals(t, r.Form["project_id"], filterProjectID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, ListBGPVPNsResult)
		})
	count := 0

	err := bgpvpns.List(fake.ServiceClient(), listOpts).EachPage(
		context.TODO(),
		func(_ context.Context, page pagination.Page) (bool, error) {
			count++
			actual, err := bgpvpns.ExtractBGPVPNs(page)
			if err != nil {
				t.Errorf("Failed to extract BGP VPNs: %v", err)
				return false, nil
			}

			expected := []bgpvpns.BGPVPN{BGPVPN}
			th.CheckDeepEquals(t, expected, actual)

			return true, nil
		})
	th.AssertNoErr(t, err)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	bgpVpnID := "460ac411-3dfb-45bb-8116-ed1a7233d143"
	th.Mux.HandleFunc("/v2.0/bgpvpn/bgpvpns/"+bgpVpnID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetBGPVPNResult)
	})

	r, err := bgpvpns.Get(context.TODO(), fake.ServiceClient(), bgpVpnID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, GetBGPVPN, *r)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/bgpvpn/bgpvpns", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, CreateResponse)
	})

	opts := bgpvpns.CreateOpts{
		TenantID: "b7549121395844bea941bb92feb3fad9",
		RouteTargets: []string{
			"64512:1444",
		},
		ImportTargets: []string{
			"64512:1555",
		},
		ExportTargets: []string{
			"64512:1666",
		},
		RouteDistinguishers: []string{
			"64512:1777",
			"64512:1888",
			"64512:1999",
		},
		Type: "l3",
		VNI:  1000,
	}

	r, err := bgpvpns.Create(context.TODO(), fake.ServiceClient(), opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, CreateBGPVPN, *r)
}

func TestDelete(t *testing.T) {
	bgpVpnID := "0f9d472a-908f-40f5-8574-b4e8a63ccbf0"
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/bgpvpn/bgpvpns/"+bgpVpnID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	err := bgpvpns.Delete(context.TODO(), fake.ServiceClient(), bgpVpnID).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestUpdate(t *testing.T) {
	bgpVpnID := "4d627abf-06dd-45ab-920b-8e61422bb984"
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/bgpvpn/bgpvpns/"+bgpVpnID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdateBGPVPNRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, UpdateBGPVPNResponse)
	})

	name := "foo"
	routeTargets := []string{"64512:1444"}
	emptyTarget := []string{}
	opts := bgpvpns.UpdateOpts{
		Name:          &name,
		RouteTargets:  &routeTargets,
		ImportTargets: &emptyTarget,
		ExportTargets: &emptyTarget,
	}

	r, err := bgpvpns.Update(context.TODO(), fake.ServiceClient(), bgpVpnID, opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, *opts.Name, r.Name)
}

func TestListNetworkAssociations(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	bgpVpnID := "460ac411-3dfb-45bb-8116-ed1a7233d143"
	fields := []string{"id", "name"}
	listOpts := bgpvpns.ListNetworkAssociationsOpts{
		Fields: fields,
	}
	th.Mux.HandleFunc("/v2.0/bgpvpn/bgpvpns/"+bgpVpnID+"/network_associations", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		r.ParseForm()
		th.AssertDeepEquals(t, fields, r.Form["fields"])

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListNetworkAssociationsResult)
	})

	count := 0
	err := bgpvpns.ListNetworkAssociations(fake.ServiceClient(), bgpVpnID, listOpts).EachPage(
		context.TODO(),
		func(_ context.Context, page pagination.Page) (bool, error) {
			count++
			actual, err := bgpvpns.ExtractNetworkAssociations(page)
			if err != nil {
				t.Errorf("Failed to extract network associations: %v", err)
				return false, nil
			}

			expected := []bgpvpns.NetworkAssociation{NetworkAssociation}
			th.CheckDeepEquals(t, expected, actual)

			return true, nil
		})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, count)
}

func TestCreateNetworkAssociation(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	bgpVpnID := "460ac411-3dfb-45bb-8116-ed1a7233d143"
	th.Mux.HandleFunc("/v2.0/bgpvpn/bgpvpns/"+bgpVpnID+"/network_associations", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateNetworkAssociationRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, CreateNetworkAssociationResponse)
	})

	opts := bgpvpns.CreateNetworkAssociationOpts{
		NetworkID: "8c5d88dc-60ac-4b02-a65a-36b65888ddcd",
	}
	r, err := bgpvpns.CreateNetworkAssociation(context.TODO(), fake.ServiceClient(), bgpVpnID, opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, CreateNetworkAssociation, *r)
}

func TestGetNetworkAssociation(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	bgpVpnID := "460ac411-3dfb-45bb-8116-ed1a7233d143"
	networkAssociationID := "73238ca1-e05d-4c7a-b4d4-70407b4b8730"
	th.Mux.HandleFunc("/v2.0/bgpvpn/bgpvpns/"+bgpVpnID+"/network_associations/"+networkAssociationID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetNetworkAssociationResult)
	})

	r, err := bgpvpns.GetNetworkAssociation(context.TODO(), fake.ServiceClient(), bgpVpnID, networkAssociationID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, GetNetworkAssociation, *r)
}

func TestDeleteNetworkAssociation(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	bgpVpnID := "460ac411-3dfb-45bb-8116-ed1a7233d143"
	networkAssociationID := "73238ca1-e05d-4c7a-b4d4-70407b4b8730"
	th.Mux.HandleFunc("/v2.0/bgpvpn/bgpvpns/"+bgpVpnID+"/network_associations/"+networkAssociationID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	err := bgpvpns.DeleteNetworkAssociation(context.TODO(), fake.ServiceClient(), bgpVpnID, networkAssociationID).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestListRouterAssociations(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	bgpVpnID := "460ac411-3dfb-45bb-8116-ed1a7233d143"
	fields := []string{"id", "name"}
	listOpts := bgpvpns.ListRouterAssociationsOpts{
		Fields: fields,
	}
	th.Mux.HandleFunc("/v2.0/bgpvpn/bgpvpns/"+bgpVpnID+"/router_associations", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		r.ParseForm()
		th.AssertDeepEquals(t, fields, r.Form["fields"])

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListRouterAssociationsResult)
	})

	count := 0
	err := bgpvpns.ListRouterAssociations(fake.ServiceClient(), bgpVpnID, listOpts).EachPage(
		context.TODO(),
		func(_ context.Context, page pagination.Page) (bool, error) {
			count++
			actual, err := bgpvpns.ExtractRouterAssociations(page)
			if err != nil {
				t.Errorf("Failed to extract router associations: %v", err)
				return false, nil
			}

			expected := []bgpvpns.RouterAssociation{RouterAssociation}
			th.CheckDeepEquals(t, expected, actual)

			return true, nil
		})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, count)
}

func TestCreateRouterAssociation(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	bgpVpnID := "460ac411-3dfb-45bb-8116-ed1a7233d143"
	th.Mux.HandleFunc("/v2.0/bgpvpn/bgpvpns/"+bgpVpnID+"/router_associations", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateRouterAssociationRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, CreateRouterAssociationResponse)
	})

	opts := bgpvpns.CreateRouterAssociationOpts{
		RouterID: "8c5d88dc-60ac-4b02-a65a-36b65888ddcd",
	}
	r, err := bgpvpns.CreateRouterAssociation(context.TODO(), fake.ServiceClient(), bgpVpnID, opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, CreateRouterAssociation, *r)
}

func TestGetRouterAssociation(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	bgpVpnID := "460ac411-3dfb-45bb-8116-ed1a7233d143"
	routerAssociationID := "73238ca1-e05d-4c7a-b4d4-70407b4b8730"
	th.Mux.HandleFunc("/v2.0/bgpvpn/bgpvpns/"+bgpVpnID+"/router_associations/"+routerAssociationID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetRouterAssociationResult)
	})

	r, err := bgpvpns.GetRouterAssociation(context.TODO(), fake.ServiceClient(), bgpVpnID, routerAssociationID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, GetRouterAssociation, *r)
}

func TestUpdateRouterAssociation(t *testing.T) {
	bgpVpnID := "4d627abf-06dd-45ab-920b-8e61422bb984"
	routerAssociationID := "73238ca1-e05d-4c7a-b4d4-70407b4b8730"
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/bgpvpn/bgpvpns/"+bgpVpnID+"/router_associations/"+routerAssociationID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdateRouterAssociationRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, UpdateRouterAssociationResponse)
	})

	opts := bgpvpns.UpdateRouterAssociationOpts{
		AdvertiseExtraRoutes: new(bool),
	}
	r, err := bgpvpns.UpdateRouterAssociation(context.TODO(), fake.ServiceClient(), bgpVpnID, routerAssociationID, opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, UpdateRouterAssociation, *r)
}

func TestDeleteRouterAssociation(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	bgpVpnID := "460ac411-3dfb-45bb-8116-ed1a7233d143"
	routerAssociationID := "73238ca1-e05d-4c7a-b4d4-70407b4b8730"
	th.Mux.HandleFunc("/v2.0/bgpvpn/bgpvpns/"+bgpVpnID+"/router_associations/"+routerAssociationID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	err := bgpvpns.DeleteRouterAssociation(context.TODO(), fake.ServiceClient(), bgpVpnID, routerAssociationID).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestListPortAssociations(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	bgpVpnID := "460ac411-3dfb-45bb-8116-ed1a7233d143"
	fields := []string{"id", "name"}
	listOpts := bgpvpns.ListPortAssociationsOpts{
		Fields: fields,
	}
	th.Mux.HandleFunc("/v2.0/bgpvpn/bgpvpns/"+bgpVpnID+"/port_associations", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		r.ParseForm()
		th.AssertDeepEquals(t, fields, r.Form["fields"])

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListPortAssociationsResult)
	})

	count := 0
	err := bgpvpns.ListPortAssociations(fake.ServiceClient(), bgpVpnID, listOpts).EachPage(
		context.TODO(),
		func(_ context.Context, page pagination.Page) (bool, error) {
			count++
			actual, err := bgpvpns.ExtractPortAssociations(page)
			if err != nil {
				t.Errorf("Failed to extract port associations: %v", err)
				return false, nil
			}

			expected := []bgpvpns.PortAssociation{PortAssociation}
			th.CheckDeepEquals(t, expected, actual)

			return true, nil
		})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, count)
}

func TestCreatePortAssociation(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	bgpVpnID := "460ac411-3dfb-45bb-8116-ed1a7233d143"
	th.Mux.HandleFunc("/v2.0/bgpvpn/bgpvpns/"+bgpVpnID+"/port_associations", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreatePortAssociationRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, CreatePortAssociationResponse)
	})

	opts := bgpvpns.CreatePortAssociationOpts{
		PortID: "8c5d88dc-60ac-4b02-a65a-36b65888ddcd",
	}
	r, err := bgpvpns.CreatePortAssociation(context.TODO(), fake.ServiceClient(), bgpVpnID, opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, CreatePortAssociation, *r)
}

func TestGetPortAssociation(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	bgpVpnID := "460ac411-3dfb-45bb-8116-ed1a7233d143"
	portAssociationID := "73238ca1-e05d-4c7a-b4d4-70407b4b8730"
	th.Mux.HandleFunc("/v2.0/bgpvpn/bgpvpns/"+bgpVpnID+"/port_associations/"+portAssociationID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetPortAssociationResult)
	})

	r, err := bgpvpns.GetPortAssociation(context.TODO(), fake.ServiceClient(), bgpVpnID, portAssociationID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, GetPortAssociation, *r)
}

func TestUpdatePortAssociation(t *testing.T) {
	bgpVpnID := "4d627abf-06dd-45ab-920b-8e61422bb984"
	portAssociationID := "73238ca1-e05d-4c7a-b4d4-70407b4b8730"
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/bgpvpn/bgpvpns/"+bgpVpnID+"/port_associations/"+portAssociationID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdatePortAssociationRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, UpdatePortAssociationResponse)
	})

	opts := bgpvpns.UpdatePortAssociationOpts{
		AdvertiseFixedIPs: new(bool),
	}
	r, err := bgpvpns.UpdatePortAssociation(context.TODO(), fake.ServiceClient(), bgpVpnID, portAssociationID, opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, UpdatePortAssociation, *r)
}

func TestDeletePortAssociation(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	bgpVpnID := "460ac411-3dfb-45bb-8116-ed1a7233d143"
	portAssociationID := "73238ca1-e05d-4c7a-b4d4-70407b4b8730"
	th.Mux.HandleFunc("/v2.0/bgpvpn/bgpvpns/"+bgpVpnID+"/port_associations/"+portAssociationID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	err := bgpvpns.DeletePortAssociation(context.TODO(), fake.ServiceClient(), bgpVpnID, portAssociationID).ExtractErr()
	th.AssertNoErr(t, err)
}
