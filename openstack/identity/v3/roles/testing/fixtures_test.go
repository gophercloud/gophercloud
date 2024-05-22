package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/roles"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// ListOutput provides a single page of Role results.
const ListOutput = `
{
    "links": {
        "next": null,
        "previous": null,
        "self": "http://example.com/identity/v3/roles"
    },
    "roles": [
        {
            "domain_id": "default",
            "id": "2844b2a08be147a08ef58317d6471f1f",
            "links": {
                "self": "http://example.com/identity/v3/roles/2844b2a08be147a08ef58317d6471f1f"
            },
            "name": "admin-read-only"
        },
        {
            "domain_id": "1789d1",
            "id": "9fe1d3",
            "links": {
                "self": "https://example.com/identity/v3/roles/9fe1d3"
            },
            "name": "support",
            "extra": {
                "description": "read-only support role"
            }
        }
    ]
}
`

// GetOutput provides a Get result.
const GetOutput = `
{
    "role": {
        "domain_id": "1789d1",
        "id": "9fe1d3",
        "links": {
            "self": "https://example.com/identity/v3/roles/9fe1d3"
        },
        "name": "support",
        "extra": {
            "description": "read-only support role"
        }
    }
}
`

// CreateRequest provides the input to a Create request.
const CreateRequest = `
{
    "role": {
        "domain_id": "1789d1",
        "name": "support",
        "description": "read-only support role"
    }
}
`

// UpdateRequest provides the input to as Update request.
const UpdateRequest = `
{
    "role": {
        "description": "admin read-only support role"
    }
}
`

// UpdateOutput provides an update result.
const UpdateOutput = `
{
    "role": {
        "domain_id": "1789d1",
        "id": "9fe1d3",
        "links": {
            "self": "https://example.com/identity/v3/roles/9fe1d3"
        },
        "name": "support",
        "extra": {
            "description": "admin read-only support role"
        }
    }
}
`

// ListAssignmentOutput provides a result of ListAssignment request.
const ListAssignmentOutput = `
{
    "role_assignments": [
        {
            "links": {
                "assignment": "http://identity:35357/v3/domains/161718/users/313233/roles/123456"
            },
            "role": {
                "id": "123456"
            },
            "scope": {
                "domain": {
                    "id": "161718"
                }
            },
            "user": {
                "domain": {
                  "id": "161718"
                },
                "id": "313233"
            }
        },
        {
            "links": {
                "assignment": "http://identity:35357/v3/projects/456789/groups/101112/roles/123456",
                "membership": "http://identity:35357/v3/groups/101112/users/313233"
            },
            "role": {
                "id": "123456"
            },
            "scope": {
                "project": {
                    "domain": {
                      "id": "161718"
                    },
                    "id": "456789"
                }
            },
            "user": {
                "domain": {
                  "id": "161718"
                },
                "id": "313233"
            }
        }
    ],
    "links": {
        "self": "http://identity:35357/v3/role_assignments?effective",
        "previous": null,
        "next": null
    }
}
`

// ListAssignmentWithNamesOutput provides a result of ListAssignment request with IncludeNames option.
const ListAssignmentWithNamesOutput = `
{
    "role_assignments": [
        {
            "links": {
                "assignment": "http://identity:35357/v3/domains/161718/users/313233/roles/123456"
            },
            "role": {
                "id": "123456",
                "name": "include_names_role"
            },
            "scope": {
                "domain": {
                    "id": "161718",
                    "name": "52833"
                }
            },
            "user": {
                "domain": {
                    "id": "161718",
                    "name": "52833"
                },
                "id": "313233",
                "name": "example-user-name"
            }
        }
    ],
    "links": {
        "self": "http://identity:35357/v3/role_assignments?include_names=True",
        "previous": null,
        "next": null
    }
}
`

// ListAssignmentsOnResourceOutput provides a result of ListAssignmentsOnResource request.
const ListAssignmentsOnResourceOutput = `
{
    "links": {
        "next": null,
        "previous": null,
        "self": "http://example.com/identity/v3/projects/9e5a15/users/b964a9/roles"
    },
    "roles": [
        {
            "id": "9fe1d3",
            "links": {
                "self": "https://example.com/identity/v3/roles/9fe1d3"
            },
            "name": "support",
            "extra": {
                "description": "read-only support role"
            }
        }
    ]
}
`

const CreateRoleInferenceRuleOutput = `
{
    "role_inference": {
        "prior_role": {
            "id": "7ceab6192ea34a548cc71b24f72e762c",
            "links": {
                "self": "http://example.com/identity/v3/roles/7ceab6192ea34a548cc71b24f72e762c"
            },
            "name": "prior role name"
        },
        "implies": {
            "id": "97e2f5d38bc94842bc3da818c16762ed",
            "links": {
                "self": "http://example.com/identity/v3/roles/97e2f5d38bc94842bc3da818c16762ed"
            },
            "name": "implied role name"
        }
    },
    "links": {
        "self": "http://example.com/identity/v3/roles/7ceab6192ea34a548cc71b24f72e762c/implies/97e2f5d38bc94842bc3da818c16762ed"
    }
}
`

const ListRoleInferenceRulesOutput = `
{
    "role_inferences": [
        {
            "prior_role": {
                "id": "1acd3c5aa0e246b9a7427d252160dcd1",
                "links": {
                    "self": "http://example.com/identity/v3/roles/1acd3c5aa0e246b9a7427d252160dcd1"
                },
                "description": "My new role",
                "name": "prior role name"
            },
            "implies": [
                {
                    "id": "3602510e2e1f499589f78a0724dcf614",
                    "links": {
                        "self": "http://example.com/identity/v3/roles/3602510e2e1f499589f78a0724dcf614"
                    },
                    "description": "My new role",
                    "name": "implied role1 name"
                },
                {
                    "id": "738289aeef684e73a987f7cf2ec6d925",
                    "links": {
                        "self": "http://example.com/identity/v3/roles/738289aeef684e73a987f7cf2ec6d925"
                    },
                    "description": "My new role",
                    "name": "implied role2 name"
                }
            ]
        },
        {
            "prior_role": {
                "id": "bbf7a5098bb34407b7164eb6ff9f144e",
                "links": {
                    "self" : "http://example.com/identity/v3/roles/bbf7a5098bb34407b7164eb6ff9f144e"
                },
                "description": "My new role",
                "name": "prior role name"
            },
            "implies": [
                {
                    "id": "872b20ad124c4c1bafaef2b1aae316ab",
                    "links": {
                        "self": "http://example.com/identity/v3/roles/872b20ad124c4c1bafaef2b1aae316ab"
                    },
                    "description": null,
                    "name": "implied role1 name"
                },
                {
                    "id": "1d865b1b2da14cb7b05254677e5f36a2",
                    "links": {
                        "self": "http://example.com/identity/v3/roles/1d865b1b2da14cb7b05254677e5f36a2"
                    },
                    "description": null,
                    "name": "implied role2 name"
                }
            ]
        }
    ],
    "links": {
        "self": "http://example.com/identity/v3/role_inferences"
    }
}
`

// FirstRole is the first role in the List request.
var FirstRole = roles.Role{
	DomainID: "default",
	ID:       "2844b2a08be147a08ef58317d6471f1f",
	Links: map[string]any{
		"self": "http://example.com/identity/v3/roles/2844b2a08be147a08ef58317d6471f1f",
	},
	Name:  "admin-read-only",
	Extra: map[string]any{},
}

// SecondRole is the second role in the List request.
var SecondRole = roles.Role{
	DomainID: "1789d1",
	ID:       "9fe1d3",
	Links: map[string]any{
		"self": "https://example.com/identity/v3/roles/9fe1d3",
	},
	Name: "support",
	Extra: map[string]any{
		"description": "read-only support role",
	},
}

// SecondRoleUpdated is how SecondRole should look after an Update.
var SecondRoleUpdated = roles.Role{
	DomainID: "1789d1",
	ID:       "9fe1d3",
	Links: map[string]any{
		"self": "https://example.com/identity/v3/roles/9fe1d3",
	},
	Name: "support",
	Extra: map[string]any{
		"description": "admin read-only support role",
	},
}

// ExpectedRolesSlice is the slice of roles expected to be returned from ListOutput.
var ExpectedRolesSlice = []roles.Role{FirstRole, SecondRole}

// HandleListRolesSuccessfully creates an HTTP handler at `/roles` on the
// test handler mux that responds with a list of two roles.
func HandleListRolesSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/roles", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListOutput)
	})
}

// HandleGetRoleSuccessfully creates an HTTP handler at `/roles` on the
// test handler mux that responds with a single role.
func HandleGetRoleSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/roles/9fe1d3", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetOutput)
	})
}

// HandleCreateRoleSuccessfully creates an HTTP handler at `/roles` on the
// test handler mux that tests role creation.
func HandleCreateRoleSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/roles", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestJSONRequest(t, r, CreateRequest)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, GetOutput)
	})
}

// HandleUpdateRoleSuccessfully creates an HTTP handler at `/roles` on the
// test handler mux that tests role update.
func HandleUpdateRoleSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/roles/9fe1d3", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestJSONRequest(t, r, UpdateRequest)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, UpdateOutput)
	})
}

// HandleDeleteRoleSuccessfully creates an HTTP handler at `/roles` on the
// test handler mux that tests role deletion.
func HandleDeleteRoleSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/roles/9fe1d3", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

func HandleAssignSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/projects/{project_id}/users/{user_id}/roles/{role_id}", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	th.Mux.HandleFunc("/projects/{project_id}/groups/{group_id}/roles/{role_id}", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	th.Mux.HandleFunc("/domains/{domain_id}/users/{user_id}/roles/{role_id}", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	th.Mux.HandleFunc("/domains/{domain_id}/groups/{group_id}/roles/{role_id}", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})
}

func HandleUnassignSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/projects/{project_id}/users/{user_id}/roles/{role_id}", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	th.Mux.HandleFunc("/projects/{project_id}/groups/{group_id}/roles/{role_id}", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	th.Mux.HandleFunc("/domains/{domain_id}/users/{user_id}/roles/{role_id}", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	th.Mux.HandleFunc("/domains/{domain_id}/groups/{group_id}/roles/{role_id}", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})
}

// FirstRoleAssignment is the first role assignment in the List request.
var FirstRoleAssignment = roles.RoleAssignment{
	Role:  roles.AssignedRole{ID: "123456"},
	Scope: roles.Scope{Domain: roles.Domain{ID: "161718"}},
	User:  roles.User{Domain: roles.Domain{ID: "161718"}, ID: "313233"},
	Group: roles.Group{},
}

// SecondRoleAssignemnt is the second role assignemnt in the List request.
var SecondRoleAssignment = roles.RoleAssignment{
	Role:  roles.AssignedRole{ID: "123456"},
	Scope: roles.Scope{Project: roles.Project{Domain: roles.Domain{ID: "161718"}, ID: "456789"}},
	User:  roles.User{Domain: roles.Domain{ID: "161718"}, ID: "313233"},
	Group: roles.Group{},
}

// ThirdRoleAssignment is the third role assignment that has entity names in the List request.
var ThirdRoleAssignment = roles.RoleAssignment{
	Role:  roles.AssignedRole{ID: "123456", Name: "include_names_role"},
	Scope: roles.Scope{Domain: roles.Domain{ID: "161718", Name: "52833"}},
	User:  roles.User{Domain: roles.Domain{ID: "161718", Name: "52833"}, ID: "313233", Name: "example-user-name"},
	Group: roles.Group{},
}

// ExpectedRoleAssignmentsSlice is the slice of role assignments expected to be
// returned from ListAssignmentOutput.
var ExpectedRoleAssignmentsSlice = []roles.RoleAssignment{FirstRoleAssignment, SecondRoleAssignment}

// ExpectedRoleAssignmentsWithNamesSlice is the slice of role assignments expected to be
// returned from ListAssignmentWithNamesOutput.
var ExpectedRoleAssignmentsWithNamesSlice = []roles.RoleAssignment{ThirdRoleAssignment}

// HandleListRoleAssignmentsSuccessfully creates an HTTP handler at `/role_assignments` on the
// test handler mux that responds with a list of two role assignments.
func HandleListRoleAssignmentsSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/role_assignments", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListAssignmentOutput)
	})
}

// HandleListRoleAssignmentsSuccessfully creates an HTTP handler at `/role_assignments` on the
// test handler mux that responds with a list of two role assignments.
func HandleListRoleAssignmentsWithNamesSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/role_assignments", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.AssertEquals(t, "include_names=true", r.URL.RawQuery)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListAssignmentWithNamesOutput)
	})
}

// HandleListRoleAssignmentsWithSubtreeSuccessfully creates an HTTP handler at `/role_assignments` on the
// test handler mux that responds with a list of two role assignments.
func HandleListRoleAssignmentsWithSubtreeSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/role_assignments", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.AssertEquals(t, "include_subtree=true", r.URL.RawQuery)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListAssignmentOutput)
	})
}

// RoleOnResource is the role in the ListAssignmentsOnResource request.
var RoleOnResource = roles.Role{
	ID: "9fe1d3",
	Links: map[string]any{
		"self": "https://example.com/identity/v3/roles/9fe1d3",
	},
	Name: "support",
	Extra: map[string]any{
		"description": "read-only support role",
	},
}

// ExpectedRolesOnResourceSlice is the slice of roles expected to be returned
// from ListAssignmentsOnResourceOutput.
var ExpectedRolesOnResourceSlice = []roles.Role{RoleOnResource}

func HandleListAssignmentsOnResourceSuccessfully_ProjectsUsers(t *testing.T) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListAssignmentsOnResourceOutput)
	}

	th.Mux.HandleFunc("/projects/{project_id}/users/{user_id}/roles", fn)
}

func HandleListAssignmentsOnResourceSuccessfully_ProjectsGroups(t *testing.T) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListAssignmentsOnResourceOutput)
	}

	th.Mux.HandleFunc("/projects/{project_id}/groups/{group_id}/roles", fn)
}

func HandleListAssignmentsOnResourceSuccessfully_DomainsUsers(t *testing.T) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListAssignmentsOnResourceOutput)
	}

	th.Mux.HandleFunc("/domains/{domain_id}/users/{user_id}/roles", fn)
}

func HandleListAssignmentsOnResourceSuccessfully_DomainsGroups(t *testing.T) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListAssignmentsOnResourceOutput)
	}

	th.Mux.HandleFunc("/domains/{domain_id}/groups/{group_id}/roles", fn)
}

var expectedRoleInferenceRule = roles.RoleInferenceRule{
	RoleInference: roles.RoleInference{
		PriorRole: roles.PriorRole{
			ID: "7ceab6192ea34a548cc71b24f72e762c",
			Links: map[string]any{
				"self": "http://example.com/identity/v3/roles/7ceab6192ea34a548cc71b24f72e762c",
			},
			Name: "prior role name",
		},
		ImpliedRole: roles.ImpliedRole{
			ID: "97e2f5d38bc94842bc3da818c16762ed",
			Links: map[string]any{
				"self": "http://example.com/identity/v3/roles/97e2f5d38bc94842bc3da818c16762ed",
			},
			Name: "implied role name",
		},
	},
	Links: map[string]any{
		"self": "http://example.com/identity/v3/roles/7ceab6192ea34a548cc71b24f72e762c/implies/97e2f5d38bc94842bc3da818c16762ed",
	},
}

func HandleCreateRoleInferenceRule(t *testing.T) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, CreateRoleInferenceRuleOutput)
	}

	th.Mux.HandleFunc("/roles/7ceab6192ea34a548cc71b24f72e762c/implies/97e2f5d38bc94842bc3da818c16762ed", fn)
}

var expectedRoleInferenceRuleList = roles.RoleInferenceRuleList{
	RoleInferenceRuleList: []roles.RoleInferenceRules{
		{
			PriorRole: roles.PriorRoleObject{
				ID: "1acd3c5aa0e246b9a7427d252160dcd1",
				Links: map[string]any{
					"self": "http://example.com/identity/v3/roles/1acd3c5aa0e246b9a7427d252160dcd1",
				},
				Name:        "prior role name",
				Description: "My new role",
			},
			ImpliedRoles: []roles.ImpliedRoleObject{
				{
					ID: "3602510e2e1f499589f78a0724dcf614",
					Links: map[string]any{
						"self": "http://example.com/identity/v3/roles/3602510e2e1f499589f78a0724dcf614",
					},
					Name:        "implied role1 name",
					Description: "My new role",
				},
				{
					ID: "738289aeef684e73a987f7cf2ec6d925",
					Links: map[string]any{
						"self": "http://example.com/identity/v3/roles/738289aeef684e73a987f7cf2ec6d925",
					},
					Name:        "implied role2 name",
					Description: "My new role",
				},
			},
		},
		{
			PriorRole: roles.PriorRoleObject{
				ID: "bbf7a5098bb34407b7164eb6ff9f144e",
				Links: map[string]any{
					"self": "http://example.com/identity/v3/roles/bbf7a5098bb34407b7164eb6ff9f144e",
				},
				Name:        "prior role name",
				Description: "My new role",
			},
			ImpliedRoles: []roles.ImpliedRoleObject{
				{
					ID: "872b20ad124c4c1bafaef2b1aae316ab",
					Links: map[string]any{
						"self": "http://example.com/identity/v3/roles/872b20ad124c4c1bafaef2b1aae316ab",
					},
					Name:        "implied role1 name",
					Description: "",
				},
				{
					ID: "1d865b1b2da14cb7b05254677e5f36a2",
					Links: map[string]any{
						"self": "http://example.com/identity/v3/roles/1d865b1b2da14cb7b05254677e5f36a2",
					},
					Name:        "implied role2 name",
					Description: "",
				},
			},
		},
	},
	Links: map[string]any{
		"self": "http://example.com/identity/v3/role_inferences",
	},
}

func HandleListRoleInferenceRules(t *testing.T) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListRoleInferenceRulesOutput)
	}

	th.Mux.HandleFunc("/role_inferences", fn)
}

func HandleDeleteRoleInferenceRule(t *testing.T) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.WriteHeader(http.StatusNoContent)
	}
	th.Mux.HandleFunc("/roles/7ceab6192ea34a548cc71b24f72e762c/implies/97e2f5d38bc94842bc3da818c16762ed", fn)
}

func HandleGetRoleInferenceRule(t *testing.T) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, CreateRoleInferenceRuleOutput)
	}

	th.Mux.HandleFunc("/roles/7ceab6192ea34a548cc71b24f72e762c/implies/97e2f5d38bc94842bc3da818c16762ed", fn)
}
