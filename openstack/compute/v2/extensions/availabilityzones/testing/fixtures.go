package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	az "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/availabilityzones"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

const GetOutput = `
{
    "availabilityZoneInfo": [
        {
            "hosts": null,
            "zoneName": "nova",
            "zoneState": {
                "available": true
            }
        }
    ]
}
`

const GetDetailOutput = `
{
    "availabilityZoneInfo": [
        {
            "hosts": {
                "conductor": {
                    "nova-conductor": {
                        "active": true,
                        "available": true,
                        "updated_at": null
                    }
                },
                "consoleauth": {
                    "nova-consoleauth": {
                        "active": true,
                        "available": true,
                        "updated_at": null
                    }
                },
                "network": {
                    "nova-network": {
                        "active": true,
                        "available": true,
                        "updated_at": "2015-09-18T09:50:55.000000"
                    }
                },
                "scheduler": {
                    "nova-scheduler": {
                        "active": true,
                        "available": true,
                        "updated_at": null
                    }
                }
            },
            "zoneName": "internal",
            "zoneState": {
                "available": true
            }
        },
        {
            "hosts": {
                "compute": {
                    "nova-compute": {
                        "active": true,
                        "available": true,
                        "updated_at": null
                    }
                }
            },
            "zoneName": "nova",
            "zoneState": {
                "available": true
            }
        }
    ]
}
`

var OSAZResult = az.OSAvailabilityZone{
	AvailabilityZoneInfo: []az.AvailabilityZone{
		{
			Hosts:     nil,
			ZoneName:  "nova",
			ZoneState: az.ZoneState{Available: true},
		},
	},
}

var nilTime time.Time
var OSAZDetailResult = az.OSAvailabilityZone{
	AvailabilityZoneInfo: []az.AvailabilityZone{
		{
			Hosts: az.Hosts{
				"conductor": az.Services{
					NovaConductor: az.StateofService{
						Active:    true,
						Available: true,
						UpdatedAt: nilTime,
					},
				},
				"consoleauth": az.Services{
					NovaConsoleauth: az.StateofService{
						Active:    true,
						Available: true,
						UpdatedAt: nilTime,
					},
				},
				"network": az.Services{
					NovaNetwork: az.StateofService{
						Active:    true,
						Available: true,
						UpdatedAt: time.Date(2015, 9, 18, 9, 50, 55, 0, time.UTC),
					},
				},
				"scheduler": az.Services{
					NovaScheduler: az.StateofService{
						Active:    true,
						Available: true,
						UpdatedAt: nilTime,
					},
				},
			},
			ZoneName:  "internal",
			ZoneState: az.ZoneState{Available: true},
		},
		{
			Hosts: az.Hosts{
				"compute": az.Services{
					NovaCompute: az.StateofService{
						Active:    true,
						Available: true,
						UpdatedAt: nilTime,
					},
				},
			},
			ZoneName:  "nova",
			ZoneState: az.ZoneState{Available: true},
		},
	},
}

// HandleGetSuccessfully configures the test server to respond to a Get request
// for availability zone information.
func HandleGetSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/os-availability-zone", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, GetOutput)
	})
}

// HandleGetDetailSuccessfully configures the test server to respond to a Get request
// for detailed availability zone information.
func HandleGetDetailSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/os-availability-zone/detail", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, GetDetailOutput)
	})
}
