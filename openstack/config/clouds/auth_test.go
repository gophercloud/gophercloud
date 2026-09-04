package clouds_test

import (
	"fmt"
	"maps"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	"github.com/gophercloud/gophercloud/v2/openstack/config/clouds"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestAuthOptionsFromCloudV3PasswordInferred(t *testing.T) {
	cloud := clouds.Cloud{
		Auth: map[string]any{
			"auth_url": "http://example.com:5000",
			"username": "testuser",
			"password": "testpass",
		},
	}

	opts, err := cloud.AuthOptions()
	th.AssertNoErr(t, err)

	v3Opts, ok := opts.(*auth.AuthOptionsV3)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000", v3Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V3PasswordOpts{Username: "testuser", Password: "testpass", Scope: &auth.Scope{}}, v3Opts.Auth)
}

func TestAuthOptionsFromCloudAllowReauth(t *testing.T) {
	cloud := clouds.Cloud{
		Auth: map[string]any{
			"auth_url":     "http://example.com:5000",
			"username":     "testuser",
			"password":     "testpass",
			"allow_reauth": true,
		},
	}

	opts, err := cloud.AuthOptions()
	th.AssertNoErr(t, err)

	v3Opts, ok := opts.(*auth.AuthOptionsV3)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000", v3Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V3PasswordOpts{Username: "testuser", Password: "testpass", Scope: &auth.Scope{}, AllowReauth: true}, v3Opts.Auth)
	th.AssertEquals(t, true, v3Opts.Auth.CanReauth())
}

func TestAuthOptionsFromCloudV2ViaIdentityAPIVersion(t *testing.T) {
	cloud := clouds.Cloud{
		IdentityAPIVersion: "2.0",
		Auth: map[string]any{
			"auth_url": "http://example.com:5000",
			"username": "testuser",
			"password": "testpass",
		},
	}

	opts, err := cloud.AuthOptions()
	th.AssertNoErr(t, err)

	v2Opts, ok := opts.(*auth.AuthOptionsV2)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000", v2Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V2PasswordOpts{Username: "testuser", Password: "testpass"}, v2Opts.Auth)
}

func TestAuthOptionsFromCloudMissingAuthURL(t *testing.T) {
	cloud := clouds.Cloud{
		Auth: map[string]any{
			"username": "testuser",
			"password": "testpass",
		},
	}

	_, err := cloud.AuthOptions()
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrMissingInput)
	th.AssertEquals(t, true, ok)
}

func TestAuthOptionsFromCloudNilAuth(t *testing.T) {
	cloud := clouds.Cloud{}

	_, err := cloud.AuthOptions()
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrMissingInput)
	th.AssertEquals(t, true, ok)
}

func TestAuthOptionsFromCloudNoUsableCredentials(t *testing.T) {
	cloud := clouds.Cloud{
		Auth: map[string]any{
			"auth_url": "http://example.com:5000",
		},
	}

	_, err := cloud.AuthOptions()
	th.AssertErr(t, err)

	missingInput, ok := err.(gophercloud.ErrMissingInput)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "Auth", missingInput.Argument)
}

func TestAuthOptionsFromCloudExplicitAuthTypes(t *testing.T) {
	tests := []struct {
		name     string
		cloud    clouds.Cloud
		wantV2   bool
		wantAuth any
	}{
		{
			name: "v2password",
			cloud: clouds.Cloud{
				AuthType: auth.AuthV2Password,
				Auth: map[string]any{
					"auth_url": "http://example.com:5000",
					"username": "testuser",
					"password": "testpass",
				},
			},
			wantV2:   true,
			wantAuth: auth.V2PasswordOpts{Username: "testuser", Password: "testpass"},
		},
		{
			name: "v2token",
			cloud: clouds.Cloud{
				AuthType: auth.AuthV2Token,
				Auth: map[string]any{
					"auth_url": "http://example.com:5000",
					"token":    "testtoken",
				},
			},
			wantV2:   true,
			wantAuth: auth.V2TokenOpts{Token: "testtoken"},
		},
		{
			name: "v3password",
			cloud: clouds.Cloud{
				AuthType: auth.AuthV3Password,
				Auth: map[string]any{
					"auth_url": "http://example.com:5000",
					"username": "testuser",
					"password": "testpass",
				},
			},
			wantAuth: auth.V3PasswordOpts{Username: "testuser", Password: "testpass", Scope: &auth.Scope{}},
		},
		{
			name: "v3token",
			cloud: clouds.Cloud{
				AuthType: auth.AuthV3Token,
				Auth: map[string]any{
					"auth_url": "http://example.com:5000",
					"token":    "testtoken",
				},
			},
			wantAuth: auth.V3TokenOpts{Token: "testtoken", Scope: &auth.Scope{}},
		},
		{
			name: "v3applicationcredential",
			cloud: clouds.Cloud{
				AuthType: auth.AuthV3ApplicationCredential,
				Auth: map[string]any{
					"auth_url":                      "http://example.com:5000",
					"application_credential_id":     "app-cred-id",
					"application_credential_secret": "app-cred-secret",
				},
			},
			wantAuth: auth.V3ApplicationCredentialOpts{ApplicationCredentialID: "app-cred-id", ApplicationCredentialSecret: "app-cred-secret"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := tt.cloud.AuthOptions()
			th.AssertNoErr(t, err)

			if tt.wantV2 {
				v2Opts, ok := opts.(*auth.AuthOptionsV2)
				th.AssertEquals(t, true, ok)
				th.AssertDeepEquals(t, tt.wantAuth, v2Opts.Auth)
			} else {
				v3Opts, ok := opts.(*auth.AuthOptionsV3)
				th.AssertEquals(t, true, ok)
				th.AssertDeepEquals(t, tt.wantAuth, v3Opts.Auth)
			}
		})
	}
}

func TestAuthOptionsFromCloudV3TOTPExplicit(t *testing.T) {
	cloud := clouds.Cloud{
		AuthType: "v3totp",
		Auth: map[string]any{
			"auth_url": "http://example.com:5000",
			"username": "testuser",
		},
	}

	opts, err := cloud.AuthOptions(auth.WithPasscode("123456"))
	th.AssertNoErr(t, err)

	v3Opts, ok := opts.(*auth.AuthOptionsV3)
	th.AssertEquals(t, true, ok)
	th.AssertDeepEquals(t, auth.V3TOTPOpts{Username: "testuser", Passcode: "123456", Scope: &auth.Scope{}}, v3Opts.Auth)
}

func TestAuthOptionsFromCloudMechanismInferencePrecedence(t *testing.T) {
	base := map[string]any{
		"auth_url":                  "http://example.com:5000",
		"username":                  "testuser",
		"password":                  "testpass",
		"token":                     "testtoken",
		"application_credential_id": "app-cred-id",
	}

	t.Run("password beats passcode, token, and appcred", func(t *testing.T) {
		cloud := clouds.Cloud{Auth: maps.Clone(base)}
		opts, err := cloud.AuthOptions(auth.WithPasscode("123456"))
		th.AssertNoErr(t, err)
		v3Opts := opts.(*auth.AuthOptionsV3)
		_, ok := v3Opts.Auth.(auth.V3PasswordOpts)
		th.AssertEquals(t, true, ok)
	})

	t.Run("passcode beats token and appcred", func(t *testing.T) {
		m := maps.Clone(base)
		delete(m, "password")
		cloud := clouds.Cloud{Auth: m}
		opts, err := cloud.AuthOptions(auth.WithPasscode("123456"))
		th.AssertNoErr(t, err)
		v3Opts := opts.(*auth.AuthOptionsV3)
		_, ok := v3Opts.Auth.(auth.V3TOTPOpts)
		th.AssertEquals(t, true, ok)
	})

	t.Run("token beats appcred", func(t *testing.T) {
		m := maps.Clone(base)
		delete(m, "password")
		cloud := clouds.Cloud{Auth: m}
		opts, err := cloud.AuthOptions()
		th.AssertNoErr(t, err)
		v3Opts := opts.(*auth.AuthOptionsV3)
		_, ok := v3Opts.Auth.(auth.V3TokenOpts)
		th.AssertEquals(t, true, ok)
	})

	t.Run("appcred alone resolves", func(t *testing.T) {
		cloud := clouds.Cloud{Auth: map[string]any{
			"auth_url":                  "http://example.com:5000",
			"application_credential_id": "app-cred-id",
		}}
		opts, err := cloud.AuthOptions()
		th.AssertNoErr(t, err)
		v3Opts := opts.(*auth.AuthOptionsV3)
		_, ok := v3Opts.Auth.(auth.V3ApplicationCredentialOpts)
		th.AssertEquals(t, true, ok)
	})
}

func TestAuthOptionsFromCloudUnsupportedAuthType(t *testing.T) {
	cloud := clouds.Cloud{
		AuthType: "v3federated",
		Auth: map[string]any{
			"auth_url": "http://example.com:5000",
		},
	}

	_, err := cloud.AuthOptions()
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrUnsupportedAuthType)
	th.AssertEquals(t, true, ok)
}

func TestAuthOptionsFromCloudScopeResolution(t *testing.T) {
	t.Run("keeps user and project domain distinct", func(t *testing.T) {
		cloud := clouds.Cloud{
			Auth: map[string]any{
				"auth_url":            "http://example.com:5000",
				"username":            "testuser",
				"password":            "testpass",
				"user_domain_name":    "userdomain",
				"project_domain_name": "projectdomain",
				"project_name":        "testproject",
			},
		}

		opts, err := cloud.AuthOptions()
		th.AssertNoErr(t, err)

		v3Opts := opts.(*auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3PasswordOpts{
			Username:       "testuser",
			Password:       "testpass",
			UserDomainName: "userdomain",
			Scope: &auth.Scope{
				ProjectDomainName: "projectdomain",
				ProjectName:       "testproject",
			},
		}, v3Opts.Auth)
	})

	t.Run("falls back to DefaultDomain when no domain is set", func(t *testing.T) {
		cloud := clouds.Cloud{
			Auth: map[string]any{
				"auth_url":       "http://example.com:5000",
				"username":       "testuser",
				"password":       "testpass",
				"default_domain": "default",
				"project_name":   "testproject",
			},
		}

		opts, err := cloud.AuthOptions()
		th.AssertNoErr(t, err)

		v3Opts := opts.(*auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3PasswordOpts{
			Username:     "testuser",
			Password:     "testpass",
			UserDomainID: "default",
			Scope: &auth.Scope{
				ProjectDomainID: "default",
				ProjectName:     "testproject",
			},
		}, v3Opts.Auth)
	})

	t.Run("generic DomainName seeds both user and project domain", func(t *testing.T) {
		cloud := clouds.Cloud{
			Auth: map[string]any{
				"auth_url":     "http://example.com:5000",
				"username":     "testuser",
				"password":     "testpass",
				"domain_name":  "shared",
				"project_name": "testproject",
			},
		}

		opts, err := cloud.AuthOptions()
		th.AssertNoErr(t, err)

		v3Opts := opts.(*auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3PasswordOpts{
			Username:       "testuser",
			Password:       "testpass",
			UserDomainName: "shared",
			Scope: &auth.Scope{
				ProjectDomainName: "shared",
				ProjectName:       "testproject",
			},
		}, v3Opts.Auth)
	})

	t.Run("system scope", func(t *testing.T) {
		cloud := clouds.Cloud{
			Auth: map[string]any{
				"auth_url":     "http://example.com:5000",
				"username":     "testuser",
				"password":     "testpass",
				"system_scope": "all",
			},
		}

		opts, err := cloud.AuthOptions()
		th.AssertNoErr(t, err)

		v3Opts := opts.(*auth.AuthOptionsV3)
		passwordOpts := v3Opts.Auth.(auth.V3PasswordOpts)
		th.AssertEquals(t, true, passwordOpts.Scope.System)
	})

	t.Run("trust ID", func(t *testing.T) {
		cloud := clouds.Cloud{
			Auth: map[string]any{
				"auth_url": "http://example.com:5000",
				"username": "testuser",
				"password": "testpass",
				"trust_id": "trust-id",
			},
		}

		opts, err := cloud.AuthOptions()
		th.AssertNoErr(t, err)

		v3Opts := opts.(*auth.AuthOptionsV3)
		passwordOpts := v3Opts.Auth.(auth.V3PasswordOpts)
		th.AssertEquals(t, "trust-id", passwordOpts.Scope.TrustID)
	})

	t.Run("application credential auth builds no scope", func(t *testing.T) {
		cloud := clouds.Cloud{
			Auth: map[string]any{
				"auth_url":                      "http://example.com:5000",
				"application_credential_id":     "app-cred-id",
				"application_credential_secret": "app-cred-secret",
				"project_name":                  "testproject",
			},
		}

		opts, err := cloud.AuthOptions()
		th.AssertNoErr(t, err)

		v3Opts := opts.(*auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3ApplicationCredentialOpts{
			ApplicationCredentialID:     "app-cred-id",
			ApplicationCredentialSecret: "app-cred-secret",
		}, v3Opts.Auth)
	})

	t.Run("WithScope bypasses computed scope", func(t *testing.T) {
		cloud := clouds.Cloud{
			Auth: map[string]any{
				"auth_url":     "http://example.com:5000",
				"username":     "testuser",
				"password":     "testpass",
				"project_name": "testproject",
			},
		}

		explicitScope := &auth.Scope{DomainID: "explicit-domain"}
		opts, err := cloud.AuthOptions(auth.WithScope(explicitScope))
		th.AssertNoErr(t, err)

		v3Opts := opts.(*auth.AuthOptionsV3)
		passwordOpts := v3Opts.Auth.(auth.V3PasswordOpts)
		th.AssertDeepEquals(t, explicitScope, passwordOpts.Scope)
	})
}

func TestCloudOptionOverrides(t *testing.T) {
	baseCloud := func() clouds.Cloud {
		return clouds.Cloud{
			Auth: map[string]any{
				"auth_url": "http://example.com:5000",
			},
		}
	}

	t.Run("WithUsername and WithPassword", func(t *testing.T) {
		opts, err := baseCloud().AuthOptions(auth.WithUsername("override-user"), auth.WithPassword("override-pass"))
		th.AssertNoErr(t, err)
		v3Opts := opts.(*auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3PasswordOpts{Username: "override-user", Password: "override-pass", Scope: &auth.Scope{}}, v3Opts.Auth)
	})

	t.Run("WithUserID", func(t *testing.T) {
		opts, err := baseCloud().AuthOptions(auth.WithUserID("override-user-id"), auth.WithPassword("override-pass"))
		th.AssertNoErr(t, err)
		v3Opts := opts.(*auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3PasswordOpts{UserID: "override-user-id", Password: "override-pass", Scope: &auth.Scope{}}, v3Opts.Auth)
	})

	t.Run("WithToken", func(t *testing.T) {
		opts, err := baseCloud().AuthOptions(auth.WithToken("override-token"))
		th.AssertNoErr(t, err)
		v3Opts := opts.(*auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3TokenOpts{Token: "override-token", Scope: &auth.Scope{}}, v3Opts.Auth)
	})

	t.Run("WithDomainID", func(t *testing.T) {
		opts, err := baseCloud().AuthOptions(auth.WithUsername("u"), auth.WithPassword("p"), auth.WithDomainID("override-domain-id"))
		th.AssertNoErr(t, err)
		v3Opts := opts.(*auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3PasswordOpts{
			Username:     "u",
			Password:     "p",
			UserDomainID: "override-domain-id",
			Scope:        &auth.Scope{ProjectDomainID: "override-domain-id"},
		}, v3Opts.Auth)
	})

	t.Run("WithDomainName", func(t *testing.T) {
		opts, err := baseCloud().AuthOptions(auth.WithUsername("u"), auth.WithPassword("p"), auth.WithDomainName("override-domain"))
		th.AssertNoErr(t, err)
		v3Opts := opts.(*auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3PasswordOpts{
			Username:       "u",
			Password:       "p",
			UserDomainName: "override-domain",
			Scope:          &auth.Scope{ProjectDomainName: "override-domain"},
		}, v3Opts.Auth)
	})

	t.Run("WithProjectID", func(t *testing.T) {
		opts, err := baseCloud().AuthOptions(auth.WithUsername("u"), auth.WithPassword("p"), auth.WithProjectID("override-project-id"))
		th.AssertNoErr(t, err)
		v3Opts := opts.(*auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3PasswordOpts{
			Username: "u",
			Password: "p",
			Scope:    &auth.Scope{ProjectID: "override-project-id"},
		}, v3Opts.Auth)
	})

	t.Run("WithProjectName", func(t *testing.T) {
		opts, err := baseCloud().AuthOptions(auth.WithUsername("u"), auth.WithPassword("p"), auth.WithProjectName("override-project-name"))
		th.AssertNoErr(t, err)
		v3Opts := opts.(*auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3PasswordOpts{
			Username: "u",
			Password: "p",
			Scope:    &auth.Scope{ProjectName: "override-project-name"},
		}, v3Opts.Auth)
	})

	t.Run("WithApplicationCredentialID and Secret", func(t *testing.T) {
		opts, err := baseCloud().AuthOptions(
			auth.WithApplicationCredentialID("override-appcred-id"),
			auth.WithApplicationCredentialSecret("override-appcred-secret"),
		)
		th.AssertNoErr(t, err)
		v3Opts := opts.(*auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3ApplicationCredentialOpts{
			ApplicationCredentialID:     "override-appcred-id",
			ApplicationCredentialSecret: "override-appcred-secret",
		}, v3Opts.Auth)
	})

	t.Run("WithApplicationCredentialName and Secret", func(t *testing.T) {
		opts, err := baseCloud().AuthOptions(
			auth.WithUsername("u"),
			auth.WithApplicationCredentialName("override-appcred-name"),
			auth.WithApplicationCredentialSecret("override-appcred-secret"),
		)
		th.AssertNoErr(t, err)
		v3Opts := opts.(*auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3ApplicationCredentialOpts{
			Username:                    "u",
			ApplicationCredentialName:   "override-appcred-name",
			ApplicationCredentialSecret: "override-appcred-secret",
		}, v3Opts.Auth)
	})
}

func TestAuthOptionsFromCloudVersionAgnosticAuthType(t *testing.T) {
	t.Run("AuthPassword with IdentityAPIVersion 2.0 resolves to V2", func(t *testing.T) {
		cloud := clouds.Cloud{
			AuthType:           auth.AuthPassword,
			IdentityAPIVersion: "2.0",
			Auth: map[string]any{
				"auth_url": "http://example.com:5000",
				"username": "testuser",
				"password": "testpass",
			},
		}

		opts, err := cloud.AuthOptions()
		th.AssertNoErr(t, err)

		v2Opts, ok := opts.(*auth.AuthOptionsV2)
		th.AssertEquals(t, true, ok)
		th.AssertDeepEquals(t, auth.V2PasswordOpts{Username: "testuser", Password: "testpass"}, v2Opts.Auth)
	})

	t.Run("AuthPassword without IdentityAPIVersion resolves to V3", func(t *testing.T) {
		cloud := clouds.Cloud{
			AuthType: auth.AuthPassword,
			Auth: map[string]any{
				"auth_url": "http://example.com:5000",
				"username": "testuser",
				"password": "testpass",
			},
		}

		opts, err := cloud.AuthOptions()
		th.AssertNoErr(t, err)

		v3Opts, ok := opts.(*auth.AuthOptionsV3)
		th.AssertEquals(t, true, ok)
		th.AssertDeepEquals(t, auth.V3PasswordOpts{Username: "testuser", Password: "testpass", Scope: &auth.Scope{}}, v3Opts.Auth)
	})

	t.Run("AuthToken with IdentityAPIVersion 2.0 resolves to V2", func(t *testing.T) {
		cloud := clouds.Cloud{
			AuthType:           auth.AuthToken,
			IdentityAPIVersion: "2.0",
			Auth: map[string]any{
				"auth_url": "http://example.com:5000",
				"token":    "testtoken",
			},
		}

		opts, err := cloud.AuthOptions()
		th.AssertNoErr(t, err)

		v2Opts, ok := opts.(*auth.AuthOptionsV2)
		th.AssertEquals(t, true, ok)
		th.AssertDeepEquals(t, auth.V2TokenOpts{Token: "testtoken"}, v2Opts.Auth)
	})

	t.Run("AuthToken without IdentityAPIVersion resolves to V3", func(t *testing.T) {
		cloud := clouds.Cloud{
			AuthType: auth.AuthToken,
			Auth: map[string]any{
				"auth_url": "http://example.com:5000",
				"token":    "testtoken",
			},
		}

		opts, err := cloud.AuthOptions()
		th.AssertNoErr(t, err)

		v3Opts, ok := opts.(*auth.AuthOptionsV3)
		th.AssertEquals(t, true, ok)
		th.AssertDeepEquals(t, auth.V3TokenOpts{Token: "testtoken", Scope: &auth.Scope{}}, v3Opts.Auth)
	})
}

func TestAuthOptionsFromCloudAuthTypeOverridesIdentityAPIVersion(t *testing.T) {
	t.Run("explicit AuthV3Password overrides IdentityAPIVersion 2.0", func(t *testing.T) {
		cloud := clouds.Cloud{
			AuthType:           auth.AuthV3Password,
			IdentityAPIVersion: "2.0",
			Auth: map[string]any{
				"auth_url": "http://example.com:5000",
				"username": "testuser",
				"password": "testpass",
			},
		}

		opts, err := cloud.AuthOptions()
		th.AssertNoErr(t, err)

		v3Opts, ok := opts.(*auth.AuthOptionsV3)
		th.AssertEquals(t, true, ok)
		th.AssertDeepEquals(t, auth.V3PasswordOpts{Username: "testuser", Password: "testpass", Scope: &auth.Scope{}}, v3Opts.Auth)
	})

	t.Run("explicit AuthV2Password takes precedence when IdentityAPIVersion is unset", func(t *testing.T) {
		cloud := clouds.Cloud{
			AuthType: auth.AuthV2Password,
			Auth: map[string]any{
				"auth_url": "http://example.com:5000",
				"username": "testuser",
				"password": "testpass",
			},
		}

		opts, err := cloud.AuthOptions()
		th.AssertNoErr(t, err)

		v2Opts, ok := opts.(*auth.AuthOptionsV2)
		th.AssertEquals(t, true, ok)
		th.AssertDeepEquals(t, auth.V2PasswordOpts{Username: "testuser", Password: "testpass"}, v2Opts.Auth)
	})
}

func ExampleWithUsername() {
	cloud := clouds.Cloud{
		Auth: map[string]any{
			"auth_url": "https://example.com:13000",
			"password": "secret",
		},
	}

	authOpts, err := cloud.AuthOptionsV3(auth.WithUsername("Kris"))
	if err != nil {
		panic(err)
	}

	fmt.Println(authOpts.Auth.(auth.V3PasswordOpts).Username)
	// Output: Kris
}
