package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/internal/ptr"
	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/defaultrules"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/rules"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/default-security-group-rules", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `
{
    "default_security_group_rules": [
        {
            "description": "",
            "direction": "egress",
            "ethertype": "IPv6",
            "id": "3c0e45ff-adaf-4124-b083-bf390e5482ff",
            "port_range_max": null,
            "port_range_min": null,
            "protocol": null,
            "remote_group_id": null,
            "remote_address_group_id": null,
            "remote_ip_prefix": null,
            "used_in_default_sg": true,
            "used_in_non_default_sg": true,
            "created_at": "2017-12-28T07:21:40Z",
            "updated_at": "2017-12-28T07:21:40Z",
            "revision_number": 0
        },
        {
            "description": "Allow HTTPS from the parent security group",
            "direction": "ingress",
            "ethertype": "IPv4",
            "id": "93aa42e5-80db-4581-9391-3a608bd0e448",
            "port_range_max": 443,
            "port_range_min": 443,
            "protocol": "tcp",
            "remote_group_id": "PARENT",
            "remote_address_group_id": null,
            "remote_ip_prefix": null,
            "used_in_default_sg": false,
            "used_in_non_default_sg": true,
            "created_at": "2017-12-28T07:21:40",
            "updated_at": "2017-12-28T07:21:40",
            "revision_number": 1
        }
    ]
}
      `)
	})

	count := 0

	err := defaultrules.List(fake.ServiceClient(fakeServer), defaultrules.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := defaultrules.ExtractDefaultRules(page)
		if err != nil {
			t.Errorf("Failed to extract default security group rules: %v", err)
			return false, err
		}

		expected := []defaultrules.DefaultSecGroupRule{
			{
				Description:        "",
				Direction:          "egress",
				EtherType:          "IPv6",
				ID:                 "3c0e45ff-adaf-4124-b083-bf390e5482ff",
				PortRangeMax:       0,
				PortRangeMin:       0,
				Protocol:           "",
				RemoteGroupID:      "",
				RemoteIPPrefix:     "",
				UsedInDefaultSG:    true,
				UsedInNonDefaultSG: true,
				CreatedAt:          time.Date(2017, 12, 28, 07, 21, 40, 0, time.UTC),
				UpdatedAt:          time.Date(2017, 12, 28, 07, 21, 40, 0, time.UTC),
			},
			{
				Description:        "Allow HTTPS from the parent security group",
				Direction:          "ingress",
				EtherType:          "IPv4",
				ID:                 "93aa42e5-80db-4581-9391-3a608bd0e448",
				PortRangeMax:       443,
				PortRangeMin:       443,
				Protocol:           "tcp",
				RemoteGroupID:      "PARENT",
				RemoteIPPrefix:     "",
				UsedInDefaultSG:    false,
				UsedInNonDefaultSG: true,
				RevisionNumber:     1,
				CreatedAt:          time.Date(2017, 12, 28, 07, 21, 40, 0, time.UTC),
				UpdatedAt:          time.Date(2017, 12, 28, 07, 21, 40, 0, time.UTC),
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

	fakeServer.Mux.HandleFunc("/v2.0/default-security-group-rules", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "default_security_group_rule": {
        "description": "Allow HTTPS from the parent security group",
        "direction": "ingress",
        "ethertype": "IPv4",
        "port_range_max": 443,
        "port_range_min": 443,
        "protocol": "tcp",
        "remote_group_id": "PARENT",
        "used_in_default_sg": false,
        "used_in_non_default_sg": true
    }
}
      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, `
{
    "default_security_group_rule": {
        "description": "Allow HTTPS from the parent security group",
        "direction": "ingress",
        "ethertype": "IPv4",
        "id": "2bc0accf-312e-429a-956e-e4407625eb62",
        "port_range_max": 443,
        "port_range_min": 443,
        "protocol": "tcp",
        "remote_group_id": "PARENT",
        "remote_address_group_id": null,
        "remote_ip_prefix": null,
        "used_in_default_sg": false,
        "used_in_non_default_sg": true,
        "revision_number": 0
    }
}
    `)
	})

	options := defaultrules.CreateOpts{
		Description:        "Allow HTTPS from the parent security group",
		Direction:          rules.DirIngress,
		EtherType:          rules.EtherType4,
		PortRangeMax:       ptr.To(443),
		PortRangeMin:       ptr.To(443),
		Protocol:           rules.ProtocolTCP,
		RemoteGroupID:      "PARENT",
		UsedInDefaultSG:    ptr.To(false),
		UsedInNonDefaultSG: ptr.To(true),
	}
	rule, err := defaultrules.Create(context.TODO(), fake.ServiceClient(fakeServer), options).Extract()
	th.AssertNoErr(t, err)

	expected := &defaultrules.DefaultSecGroupRule{
		Description:        "Allow HTTPS from the parent security group",
		Direction:          "ingress",
		EtherType:          "IPv4",
		ID:                 "2bc0accf-312e-429a-956e-e4407625eb62",
		PortRangeMax:       443,
		PortRangeMin:       443,
		Protocol:           "tcp",
		RemoteGroupID:      "PARENT",
		UsedInDefaultSG:    false,
		UsedInNonDefaultSG: true,
	}
	th.AssertDeepEquals(t, expected, rule)
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/default-security-group-rules/3c0e45ff-adaf-4124-b083-bf390e5482ff", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `
{
    "default_security_group_rule": {
        "description": "",
        "direction": "egress",
        "ethertype": "IPv6",
        "id": "3c0e45ff-adaf-4124-b083-bf390e5482ff",
        "port_range_max": null,
        "port_range_min": null,
        "protocol": null,
        "remote_group_id": null,
        "remote_address_group_id": null,
        "remote_ip_prefix": null,
        "used_in_default_sg": true,
        "used_in_non_default_sg": true,
        "revision_number": 0
    }
}
      `)
	})

	rule, err := defaultrules.Get(context.TODO(), fake.ServiceClient(fakeServer), "3c0e45ff-adaf-4124-b083-bf390e5482ff").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "egress", rule.Direction)
	th.AssertEquals(t, "IPv6", rule.EtherType)
	th.AssertEquals(t, "3c0e45ff-adaf-4124-b083-bf390e5482ff", rule.ID)
	th.AssertEquals(t, true, rule.UsedInDefaultSG)
	th.AssertEquals(t, true, rule.UsedInNonDefaultSG)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/default-security-group-rules/4ec89087-d057-4e2c-911f-60a3b47ee304", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := defaultrules.Delete(context.TODO(), fake.ServiceClient(fakeServer), "4ec89087-d057-4e2c-911f-60a3b47ee304")
	th.AssertNoErr(t, res.Err)
}
