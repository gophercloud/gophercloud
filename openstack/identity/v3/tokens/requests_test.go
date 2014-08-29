package tokens

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/testhelper"
)

// authTokenPost verifies that providing certain AuthOptions and Scope results in an expected JSON structure.
func authTokenPost(t *testing.T, options gophercloud.AuthOptions, scope *Scope, requestJSON string) {
	setup()
	defer teardown()

	client := gophercloud.ServiceClient{
		Endpoint: endpoint(),
		Options:  options,
		TokenID:  "12345abcdef",
	}

	mux.HandleFunc("/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "POST")
		testhelper.TestHeader(t, r, "Content-Type", "application/json")
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestJSONRequest(t, r, requestJSON)

		fmt.Fprintf(w, `{}`)
	})

	_, err := Create(&client, scope)
	if err != nil {
		t.Errorf("Create returned an error: %v", err)
	}
}

func authTokenPostErr(t *testing.T, options gophercloud.AuthOptions, scope *Scope, includeToken bool, expectedErr error) {
	setup()
	defer teardown()

	client := gophercloud.ServiceClient{
		Endpoint: endpoint(),
		Options:  options,
	}
	if includeToken {
		client.TokenID = "abcdef123456"
	}

	_, err := Create(&client, scope)
	if err == nil {
		t.Errorf("Create did NOT return an error")
	}
	if err != expectedErr {
		t.Errorf("Create returned an unexpected error: wanted %v, got %v", expectedErr, err)
	}
}

func TestCreateUserIDAndPassword(t *testing.T) {
	authTokenPost(t, gophercloud.AuthOptions{UserID: "me", Password: "squirrel!"}, nil, `
		{
			"auth": {
				"identity": {
					"methods": ["password"],
					"password": {
						"user": { "id": "me", "password": "squirrel!" }
					}
				}
			}
		}
	`)
}

func TestCreateUsernameDomainIDPassword(t *testing.T) {
	authTokenPost(t, gophercloud.AuthOptions{Username: "fakey", Password: "notpassword", DomainID: "abc123"}, nil, `
		{
			"auth": {
				"identity": {
					"methods": ["password"],
					"password": {
						"user": {
							"domain": {
								"id": "abc123"
							},
							"name": "fakey",
							"password": "notpassword"
						}
					}
				}
			}
		}
	`)
}

func TestCreateUsernameDomainNamePassword(t *testing.T) {
	authTokenPost(t, gophercloud.AuthOptions{Username: "frank", Password: "swordfish", DomainName: "spork.net"}, nil, `
		{
			"auth": {
				"identity": {
					"methods": ["password"],
					"password": {
						"user": {
							"domain": {
								"name": "spork.net"
							},
							"name": "frank",
							"password": "swordfish"
						}
					}
				}
			}
		}
	`)
}

func TestCreateTokenID(t *testing.T) {
	authTokenPost(t, gophercloud.AuthOptions{}, nil, `
		{
			"auth": {
				"identity": {
					"methods": ["token"],
					"token": {
						"id": "12345abcdef"
					}
				}
			}
		}
	`)
}

func TestCreateProjectIDScope(t *testing.T) {
	options := gophercloud.AuthOptions{UserID: "fenris", Password: "g0t0h311"}
	scope := &Scope{ProjectID: "123456"}
	authTokenPost(t, options, scope, `
		{
			"auth": {
				"identity": {
					"methods": ["password"],
					"password": {
						"user": {
							"id": "fenris",
							"password": "g0t0h311"
						}
					}
				},
				"scope": {
					"project": {
						"id": "123456"
					}
				}
			}
		}
	`)
}

func TestCreateDomainIDScope(t *testing.T) {
	options := gophercloud.AuthOptions{UserID: "fenris", Password: "g0t0h311"}
	scope := &Scope{DomainID: "1000"}
	authTokenPost(t, options, scope, `
		{
			"auth": {
				"identity": {
					"methods": ["password"],
					"password": {
						"user": {
							"id": "fenris",
							"password": "g0t0h311"
						}
					}
				},
				"scope": {
					"domain": {
						"id": "1000"
					}
				}
			}
		}
	`)
}

func TestCreateProjectNameAndDomainIDScope(t *testing.T) {
	options := gophercloud.AuthOptions{UserID: "fenris", Password: "g0t0h311"}
	scope := &Scope{ProjectName: "world-domination", DomainID: "1000"}
	authTokenPost(t, options, scope, `
		{
			"auth": {
				"identity": {
					"methods": ["password"],
					"password": {
						"user": {
							"id": "fenris",
							"password": "g0t0h311"
						}
					}
				},
				"scope": {
					"project": {
						"domain": {
							"id": "1000"
						},
						"name": "world-domination"
					}
				}
			}
		}
	`)
}

func TestCreateProjectNameAndDomainNameScope(t *testing.T) {
	options := gophercloud.AuthOptions{UserID: "fenris", Password: "g0t0h311"}
	scope := &Scope{ProjectName: "world-domination", DomainName: "evil-plans"}
	authTokenPost(t, options, scope, `
		{
			"auth": {
				"identity": {
					"methods": ["password"],
					"password": {
						"user": {
							"id": "fenris",
							"password": "g0t0h311"
						}
					}
				},
				"scope": {
					"project": {
						"domain": {
							"name": "evil-plans"
						},
						"name": "world-domination"
					}
				}
			}
		}
	`)
}

func TestCreateFailureEmptyAuth(t *testing.T) {
	authTokenPostErr(t, gophercloud.AuthOptions{}, nil, false, ErrMissingPassword)
}

func TestCreateFailureAPIKey(t *testing.T) {
	authTokenPostErr(t, gophercloud.AuthOptions{APIKey: "something"}, nil, false, ErrAPIKeyProvided)
}

func TestCreateFailureTenantID(t *testing.T) {
	authTokenPostErr(t, gophercloud.AuthOptions{TenantID: "something"}, nil, false, ErrTenantIDProvided)
}

func TestCreateFailureTenantName(t *testing.T) {
	authTokenPostErr(t, gophercloud.AuthOptions{TenantName: "something"}, nil, false, ErrTenantNameProvided)
}

func TestCreateFailureTokenIDUsername(t *testing.T) {
	authTokenPostErr(t, gophercloud.AuthOptions{Username: "something"}, nil, true, ErrUsernameWithToken)
}

func TestCreateFailureTokenIDUserID(t *testing.T) {
	authTokenPostErr(t, gophercloud.AuthOptions{UserID: "something"}, nil, true, ErrUserIDWithToken)
}

func TestCreateFailureTokenIDDomainID(t *testing.T) {
	authTokenPostErr(t, gophercloud.AuthOptions{DomainID: "something"}, nil, true, ErrDomainIDWithToken)
}

func TestCreateFailureTokenIDDomainName(t *testing.T) {
	authTokenPostErr(t, gophercloud.AuthOptions{DomainName: "something"}, nil, true, ErrDomainNameWithToken)
}

func TestCreateFailureMissingUser(t *testing.T) {
	options := gophercloud.AuthOptions{Password: "supersecure"}
	authTokenPostErr(t, options, nil, false, ErrUsernameOrUserID)
}

func TestCreateFailureBothUser(t *testing.T) {
	options := gophercloud.AuthOptions{
		Password: "supersecure",
		Username: "oops",
		UserID:   "redundancy",
	}
	authTokenPostErr(t, options, nil, false, ErrUsernameOrUserID)
}

func TestCreateFailureMissingDomain(t *testing.T) {
	options := gophercloud.AuthOptions{
		Password: "supersecure",
		Username: "notuniqueenough",
	}
	authTokenPostErr(t, options, nil, false, ErrDomainIDOrDomainName)
}

func TestCreateFailureBothDomain(t *testing.T) {
	options := gophercloud.AuthOptions{
		Password:   "supersecure",
		Username:   "someone",
		DomainID:   "hurf",
		DomainName: "durf",
	}
	authTokenPostErr(t, options, nil, false, ErrDomainIDOrDomainName)
}

func TestCreateFailureUserIDDomainID(t *testing.T) {
	options := gophercloud.AuthOptions{
		UserID:   "100",
		Password: "stuff",
		DomainID: "oops",
	}
	authTokenPostErr(t, options, nil, false, ErrDomainIDWithUserID)
}

func TestCreateFailureUserIDDomainName(t *testing.T) {
	options := gophercloud.AuthOptions{
		UserID:     "100",
		Password:   "sssh",
		DomainName: "oops",
	}
	authTokenPostErr(t, options, nil, false, ErrDomainNameWithUserID)
}

func TestCreateFailureScopeProjectNameAlone(t *testing.T) {
	options := gophercloud.AuthOptions{UserID: "myself", Password: "swordfish"}
	scope := &Scope{ProjectName: "notenough"}
	authTokenPostErr(t, options, scope, false, ErrScopeDomainIDOrDomainName)
}

func TestCreateFailureScopeProjectNameAndID(t *testing.T) {
	options := gophercloud.AuthOptions{UserID: "myself", Password: "swordfish"}
	scope := &Scope{ProjectName: "whoops", ProjectID: "toomuch", DomainID: "1234"}
	authTokenPostErr(t, options, scope, false, ErrScopeProjectIDOrProjectName)
}

func TestCreateFailureScopeProjectIDAndDomainID(t *testing.T) {
	options := gophercloud.AuthOptions{UserID: "myself", Password: "swordfish"}
	scope := &Scope{ProjectID: "toomuch", DomainID: "notneeded"}
	authTokenPostErr(t, options, scope, false, ErrScopeProjectIDAlone)
}

func TestCreateFailureScopeProjectIDAndDomainNAme(t *testing.T) {
	options := gophercloud.AuthOptions{UserID: "myself", Password: "swordfish"}
	scope := &Scope{ProjectID: "toomuch", DomainName: "notneeded"}
	authTokenPostErr(t, options, scope, false, ErrScopeProjectIDAlone)
}

func TestCreateFailureScopeDomainIDAndDomainName(t *testing.T) {
	options := gophercloud.AuthOptions{UserID: "myself", Password: "swordfish"}
	scope := &Scope{DomainID: "toomuch", DomainName: "notneeded"}
	authTokenPostErr(t, options, scope, false, ErrScopeDomainIDOrDomainName)
}

func TestCreateFailureScopeDomainNameAlone(t *testing.T) {
	options := gophercloud.AuthOptions{UserID: "myself", Password: "swordfish"}
	scope := &Scope{DomainName: "notenough"}
	authTokenPostErr(t, options, scope, false, ErrScopeDomainName)
}

func TestCreateFailureEmptyScope(t *testing.T) {
	options := gophercloud.AuthOptions{UserID: "myself", Password: "swordfish"}
	scope := &Scope{}
	authTokenPostErr(t, options, scope, false, ErrScopeEmpty)
}
