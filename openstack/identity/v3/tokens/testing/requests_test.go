package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// authTokenPost verifies that providing certain AuthOptions and Scope results in an expected JSON structure.
func authTokenPost(t *testing.T, options auth.AuthOptionsBuilderV3, requestJSON string) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}

	fakeServer.Mux.HandleFunc("/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, requestJSON)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"token": {
				"expires_at": "2014-10-02T13:45:00.000000Z"
			}
		}`)
	})

	expected := &tokens.Token{
		ExpiresAt: time.Date(2014, 10, 2, 13, 45, 0, 0, time.UTC),
	}
	actual, err := tokens.Create(context.TODO(), &client, options).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, expected, actual)
}

func authTokenPostErr(t *testing.T, options auth.AuthOptionsBuilderV3, expectedErr error) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}

	_, err := tokens.Create(context.TODO(), &client, options).Extract()
	if err == nil {
		t.Errorf("Create did NOT return an error")
	}
	if err != expectedErr {
		t.Errorf("Create returned an unexpected error: wanted %v, got %v", expectedErr, err)
	}
}

func TestCreateUserIDAndPassword(t *testing.T) {
	opts := auth.V3PasswordOpts{
		UserID:   "me",
		Password: "squirrel!",
	}

	authTokenPost(t, opts, `
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
	opts := auth.V3PasswordOpts{
		Username:     "fakey",
		Password:     "notpassword",
		UserDomainID: "abc123",
	}

	authTokenPost(t, opts, `
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

	opts := auth.V3PasswordOpts{
		Username:       "frank",
		Password:       "swordfish",
		UserDomainName: "spork.net",
	}

	authTokenPost(t, opts, `
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

	opts := auth.V3TokenOpts{
		Token: "12345abcdef",
	}

	authTokenPost(t, opts, `
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

	opts := auth.V3PasswordOpts{
		UserID:   "someuser",
		Password: "somepassword",
		Scope: &auth.Scope{
			ProjectID: "123456",
		},
	}

	authTokenPost(t, opts, `
		{
			"auth": {
				"identity": {
					"methods": ["password"],
					"password": {
						"user": {
							"id": "someuser",
							"password": "somepassword"
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
	opts := auth.V3PasswordOpts{
		UserID:   "someuser",
		Password: "somepassword",
		Scope: &auth.Scope{
			DomainID: "1000",
		},
	}

	authTokenPost(t, opts, `
		{
			"auth": {
				"identity": {
					"methods": ["password"],
					"password": {
						"user": {
							"id": "someuser",
							"password": "somepassword"
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

func TestCreateDomainNameScope(t *testing.T) {
	opts := auth.V3PasswordOpts{
		UserID:   "someuser",
		Password: "somepassword",
		Scope: &auth.Scope{
			DomainName: "evil-plans",
		},
	}

	authTokenPost(t, opts, `
		{
			"auth": {
				"identity": {
					"methods": ["password"],
					"password": {
						"user": {
							"id": "someuser",
							"password": "somepassword"
						}
					}
				},
				"scope": {
					"domain": {
						"name": "evil-plans"
					}
				}
			}
		}
	`)
}

func TestCreateProjectNameAndDomainIDScope(t *testing.T) {
	opts := auth.V3PasswordOpts{
		UserID:   "someuser",
		Password: "somepassword",
		Scope: &auth.Scope{
			ProjectName:     "world-domination",
			ProjectDomainID: "1000",
		},
	}

	authTokenPost(t, opts, `
		{
			"auth": {
				"identity": {
					"methods": ["password"],
					"password": {
						"user": {
							"id": "someuser",
							"password": "somepassword"
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
	opts := auth.V3PasswordOpts{
		UserID:   "someuser",
		Password: "somepassword",
		Scope: &auth.Scope{
			ProjectName:       "world-domination",
			ProjectDomainName: "evil-plans",
		},
	}

	authTokenPost(t, opts, `
		{
			"auth": {
				"identity": {
					"methods": ["password"],
					"password": {
						"user": {
							"id": "someuser",
							"password": "somepassword"
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

func TestCreateSystemScope(t *testing.T) {
	opts := auth.V3PasswordOpts{
		UserID:   "someuser",
		Password: "somepassword",
		Scope: &auth.Scope{
			System: true,
		},
	}

	authTokenPost(t, opts, `
		{
			"auth": {
				"identity": {
					"methods": ["password"],
					"password": {
						"user": {
							"id": "someuser",
							"password": "somepassword"
						}
					}
				},
				"scope": {
					"system": {
						"all": true
					}
				}
			}
		}
	`)
}

func TestCreateUserIDPasswordTrustID(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	requestJSON := `{
		"auth": {
			"identity": {
				"methods": ["password"],
				"password": {
					"user": { "id": "demo", "password": "squirrel!" }
				}
			},
			"scope": {
				"OS-TRUST:trust": {
					"id": "95946f9eef864fdc993079d8fe3e5747"
				}
			}
		}
	}`
	responseJSON := `{
		"token": {
			"OS-TRUST:trust": {
				"id": "95946f9eef864fdc993079d8fe3e5747",
				"impersonation": false,
				"trustee_user": {
					"id": "64f9caa2872b442c98d42a986ee3b37a"
				},
				"trustor_user": {
					"id": "c88693b7c81c408e9084ac1e51082bfb"
				}
			},
			"audit_ids": [
				"wwcoUZGPR6mCIIl-COn8Kg"
			],
			"catalog": [],
			"expires_at": "2024-02-28T12:10:39.000000Z",
			"issued_at": "2024-02-28T11:10:39.000000Z",
			"methods": [
				"password"
			],
			"project": {
				"domain": {
					"id": "default",
					"name": "Default"
				},
				"id": "1fd93a4455c74d2ea94b929fc5f0e488",
				"name": "admin"
			},
			"roles": [],
			"user": {
				"domain": {
					"id": "default",
					"name": "Default"
				},
				"id": "64f9caa2872b442c98d42a986ee3b37a",
				"name": "demo",
				"password_expires_at": null
			}
		}
	}`
	fakeServer.Mux.HandleFunc("/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, requestJSON)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, responseJSON)
	})

	opts := auth.V3PasswordOpts{
		UserID:   "demo",
		Password: "squirrel!",
		Scope: &auth.Scope{
			TrustID: "95946f9eef864fdc993079d8fe3e5747",
		},
	}

	rsp := tokens.Create(context.TODO(), client.ServiceClient(fakeServer), opts)

	token, err := rsp.Extract()
	th.AssertNoErr(t, err)
	expectedToken := &tokens.Token{
		ExpiresAt: time.Date(2024, 02, 28, 12, 10, 39, 0, time.UTC),
	}
	th.AssertDeepEquals(t, expectedToken, token)

	trust, err := rsp.ExtractTrust()
	th.AssertNoErr(t, err)
	expectedTrust := &tokens.Trust{
		ID:            "95946f9eef864fdc993079d8fe3e5747",
		Impersonation: false,
		TrusteeUserID: tokens.TrustUser{
			ID: "64f9caa2872b442c98d42a986ee3b37a",
		},
		TrustorUserID: tokens.TrustUser{
			ID: "c88693b7c81c408e9084ac1e51082bfb",
		},
	}
	th.AssertDeepEquals(t, expectedTrust, trust)
}

func TestCreateApplicationCredentialIDAndSecret(t *testing.T) {
	opts := auth.V3ApplicationCredentialOpts{
		ApplicationCredentialID:     "12345abcdef",
		ApplicationCredentialSecret: "mysecret",
	}

	authTokenPost(t, opts, `
		{
			"auth": {
				"identity": {
					"application_credential": {
						"id": "12345abcdef",
						"secret": "mysecret"
					},
					"methods": [
						"application_credential"
					]
				}
			}
		}
	`)
}

func TestCreateApplicationCredentialNameAndSecret(t *testing.T) {
	opts := auth.V3ApplicationCredentialOpts{
		ApplicationCredentialName:   "myappcred",
		ApplicationCredentialSecret: "mysecret",
		Username:                    "someuser",
		UserDomainName:              "evil-plans",
	}

	authTokenPost(t, opts, `
		{
			"auth": {
				"identity": {
					"application_credential": {
						"name": "myappcred",
						"secret": "mysecret",
						"user": {
							"name": "someuser",
							"domain": {
								"name": "evil-plans"
							}
						}
					},
					"methods": [
						"application_credential"
					]
				}
			}
		}
	`)
}

func TestCreateTOTPProjectNameAndDomainNameScope(t *testing.T) {
	opts := auth.V3TOTPOpts{
		UserID:   "someuser",
		Passcode: "12345678",
		Scope: &auth.Scope{
			ProjectName:       "world-domination",
			ProjectDomainName: "evil-plans",
		},
	}

	authTokenPost(t, opts, `
		{
			"auth": {
				"identity": {
					"methods": ["totp"],
					"totp": {
						"user": {
							"id": "someuser",
							"passcode": "12345678"
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

func TestCreatePasswordTOTPProjectNameAndDomainNameScope(t *testing.T) {
	opts := auth.V3MultifactorOpts{
		AuthMethods: []auth.AuthOptionsBuilderV3{
			auth.V3PasswordOpts{
				UserID:   "someuser",
				Password: "somepassword",
			},
			auth.V3TOTPOpts{
				UserID:   "someuser",
				Passcode: "12345678",
			},
		},
		Scope: &auth.Scope{
			ProjectName:       "world-domination",
			ProjectDomainName: "evil-plans",
		},
	}

	authTokenPost(t, opts, `
		{
			"auth": {
				"identity": {
					"methods": ["password","totp"],
					"password": {
						"user": {
							"id": "someuser",
							"password": "somepassword"
						}
					},
					"totp": {
						"user": {
							"id": "someuser",
							"passcode": "12345678"
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
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}

	fakeServer.Mux.HandleFunc("/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Subject-Token", "aaa111")

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"token": {
				"expires_at": "2014-10-02T13:45:00.000000Z"
			}
		}`)
	})

	opts := auth.V3PasswordOpts{
		UserID:   "me",
		Password: "shhh",
	}

	token, err := tokens.Create(context.TODO(), &client, opts).Extract()
	if err != nil {
		t.Fatalf("Create returned an error: %v", err)
	}

	if token.ID != "aaa111" {
		t.Errorf("Expected token to be aaa111, but was %s", token.ID)
	}
}

func TestCreateFailureEmptyAuth(t *testing.T) {
	authTokenPostErr(t, auth.V3PasswordOpts{}, gophercloud.ErrMissingPassword{})
}

func TestCreateFailureMissingUser(t *testing.T) {
	opts := auth.V3PasswordOpts{
		Password: "supersecure",
	}

	authTokenPostErr(t, opts, gophercloud.ErrUsernameOrUserID{})
}

func TestCreateFailureBothUser(t *testing.T) {
	opts := auth.V3PasswordOpts{
		Password: "supersecure",
		Username: "oops",
		UserID:   "redundancy",
	}
	authTokenPostErr(t, opts, gophercloud.ErrUsernameOrUserID{})
}

func TestCreateFailureMissingDomain(t *testing.T) {
	opts := auth.V3PasswordOpts{
		Password: "supersecure",
		Username: "notuniqueenough",
	}
	authTokenPostErr(t, opts, gophercloud.ErrDomainIDOrDomainName{})
}

func TestCreateFailureBothDomain(t *testing.T) {
	opts := auth.V3PasswordOpts{
		Password:       "supersecure",
		Username:       "someone",
		UserDomainID:   "hurf",
		UserDomainName: "durf",
	}

	authTokenPostErr(t, opts, gophercloud.ErrDomainIDOrDomainName{})
}

func TestCreateFailureUserIDDomainID(t *testing.T) {
	opts := auth.V3PasswordOpts{
		UserID:       "100",
		Password:     "stuff",
		UserDomainID: "oops",
	}

	authTokenPostErr(t, opts, gophercloud.ErrDomainIDWithUserID{})
}

func TestCreateFailureUserIDDomainName(t *testing.T) {
	opts := auth.V3PasswordOpts{
		UserID:         "100",
		Password:       "sssh",
		UserDomainName: "oops",
	}

	authTokenPostErr(t, opts, gophercloud.ErrDomainNameWithUserID{})
}

func TestCreateFailureScopeProjectNameAlone(t *testing.T) {
	opts := auth.V3PasswordOpts{
		UserID:   "myself",
		Password: "swordfish",
		Scope: &auth.Scope{
			ProjectName: "notenough",
		},
	}

	authTokenPostErr(t, opts, gophercloud.ErrScopeProjectDomainIDOrProjectDomainName{})
}

func TestCreateFailureScopeProjectNameAndDomainID(t *testing.T) {
	opts := auth.V3PasswordOpts{
		UserID:   "myself",
		Password: "swordfish",
		Scope: &auth.Scope{
			ProjectName: "whoops",
			DomainID:    "1234",
		},
	}

	authTokenPostErr(t, opts, gophercloud.ErrProjectScopeAndDomainScope{
		DomainOption:  "DomainID",
		ProjectOption: "ProjectName",
	})
}

func TestCreateFailureScopeProjectIDAndDomainID(t *testing.T) {
	opts := auth.V3PasswordOpts{
		UserID:   "myself",
		Password: "swordfish",
		Scope: &auth.Scope{
			ProjectID: "toomuch",
			DomainID:  "notneeded",
		},
	}

	authTokenPostErr(t, opts, gophercloud.ErrProjectScopeAndDomainScope{
		DomainOption:  "DomainID",
		ProjectOption: "ProjectID",
	})
}

func TestCreateFailureScopeProjectIDAndDomainName(t *testing.T) {
	opts := auth.V3PasswordOpts{
		UserID:   "myself",
		Password: "swordfish",
		Scope: &auth.Scope{
			ProjectID:  "toomuch",
			DomainName: "notneeded",
		},
	}

	authTokenPostErr(t, opts, gophercloud.ErrProjectScopeAndDomainScope{
		DomainOption:  "DomainName",
		ProjectOption: "ProjectID",
	})
}

func TestGetRequest(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{
			TokenID: "12345abcdef",
		},
		Endpoint: fakeServer.Endpoint(),
	}

	fakeServer.Mux.HandleFunc("/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeaderUnset(t, r, "Content-Type")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", "12345abcdef")
		th.TestHeader(t, r, "X-Subject-Token", "abcdef12345")

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `
			{ "token": { "expires_at": "2014-08-29T13:10:01.000000Z" } }
		`)
	})

	token, err := tokens.Get(context.TODO(), &client, "abcdef12345", nil).Extract()
	th.AssertNoErr(t, err)

	expected, _ := time.Parse(time.UnixDate, "Fri Aug 29 13:10:01 UTC 2014")
	if token.ExpiresAt != expected {
		t.Errorf("Expected expiration time %s, but was %s", expected.Format(time.UnixDate), time.Time(token.ExpiresAt).Format(time.UnixDate))
	}
}

func prepareAuthTokenHandler(t *testing.T, fakeServer th.FakeServer, expectedMethod string, status int) gophercloud.ServiceClient {
	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{
			TokenID: "12345abcdef",
		},
		Endpoint: fakeServer.Endpoint(),
	}

	fakeServer.Mux.HandleFunc("/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, expectedMethod)
		th.TestHeaderUnset(t, r, "Content-Type")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", "12345abcdef")
		th.TestHeader(t, r, "X-Subject-Token", "abcdef12345")

		w.WriteHeader(status)
	})

	return client
}

func TestValidateRequestSuccessful(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	client := prepareAuthTokenHandler(t, fakeServer, "HEAD", http.StatusNoContent)

	ok, err := tokens.Validate(context.TODO(), &client, "abcdef12345", nil)
	if err != nil {
		t.Errorf("Unexpected error from Validate: %v", err)
	}

	if !ok {
		t.Errorf("Validate returned false for a valid token")
	}
}

func TestValidateRequestFailure(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	client := prepareAuthTokenHandler(t, fakeServer, "HEAD", http.StatusNotFound)

	ok, err := tokens.Validate(context.TODO(), &client, "abcdef12345", nil)
	if err != nil {
		t.Errorf("Unexpected error from Validate: %v", err)
	}

	if ok {
		t.Errorf("Validate returned true for an invalid token")
	}
}

func TestValidateRequestError(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	client := prepareAuthTokenHandler(t, fakeServer, "HEAD", http.StatusMethodNotAllowed)

	_, err := tokens.Validate(context.TODO(), &client, "abcdef12345", nil)
	if err == nil {
		t.Errorf("Missing expected error from Validate")
	}
}

func TestRevokeRequestSuccessful(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	client := prepareAuthTokenHandler(t, fakeServer, "DELETE", http.StatusNoContent)

	res := tokens.Revoke(context.TODO(), &client, "abcdef12345")
	th.AssertNoErr(t, res.Err)
}

func TestRevokeRequestError(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	client := prepareAuthTokenHandler(t, fakeServer, "DELETE", http.StatusNotFound)

	res := tokens.Revoke(context.TODO(), &client, "abcdef12345")
	if res.Err == nil {
		t.Errorf("Missing expected error from Revoke")
	}
}

func TestGetRequestWithAccessRules(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{
			TokenID: "12345abcdef",
		},
		Endpoint: fakeServer.Endpoint(),
	}

	fakeServer.Mux.HandleFunc("/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", "12345abcdef")
		th.TestHeader(t, r, "X-Subject-Token", "abcdef12345")
		th.TestHeader(t, r, "OpenStack-Identity-Access-Rules", "1")

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `
			{ "token": { "expires_at": "2014-08-29T13:10:01.000000Z" } }
		`)
	})

	getOpts := tokens.GetOpts{
		AccessRulesVersion: "1",
	}
	token, err := tokens.Get(context.TODO(), &client, "abcdef12345", getOpts).Extract()
	th.AssertNoErr(t, err)

	expected, _ := time.Parse(time.UnixDate, "Fri Aug 29 13:10:01 UTC 2014")
	if token.ExpiresAt != expected {
		t.Errorf("Expected expiration time %s, but was %s", expected.Format(time.UnixDate), time.Time(token.ExpiresAt).Format(time.UnixDate))
	}
}

func TestValidateRequestWithAccessRules(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{
			TokenID: "12345abcdef",
		},
		Endpoint: fakeServer.Endpoint(),
	}

	fakeServer.Mux.HandleFunc("/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "HEAD")
		th.TestHeader(t, r, "X-Auth-Token", "12345abcdef")
		th.TestHeader(t, r, "X-Subject-Token", "abcdef12345")
		th.TestHeader(t, r, "OpenStack-Identity-Access-Rules", "1")

		w.WriteHeader(http.StatusNoContent)
	})

	validateOpts := tokens.ValidateOpts{
		AccessRulesVersion: "1",
	}
	ok, err := tokens.Validate(context.TODO(), &client, "abcdef12345", validateOpts)
	if err != nil {
		t.Errorf("Unexpected error from Validate: %v", err)
	}

	if !ok {
		t.Errorf("Validate returned false for a valid token")
	}
}

func TestNoTokenInResponse(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}

	fakeServer.Mux.HandleFunc("/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{}`)
	})

	opts := auth.V3PasswordOpts{
		UserID:   "me",
		Password: "squirrel!",
	}

	_, err := tokens.Create(context.TODO(), &client, opts).Extract()
	th.AssertNoErr(t, err)
}
