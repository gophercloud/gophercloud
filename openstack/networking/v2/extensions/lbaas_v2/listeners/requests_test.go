package listeners

import (
	"fmt"
	"net/http"
	"testing"
	//"github.com/davecgh/go-spew/spew"
	fake "github.com/rackspace/gophercloud/openstack/networking/v2/common"
	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
)

func TestURLs(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.AssertEquals(t, th.Endpoint()+"v2.0/lbaas/listeners", rootURL(fake.ServiceClient()))
	th.AssertEquals(t, th.Endpoint()+"v2.0/lbaas/listeners/foo", resourceURL(fake.ServiceClient(), "foo"))
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/lbaas/listeners", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "listeners":[
         {
           "id": "db902c0c-d5ff-4753-b465-668ad9656918",
           "tenant_id": "310df60f-2a10-4ee5-9554-98393092194c",
           "name": "web_listener",
           "description": "listener config for the web tier",
           "loadbalancers": [{"id": "53306cda-815d-4354-9444-59e09da9c3c5"}],
           "protocol": "HTTP",
           "protocol_port": 80,
           "default_pool_id": "fad389a3-9a4a-4762-a365-8c7038508b5d",
           "admin_state_up": true,
           "default_tls_container_ref": "2c433435-20de-4411-84ae-9cc8917def76",
           "sni_container_refs": ["3d328d82-2547-4921-ac2f-61c3b452b5ff", "b3cfd7e3-8c19-455c-8ebb-d78dfd8f7e7d"]
         },
         {
           "id": "36e08a3e-a78f-4b40-a229-1e7e23eee1ab",
           "tenant_id": "310df60f-2a10-4ee5-9554-98393092194c",
           "name": "db_listener",
	   "description": "listener config for the db tier",
           "loadbalancers": [{"id": "79e05663-7f03-45d2-a092-8b94062f22ab"}],
           "protocol": "TCP",
           "protocol_port": 3306,
           "default_pool_id": "41efe233-7591-43c5-9cf7-923964759f9e",
           "connection_limit": 2000,
           "admin_state_up": true,
           "default_tls_container_ref": "2c433435-20de-4411-84ae-9cc8917def76",
           "sni_container_refs": ["3d328d82-2547-4921-ac2f-61c3b452b5ff", "b3cfd7e3-8c19-455c-8ebb-d78dfd8f7e7d"]
         }
      ]
}
			`)
	})

	count := 0

	List(fake.ServiceClient(), ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractListeners(page)
		if err != nil {
			t.Errorf("Failed to extract LBs: %v", err)
			return false, err
		}

		expected := []Listener{
			{
				ID:                     "db902c0c-d5ff-4753-b465-668ad9656918",
				TenantID:               "310df60f-2a10-4ee5-9554-98393092194c",
				Name:                   "web_listener",
				Description:            "listener config for the web tier",
				Loadbalancers:          []map[string]interface{}{{"id": "53306cda-815d-4354-9444-59e09da9c3c5"}},
				Protocol:               "HTTP",
				ProtocolPort:           80,
				DefaultPoolID:          "fad389a3-9a4a-4762-a365-8c7038508b5d",
				AdminStateUp:           true,
				DefaultTlsContainerRef: "2c433435-20de-4411-84ae-9cc8917def76",
				SniContainerRefs:       []string{"3d328d82-2547-4921-ac2f-61c3b452b5ff", "b3cfd7e3-8c19-455c-8ebb-d78dfd8f7e7d"},
			},
			{
				ID:                     "36e08a3e-a78f-4b40-a229-1e7e23eee1ab",
				TenantID:               "310df60f-2a10-4ee5-9554-98393092194c",
				Name:                   "db_listener",
				Description:            "listener config for the db tier",
				Loadbalancers:          []map[string]interface{}{{"id": "79e05663-7f03-45d2-a092-8b94062f22ab"}},
				Protocol:               "TCP",
				ProtocolPort:           3306,
				DefaultPoolID:          "41efe233-7591-43c5-9cf7-923964759f9e",
				ConnLimit:              2000,
				AdminStateUp:           true,
				DefaultTlsContainerRef: "2c433435-20de-4411-84ae-9cc8917def76",
				SniContainerRefs:       []string{"3d328d82-2547-4921-ac2f-61c3b452b5ff", "b3cfd7e3-8c19-455c-8ebb-d78dfd8f7e7d"},
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/lbaas/listeners", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "listener": {
        "loadbalancer_id": "53306cda-815d-4354-9444-59e09da9c3c5",
        "protocol": "HTTP",
        "name": "NewListener",
        "admin_state_up": true,
        "default_tls_container_ref": "8032909d-47a1-4715-90af-5153ffe39861",
        "default_pool_id": "61b1f87a-7a21-4ad3-9dda-7f81d249944f",
        "protocol_port": 80
    }
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
    "listener": {
        "id": "36e08a3e-a78f-4b40-a229-1e7e23eee1ab",
	"tenant_id": "83657cfcdfe44cd5920adaf26c48ceea",
	"name": "NewListener",
	"description": "",
	"loadbalancers": [{"id": "53306cda-815d-4354-9444-59e09da9c3c5"}],
	"protocol": "HTTP",
	"protocol_port": 80,
	"connection_limit": -1,
	"default_pool_id": "61b1f87a-7a21-4ad3-9dda-7f81d249944f",
	"admin_state_up": true,
	"default_tls_container_ref": "8032909d-47a1-4715-90af-5153ffe39861"
    }
}
		`)
	})

	opts := CreateOpts{
		Protocol:               "HTTP",
		Name:                   "NewListener",
		LoadbalancerID:         "53306cda-815d-4354-9444-59e09da9c3c5",
		AdminStateUp:           Up,
		DefaultTlsContainerRef: "8032909d-47a1-4715-90af-5153ffe39861",
		DefaultPoolID:          "61b1f87a-7a21-4ad3-9dda-7f81d249944f",
		ProtocolPort:           80,
	}

	r, err := Create(fake.ServiceClient(), opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "HTTP", r.Protocol)
	th.AssertEquals(t, "", r.Description)
	th.AssertEquals(t, true, r.AdminStateUp)
	th.AssertEquals(t, "8032909d-47a1-4715-90af-5153ffe39861", r.DefaultTlsContainerRef)
	th.AssertEquals(t, "83657cfcdfe44cd5920adaf26c48ceea", r.TenantID)
	th.AssertEquals(t, -1, r.ConnLimit)
	th.AssertEquals(t, "61b1f87a-7a21-4ad3-9dda-7f81d249944f", r.DefaultPoolID)
	th.AssertEquals(t, 80, r.ProtocolPort)
	th.AssertEquals(t, "36e08a3e-a78f-4b40-a229-1e7e23eee1ab", r.ID)
	th.AssertEquals(t, "53306cda-815d-4354-9444-59e09da9c3c5", r.Loadbalancers[0]["id"])
	th.AssertEquals(t, "NewListener", r.Name)
}

func TestRequiredCreateOpts(t *testing.T) {
	res := Create(fake.ServiceClient(), CreateOpts{})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
	res = Create(fake.ServiceClient(), CreateOpts{Name: "foo"})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
	res = Create(fake.ServiceClient(), CreateOpts{Name: "foo", TenantID: "bar"})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
	res = Create(fake.ServiceClient(), CreateOpts{Name: "foo", TenantID: "bar", Protocol: "bar"})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
	res = Create(fake.ServiceClient(), CreateOpts{Name: "foo", TenantID: "bar", Protocol: "bar", ProtocolPort: 80})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/lbaas/listeners/4ec89087-d057-4e2c-911f-60a3b47ee304", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "listener": {
        "id": "4ec89087-d057-4e2c-911f-60a3b47ee304",
	"tenant_id": "83657cfcdfe44cd5920adaf26c48ceea",
	"name": "NewListener",
	"description": "",
	"loadbalancers": [{"id": "53306cda-815d-4354-9444-59e09da9c3c5"}],
	"protocol": "HTTP",
	"protocol_port": 80,
	"connection_limit": -1,
	"default_pool_id": "61b1f87a-7a21-4ad3-9dda-7f81d249944f",
	"admin_state_up": true,
	"default_tls_container_ref": "8032909d-47a1-4715-90af-5153ffe39861"
    }
}
			`)
	})

	l, err := Get(fake.ServiceClient(), "4ec89087-d057-4e2c-911f-60a3b47ee304").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "HTTP", l.Protocol)
	th.AssertEquals(t, "", l.Description)
	th.AssertEquals(t, true, l.AdminStateUp)
	th.AssertEquals(t, "8032909d-47a1-4715-90af-5153ffe39861", l.DefaultTlsContainerRef)
	th.AssertEquals(t, "83657cfcdfe44cd5920adaf26c48ceea", l.TenantID)
	th.AssertEquals(t, -1, l.ConnLimit)
	th.AssertEquals(t, "61b1f87a-7a21-4ad3-9dda-7f81d249944f", l.DefaultPoolID)
	th.AssertEquals(t, 80, l.ProtocolPort)
	th.AssertEquals(t, "4ec89087-d057-4e2c-911f-60a3b47ee304", l.ID)
	th.AssertEquals(t, "53306cda-815d-4354-9444-59e09da9c3c5", l.Loadbalancers[0]["id"])
	th.AssertEquals(t, "NewListener", l.Name)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/lbaas/listeners/4ec89087-d057-4e2c-911f-60a3b47ee304", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "listener": {
        "name": "NewListenerName",
        "connection_limit": 1001
    }
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

		fmt.Fprintf(w, `
{
    "listener": {
        "id": "4ec89087-d057-4e2c-911f-60a3b47ee304",
	"tenant_id": "83657cfcdfe44cd5920adaf26c48ceea",
	"name": "NewListenerName",
	"description": "",
	"loadbalancers": [{"id": "53306cda-815d-4354-9444-59e09da9c3c5"}],
	"protocol": "HTTP",
	"protocol_port": 80,
	"connection_limit": 1001,
	"default_pool_id": "61b1f87a-7a21-4ad3-9dda-7f81d249944f",
	"admin_state_up": true,
	"default_tls_container_ref": "8032909d-47a1-4715-90af-5153ffe39861"
    }
}
		`)
	})

	i1001 := 1001
	options := UpdateOpts{
		Name:      "NewListenerName",
		ConnLimit: &i1001,
	}

	l, err := Update(fake.ServiceClient(), "4ec89087-d057-4e2c-911f-60a3b47ee304", options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, *(options.ConnLimit), l.ConnLimit)
	th.AssertEquals(t, options.Name, l.Name)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/lbaas/listeners/4ec89087-d057-4e2c-911f-60a3b47ee304", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := Delete(fake.ServiceClient(), "4ec89087-d057-4e2c-911f-60a3b47ee304")
	th.AssertNoErr(t, res.Err)
}
