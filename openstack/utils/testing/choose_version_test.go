package testing

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/utils"
	"github.com/gophercloud/gophercloud/v2/testhelper"
)

func setupVersionHandler() {
	testhelper.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			{
				"versions": {
					"values": [
						{
							"status": "stable",
							"id": "v3.0",
							"links": [
								{ "href": "%s/v3.0", "rel": "self" }
							]
						},
						{
							"status": "stable",
							"id": "v2.0",
							"links": [
								{ "href": "%s/v2.0", "rel": "self" }
							]
						}
					]
				}
			}
		`, testhelper.Server.URL, testhelper.Server.URL)
	})
	// Compute v2.1 API
	testhelper.Mux.HandleFunc("/compute/v2.1/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			{
				"version": {
					"id": "v2.1",
					"status": "CURRENT",
					"version": "2.90",
					"min_version": "2.1",
					"updated": "2013-07-23T11:33:21Z",
					"links": [
						{
							"rel": "self",
							"href": "%s/compute/v2.1/"
						},
						{
							"rel": "describedby",
							"type": "text/html",
							"href": "http://docs.openstack.org/"
						}
					],
					"media-types": [
						{
							"base": "application/json",
							"type": "application/vnd.openstack.compute+json;version=2.1"
						}
					]
				}
			}
		`, testhelper.Server.URL)
	})
	// Compute v2 API
	testhelper.Mux.HandleFunc("/compute/v2/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			{
				"version": {
					"id": "v2.0",
					"status": "SUPPORTED",
					"version": "",
					"min_version": "",
					"updated": "2011-01-21T11:33:21Z",
					"links": [
						{
							"rel": "self",
							"href": "%s/compute/v2/"
						},
						{
							"rel": "describedby",
							"type": "text/html",
							"href": "http://docs.openstack.org/"
						}
					],
					"media-types": [
						{
							"base": "application/json",
							"type": "application/vnd.openstack.compute+json;version=2"
						}
					]
				}
			}
		`, testhelper.Server.URL)
	})
	// Ironic API
	testhelper.Mux.HandleFunc("/ironic/v1/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
		{
			"name": "OpenStack Ironic API",
			"description": "Ironic is an OpenStack project which enables the provision and management of baremetal machines.",
			"default_version": {
				"id": "v1",
				"links": [
					{
						"href": "%s/ironic/v1/",
						"rel": "self"
					}
				],
				"status": "CURRENT",
				"min_version": "1.1",
				"version": "1.87"
			},
			"versions": [
				{
					"id": "v1",
					"links": [
						{
							"href": "%s/ironic/v1/",
							"rel": "self"
						}
					],
					"status": "CURRENT",
					"min_version": "1.1",
					"version": "1.87"
				}
			]
		}
		`, testhelper.Server.URL, testhelper.Server.URL)
	})
	// Ironic multi-version
	testhelper.Mux.HandleFunc("/ironic/v1.2/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
		{
			"name": "OpenStack Ironic API",
			"description": "Ironic is an OpenStack project which enables the provision and management of baremetal machines.",
			"default_version": {
				"id": "v1",
				"links": [
					{
						"href": "%s/ironic/v1/",
						"rel": "self"
					}
				],
				"status": "CURRENT",
				"min_version": "1.1",
				"version": "1.87"
			},
			"versions": [
				{
					"id": "v1",
					"links": [
						{
							"href": "%s/ironic/v1/",
							"rel": "self"
						}
					],
					"status": "CURRENT",
					"min_version": "1.1",
					"version": "1.87"
				},
				{
					"id": "v1.2",
					"links": [
						{
							"href": "%s/ironic/v1/",
							"rel": "self"
						}
					],
					"status": "CURRENT",
					"min_version": "1.2",
					"version": "1.90"
				}
			]
		}
		`, testhelper.Server.URL, testhelper.Server.URL, testhelper.Server.URL)
	})
}

func TestChooseVersion(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	setupVersionHandler()

	v2 := &utils.Version{ID: "v2.0", Priority: 2, Suffix: "blarg"}
	v3 := &utils.Version{ID: "v3.0", Priority: 3, Suffix: "hargl"}

	c := &gophercloud.ProviderClient{
		IdentityBase:     testhelper.Endpoint(),
		IdentityEndpoint: "",
	}
	v, endpoint, err := utils.ChooseVersion(context.TODO(), c, []*utils.Version{v2, v3})

	if err != nil {
		t.Fatalf("Unexpected error from ChooseVersion: %v", err)
	}

	if v != v3 {
		t.Errorf("Expected %#v to win, but %#v did instead", v3, v)
	}

	expected := testhelper.Endpoint() + "v3.0/"
	if endpoint != expected {
		t.Errorf("Expected endpoint [%s], but was [%s] instead", expected, endpoint)
	}
}

func TestChooseVersionOpinionatedLink(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	setupVersionHandler()

	v2 := &utils.Version{ID: "v2.0", Priority: 2, Suffix: "nope"}
	v3 := &utils.Version{ID: "v3.0", Priority: 3, Suffix: "northis"}

	c := &gophercloud.ProviderClient{
		IdentityBase:     testhelper.Endpoint(),
		IdentityEndpoint: testhelper.Endpoint() + "v2.0/",
	}
	v, endpoint, err := utils.ChooseVersion(context.TODO(), c, []*utils.Version{v2, v3})
	if err != nil {
		t.Fatalf("Unexpected error from ChooseVersion: %v", err)
	}

	if v != v2 {
		t.Errorf("Expected %#v to win, but %#v did instead", v2, v)
	}

	expected := testhelper.Endpoint() + "v2.0/"
	if endpoint != expected {
		t.Errorf("Expected endpoint [%s], but was [%s] instead", expected, endpoint)
	}
}

func TestChooseVersionFromSuffix(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	v2 := &utils.Version{ID: "v2.0", Priority: 2, Suffix: "/v2.0/"}
	v3 := &utils.Version{ID: "v3.0", Priority: 3, Suffix: "/v3.0/"}

	c := &gophercloud.ProviderClient{
		IdentityBase:     testhelper.Endpoint(),
		IdentityEndpoint: testhelper.Endpoint() + "v2.0/",
	}
	v, endpoint, err := utils.ChooseVersion(context.TODO(), c, []*utils.Version{v2, v3})
	if err != nil {
		t.Fatalf("Unexpected error from ChooseVersion: %v", err)
	}

	if v != v2 {
		t.Errorf("Expected %#v to win, but %#v did instead", v2, v)
	}

	expected := testhelper.Endpoint() + "v2.0/"
	if endpoint != expected {
		t.Errorf("Expected endpoint [%s], but was [%s] instead", expected, endpoint)
	}
}

type getSupportedServiceMicroversions struct {
	Endpoint    string
	ExpectedMax string
	ExpectedMin string
	ExpectedErr bool
}

func TestGetSupportedVersions(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	setupVersionHandler()

	tests := []getSupportedServiceMicroversions{
		{
			// v2 does not support microversions and returns error
			Endpoint:    testhelper.Endpoint() + "compute/v2/",
			ExpectedMax: "",
			ExpectedMin: "",
			ExpectedErr: true,
		},
		{
			Endpoint:    testhelper.Endpoint() + "compute/v2.1/",
			ExpectedMax: "2.90",
			ExpectedMin: "2.1",
			ExpectedErr: false,
		},
		{
			Endpoint:    testhelper.Endpoint() + "ironic/v1/",
			ExpectedMax: "1.87",
			ExpectedMin: "1.1",
			ExpectedErr: false,
		},
		{
			// This endpoint returns multiple versions, which is not supported
			Endpoint:    testhelper.Endpoint() + "ironic/v1.2/",
			ExpectedMax: "not-relevant",
			ExpectedMin: "not-relevant",
			ExpectedErr: true,
		},
	}

	for _, test := range tests {
		c := &gophercloud.ProviderClient{
			IdentityBase:     testhelper.Endpoint(),
			IdentityEndpoint: testhelper.Endpoint() + "v2.0/",
		}

		client := &gophercloud.ServiceClient{
			ProviderClient: c,
			Endpoint:       test.Endpoint,
		}

		supported, err := utils.GetSupportedMicroversions(context.TODO(), client)

		if test.ExpectedErr {
			if err == nil {
				t.Error("Expected error but got none!")
			}
			// Check for reasonable error message
			if !strings.Contains(err.Error(), "not supported") {
				t.Error("Expected error to contain 'not supported' but it did not!")
			}
			// No point parsing and comparing versions after error, so continue to next test case
			continue
		} else {
			if err != nil {
				t.Errorf("Expected no error but got %s", err.Error())
			}
		}

		min := fmt.Sprintf("%d.%d", supported.MinMajor, supported.MinMinor)
		max := fmt.Sprintf("%d.%d", supported.MaxMajor, supported.MaxMinor)

		if (min != test.ExpectedMin) || (max != test.ExpectedMax) {
			t.Errorf("Expected min=%s and max=%s but got min=%s and max=%s", test.ExpectedMin, test.ExpectedMax, min, max)
		}
	}
}

type microversionSupported struct {
	Version    string
	MinVersion string
	MaxVersion string
	Supported  bool
	Error      bool
}

func TestMicroversionSupported(t *testing.T) {
	tests := []microversionSupported{
		{
			// Checking min version
			Version:    "2.1",
			MinVersion: "2.1",
			MaxVersion: "2.90",
			Supported:  true,
			Error:      false,
		},
		{
			// Checking max version
			Version:    "2.90",
			MinVersion: "2.1",
			MaxVersion: "2.90",
			Supported:  true,
			Error:      false,
		},
		{
			// Checking too high version
			Version:    "2.95",
			MinVersion: "2.1",
			MaxVersion: "2.90",
			Supported:  false,
			Error:      false,
		},
		{
			// Checking too low version
			Version:    "2.1",
			MinVersion: "2.53",
			MaxVersion: "2.90",
			Supported:  false,
			Error:      false,
		},
		{
			// Invalid version
			Version:    "2.1.53",
			MinVersion: "2.53",
			MaxVersion: "2.90",
			Supported:  false,
			Error:      true,
		},
	}

	for _, test := range tests {
		var err error
		var supportedVersions utils.SupportedMicroversions
		supportedVersions.MaxMajor, supportedVersions.MaxMinor, err = utils.ParseMicroversion(test.MaxVersion)
		if err != nil {
			t.Error("Error parsing MaxVersion!")
		}
		supportedVersions.MinMajor, supportedVersions.MinMinor, err = utils.ParseMicroversion(test.MinVersion)
		if err != nil {
			t.Error("Error parsing MinVersion!")
		}

		supported, err := supportedVersions.IsSupported(test.Version)
		if test.Error {
			if err == nil {
				t.Error("Expected error but got none!")
			}
		} else {
			if err != nil {
				t.Errorf("Expected no error but got %s", err.Error())
			}
		}
		if test.Supported != supported {
			t.Errorf("Expected supported=%t to be %t, when version=%s, min=%s and max=%s",
				supported, test.Supported, test.Version, test.MinVersion, test.MaxVersion)
		}
	}
}
