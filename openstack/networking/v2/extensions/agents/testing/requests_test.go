package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/agents"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/agents", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, AgentsListResult)
	})

	count := 0

	err := agents.List(fake.ServiceClient(), agents.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := agents.ExtractAgents(page)

		if err != nil {
			t.Errorf("Failed to extract agents: %v", err)
			return false, nil
		}

		expected := []agents.Agent{
			Agent1,
			Agent2,
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/agents/43583cf5-472e-4dc8-af5b-6aed4c94ee3a", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, AgentsGetResult)
	})

	s, err := agents.Get(context.TODO(), fake.ServiceClient(), "43583cf5-472e-4dc8-af5b-6aed4c94ee3a").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.ID, "43583cf5-472e-4dc8-af5b-6aed4c94ee3a")
	th.AssertEquals(t, s.Binary, "neutron-openvswitch-agent")
	th.AssertEquals(t, s.AdminStateUp, true)
	th.AssertEquals(t, s.Alive, true)
	th.AssertEquals(t, s.Topic, "N/A")
	th.AssertEquals(t, s.Host, "compute3")
	th.AssertEquals(t, s.AgentType, "Open vSwitch agent")
	th.AssertEquals(t, s.HeartbeatTimestamp, time.Date(2019, 1, 9, 11, 43, 01, 0, time.UTC))
	th.AssertEquals(t, s.StartedAt, time.Date(2018, 6, 26, 21, 46, 20, 0, time.UTC))
	th.AssertEquals(t, s.CreatedAt, time.Date(2017, 7, 26, 23, 2, 5, 0, time.UTC))
	th.AssertDeepEquals(t, s.Configurations, map[string]any{
		"ovs_hybrid_plug":            false,
		"datapath_type":              "system",
		"vhostuser_socket_dir":       "/var/run/openvswitch",
		"log_agent_heartbeats":       false,
		"l2_population":              true,
		"enable_distributed_routing": false,
	})
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/agents/43583cf5-472e-4dc8-af5b-6aed4c94ee3a", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, AgentUpdateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, AgentsUpdateResult)
	})

	iTrue := true
	description := "My OVS agent for OpenStack"
	updateOpts := &agents.UpdateOpts{
		Description:  &description,
		AdminStateUp: &iTrue,
	}
	s, err := agents.Update(context.TODO(), fake.ServiceClient(), "43583cf5-472e-4dc8-af5b-6aed4c94ee3a", updateOpts).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, *s, Agent)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/agents/43583cf5-472e-4dc8-af5b-6aed4c94ee3a", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	err := agents.Delete(context.TODO(), fake.ServiceClient(), "43583cf5-472e-4dc8-af5b-6aed4c94ee3a").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestListDHCPNetworks(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/agents/43583cf5-472e-4dc8-af5b-6aed4c94ee3a/dhcp-networks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, AgentDHCPNetworksListResult)
	})

	s, err := agents.ListDHCPNetworks(context.TODO(), fake.ServiceClient(), "43583cf5-472e-4dc8-af5b-6aed4c94ee3a").Extract()
	th.AssertNoErr(t, err)

	var nilSlice []string
	th.AssertEquals(t, len(s), 1)
	th.AssertEquals(t, s[0].ID, "d32019d3-bc6e-4319-9c1d-6722fc136a22")
	th.AssertEquals(t, s[0].AdminStateUp, true)
	th.AssertEquals(t, s[0].ProjectID, "4fd44f30292945e481c7b8a0c8908869")
	th.AssertEquals(t, s[0].Shared, false)
	th.AssertEquals(t, s[0].Name, "net1")
	th.AssertEquals(t, s[0].Status, "ACTIVE")
	th.AssertDeepEquals(t, s[0].Tags, nilSlice)
	th.AssertEquals(t, s[0].TenantID, "4fd44f30292945e481c7b8a0c8908869")
	th.AssertDeepEquals(t, s[0].AvailabilityZoneHints, []string{})
	th.AssertDeepEquals(t, s[0].Subnets, []string{"54d6f61d-db07-451c-9ab3-b9609b6b6f0b"})

}

func TestScheduleDHCPNetwork(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/agents/43583cf5-472e-4dc8-af5b-6aed4c94ee3a/dhcp-networks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, ScheduleDHCPNetworkRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	})

	opts := &agents.ScheduleDHCPNetworkOpts{
		NetworkID: "1ae075ca-708b-4e66-b4a7-b7698632f05f",
	}
	err := agents.ScheduleDHCPNetwork(context.TODO(), fake.ServiceClient(), "43583cf5-472e-4dc8-af5b-6aed4c94ee3a", opts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestRemoveDHCPNetwork(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/agents/43583cf5-472e-4dc8-af5b-6aed4c94ee3a/dhcp-networks/1ae075ca-708b-4e66-b4a7-b7698632f05f", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	err := agents.RemoveDHCPNetwork(context.TODO(), fake.ServiceClient(), "43583cf5-472e-4dc8-af5b-6aed4c94ee3a", "1ae075ca-708b-4e66-b4a7-b7698632f05f").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestListBGPSpeakers(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	agentID := "30d76012-46de-4215-aaa1-a1630d01d891"

	th.Mux.HandleFunc("/v2.0/agents/"+agentID+"/bgp-drinstances",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprintf(w, ListBGPSpeakersResult)
		})

	count := 0
	err := agents.ListBGPSpeakers(fake.ServiceClient(), agentID).EachPage(
		context.TODO(),
		func(_ context.Context, page pagination.Page) (bool, error) {
			count++
			actual, err := agents.ExtractBGPSpeakers(page)

			th.AssertNoErr(t, err)
			th.AssertEquals(t, len(actual), 1)
			th.AssertEquals(t, actual[0].ID, "cab00464-284d-4251-9798-2b27db7b1668")
			th.AssertEquals(t, actual[0].Name, "gophercloud-testing-speaker")
			th.AssertEquals(t, actual[0].LocalAS, 12345)
			th.AssertEquals(t, actual[0].IPVersion, 4)
			return true, nil
		})
	th.AssertNoErr(t, err)
	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestScheduleBGPSpeaker(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	agentID := "30d76012-46de-4215-aaa1-a1630d01d891"
	speakerID := "8edb2c68-0654-49a9-b3fe-030f92e3ddf6"

	th.Mux.HandleFunc("/v2.0/agents/"+agentID+"/bgp-drinstances",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "POST")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, ScheduleBGPSpeakerRequest)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
		})

	var opts agents.ScheduleBGPSpeakerOpts
	opts.SpeakerID = speakerID
	err := agents.ScheduleBGPSpeaker(context.TODO(), fake.ServiceClient(), agentID, opts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestRemoveBGPSpeaker(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	agentID := "30d76012-46de-4215-aaa1-a1630d01d891"
	speakerID := "8edb2c68-0654-49a9-b3fe-030f92e3ddf6"

	th.Mux.HandleFunc("/v2.0/agents/"+agentID+"/bgp-drinstances/"+speakerID,
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "DELETE")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
			th.TestHeader(t, r, "Accept", "application/json")

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
		})

	err := agents.RemoveBGPSpeaker(context.TODO(), fake.ServiceClient(), agentID, speakerID).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestListDRAgentHostingBGPSpeakers(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	speakerID := "3f511b1b-d541-45f1-aa98-2e44e8183d4c"
	th.Mux.HandleFunc("/v2.0/bgp-speakers/"+speakerID+"/bgp-dragents",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, ListDRAgentHostingBGPSpeakersResult)
		})

	count := 0
	err := agents.ListDRAgentHostingBGPSpeakers(fake.ServiceClient(), speakerID).EachPage(
		context.TODO(),
		func(_ context.Context, page pagination.Page) (bool, error) {
			count++
			actual, err := agents.ExtractAgents(page)

			if err != nil {
				t.Errorf("Failed to extract agents: %v", err)
				return false, nil
			}

			expected := []agents.Agent{BGPAgent1, BGPAgent2}
			th.CheckDeepEquals(t, expected, actual)
			return true, nil
		})
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestListL3Routers(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/agents/43583cf5-472e-4dc8-af5b-6aed4c94ee3a/l3-routers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, AgentL3RoutersListResult)
	})

	s, err := agents.ListL3Routers(context.TODO(), fake.ServiceClient(), "43583cf5-472e-4dc8-af5b-6aed4c94ee3a").Extract()
	th.AssertNoErr(t, err)

	routes := []routers.Route{
		{
			NextHop:         "172.24.3.99",
			DestinationCIDR: "179.24.1.0/24",
		},
	}

	var snat bool = true
	gw := routers.GatewayInfo{
		EnableSNAT: &snat,
		NetworkID:  "ae34051f-aa6c-4c75-abf5-50dc9ac99ef3",
		ExternalFixedIPs: []routers.ExternalFixedIP{
			{
				IPAddress: "172.24.4.3",
				SubnetID:  "b930d7f6-ceb7-40a0-8b81-a425dd994ccf",
			},

			{
				IPAddress: "2001:db8::c",
				SubnetID:  "0c56df5d-ace5-46c8-8f4c-45fa4e334d18",
			},
		},
	}

	var nilSlice []string
	th.AssertEquals(t, len(s), 2)
	th.AssertEquals(t, s[0].ID, "915a14a6-867b-4af7-83d1-70efceb146f9")
	th.AssertEquals(t, s[0].AdminStateUp, true)
	th.AssertEquals(t, s[0].ProjectID, "0bd18306d801447bb457a46252d82d13")
	th.AssertEquals(t, s[0].Name, "router2")
	th.AssertEquals(t, s[0].Status, "ACTIVE")
	th.AssertEquals(t, s[0].TenantID, "0bd18306d801447bb457a46252d82d13")
	th.AssertDeepEquals(t, s[0].AvailabilityZoneHints, []string{})
	th.AssertDeepEquals(t, s[0].Routes, routes)
	th.AssertDeepEquals(t, s[0].GatewayInfo, gw)
	th.AssertDeepEquals(t, s[0].Tags, nilSlice)
	th.AssertEquals(t, s[1].ID, "f8a44de0-fc8e-45df-93c7-f79bf3b01c95")
	th.AssertEquals(t, s[1].Name, "router1")

}

func TestScheduleL3Router(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/agents/43583cf5-472e-4dc8-af5b-6aed4c94ee3a/l3-routers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, ScheduleL3RouterRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	})

	opts := &agents.ScheduleL3RouterOpts{
		RouterID: "43e66290-79a4-415d-9eb9-7ff7919839e1",
	}
	err := agents.ScheduleL3Router(context.TODO(), fake.ServiceClient(), "43583cf5-472e-4dc8-af5b-6aed4c94ee3a", opts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestRemoveL3Router(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/agents/43583cf5-472e-4dc8-af5b-6aed4c94ee3a/l3-routers/43e66290-79a4-415d-9eb9-7ff7919839e1", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	err := agents.RemoveL3Router(context.TODO(), fake.ServiceClient(), "43583cf5-472e-4dc8-af5b-6aed4c94ee3a", "43e66290-79a4-415d-9eb9-7ff7919839e1").ExtractErr()
	th.AssertNoErr(t, err)
}
