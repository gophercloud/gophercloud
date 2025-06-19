package testing

import (
	"context"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/utils"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestGetServiceVersions(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	setupMultiServiceVersionHandler(fakeServer)

	tests := []struct {
		name             string
		endpoint         string
		discoverVersions bool
		expectedVersions []utils.SupportedVersion
		expectedErr      string
	}{
		{
			name:             "identity unversioned endpoint",
			endpoint:         fakeServer.Endpoint() + "identity/",
			discoverVersions: true,
			expectedVersions: []utils.SupportedVersion{
				{
					Major:  3,
					Minor:  14,
					Status: utils.StatusCurrent,
				},
			},
		},
		{
			name:             "identity unversioned endpoint without discovery",
			endpoint:         fakeServer.Endpoint() + "identity/",
			discoverVersions: false,
			// we will still run discovery since we can't extract the version from the URL
			expectedVersions: []utils.SupportedVersion{
				{
					Major:  3,
					Minor:  14,
					Status: utils.StatusCurrent,
				},
			},
		},
		{
			name:             "identity versioned endpoint",
			endpoint:         fakeServer.Endpoint() + "identity/v3/",
			discoverVersions: true,
			expectedVersions: []utils.SupportedVersion{
				{
					Major:  3,
					Minor:  14,
					Status: utils.StatusCurrent,
				},
			},
		},
		{
			name:             "identity versioned endpoint without discovery",
			endpoint:         fakeServer.Endpoint() + "identity/v3/",
			discoverVersions: false,
			// we will skip discovery since we can extract a version from the URL
			expectedVersions: []utils.SupportedVersion{
				{
					Major:  3,
					Minor:  0,
					Status: utils.StatusUnknown,
				},
			},
		},
		{
			name:             "compute unversioned endpoint",
			endpoint:         fakeServer.Endpoint() + "compute/",
			discoverVersions: true,
			expectedVersions: []utils.SupportedVersion{
				{
					Major:  2,
					Minor:  1,
					Status: utils.StatusCurrent,
					SupportedMicroversions: utils.SupportedMicroversions{
						MaxMajor: 2, MaxMinor: 90, MinMajor: 2, MinMinor: 1,
					},
				},
				{
					Major:  2,
					Minor:  0,
					Status: utils.StatusSupported,
				},
			},
		},
		{
			name:             "compute legacy endpoint",
			endpoint:         fakeServer.Endpoint() + "compute/v2/",
			discoverVersions: true,
			expectedVersions: []utils.SupportedVersion{
				{
					Major:  2,
					Minor:  0,
					Status: utils.StatusSupported,
				},
			},
		},
		{
			name:             "compute versioned endpoint",
			endpoint:         fakeServer.Endpoint() + "compute/v2.1/",
			discoverVersions: true,
			expectedVersions: []utils.SupportedVersion{
				{
					Major:  2,
					Minor:  1,
					Status: utils.StatusCurrent,
					SupportedMicroversions: utils.SupportedMicroversions{
						MaxMajor: 2, MaxMinor: 90, MinMajor: 2, MinMinor: 1,
					},
				},
			},
		},
		{
			name:             "container-infra unversioned endpoint",
			endpoint:         fakeServer.Endpoint() + "container-infra/",
			discoverVersions: true,
			expectedVersions: []utils.SupportedVersion{
				{
					Major:  1,
					Minor:  0,
					Status: utils.StatusCurrent,
					SupportedMicroversions: utils.SupportedMicroversions{
						MaxMajor: 1, MaxMinor: 11, MinMajor: 1, MinMinor: 1,
					},
				},
			},
		},
		{
			name:             "container-infra versioned endpoint",
			endpoint:         fakeServer.Endpoint() + "container-infra/v1/",
			discoverVersions: true,
			expectedVersions: []utils.SupportedVersion{
				{
					Major:  1,
					Minor:  0,
					Status: utils.StatusUnknown,
				},
			},
		},
		{
			name:             "orchestration unversioned endpoint",
			endpoint:         fakeServer.Endpoint() + "heat-api/",
			discoverVersions: true,
			expectedVersions: []utils.SupportedVersion{
				{
					Major:  1,
					Minor:  0,
					Status: utils.StatusCurrent,
				},
			},
		},
		{
			name:             "orchestration versioned endpoint",
			endpoint:         fakeServer.Endpoint() + "heat-api/v1/",
			discoverVersions: true,
			expectedVersions: []utils.SupportedVersion{
				{
					Major:  1,
					Minor:  0,
					Status: utils.StatusUnknown,
				},
			},
		},
		{
			name:             "workflow unversioned endpoint",
			endpoint:         fakeServer.Endpoint() + "workflow/",
			discoverVersions: true,
			expectedVersions: []utils.SupportedVersion{
				{
					Major:  2,
					Minor:  0,
					Status: utils.StatusCurrent,
				},
			},
		},
		{
			name:             "workflow versioned endpoint",
			endpoint:         fakeServer.Endpoint() + "workflow/v2/",
			discoverVersions: true,
			expectedVersions: []utils.SupportedVersion{
				{
					Major:  2,
					Minor:  0,
					Status: utils.StatusUnknown,
				},
			},
		},
		{
			name:             "baremetal unversioned endpoint",
			endpoint:         fakeServer.Endpoint() + "baremetal/",
			discoverVersions: true,
			expectedVersions: []utils.SupportedVersion{
				{
					Major:  1,
					Minor:  0,
					Status: utils.StatusCurrent,
					SupportedMicroversions: utils.SupportedMicroversions{
						MaxMajor: 1, MaxMinor: 87, MinMajor: 1, MinMinor: 1,
					},
				},
			},
		},
		{
			name:             "baremetal versioned endpoint",
			endpoint:         fakeServer.Endpoint() + "baremetal/v1/",
			discoverVersions: true,
			expectedVersions: []utils.SupportedVersion{
				{
					Major:  1,
					Minor:  0,
					Status: utils.StatusCurrent,
					SupportedMicroversions: utils.SupportedMicroversions{
						MaxMajor: 1, MaxMinor: 87, MinMajor: 1, MinMinor: 1,
					},
				},
			},
		},
		{
			name:             "baremetal versioned endpoint",
			endpoint:         fakeServer.Endpoint() + "baremetal/v1/",
			discoverVersions: true,
			expectedVersions: []utils.SupportedVersion{
				{
					Major:  1,
					Minor:  0,
					Status: utils.StatusCurrent,
					SupportedMicroversions: utils.SupportedMicroversions{
						MaxMajor: 1, MaxMinor: 87, MinMajor: 1, MinMinor: 1,
					},
				},
			},
		},
		{
			name:             "fictional multi-version endpoint",
			endpoint:         fakeServer.Endpoint() + "multi-version/v1.2/",
			discoverVersions: true,
			expectedVersions: []utils.SupportedVersion{
				{
					Major:  1,
					Minor:  2,
					Status: utils.StatusCurrent,
					SupportedMicroversions: utils.SupportedMicroversions{
						MaxMajor: 1, MaxMinor: 90, MinMajor: 1, MinMinor: 2,
					},
				},
				{
					Major:  1,
					Minor:  0,
					Status: utils.StatusCurrent,
					SupportedMicroversions: utils.SupportedMicroversions{
						MaxMajor: 1, MaxMinor: 87, MinMajor: 1, MinMinor: 1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &gophercloud.ProviderClient{}

			actualVersions, err := utils.GetServiceVersions(context.TODO(), client, tt.endpoint, tt.discoverVersions)

			if tt.expectedErr != "" {
				th.AssertErr(t, err)
				if !strings.Contains(err.Error(), tt.expectedErr) {
					t.Errorf("Expected error to contain '%s', got '%s': %+v", tt.expectedErr, err, tt)
				}
				return
			} else {
				th.AssertNoErr(t, err)
			}

			th.AssertDeepEquals(t, tt.expectedVersions, actualVersions)
		})
	}
}

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
			name:        "identity unversioned endpoint",
			endpoint:    fakeServer.Endpoint() + "identity/",
			expectedErr: "not supported",
		},
		{
			// identity does not support microversions and returns error
			name:        "identity versioned endpoint",
			endpoint:    fakeServer.Endpoint() + "identity/v3/",
			expectedErr: "not supported",
		},
		{
			// compute root API does not expose microversion info and returns error
			name:        "compute unversioned endpoint",
			endpoint:    fakeServer.Endpoint() + "compute/",
			expectedErr: "not supported",
		},
		{
			// compute v2 does not support microversions and returns error
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
			name:     "container-infra unversioned endpoint",
			endpoint: fakeServer.Endpoint() + "container-infra/",
			expectedVersions: utils.SupportedMicroversions{
				MaxMajor: 1, MaxMinor: 11, MinMajor: 1, MinMinor: 1,
			},
		},
		{
			// container-infra does not expose proper discovery information
			name:        "container-infra versioned endpoint",
			endpoint:    fakeServer.Endpoint() + "container-infra/v1/",
			expectedErr: "not supported",
		},
		{
			// orchestration does not support microversions and returns error
			name:        "orchestration unversioned endpoint",
			endpoint:    fakeServer.Endpoint() + "heat-api/",
			expectedErr: "not supported",
		},
		{
			// orchestration does not support microversions and returns error
			name:        "orchestration versioned endpoint",
			endpoint:    fakeServer.Endpoint() + "heat-api/v1/",
			expectedErr: "not supported",
		},
		{
			// workflow does not support microversions and returns error
			name:        "workflow unversioned endpoint",
			endpoint:    fakeServer.Endpoint() + "workflow/",
			expectedErr: "not supported",
		},
		{
			// workflow does not support microversions and returns error
			name:        "workflow versioned endpoint",
			endpoint:    fakeServer.Endpoint() + "workflow/v2/",
			expectedErr: "not supported",
		},
		{
			name:     "baremetal unversioned endpoint",
			endpoint: fakeServer.Endpoint() + "baremetal/",
			expectedVersions: utils.SupportedMicroversions{
				MaxMajor: 1, MaxMinor: 87, MinMajor: 1, MinMinor: 1,
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
				// No point parsing and comparing versions after error, so continue to next test case
				return
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
