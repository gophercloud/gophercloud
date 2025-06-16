package testing

import (
	"context"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/utils"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestGetSupportedMicroversions(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	setupMultiServiceVersionHandler(fakeServer)

	tests := []struct {
		name             string
		endpoint         string
		expectedVersions utils.SupportedMicroversions
		expectedErr      string
	}{
		{
			// v2 does not support microversions and returns error
			name:        "compute legacy endpoint",
			endpoint:    fakeServer.Endpoint() + "compute/v2/",
			expectedErr: "not supported",
		},
		{
			name:     "compute versioned endpoint",
			endpoint: fakeServer.Endpoint() + "compute/v2.1/",
			expectedVersions: utils.SupportedMicroversions{
				MaxMajor: 2, MaxMinor: 90, MinMajor: 2, MinMinor: 1,
			},
		},
		{
			name:     "baremetal versioned endpoint",
			endpoint: fakeServer.Endpoint() + "baremetal/v1/",
			expectedVersions: utils.SupportedMicroversions{
				MaxMajor: 1, MaxMinor: 87, MinMajor: 1, MinMinor: 1,
			},
		},
		{
			// This endpoint returns multiple versions, which is not supported
			name:        "fictional multi-version endpoint",
			endpoint:    fakeServer.Endpoint() + "multi-version/v1.2/",
			expectedErr: "not supported",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &gophercloud.ProviderClient{}
			client := &gophercloud.ServiceClient{
				ProviderClient: c,
				Endpoint:       tt.endpoint,
			}

			actualVersions, err := utils.GetSupportedMicroversions(context.TODO(), client)

			if tt.expectedErr != "" {
				th.AssertErr(t, err)
				if !strings.Contains(err.Error(), tt.expectedErr) {
					t.Fatalf("Expected error to contain '%s', got '%s': %+v", tt.expectedErr, err, tt)
				}
			} else {
				th.AssertNoErr(t, err)
			}

			th.AssertDeepEquals(t, tt.expectedVersions, actualVersions)
		})
	}
}

func TestMicroversionSupported(t *testing.T) {
	tests := []struct {
		name          string
		version       string
		minVersion    string
		maxVersion    string
		supported     bool
		expectedError bool
	}{
		{
			name:          "Checking min version",
			version:       "2.1",
			minVersion:    "2.1",
			maxVersion:    "2.90",
			supported:     true,
			expectedError: false,
		},
		{
			name:          "Checking max version",
			version:       "2.90",
			minVersion:    "2.1",
			maxVersion:    "2.90",
			supported:     true,
			expectedError: false,
		},
		{
			name:          "Checking too high version",
			version:       "2.95",
			minVersion:    "2.1",
			maxVersion:    "2.90",
			supported:     false,
			expectedError: false,
		},
		{
			name:          "Checking too low version",
			version:       "2.1",
			minVersion:    "2.53",
			maxVersion:    "2.90",
			supported:     false,
			expectedError: false,
		},
		{
			name:          "Invalid version",
			version:       "2.1.53",
			minVersion:    "2.53",
			maxVersion:    "2.90",
			supported:     false,
			expectedError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			var supportedVersions utils.SupportedMicroversions

			supportedVersions.MaxMajor, supportedVersions.MaxMinor, err = utils.ParseMicroversion(tt.maxVersion)
			if err != nil {
				t.Fatal("Error parsing MaxVersion!")
			}

			supportedVersions.MinMajor, supportedVersions.MinMinor, err = utils.ParseMicroversion(tt.minVersion)
			if err != nil {
				t.Fatal("Error parsing MinVersion!")
			}

			supported, err := supportedVersions.IsSupported(tt.version)
			if tt.expectedError {
				th.AssertErr(t, err)
			} else {
				th.AssertNoErr(t, err)
			}

			if tt.supported != supported {
				t.Fatalf("Expected supported=%t to be %t, when version=%s, min=%s and max=%s",
					supported, tt.supported, tt.version, tt.minVersion, tt.maxVersion)
			}
		})
	}
}
