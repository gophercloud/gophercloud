package tokens

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/testhelper"
)

// authTokenPost verifies that providing certain AuthOptions and Scope results in an expected JSON structure.
func authTokenPost(t *testing.T, options AuthOptionsBuilder, scope *gophercloud.ScopeOptsV3, requestJSON string) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{
			TokenID: "12345abcdef",
		},
		Endpoint: testhelper.Endpoint(),
	}

	testhelper.Mux.HandleFunc("/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "POST")
		testhelper.TestHeader(t, r, "Content-Type", "application/json")
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestJSONRequest(t, r, requestJSON)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{
			"token": {
				"expires_at": "2014-10-02T13:45:00.000000Z"
			}
		}`)
	})

	_, err := Create(&client, options, scope).Extract()
	if err != nil {
		t.Errorf("Create returned an error: %v", err)
	}
}

func authTokenPostErr(t *testing.T, options AuthOptionsBuilder, scope *gophercloud.ScopeOptsV3, includeToken bool, expectedErr error) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       testhelper.Endpoint(),
	}
	if includeToken {
		client.TokenID = "abcdef123456"
	}

	_, err := Create(&client, options, scope).Extract()
	if err == nil {
		t.Errorf("Create did NOT return an error")
	}
	if err != expectedErr {
		t.Errorf("Create returned an unexpected error: wanted %v, got %v", expectedErr, err)
	}
}

func TestCreateUserIDAndPassword(t *testing.T) {
	ao := gophercloud.AuthOptions{}
	ao.UserID = "me"
	ao.Password = "squirrel!"
	authTokenPost(t, ao, nil, `
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
	ao := gophercloud.AuthOptions{DomainID: "abc123"}
	ao.Username = "fakey"
	ao.Password = "notpassword"
	authTokenPost(t, ao, nil, `
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
	ao := gophercloud.AuthOptions{DomainName: "spork.net"}
	ao.Username = "frank"
	ao.Password = "swordfish"
	authTokenPost(t, ao, nil, `
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
	authTokenPost(t, gophercloud.AuthOptions{TokenID: "12345abcdef"}, nil, `
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
	ao := gophercloud.AuthOptions{}
	ao.UserID = "fenris"
	ao.Password = "g0t0h311"
	scope := &gophercloud.ScopeOptsV3{ProjectID: "123456"}
	authTokenPost(t, ao, scope, `
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
	ao := gophercloud.AuthOptions{}
	ao.UserID = "fenris"
	ao.Password = "g0t0h311"
	scope := &gophercloud.ScopeOptsV3{DomainID: "1000"}
	authTokenPost(t, ao, scope, `
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
	ao := gophercloud.AuthOptions{}
	ao.UserID = "fenris"
	ao.Password = "g0t0h311"
	scope := &gophercloud.ScopeOptsV3{ProjectName: "world-domination", DomainID: "1000"}
	authTokenPost(t, ao, scope, `
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
	ao := gophercloud.AuthOptions{}
	ao.UserID = "fenris"
	ao.Password = "g0t0h311"
	scope := &gophercloud.ScopeOptsV3{ProjectName: "world-domination", DomainName: "evil-plans"}
	authTokenPost(t, ao, scope, `
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

func TestCreateExtractsTokenFromResponse(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       testhelper.Endpoint(),
	}

	testhelper.Mux.HandleFunc("/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Subject-Token", "aaa111")

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{
			"token": {
				"expires_at": "2014-10-02T13:45:00.000000Z"
			}
		}`)
	})

	ao := gophercloud.AuthOptions{}
	ao.UserID = "me"
	ao.Password = "shhh"
	token, err := Create(&client, ao, nil).Extract()
	if err != nil {
		t.Fatalf("Create returned an error: %v", err)
	}

	if token.ID != "aaa111" {
		t.Errorf("Expected token to be aaa111, but was %s", token.ID)
	}
}

func TestCreateFailureEmptyAuth(t *testing.T) {
	authTokenPostErr(t, gophercloud.AuthOptions{}, nil, false, ErrMissingPassword)
}

func TestCreateFailureTokenIDUsername(t *testing.T) {
	ao := gophercloud.AuthOptions{}
	ao.Username = "somthing"
	authTokenPostErr(t, ao, nil, true, ErrUsernameWithToken)
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
	ao := gophercloud.AuthOptions{}
	ao.Password = "supersecure"
	authTokenPostErr(t, ao, nil, false, ErrUsernameOrUserID)
}

func TestCreateFailureBothUser(t *testing.T) {
	ao := gophercloud.AuthOptions{}
	ao.UserID = "redundancy"
	ao.Username = "oops"
	ao.Password = "supersecure"
	authTokenPostErr(t, ao, nil, false, ErrUsernameOrUserID)
}

func TestCreateFailureMissingDomain(t *testing.T) {
	ao := gophercloud.AuthOptions{}
	ao.Username = "notuniqueenough"
	ao.Password = "supersecure"
	authTokenPostErr(t, ao, nil, false, ErrDomainIDOrDomainName)
}

func TestCreateFailureBothDomain(t *testing.T) {
	ao := gophercloud.AuthOptions{}
	ao.Username = "someone"
	ao.Password = "supersecure"
	ao.DomainID = "hurf"
	ao.DomainName = "durf"
	authTokenPostErr(t, ao, nil, false, ErrDomainIDOrDomainName)
}

func TestCreateFailureUserIDDomainID(t *testing.T) {
	ao := gophercloud.AuthOptions{}
	ao.UserID = "100"
	ao.Password = "stuff"
	ao.DomainID = "oops"
	authTokenPostErr(t, ao, nil, false, ErrDomainIDWithUserID)
}

func TestCreateFailureUserIDDomainName(t *testing.T) {
	ao := gophercloud.AuthOptions{}
	ao.UserID = "100"
	ao.Password = "sssh"
	ao.DomainName = "oops"
	authTokenPostErr(t, ao, nil, false, ErrDomainNameWithUserID)
}

func TestCreateFailureScopeProjectNameAlone(t *testing.T) {
	ao := gophercloud.AuthOptions{}
	ao.UserID = "myself"
	ao.Password = "swordfish"
	scope := &gophercloud.ScopeOptsV3{ProjectName: "notenough"}
	authTokenPostErr(t, ao, scope, false, ErrScopeDomainIDOrDomainName)
}

func TestCreateFailureScopeProjectNameAndID(t *testing.T) {
	ao := gophercloud.AuthOptions{}
	ao.UserID = "myself"
	ao.Password = "swordfish"
	scope := &gophercloud.ScopeOptsV3{ProjectName: "whoops", ProjectID: "toomuch", DomainID: "1234"}
	authTokenPostErr(t, ao, scope, false, ErrScopeProjectIDOrProjectName)
}

func TestCreateFailureScopeProjectIDAndDomainID(t *testing.T) {
	ao := gophercloud.AuthOptions{}
	ao.UserID = "myself"
	ao.Password = "swordfish"
	scope := &gophercloud.ScopeOptsV3{ProjectID: "toomuch", DomainID: "notneeded"}
	authTokenPostErr(t, ao, scope, false, ErrScopeProjectIDAlone)
}

func TestCreateFailureScopeProjectIDAndDomainNAme(t *testing.T) {
	ao := gophercloud.AuthOptions{}
	ao.UserID = "myself"
	ao.Password = "swordfish"
	scope := &gophercloud.ScopeOptsV3{ProjectID: "toomuch", DomainName: "notneeded"}
	authTokenPostErr(t, ao, scope, false, ErrScopeProjectIDAlone)
}

func TestCreateFailureScopeDomainIDAndDomainName(t *testing.T) {
	ao := gophercloud.AuthOptions{}
	ao.UserID = "myself"
	ao.Password = "swordfish"
	scope := &gophercloud.ScopeOptsV3{DomainID: "toomuch", DomainName: "notneeded"}
	authTokenPostErr(t, ao, scope, false, ErrScopeDomainIDOrDomainName)
}

func TestCreateFailureScopeDomainNameAlone(t *testing.T) {
	ao := gophercloud.AuthOptions{}
	ao.UserID = "myself"
	ao.Password = "swordfish"
	scope := &gophercloud.ScopeOptsV3{DomainName: "notenough"}
	authTokenPostErr(t, ao, scope, false, ErrScopeDomainName)
}

func TestCreateFailureEmptyScope(t *testing.T) {
	ao := gophercloud.AuthOptions{}
	ao.UserID = "myself"
	ao.Password = "swordfish"
	scope := &gophercloud.ScopeOptsV3{}
	authTokenPostErr(t, ao, scope, false, ErrScopeEmpty)
}

func TestGetRequest(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{
			TokenID: "12345abcdef",
		},
		Endpoint: testhelper.Endpoint(),
	}

	testhelper.Mux.HandleFunc("/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "Content-Type", "")
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestHeader(t, r, "X-Auth-Token", "12345abcdef")
		testhelper.TestHeader(t, r, "X-Subject-Token", "abcdef12345")

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `
			{ "token": { "expires_at": "2014-08-29T13:10:01.000000Z" } }
		`)
	})

	token, err := Get(&client, "abcdef12345").Extract()
	if err != nil {
		t.Errorf("Info returned an error: %v", err)
	}

	expected, _ := time.Parse(time.UnixDate, "Fri Aug 29 13:10:01 UTC 2014")
	if token.ExpiresAt != expected {
		t.Errorf("Expected expiration time %s, but was %s", expected.Format(time.UnixDate), token.ExpiresAt.Format(time.UnixDate))
	}
}

func prepareAuthTokenHandler(t *testing.T, expectedMethod string, status int) gophercloud.ServiceClient {
	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{
			TokenID: "12345abcdef",
		},
		Endpoint: testhelper.Endpoint(),
	}

	testhelper.Mux.HandleFunc("/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, expectedMethod)
		testhelper.TestHeader(t, r, "Content-Type", "")
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestHeader(t, r, "X-Auth-Token", "12345abcdef")
		testhelper.TestHeader(t, r, "X-Subject-Token", "abcdef12345")

		w.WriteHeader(status)
	})

	return client
}

func TestValidateRequestSuccessful(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	client := prepareAuthTokenHandler(t, "HEAD", http.StatusNoContent)

	ok, err := Validate(&client, "abcdef12345")
	if err != nil {
		t.Errorf("Unexpected error from Validate: %v", err)
	}

	if !ok {
		t.Errorf("Validate returned false for a valid token")
	}
}

func TestValidateRequestFailure(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	client := prepareAuthTokenHandler(t, "HEAD", http.StatusNotFound)

	ok, err := Validate(&client, "abcdef12345")
	if err != nil {
		t.Errorf("Unexpected error from Validate: %v", err)
	}

	if ok {
		t.Errorf("Validate returned true for an invalid token")
	}
}

func TestValidateRequestError(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	client := prepareAuthTokenHandler(t, "HEAD", http.StatusUnauthorized)

	_, err := Validate(&client, "abcdef12345")
	if err == nil {
		t.Errorf("Missing expected error from Validate")
	}
}

func TestRevokeRequestSuccessful(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	client := prepareAuthTokenHandler(t, "DELETE", http.StatusNoContent)

	res := Revoke(&client, "abcdef12345")
	testhelper.AssertNoErr(t, res.Err)
}

func TestRevokeRequestError(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	client := prepareAuthTokenHandler(t, "DELETE", http.StatusNotFound)

	res := Revoke(&client, "abcdef12345")
	if res.Err == nil {
		t.Errorf("Missing expected error from Revoke")
	}
}
