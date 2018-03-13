package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/clustering/v1/policytypes"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListPolicyTypes(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
			"policy_type": [
				{
					"name": "senlin.policy.affinity",
					"schema": { "availability_zone": {
									"description": "Name of the availability zone to place the node",
									"required": false,
									"type": "String",
									"updatable": false
 								},
							    "enable_drs_extension": {
									"default": false,
									"description": "Enable vSphere DRS extension.",
									"required": false,
									"type": "Boolean",
									"updatable": false
								},
								"servergroup": {
									"description": "Properties of the VM server group",
									"required": false,
									"schema": {
										"name": {
											"description": "The name of the server group",
											"required": false,
											"type": "String",
											"updatable": false
										},
										"policies": {
											"constraints": [
												{
													"constraint": [
														"affinity",
														"anti-affinity"
													],
													"type": "AllowedValues"
												}
											],
											"default": "anti-affinity",
											"description": "The server group policies.",
											"required": false,
											"type": "String",
											"updatable": false
										}
									},
									"type": "Map",
									"updatable": false
								}
							  },
					"support_status": {
						"1.0": [{
							"status": "SUPPORTED",
							"since": "2016.10"
						}]
					}
				}
			]
		}`)
	})

	count := 0

	policytypes.ListDetail(fake.ServiceClient()).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := policytypes.ExtractPolicyTypes(page)
		if err != nil {
			t.Errorf("Failed to extract policy types: %v", err)
			return false, err
		}

		expected := []policytypes.PolicyType{
			{
				Name: "senlin.policy.affinity",
				Schema: policytypes.SchemaType{
					AvailabilityZone: map[string]interface{}{
						"description": "Name of the availability zone to place the node",
						"required":    false,
						"type":        "String",
						"updatable":   false,
					},
					EnableDrsExtension: map[string]interface{}{
						"default":     false,
						"description": "Enable vSphere DRS extension.",
						"required":    false,
						"type":        "Boolean",
						"updatable":   false,
					},
					Servergroup: policytypes.ServerGroupType{
						Description: "Properties of the VM server group",
						Required:    false,
						Schema: policytypes.ServerGroupSchemaType{
							policytypes.ServerGroupSchemaName{
								Description: "The name of the server group",
								Required:    false,
								Type:        "String",
								Updatable:   false,
							},
							policytypes.ServerGroupSchemaPolicies{
								Constraints: []policytypes.ServerGroupSchemaPoliciesConstraints{
									{
										Constraint: []string{"affinity", "anti-affinity"},
										Type:       "AllowedValues",
									},
								},
								Default:     "anti-affinity",
								Description: "The server group policies.",
								Required:    false,
								Type:        "String",
								Updatable:   false,
							},
						},
						Type:      "Map",
						Updatable: false,
					},
				},
				SupportStatus: policytypes.SupportStatusType{
					SupportVersion: map[string]interface{}{
						"1.0": []interface{}{
							map[string]interface{}{
								"status": "SUPPORTED",
								"since":  "2016.10",
							},
						},
					},
				},
			},
		}

		th.AssertDeepEquals(t, expected, actual)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestNonJSONCannotBeExtractedIntoPolicyTypes(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	policytypes.ListDetail(fake.ServiceClient()).EachPage(func(page pagination.Page) (bool, error) {
		if _, err := policytypes.ExtractPolicyTypes(page); err == nil {
			t.Fatalf("Expected error, got nil")
		}
		return true, nil
	})
}
