package testing

import (
	"github.com/gophercloud/gophercloud/utils/gnocchi/metric/v1/archivepolicies"
)

// ArchivePoliciesListResult represents raw server response from a server to a list call.
const ArchivePoliciesListResult = `[
    {
        "aggregation_methods": [
					"max",
					"min"
        ],
        "back_window": 0,
        "definition": [
            {
                "granularity": "1:00:00",
                "points": 2304,
                "timespan": "96 days, 0:00:00"
            },
            {
                "granularity": "0:05:00",
                "points": 9216,
                "timespan": "32 days, 0:00:00"
            },
            {
                "granularity": "1 day, 0:00:00",
                "points": 400,
                "timespan": "400 days, 0:00:00"
            }
        ],
        "name": "precise"
    },
    {
	        "aggregation_methods": [
							"mean",
							"sum"
	        ],
	        "back_window": 12,
	        "definition": [
			    {
					    "granularity": "1:00:00",
					    "points": 2160,
					    "timespan": "90 days, 0:00:00"
			    },
			    {
				      "granularity": "1 day, 0:00:00",
				      "points": 200,
				      "timespan": "200 days, 0:00:00"
			    }
	        ],
	        "name": "not_so_precise"
    }
]`

// ListArchivePoliciesExpected represents an expected repsonse from a List request.
var ListArchivePoliciesExpected = []archivepolicies.ArchivePolicy{
	{
		AggregationMethods: []string{
			"max",
			"min",
		},
		BackWindow: 0,
		Definitions: []archivepolicies.ArchivePolicyDefinition{
			{
				Granularity: "1:00:00",
				Points:      2304,
				TimeSpan:    "96 days, 0:00:00",
			},
			{
				Granularity: "0:05:00",
				Points:      9216,
				TimeSpan:    "32 days, 0:00:00",
			},
			{
				Granularity: "1 day, 0:00:00",
				Points:      400,
				TimeSpan:    "400 days, 0:00:00",
			},
		},
		Name: "precise",
	},
	{
		AggregationMethods: []string{
			"mean",
			"sum",
		},
		BackWindow: 12,
		Definitions: []archivepolicies.ArchivePolicyDefinition{
			{
				Granularity: "1:00:00",
				Points:      2160,
				TimeSpan:    "90 days, 0:00:00",
			},
			{
				Granularity: "1 day, 0:00:00",
				Points:      200,
				TimeSpan:    "200 days, 0:00:00",
			},
		},
		Name: "not_so_precise",
	},
}
