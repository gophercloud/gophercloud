package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/rules"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/security-group-rules", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `
{
    "security_group_rules": [
        {
            "direction": "egress",
            "ethertype": "IPv6",
            "id": "3c0e45ff-adaf-4124-b083-bf390e5482ff",
            "port_range_max": null,
            "port_range_min": null,
            "protocol": null,
            "remote_group_id": null,
            "remote_ip_prefix": null,
            "security_group_id": "85cc3048-abc3-43cc-89b3-377341426ac5",
            "tenant_id": "e4f50856753b4dc6afee5fa6b9b6c550"
        },
        {
            "direction": "egress",
            "ethertype": "IPv4",
            "id": "93aa42e5-80db-4581-9391-3a608bd0e448",
            "port_range_max": null,
            "port_range_min": null,
            "protocol": null,
            "remote_group_id": null,
            "remote_ip_prefix": null,
            "security_group_id": "85cc3048-abc3-43cc-89b3-377341426ac5",
            "tenant_id": "e4f50856753b4dc6afee5fa6b9b6c550"
        }
    ]
}
      `)
	})

	count := 0

	err := rules.List(fake.ServiceClient(fakeServer), rules.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := rules.ExtractRules(page)
		if err != nil {
			t.Errorf("Failed to extract secrules: %v", err)
			return false, err
		}

		expected := []rules.SecGroupRule{
			{
				Description:    "",
				Direction:      "egress",
				EtherType:      "IPv6",
				ID:             "3c0e45ff-adaf-4124-b083-bf390e5482ff",
				PortRangeMax:   0,
				PortRangeMin:   0,
				Protocol:       "",
				RemoteGroupID:  "",
				RemoteIPPrefix: "",
				SecGroupID:     "85cc3048-abc3-43cc-89b3-377341426ac5",
				TenantID:       "e4f50856753b4dc6afee5fa6b9b6c550",
			},
			{
				Direction:      "egress",
				EtherType:      "IPv4",
				ID:             "93aa42e5-80db-4581-9391-3a608bd0e448",
				PortRangeMax:   0,
				PortRangeMin:   0,
				Protocol:       "",
				RemoteGroupID:  "",
				RemoteIPPrefix: "",
				SecGroupID:     "85cc3048-abc3-43cc-89b3-377341426ac5",
				TenantID:       "e4f50856753b4dc6afee5fa6b9b6c550",
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/security-group-rules", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "security_group_rule": {
        "description": "test description of rule",
        "direction": "ingress",
        "port_range_min": 80,
        "ethertype": "IPv4",
        "port_range_max": 80,
        "protocol": "tcp",
        "remote_group_id": "85cc3048-abc3-43cc-89b3-377341426ac5",
        "security_group_id": "a7734e61-b545-452d-a3cd-0189cbd9747a"
    }
}
      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, `
{
    "security_group_rule": {
        "description": "test description of rule",
        "direction": "ingress",
        "ethertype": "IPv4",
        "id": "2bc0accf-312e-429a-956e-e4407625eb62",
        "port_range_max": 80,
        "port_range_min": 80,
        "protocol": "tcp",
        "remote_group_id": "85cc3048-abc3-43cc-89b3-377341426ac5",
        "remote_ip_prefix": null,
        "security_group_id": "a7734e61-b545-452d-a3cd-0189cbd9747a",
        "tenant_id": "e4f50856753b4dc6afee5fa6b9b6c550"
    }
}
    `)
	})

	opts := rules.CreateOpts{
		Description:   "test description of rule",
		Direction:     "ingress",
		PortRangeMin:  80,
		EtherType:     rules.EtherType4,
		PortRangeMax:  80,
		Protocol:      "tcp",
		RemoteGroupID: "85cc3048-abc3-43cc-89b3-377341426ac5",
		SecGroupID:    "a7734e61-b545-452d-a3cd-0189cbd9747a",
	}
	_, err := rules.Create(context.TODO(), fake.ServiceClient(fakeServer), opts).Extract()
	th.AssertNoErr(t, err)
}

func TestCreateAnyProtocol(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/security-group-rules", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "security_group_rule": {
        "description": "test description of rule",
        "direction": "ingress",
        "port_range_min": 80,
        "ethertype": "IPv4",
        "port_range_max": 80,
        "remote_group_id": "85cc3048-abc3-43cc-89b3-377341426ac5",
        "security_group_id": "a7734e61-b545-452d-a3cd-0189cbd9747a"
    }
}
      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, `
{
    "security_group_rule": {
        "description": "test description of rule",
        "direction": "ingress",
        "ethertype": "IPv4",
        "id": "2bc0accf-312e-429a-956e-e4407625eb62",
        "port_range_max": 80,
        "port_range_min": 80,
        "remote_group_id": "85cc3048-abc3-43cc-89b3-377341426ac5",
        "remote_ip_prefix": null,
        "security_group_id": "a7734e61-b545-452d-a3cd-0189cbd9747a",
        "tenant_id": "e4f50856753b4dc6afee5fa6b9b6c550"
    }
}
    `)
	})

	opts := rules.CreateOpts{
		Description:   "test description of rule",
		Direction:     "ingress",
		PortRangeMin:  80,
		EtherType:     rules.EtherType4,
		PortRangeMax:  80,
		Protocol:      rules.ProtocolAny,
		RemoteGroupID: "85cc3048-abc3-43cc-89b3-377341426ac5",
		SecGroupID:    "a7734e61-b545-452d-a3cd-0189cbd9747a",
	}
	_, err := rules.Create(context.TODO(), fake.ServiceClient(fakeServer), opts).Extract()
	th.AssertNoErr(t, err)
}

func TestCreateBulk(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/security-group-rules", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "security_group_rules": [
        {
            "description": "test description of rule",
            "direction": "ingress",
            "port_range_min": 80,
            "ethertype": "IPv4",
            "port_range_max": 80,
            "protocol": "tcp",
            "remote_group_id": "85cc3048-abc3-43cc-89b3-377341426ac5",
            "security_group_id": "a7734e61-b545-452d-a3cd-0189cbd9747a"
        },
        {
            "description": "test description of rule",
            "direction": "ingress",
            "port_range_min": 443,
            "ethertype": "IPv4",
            "port_range_max": 443,
            "protocol": "tcp",
            "security_group_id": "a7734e61-b545-452d-a3cd-0189cbd9747a"
        }
    ]
}
      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, `
{
    "security_group_rules": [
        {
            "description": "test description of rule",
            "direction": "ingress",
            "ethertype": "IPv4",
            "port_range_max": 80,
            "port_range_min": 80,
            "protocol": "tcp",
            "remote_group_id": "85cc3048-abc3-43cc-89b3-377341426ac5",
            "remote_ip_prefix": null,
            "security_group_id": "a7734e61-b545-452d-a3cd-0189cbd9747a",
            "tenant_id": "e4f50856753b4dc6afee5fa6b9b6c550"
        },
        {
            "description": "test description of rule",
            "direction": "ingress",
            "ethertype": "IPv4",
            "port_range_max": 443,
            "port_range_min": 443,
            "protocol": "tcp",
            "remote_group_id": null,
            "remote_ip_prefix": null,
            "security_group_id": "a7734e61-b545-452d-a3cd-0189cbd9747a",
            "tenant_id": "e4f50856753b4dc6afee5fa6b9b6c550"
        }
    ]
}
    `)
	})

	opts := []rules.CreateOpts{
		{
			Description:   "test description of rule",
			Direction:     "ingress",
			PortRangeMin:  80,
			EtherType:     rules.EtherType4,
			PortRangeMax:  80,
			Protocol:      "tcp",
			RemoteGroupID: "85cc3048-abc3-43cc-89b3-377341426ac5",
			SecGroupID:    "a7734e61-b545-452d-a3cd-0189cbd9747a",
		},
		{
			Description:  "test description of rule",
			Direction:    "ingress",
			PortRangeMin: 443,
			EtherType:    rules.EtherType4,
			PortRangeMax: 443,
			Protocol:     "tcp",
			SecGroupID:   "a7734e61-b545-452d-a3cd-0189cbd9747a",
		},
	}
	{
		_, err := rules.CreateBulk(context.TODO(), fake.ServiceClient(fakeServer), opts).Extract()
		th.AssertNoErr(t, err)
	}

	{
		optsBuilder := make([]rules.CreateOptsBuilder, len(opts))
		for i := range opts {
			optsBuilder[i] = opts[i]
		}
		_, err := rules.CreateBulk(context.TODO(), fake.ServiceClient(fakeServer), optsBuilder).Extract()
		th.AssertNoErr(t, err)
	}
}

func TestRequiredCreateOpts(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	res := rules.Create(context.TODO(), fake.ServiceClient(fakeServer), rules.CreateOpts{Direction: rules.DirIngress})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
	res = rules.Create(context.TODO(), fake.ServiceClient(fakeServer), rules.CreateOpts{Direction: rules.DirIngress, EtherType: rules.EtherType4})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
	res = rules.Create(context.TODO(), fake.ServiceClient(fakeServer), rules.CreateOpts{Direction: rules.DirIngress, EtherType: rules.EtherType4})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
	res = rules.Create(context.TODO(), fake.ServiceClient(fakeServer), rules.CreateOpts{Direction: rules.DirIngress, EtherType: rules.EtherType4, SecGroupID: "something", Protocol: "foo"})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/security-group-rules/3c0e45ff-adaf-4124-b083-bf390e5482ff", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `
{
    "security_group_rule": {
        "direction": "egress",
        "ethertype": "IPv6",
        "id": "3c0e45ff-adaf-4124-b083-bf390e5482ff",
        "port_range_max": null,
        "port_range_min": null,
        "protocol": null,
        "remote_group_id": null,
        "remote_ip_prefix": null,
        "security_group_id": "85cc3048-abc3-43cc-89b3-377341426ac5",
        "tenant_id": "e4f50856753b4dc6afee5fa6b9b6c550"
    }
}
      `)
	})

	sr, err := rules.Get(context.TODO(), fake.ServiceClient(fakeServer), "3c0e45ff-adaf-4124-b083-bf390e5482ff").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "egress", sr.Direction)
	th.AssertEquals(t, "IPv6", sr.EtherType)
	th.AssertEquals(t, "3c0e45ff-adaf-4124-b083-bf390e5482ff", sr.ID)
	th.AssertEquals(t, 0, sr.PortRangeMax)
	th.AssertEquals(t, 0, sr.PortRangeMin)
	th.AssertEquals(t, "", sr.Protocol)
	th.AssertEquals(t, "", sr.RemoteGroupID)
	th.AssertEquals(t, "", sr.RemoteIPPrefix)
	th.AssertEquals(t, "85cc3048-abc3-43cc-89b3-377341426ac5", sr.SecGroupID)
	th.AssertEquals(t, "e4f50856753b4dc6afee5fa6b9b6c550", sr.TenantID)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/security-group-rules/4ec89087-d057-4e2c-911f-60a3b47ee304", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := rules.Delete(context.TODO(), fake.ServiceClient(fakeServer), "4ec89087-d057-4e2c-911f-60a3b47ee304")
	th.AssertNoErr(t, res.Err)
}
