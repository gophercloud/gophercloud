package clouds

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"go.yaml.in/yaml/v3"
)

func TestDefaultMicroversionSerialization(t *testing.T) {
	const input = `auth:
  auth_url: https://example.org/v3
compute_default_microversion: "2.87"
volumev3_default_microversion: "3.60"
shared_file_system_default_microversion: 2.10
future_service_default_microversion: "1.4"
unrelated_setting:
  nested: true
regions:
  - name: mars
    values:
      compute_default_microversion: "2.79"
`
	var cloud Cloud
	if err := yaml.Unmarshal([]byte(input), &cloud); err != nil {
		t.Fatal(err)
	}
	expected := map[string]string{"compute": "2.87", "volumev3": "3.60", "shared-file-system": "2.10", "future-service": "1.4"}
	if !reflect.DeepEqual(cloud.DefaultMicroversions, expected) {
		t.Fatalf("defaults: %#v", cloud.DefaultMicroversions)
	}
	for _, codec := range []struct {
		name      string
		marshal   func(any) ([]byte, error)
		unmarshal func([]byte, any) error
	}{{"yaml", yaml.Marshal, yaml.Unmarshal}, {"json", json.Marshal, json.Unmarshal}} {
		t.Run(codec.name, func(t *testing.T) {
			data, err := codec.marshal(cloud)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(data), "volumev3_default_microversion") {
				t.Fatalf("missing flat key: %s", data)
			}
			var got Cloud
			if err := codec.unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(cloud, got) {
				t.Fatalf("roundtrip: %#v", got)
			}
		})
	}
	for _, input := range []string{"compute_default_microversion: [2.87]", "compute_default_microversion: {bad: value}"} {
		if err := yaml.Unmarshal([]byte(input), &Cloud{}); err == nil {
			t.Fatalf("accepted invalid value: %s", input)
		}
	}
}

func TestDefaultMicroversionMerging(t *testing.T) {
	base := Cloud{DefaultMicroversions: map[string]string{"compute": "2.87", "volumev3": "3.60"}}
	override := Cloud{DefaultMicroversions: map[string]string{"compute": "2.79", "placement": "1.20"}}
	got, err := mergeClouds(override, base)
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]string{"compute": "2.79", "volumev3": "3.60", "placement": "1.20"}
	if !reflect.DeepEqual(got.DefaultMicroversions, expected) {
		t.Fatalf("merged defaults: %#v", got.DefaultMicroversions)
	}
}

func TestDefaultMicroversionAliases(t *testing.T) {
	for _, alias := range []string{"block-storage", "block-store", "volume", "volumev2", "volumev3"} {
		t.Run(alias, func(t *testing.T) {
			defaults := map[string]string{alias: "3.60"}
			for _, requested := range []string{"block-storage", "volume", "volumev3"} {
				if got := (gophercloud.EndpointOpts{DefaultMicroversions: defaults}).MicroversionFor(requested); got != "3.60" {
					t.Fatalf("%s: %q", requested, got)
				}
			}
		})
	}
	defaults := map[string]string{"block-storage": "3.60", "block-store": "3.50", "volume": "3.40", "volumev3": "3.30"}
	requested := "volumev3"
	if got := (gophercloud.EndpointOpts{DefaultMicroversions: defaults}).MicroversionFor(requested); got != "3.60" {
		t.Fatalf("canonical precedence: %q", got)
	}
	defaults["block-storage"] = ""
	if got := (gophercloud.EndpointOpts{DefaultMicroversions: defaults}).MicroversionFor(requested); got != "" {
		t.Fatalf("empty canonical default: %q", got)
	}
	delete(defaults, "block-storage")
	if got := (gophercloud.EndpointOpts{DefaultMicroversions: defaults}).MicroversionFor(requested); got != "3.50" {
		t.Fatalf("alias precedence: %q", got)
	}
}

func TestParsedMicroversionsReachRequests(t *testing.T) {
	const config = `clouds:
  test:
    auth:
      auth_url: https://example.org/v3
    compute_default_microversion: "2.87"
    volumev3_default_microversion: "3.60"
    share_default_microversion: "2.65"
    placement_default_microversion: "1.20"
    bare_metal_default_microversion: "1.80"
    container_default_microversion: "1.40"
    container_infra_default_microversion: "1.4"
    regions:
      - name: mars
        values:
          compute_default_microversion: "2.79"
`
	for _, tc := range []struct {
		name, version, header, value string
		factory                      func(context.Context, *gophercloud.ProviderClient, gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error)
	}{
		{"compute", "2.79", "OpenStack-API-Version", "compute 2.79", openstack.NewComputeV2},
		{"block-storage", "3.60", "X-OpenStack-Volume-API-Version", "3.60", openstack.NewBlockStorageV3},
		{"shared-file-system", "2.65", "X-OpenStack-Manila-API-Version", "2.65", openstack.NewSharedFileSystemV2},
		{"placement", "1.20", "OpenStack-API-Version", "placement 1.20", openstack.NewPlacementV1},
		{"baremetal", "1.80", "X-OpenStack-Ironic-API-Version", "1.80", openstack.NewBareMetalV1},
		{"application-container", "1.40", "OpenStack-API-Version", "container 1.40", openstack.NewContainerV1},
		{"container-infrastructure-management", "1.4", "OpenStack-API-Version", "container-infra 1.4", openstack.NewContainerInfraV1},
		{"image", "", "OpenStack-API-Version", "", openstack.NewImageV2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if got := r.Header.Get(tc.header); got != tc.value {
					t.Errorf("header: %q, want %q", got, tc.value)
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()
			_, eo, _, err := Parse(WithCloudName("test"), WithCloudsYAML(strings.NewReader(config)), WithRegion("mars"))
			if err != nil {
				t.Fatal(err)
			}
			provider := &gophercloud.ProviderClient{EndpointLocator: func(context.Context, gophercloud.EndpointOpts) (string, error) { return server.URL + "/", nil }}
			client, err := tc.factory(context.Background(), provider, eo)
			if err != nil {
				t.Fatal(err)
			}
			if client.Microversion != tc.version {
				t.Fatalf("microversion: %q", client.Microversion)
			}
			if _, err := client.Get(context.Background(), server.URL, nil, nil); err != nil {
				t.Fatal(err)
			}
			eo.Microversion = "9.9"
			client, err = tc.factory(context.Background(), provider, eo)
			if err != nil {
				t.Fatal(err)
			}
			if client.Microversion != "9.9" {
				t.Fatalf("explicit override: %q", client.Microversion)
			}
		})
	}
}

func TestSecureMicroversionDefaults(t *testing.T) {
	_, eo, _, err := Parse(WithCloudName("test"), WithCloudsYAML(strings.NewReader(`clouds:
  test:
    auth:
      auth_url: https://example.org/v3
    compute_default_microversion: "2.87"
    volumev3_default_microversion: "3.60"
`)), WithSecureYAML(strings.NewReader(`clouds:
  test:
    compute_default_microversion: "2.79"
`)))
	if err != nil {
		t.Fatal(err)
	}
	if eo.MicroversionFor("compute") != "2.79" || eo.MicroversionFor("block-storage") != "3.60" {
		t.Fatalf("defaults: %#v", eo.DefaultMicroversions)
	}
}
