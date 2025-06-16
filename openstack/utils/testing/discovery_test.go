package testing

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/utils"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

type getSupportedServiceMicroversions struct {
	Endpoint    string
	ExpectedMax string
	ExpectedMin string
	ExpectedErr bool
}

func TestGetSupportedVersions(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	setupVersionHandler(fakeServer)

	tests := []getSupportedServiceMicroversions{
		{
			// v2 does not support microversions and returns error
			Endpoint:    fakeServer.Endpoint() + "compute/v2/",
			ExpectedMax: "",
			ExpectedMin: "",
			ExpectedErr: true,
		},
		{
			Endpoint:    fakeServer.Endpoint() + "compute/v2.1/",
			ExpectedMax: "2.90",
			ExpectedMin: "2.1",
			ExpectedErr: false,
		},
		{
			Endpoint:    fakeServer.Endpoint() + "ironic/v1/",
			ExpectedMax: "1.87",
			ExpectedMin: "1.1",
			ExpectedErr: false,
		},
		{
			// This endpoint returns multiple versions, which is not supported
			Endpoint:    fakeServer.Endpoint() + "ironic/v1.2/",
			ExpectedMax: "not-relevant",
			ExpectedMin: "not-relevant",
			ExpectedErr: true,
		},
	}

	for _, test := range tests {
		c := &gophercloud.ProviderClient{
			IdentityBase:     fakeServer.Endpoint(),
			IdentityEndpoint: fakeServer.Endpoint() + "v2.0/",
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
