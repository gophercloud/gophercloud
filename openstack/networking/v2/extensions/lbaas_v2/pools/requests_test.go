package pools

import (
	"fmt"
	"net/http"
	"testing"

	fake "github.com/rackspace/gophercloud/openstack/networking/v2/common"
	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
)

func TestURLs(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.AssertEquals(t, th.Endpoint()+"v2.0/lbaas/pools", rootURL(fake.ServiceClient()))
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/lbaas/pools", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
   "pools":[
      {
         "lb_algorithm":"ROUND_ROBIN",
         "protocol":"HTTP",
         "description":"",
         "health_monitors":[
            "466c8345-28d8-4f84-a246-e04380b0461d",
            "5d4b5228-33b0-4e60-b225-9b727c1a20e7"
         ],
         "members":[{"id": "53306cda-815d-4354-9fe4-59e09da9c3c5"}],
         "listeners":[{"id": "2a280670-c202-4b0b-a562-34077415aabf"}],
         "loadbalancers":[{"id": "79e05663-7f03-45d2-a092-8b94062f22ab"}],
         "id":"72741b06-df4d-4715-b142-276b6bce75ab",
         "name":"app_pool",
         "admin_state_up":true,
         "subnet_id":"8032909d-47a1-4715-90af-5153ffe39861",
         "tenant_id":"83657cfcdfe44cd5920adaf26c48ceea",
         "provider": "haproxy"
      }
   ]
}
			`)
	})

	count := 0

	List(fake.ServiceClient(), ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractPools(page)
		if err != nil {
			t.Errorf("Failed to extract pools: %v", err)
			return false, err
		}

		expected := []Pool{
			{
				LBMethod:    "ROUND_ROBIN",
				Protocol:    "HTTP",
				Description: "",
				MonitorIDs: []string{
					"466c8345-28d8-4f84-a246-e04380b0461d",
					"5d4b5228-33b0-4e60-b225-9b727c1a20e7",
				},
				SubnetID:      "8032909d-47a1-4715-90af-5153ffe39861",
				TenantID:      "83657cfcdfe44cd5920adaf26c48ceea",
				AdminStateUp:  true,
				Name:          "app_pool",
				Members:       []map[string]interface{}{{"id": "53306cda-815d-4354-9fe4-59e09da9c3c5"}},
				ID:            "72741b06-df4d-4715-b142-276b6bce75ab",
				Loadbalancers: []map[string]interface{}{{"id": "79e05663-7f03-45d2-a092-8b94062f22ab"}},
				Listeners:     []map[string]interface{}{{"id": "2a280670-c202-4b0b-a562-34077415aabf"}},
				Provider:      "haproxy",
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

	th.Mux.HandleFunc("/v2.0/lbaas/pools", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "pool": {
        "lb_algorithm": "ROUND_ROBIN",
        "protocol": "HTTP",
        "name": "Example pool",
        "subnet_id": "1981f108-3c48-48d2-b908-30f7d28532c9",
        "tenant_id": "2ffc6e22aae24e4795f87155d24c896f",
        "loadbalancer_id": "79e05663-7f03-45d2-a092-8b94062f22ab"
    }
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
    "pool": {
        "lb_algorithm": "ROUND_ROBIN",
        "protocol": "HTTP",
        "description": "",
        "health_monitors": [],
        "members": [{}],
        "id": "69055154-f603-4a28-8951-7cc2d9e54a9a",
        "name": "Example pool",
        "admin_state_up": true,
        "subnet_id": "1981f108-3c48-48d2-b908-30f7d28532c9",
        "tenant_id": "2ffc6e22aae24e4795f87155d24c896f",
        "listeners":[{"id": "2a280670-c202-4b0b-a562-34077415aabf"}],
        "loadbalancers":[{"id": "79e05663-7f03-45d2-a092-8b94062f22ab"}]
    }
}
		`)
	})

	options := CreateOpts{
		LBMethod:       LBMethodRoundRobin,
		Protocol:       "HTTP",
		Name:           "Example pool",
		SubnetID:       "1981f108-3c48-48d2-b908-30f7d28532c9",
		TenantID:       "2ffc6e22aae24e4795f87155d24c896f",
		LoadbalancerID: "79e05663-7f03-45d2-a092-8b94062f22ab",
	}
	p, err := Create(fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "ROUND_ROBIN", p.LBMethod)
	th.AssertEquals(t, "HTTP", p.Protocol)
	th.AssertEquals(t, "", p.Description)
	th.AssertDeepEquals(t, []string{}, p.MonitorIDs)
	th.AssertDeepEquals(t, []map[string]interface{}{{}}, p.Members)
	th.AssertEquals(t, "69055154-f603-4a28-8951-7cc2d9e54a9a", p.ID)
	th.AssertEquals(t, "79e05663-7f03-45d2-a092-8b94062f22ab", p.Loadbalancers[0]["id"])
	th.AssertEquals(t, "Example pool", p.Name)
	th.AssertEquals(t, "1981f108-3c48-48d2-b908-30f7d28532c9", p.SubnetID)
	th.AssertEquals(t, "2ffc6e22aae24e4795f87155d24c896f", p.TenantID)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/lbaas/pools/332abe93-f488-41ba-870b-2ac66be7f853", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
   "pool":{
      "id":"332abe93-f488-41ba-870b-2ac66be7f853",
      "tenant_id":"19eaa775-cf5d-49bc-902e-2f85f668d995",
      "name":"Example pool",
      "description":"",
      "protocol":"tcp",
      "lb_algorithm":"ROUND_ROBIN",
      "session_persistence":{
      },
      "members":[{}],
      "admin_state_up":true
   }
}
			`)
	})

	n, err := Get(fake.ServiceClient(), "332abe93-f488-41ba-870b-2ac66be7f853").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.ID, "332abe93-f488-41ba-870b-2ac66be7f853")
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/lbaas/pools/332abe93-f488-41ba-870b-2ac66be7f853", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
   "pool":{
      "name": "SuperPool",
      "lb_algorithm": "LEAST_CONNECTIONS"
   }
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
   "pool":{
      "lb_algorithm":"LEAST_CONNECTIONS",
      "protocol":"TCP",
      "description":"",
      "health_monitors":[],
      "subnet_id":"8032909d-47a1-4715-90af-5153ffe39861",
      "tenant_id":"83657cfcdfe44cd5920adaf26c48ceea",
      "admin_state_up":true,
      "name":"SuperPool",
      "members":[{}],
      "id":"61b1f87a-7a21-4ad3-9dda-7f81d249944f"
   }
}
		`)
	})

	options := UpdateOpts{Name: "SuperPool", LBMethod: LBMethodLeastConnections}

	n, err := Update(fake.ServiceClient(), "332abe93-f488-41ba-870b-2ac66be7f853", options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "SuperPool", n.Name)
	th.AssertDeepEquals(t, "LEAST_CONNECTIONS", n.LBMethod)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/lbaas/pools/332abe93-f488-41ba-870b-2ac66be7f853", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := Delete(fake.ServiceClient(), "332abe93-f488-41ba-870b-2ac66be7f853")
	th.AssertNoErr(t, res.Err)
}

func TestListAssociateMembers(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/lbaas/pools/332abe93-f488-41ba-870b-2ac66be7f853/members", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
   "members":[
      {
        "id": "2a280670-c202-4b0b-a562-34077415aabf",
        "address": "10.0.2.10",
        "weight": 5,
        "name": "member1",
        "subnet_id": "1981f108-3c48-48d2-b908-30f7d28532c9",
        "tenant_id": "2ffc6e22aae24e4795f87155d24c896f",
        "admin_state_up":true,
        "protocol_port": 80
      },
      {
        "id": "fad389a3-9a4a-4762-a365-8c7038508b5d",
        "address": "10.0.2.11",
        "weight": 10,
        "name": "member2",
        "subnet_id": "1981f108-3c48-48d2-b908-30f7d28532c9",
        "tenant_id": "2ffc6e22aae24e4795f87155d24c896f",
        "admin_state_up":false,
        "protocol_port": 80
      }
   ]
}
			`)
	})

	count := 0

	ListAssociateMembers(fake.ServiceClient(), "332abe93-f488-41ba-870b-2ac66be7f853", MemberListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractMembers(page)
		if err != nil {
			t.Errorf("Failed to extract members: %v", err)
			return false, err
		}

		expected := []Member{
			{
				SubnetID:     "1981f108-3c48-48d2-b908-30f7d28532c9",
				TenantID:     "2ffc6e22aae24e4795f87155d24c896f",
				AdminStateUp: true,
				Name:         "member1",
				ID:           "2a280670-c202-4b0b-a562-34077415aabf",
				Address:      "10.0.2.10",
				Weight:       5,
				ProtocolPort: 80,
			},
			{
				SubnetID:     "1981f108-3c48-48d2-b908-30f7d28532c9",
				TenantID:     "2ffc6e22aae24e4795f87155d24c896f",
				AdminStateUp: false,
				Name:         "member2",
				ID:           "fad389a3-9a4a-4762-a365-8c7038508b5d",
				Address:      "10.0.2.11",
				Weight:       10,
				ProtocolPort: 80,
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestCreateAssociateMember(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/lbaas/pools/332abe93-f488-41ba-870b-2ac66be7f853/members", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "member": {
        "address": "10.0.2.10",
        "weight": 5,
        "name": "Example member",
        "subnet_id": "1981f108-3c48-48d2-b908-30f7d28532c9",
        "tenant_id": "2ffc6e22aae24e4795f87155d24c896f",
        "protocol_port": 80
    }
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
    "member": {
        "id": "2a280670-c202-4b0b-a562-34077415aabf",
        "address": "10.0.2.10",
        "weight": 5,
        "name": "Example member",
        "subnet_id": "1981f108-3c48-48d2-b908-30f7d28532c9",
        "tenant_id": "2ffc6e22aae24e4795f87155d24c896f",
        "admin_state_up":true,
        "protocol_port": 80
    }
}
		`)
	})

	options := MemberCreateOpts{
		Name:         "Example member",
		SubnetID:     "1981f108-3c48-48d2-b908-30f7d28532c9",
		TenantID:     "2ffc6e22aae24e4795f87155d24c896f",
		Address:      "10.0.2.10",
		ProtocolPort: 80,
		Weight:       5,
	}
	p, err := CreateAssociateMember(fake.ServiceClient(), "332abe93-f488-41ba-870b-2ac66be7f853", options).ExtractMember()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "2a280670-c202-4b0b-a562-34077415aabf", p.ID)
	th.AssertEquals(t, "Example member", p.Name)
	th.AssertEquals(t, "1981f108-3c48-48d2-b908-30f7d28532c9", p.SubnetID)
	th.AssertEquals(t, "2ffc6e22aae24e4795f87155d24c896f", p.TenantID)
}

func TestGetAssociateMember(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/lbaas/pools/332abe93-f488-41ba-870b-2ac66be7f853/members/2a280670-c202-4b0b-a562-34077415aabf", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
   "member": {
        "id": "2a280670-c202-4b0b-a562-34077415aabf",
        "address": "10.0.2.10",
        "weight": 5,
        "name": "Example member",
        "subnet_id": "1981f108-3c48-48d2-b908-30f7d28532c9",
        "tenant_id": "2ffc6e22aae24e4795f87155d24c896f",
        "admin_state_up":true,
        "protocol_port": 80
    }
}
			`)
	})

	n, err := GetAssociateMember(fake.ServiceClient(), "332abe93-f488-41ba-870b-2ac66be7f853", "2a280670-c202-4b0b-a562-34077415aabf").ExtractMember()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.ID, "2a280670-c202-4b0b-a562-34077415aabf")
}

func TestUpdateAssociateMember(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/lbaas/pools/332abe93-f488-41ba-870b-2ac66be7f853/members/2a280670-c202-4b0b-a562-34077415aabf", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
   "member":{
      "name": "newMemberName",
      "weight": 4
   }
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
   "member": {
        "id": "2a280670-c202-4b0b-a562-34077415aabf",
        "address": "10.0.2.10",
        "weight": 4,
        "name": "newMemberName",
        "subnet_id": "1981f108-3c48-48d2-b908-30f7d28532c9",
        "tenant_id": "2ffc6e22aae24e4795f87155d24c896f",
        "admin_state_up":true,
        "protocol_port": 80
    }
}
		`)
	})

	options := MemberUpdateOpts{Name: "newMemberName", Weight: 4}

	n, err := UpdateAssociateMember(fake.ServiceClient(), "332abe93-f488-41ba-870b-2ac66be7f853", "2a280670-c202-4b0b-a562-34077415aabf", options).ExtractMember()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "newMemberName", n.Name)
	th.AssertDeepEquals(t, 4, n.Weight)
}

func TestDeleteMember(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/lbaas/pools/332abe93-f488-41ba-870b-2ac66be7f853/members/2a280670-c202-4b0b-a562-34077415aabf", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := DeleteMember(fake.ServiceClient(), "332abe93-f488-41ba-870b-2ac66be7f853", "2a280670-c202-4b0b-a562-34077415aabf")
	th.AssertNoErr(t, res.Err)
}
