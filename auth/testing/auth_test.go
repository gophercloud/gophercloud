package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

const ID = "abc123def"

func TestAuthOptionsV2AuthenticatePassword(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestJSONRequest(t, r, `
			{
				"auth": {
					"tenantId": "tenant-id",
					"passwordCredentials": {
						"username": "testuser",
						"password": "testpass"
					}
				}
			}
		`)

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `
			{
				"access": {
					"token": {
						"id": "v2-token-id",
						"expires": "2026-09-01T12:00:00Z",
						"tenant": {"id": "tenant-id", "name": "tenant-name"}
					},
					"user": {
						"id": "user-id",
						"name": "testuser",
						"roles": [{"name": "member"}]
					},
					"serviceCatalog": [
						{
							"name": "nova",
							"type": "compute",
							"endpoints": [
								{
									"region": "RegionOne",
									"publicURL": "http://public.example.com/compute",
									"internalURL": "http://internal.example.com/compute",
									"adminURL": "http://admin.example.com/compute"
								}
							]
						}
					]
				}
			}
		`)
	})

	opts := auth.AuthOptionsV2{
		AuthURL: fakeServer.Endpoint(),
		Auth: auth.V2PasswordOpts{
			Username: "testuser",
			Password: "testpass",
			TenantID: "tenant-id",
		},
	}

	result, err := opts.Authenticate(context.TODO(), nil)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "v2-token-id", result.TokenID)
	th.AssertEquals(t, "user-id", result.User.ID)
	th.AssertEquals(t, "tenant-id", result.Project.ID)
	th.AssertEquals(t, "tenant-name", result.Project.Name)
	th.AssertEquals(t, 1, len(result.Roles))
	th.AssertEquals(t, "member", result.Roles[0].Name)
	th.AssertEquals(t, 1, len(result.Catalog.Entries))
	th.AssertEquals(t, 3, len(result.Catalog.Entries[0].Endpoints))

	var interfaces []string
	for _, ep := range result.Catalog.Entries[0].Endpoints {
		interfaces = append(interfaces, ep.Interface)
	}
	th.AssertDeepEquals(t, []string{"public", "internal", "admin"}, interfaces)

	expectedExpires := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	th.AssertEquals(t, true, result.ExpiresAt.Equal(expectedExpires))
}

func TestAuthOptionsV2AuthenticateToken(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestJSONRequest(t, r, `
			{
				"auth": {
					"tenantId": "tenant-id",
					"tokenCredentials": {"id": "existing-token"}
				}
			}
		`)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `
			{
				"access": {
					"token": {"id": "reissued-token", "expires": "2026-09-01T12:00:00Z", "tenant": {"id": "tenant-id", "name": "tenant-name"}},
					"user": {"id": "user-id", "name": "testuser", "roles": []},
					"serviceCatalog": []
				}
			}
		`)
	})

	opts := auth.AuthOptionsV2{
		AuthURL: fakeServer.Endpoint(),
		Auth:    auth.V2TokenOpts{Token: "existing-token", TenantID: "tenant-id", AllowReauth: true},
	}

	result, err := opts.Authenticate(context.TODO(), nil)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "reissued-token", result.TokenID)
	th.AssertEquals(t, true, result.CanReauth)
}

func TestAuthOptionsV2AuthenticateRescope(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestJSONRequest(t, r, `
			{
				"auth": {
					"tenantId": "new-tenant-id",
					"tokenCredentials": {"id": "existing-token"}
				}
			}
		`)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `
			{
				"access": {
					"token": {"id": "rescoped-token", "expires": "2026-09-01T12:00:00Z", "tenant": {"id": "new-tenant-id", "name": "new-tenant-name"}},
					"user": {"id": "user-id", "name": "testuser", "roles": []},
					"serviceCatalog": []
				}
			}
		`)
	})

	opts := auth.AuthOptionsV2{
		AuthURL: fakeServer.Endpoint(),
		Auth:    auth.V2RescopeTokenOpts{Token: "existing-token", TenantID: "new-tenant-id"},
	}

	result, err := opts.Authenticate(context.TODO(), nil)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "rescoped-token", result.TokenID)
	th.AssertEquals(t, "new-tenant-id", result.Project.ID)
}

func TestAuthOptionsV2AuthenticateNilAuth(t *testing.T) {
	opts := auth.AuthOptionsV2{AuthURL: "http://example.com:5000/v2.0"}

	_, err := opts.Authenticate(context.TODO(), nil)
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrMissingInput)
	th.AssertEquals(t, true, ok)
}

func TestAuthOptionsV2AuthenticateHTTPFailure(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/tokens", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"error": {"message": "bad credentials"}}`)
	})

	opts := auth.AuthOptionsV2{
		AuthURL: fakeServer.Endpoint(),
		Auth:    auth.V2PasswordOpts{Username: "testuser", Password: "wrongpass"},
	}

	_, err := opts.Authenticate(context.TODO(), nil)
	th.AssertErr(t, err)
}

func TestAuthOptionsV3AuthenticatePassword(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestJSONRequest(t, r, `
			{
				"auth": {
					"identity": {
						"methods": ["password"],
						"password": {
							"user": {
								"name": "testuser",
								"domain": {"id": "default"},
								"password": "testpass"
							}
						}
					},
					"scope": {
						"project": {"id": "project-id"}
					}
				}
			}
		`)

		w.Header().Set("X-Subject-Token", "the-token-id")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `
			{
				"token": {
					"methods": ["password"],
					"expires_at": "2026-09-01T12:00:00.000000Z",
					"issued_at": "2026-09-01T11:00:00.000000Z",
					"audit_ids": ["abc123"],
					"user": {
						"id": "user-id",
						"name": "testuser",
						"domain": {"id": "default", "name": "Default"}
					},
					"project": {
						"id": "project-id",
						"name": "testproject",
						"domain": {"id": "default", "name": "Default"}
					},
					"roles": [{"id": "role-id", "name": "member"}],
					"catalog": [
						{
							"id": "svc-id",
							"name": "nova",
							"type": "compute",
							"endpoints": [
								{"id": "ep-id", "region": "RegionOne", "region_id": "RegionOne", "interface": "public", "url": "http://example.com:8774/v2.1"}
							]
						}
					]
				}
			}
		`)
	})

	opts := auth.AuthOptionsV3{
		AuthURL: fakeServer.Endpoint(),
		Auth: auth.V3PasswordOpts{
			Username:     "testuser",
			Password:     "testpass",
			UserDomainID: "default",
			Scope:        &auth.Scope{ProjectID: "project-id"},
			AllowReauth:  true,
		},
	}

	result, err := opts.Authenticate(context.TODO(), nil)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "the-token-id", result.TokenID)
	th.AssertEquals(t, "user-id", result.User.ID)
	th.AssertEquals(t, "testproject", result.Project.Name)
	th.AssertEquals(t, "default", result.Project.Domain.ID)
	th.AssertEquals(t, 1, len(result.Roles))
	th.AssertEquals(t, "member", result.Roles[0].Name)
	th.AssertEquals(t, 1, len(result.Catalog.Entries))
	th.AssertEquals(t, "compute", result.Catalog.Entries[0].Type)
	th.AssertEquals(t, "http://example.com:8774/v2.1", result.Catalog.Entries[0].Endpoints[0].URL)
	th.AssertDeepEquals(t, []auth.AuthType{"password"}, result.Methods)
	th.AssertEquals(t, "abc123", result.AuditIDs[0])
	th.AssertEquals(t, false, result.System)
	th.AssertEquals(t, true, result.Domain == nil)
	th.AssertEquals(t, true, result.CanReauth)

	expectedExpires := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	th.AssertEquals(t, true, result.ExpiresAt.Equal(expectedExpires))
}

func TestAuthOptionsV3AuthenticateNilAuth(t *testing.T) {
	opts := auth.AuthOptionsV3{AuthURL: "http://example.com:5000/v3"}

	_, err := opts.Authenticate(context.TODO(), nil)
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrMissingInput)
	th.AssertEquals(t, true, ok)
}

func TestAuthOptionsV3AuthenticateToAuthBodyError(t *testing.T) {
	opts := auth.AuthOptionsV3{
		AuthURL: "http://example.com:5000/v3",
		Auth:    auth.V3PasswordOpts{}, // no password set
	}

	_, err := opts.Authenticate(context.TODO(), nil)
	th.AssertErr(t, err)
}

func TestAuthOptionsV3AuthenticateHTTPFailure(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"error": {"message": "bad credentials"}}`)
	})

	opts := auth.AuthOptionsV3{
		AuthURL: fakeServer.Endpoint(),
		Auth: auth.V3PasswordOpts{
			Username: "testuser",
			Password: "wrongpass",
		},
	}

	_, err := opts.Authenticate(context.TODO(), nil)
	th.AssertErr(t, err)
}

func TestAuthOptionsV3AuthenticateTOTP(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestJSONRequest(t, r, `
			{
				"auth": {
					"identity": {
						"methods": ["totp"],
						"totp": {
							"user": {
								"name": "testuser",
								"domain": {"id": "default"},
								"passcode": "123456"
							}
						}
					}
				}
			}
		`)
		w.Header().Set("X-Subject-Token", "totp-token")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"token": {"methods": ["totp"], "user": {"id": "u", "name": "testuser"}}}`)
	})

	opts := auth.AuthOptionsV3{
		AuthURL: fakeServer.Endpoint(),
		Auth: auth.V3TOTPOpts{
			Username:     "testuser",
			Passcode:     "123456",
			UserDomainID: "default",
		},
	}

	result, err := opts.Authenticate(context.TODO(), nil)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "totp-token", result.TokenID)
}

func TestAuthOptionsV3AuthenticateApplicationCredential(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestJSONRequest(t, r, `
			{
				"auth": {
					"identity": {
						"methods": ["application_credential"],
						"application_credential": {
							"id": "appcred-id",
							"secret": "appcred-secret",
							"user": {}
						}
					}
				}
			}
		`)
		w.Header().Set("X-Subject-Token", "appcred-token")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `
			{
				"token": {
					"methods": ["application_credential"],
					"user": {"id": "u", "name": "testuser"},
					"project": {"id": "project-id", "name": "testproject"},
					"application_credential": {
						"id": "appcred-id",
						"name": "my-appcred",
						"restricted": true,
						"access_rules": [{"id": "rule-1", "path": "/v2.1/servers", "method": "GET", "service": "compute"}]
					}
				}
			}
		`)
	})

	opts := auth.AuthOptionsV3{
		AuthURL: fakeServer.Endpoint(),
		Auth: auth.V3ApplicationCredentialOpts{
			ApplicationCredentialID:     "appcred-id",
			ApplicationCredentialSecret: "appcred-secret",
		},
	}

	result, err := opts.Authenticate(context.TODO(), nil)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "appcred-token", result.TokenID)
	th.AssertEquals(t, "appcred-id", result.ApplicationCredential.ID)
	th.AssertEquals(t, true, result.ApplicationCredential.Restricted)
	th.AssertEquals(t, 1, len(result.ApplicationCredential.AccessRules))
	th.AssertEquals(t, "compute", result.ApplicationCredential.AccessRules[0].Service)
}

func TestAuthOptionsV3AuthenticateToken(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestJSONRequest(t, r, `
			{
				"auth": {
					"identity": {
						"methods": ["token"],
						"token": {"id": "existing-token"}
					}
				}
			}
		`)
		w.Header().Set("X-Subject-Token", "reissued-token")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"token": {"methods": ["token"], "user": {"id": "u", "name": "testuser"}}}`)
	})

	opts := auth.AuthOptionsV3{
		AuthURL: fakeServer.Endpoint(),
		Auth:    auth.V3TokenOpts{Token: "existing-token"},
	}

	result, err := opts.Authenticate(context.TODO(), nil)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "reissued-token", result.TokenID)
	th.AssertEquals(t, false, result.CanReauth)
}

func TestAuthOptionsV3AuthenticateMultifactor(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestJSONRequest(t, r, `
			{
				"auth": {
					"identity": {
						"methods": ["password", "totp"],
						"password": {
							"user": {
								"name": "testuser",
								"domain": {"id": "default"},
								"password": "testpass"
							}
						},
						"totp": {
							"user": {
								"name": "testuser",
								"domain": {"id": "default"},
								"passcode": "123456"
							}
						}
					}
				}
			}
		`)
		w.Header().Set("X-Subject-Token", "mfa-token")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"token": {"methods": ["password", "totp"], "user": {"id": "u", "name": "testuser"}}}`)
	})

	opts := auth.AuthOptionsV3{
		AuthURL: fakeServer.Endpoint(),
		Auth: auth.V3MultifactorOpts{
			AuthMethods: []auth.AuthOptionsBuilderV3{
				auth.V3PasswordOpts{Username: "testuser", Password: "testpass", UserDomainID: "default"},
				auth.V3TOTPOpts{Username: "testuser", Passcode: "123456", UserDomainID: "default"},
			},
		},
	}

	result, err := opts.Authenticate(context.TODO(), nil)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "mfa-token", result.TokenID)
	th.AssertDeepEquals(t, []auth.AuthType{"password", "totp"}, result.Methods)
}

func TestAuthOptionsV3AuthenticateRescope(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestJSONRequest(t, r, `
			{
				"auth": {
					"identity": {
						"methods": ["token"],
						"token": {"id": "existing-token"}
					},
					"scope": {
						"project": {"id": "new-project-id"}
					}
				}
			}
		`)
		w.Header().Set("X-Subject-Token", "rescoped-token")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `
			{
				"token": {
					"methods": ["token"],
					"user": {"id": "u", "name": "testuser"},
					"project": {"id": "new-project-id", "name": "newproject"}
				}
			}
		`)
	})

	opts := auth.AuthOptionsV3{
		AuthURL: fakeServer.Endpoint(),
		Auth: auth.V3RescopeTokenOpts{
			Token: "existing-token",
			Scope: &auth.Scope{ProjectID: "new-project-id"},
		},
	}

	result, err := opts.Authenticate(context.TODO(), nil)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "rescoped-token", result.TokenID)
	th.AssertEquals(t, "new-project-id", result.Project.ID)
}

func TestAuthOptionsV3AuthenticateRescopeMissingScope(t *testing.T) {
	opts := auth.AuthOptionsV3{
		AuthURL: "http://example.com:5000/v3",
		Auth:    auth.V3RescopeTokenOpts{Token: "existing-token"},
	}

	_, err := opts.Authenticate(context.TODO(), nil)
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrScopeEmpty)
	th.AssertEquals(t, true, ok)
}

func TestAuthOptionsV3AuthenticateSystemScoped(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestJSONRequest(t, r, `
			{
				"auth": {
					"identity": {
						"methods": ["password"],
						"password": {
							"user": {
								"name": "testuser",
								"domain": {"id": "default"},
								"password": "testpass"
							}
						}
					},
					"scope": {
						"system": {"all": true}
					}
				}
			}
		`)
		w.Header().Set("X-Subject-Token", "system-token")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `
			{
				"token": {
					"methods": ["password"],
					"user": {"id": "u", "name": "testuser"},
					"system": {"all": true}
				}
			}
		`)
	})

	opts := auth.AuthOptionsV3{
		AuthURL: fakeServer.Endpoint(),
		Auth: auth.V3PasswordOpts{
			Username:     "testuser",
			Password:     "testpass",
			UserDomainID: "default",
			Scope:        &auth.Scope{System: true},
		},
	}

	result, err := opts.Authenticate(context.TODO(), nil)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, true, result.System)
	th.AssertEquals(t, true, result.Project == nil)
}

func TestAuthOptionsV3AuthenticateDomainScoped(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestJSONRequest(t, r, `
			{
				"auth": {
					"identity": {
						"methods": ["password"],
						"password": {
							"user": {
								"name": "testuser",
								"domain": {"id": "default"},
								"password": "testpass"
							}
						}
					},
					"scope": {
						"domain": {"id": "domain-id"}
					}
				}
			}
		`)
		w.Header().Set("X-Subject-Token", "domain-token")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `
			{
				"token": {
					"methods": ["password"],
					"user": {"id": "u", "name": "testuser"},
					"domain": {"id": "domain-id", "name": "domain-name"}
				}
			}
		`)
	})

	opts := auth.AuthOptionsV3{
		AuthURL: fakeServer.Endpoint(),
		Auth: auth.V3PasswordOpts{
			Username:     "testuser",
			Password:     "testpass",
			UserDomainID: "default",
			Scope:        &auth.Scope{DomainID: "domain-id"},
		},
	}

	result, err := opts.Authenticate(context.TODO(), nil)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "domain-id", result.Domain.ID)
	th.AssertEquals(t, true, result.Project == nil)
	th.AssertEquals(t, false, result.System)
}

func TestAuthResultExtractTokenID(t *testing.T) {
	result := auth.AuthResult{TokenID: "abc123"}
	id, err := result.ExtractTokenID()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "abc123", id)
}

func TestAuthOptionsV3AuthenticateUsesProvidedHTTPClient(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Subject-Token", ID)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{ "token": { "expires_at": "2013-02-02T18:30:59.000000Z" } }`)
	})

	var usedCustomTransport bool
	httpClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			usedCustomTransport = true
			return http.DefaultTransport.RoundTrip(req)
		}),
	}

	opts := auth.AuthOptionsV3{
		AuthURL: fakeServer.Endpoint(),
		Auth:    auth.V3PasswordOpts{Username: "me", Password: "secret", UserDomainName: "default"},
	}
	result, err := opts.Authenticate(context.TODO(), httpClient)
	th.AssertNoErr(t, err)
	th.CheckEquals(t, ID, result.TokenID)
	th.CheckEquals(t, true, usedCustomTransport)
}

func TestAuthOptionsV3AuthenticateNilHTTPClientDefaults(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Subject-Token", ID)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{ "token": { "expires_at": "2013-02-02T18:30:59.000000Z" } }`)
	})

	opts := auth.AuthOptionsV3{
		AuthURL: fakeServer.Endpoint(),
		Auth:    auth.V3PasswordOpts{Username: "me", Password: "secret", UserDomainName: "default"},
	}
	result, err := opts.Authenticate(context.TODO(), nil)
	th.AssertNoErr(t, err)
	th.CheckEquals(t, ID, result.TokenID)
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
