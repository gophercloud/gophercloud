package clouds_test

import (
	"fmt"
	"maps"
	"net/http"
	"strings"
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

	opts, err := cloud.ToAuthOptions()
	th.AssertNoErr(t, err)

	v3Opts, ok := opts.(auth.AuthOptionsV3)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000", v3Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V3PasswordOpts{Username: "testuser", Password: "testpass", Scope: &auth.Scope{}, AllowReauth: true}, v3Opts.Auth)
}

func TestAuthOptionsFromCloudDefaultsAllowReauthForReusableMechanisms(t *testing.T) {
	testCases := map[string]clouds.Cloud{
		"v2 password": {
			AuthType: auth.AuthV2Password,
			Auth: map[string]any{
				"auth_url": "http://example.com:5000",
				"username": "testuser",
				"password": "testpass",
			},
		},
		"v2 token": {
			AuthType: auth.AuthV2Token,
			Auth: map[string]any{
				"auth_url": "http://example.com:5000",
				"token":    "testtoken",
			},
		},
		"v3 password": {
			AuthType: auth.AuthV3Password,
			Auth: map[string]any{
				"auth_url": "http://example.com:5000",
				"username": "testuser",
				"password": "testpass",
			},
		},
		"v3 application credential": {
			AuthType: auth.AuthV3ApplicationCredential,
			Auth: map[string]any{
				"auth_url":                      "http://example.com:5000",
				"application_credential_id":     "test-id",
				"application_credential_secret": "test-secret",
			},
		},
	}

	for name, cloud := range testCases {
		t.Run(name, func(t *testing.T) {
			opts, err := cloud.ToAuthOptions()
			th.AssertNoErr(t, err)
			th.AssertEquals(t, true, canReauth(opts))
		})
	}
}

func TestAuthOptionsFromCloudExplicitFalseDisablesReauth(t *testing.T) {
	testCases := map[string]clouds.Cloud{
		"v2 password": {
			AuthType: auth.AuthV2Password,
			Auth: map[string]any{
				"auth_url":     "http://example.com:5000",
				"username":     "testuser",
				"password":     "testpass",
				"allow_reauth": false,
			},
		},
		"v2 token": {
			AuthType: auth.AuthV2Token,
			Auth: map[string]any{
				"auth_url":     "http://example.com:5000",
				"token":        "testtoken",
				"allow_reauth": false,
			},
		},
		"v3 password": {
			AuthType: auth.AuthV3Password,
			Auth: map[string]any{
				"auth_url":     "http://example.com:5000",
				"username":     "testuser",
				"password":     "testpass",
				"allow_reauth": false,
			},
		},
		"v3 application credential": {
			AuthType: auth.AuthV3ApplicationCredential,
			Auth: map[string]any{
				"auth_url":                      "http://example.com:5000",
				"application_credential_id":     "test-id",
				"application_credential_secret": "test-secret",
				"allow_reauth":                  false,
			},
		},
	}

	for name, cloud := range testCases {
		t.Run(name, func(t *testing.T) {
			opts, err := cloud.ToAuthOptions()
			th.AssertNoErr(t, err)
			th.AssertEquals(t, false, canReauth(opts))
		})
	}
}

func TestAuthOptionsFromCloudCannotEnableReauthForNonReusableMechanisms(t *testing.T) {
	testCases := map[string]clouds.Cloud{
		"v3 token": {
			AuthType: auth.AuthV3Token,
			Auth: map[string]any{
				"auth_url":     "http://example.com:5000",
				"token":        "testtoken",
				"allow_reauth": true,
			},
		},
		"v3 TOTP": {
			AuthType: auth.AuthV3Totp,
			Auth: map[string]any{
				"auth_url":     "http://example.com:5000",
				"username":     "testuser",
				"passcode":     "123456",
				"allow_reauth": true,
			},
		},
	}

	for name, cloud := range testCases {
		t.Run(name, func(t *testing.T) {
			opts, err := cloud.ToAuthOptions()
			th.AssertNoErr(t, err)
			th.AssertEquals(t, false, canReauth(opts))
		})
	}
}

func TestAuthOptionsFromCloudV3Multifactor(t *testing.T) {
	const cloudsYAML = `
clouds:
  multifactor:
    auth_type: v3multifactor
    auth:
      auth_url: http://example.com:5000/v3
      auth_methods:
        - v3password
        - v3totp
      username: testuser
      password: testpass
      user_domain_id: default
`

	cloud, _, _, err := clouds.Parse(
		clouds.WithCloudName("multifactor"),
		clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
	)
	th.AssertNoErr(t, err)

	opts, err := cloud.ToAuthOptions(auth.WithPasscode("123456"))
	th.AssertNoErr(t, err)

	v3Opts, ok := opts.(auth.AuthOptionsV3)
	th.AssertEquals(t, true, ok)
	th.AssertDeepEquals(t, auth.V3MultifactorOpts{
		AuthMethods: []auth.AuthOptionsBuilderV3{
			auth.V3PasswordOpts{
				Username:     "testuser",
				Password:     "testpass",
				UserDomainID: "default",
				Scope:        &auth.Scope{},
				AllowReauth:  true,
			},
			auth.V3TOTPOpts{
				Username:     "testuser",
				Passcode:     "123456",
				UserDomainID: "default",
				Scope:        &auth.Scope{},
			},
		},
		Scope: &auth.Scope{},
	}, v3Opts.Auth)
	th.AssertEquals(t, false, v3Opts.Auth.CanReauth())
}

func TestAuthOptionsFromCloudV3MultifactorRequiresAuthMethods(t *testing.T) {
	cloud := clouds.Cloud{
		AuthType: auth.AuthV3MultiFactor,
		Auth: map[string]any{
			"auth_url": "http://example.com:5000/v3",
		},
	}

	_, err := cloud.ToAuthOptions()
	th.AssertErr(t, err)

	missing, ok := err.(gophercloud.ErrMissingInput)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "AuthMethods", missing.Argument)
}

func TestAuthOptionsFromCloudV3MultifactorRejectsUnsupportedMethod(t *testing.T) {
	cloud := clouds.Cloud{
		AuthType: auth.AuthV3MultiFactor,
		Auth: map[string]any{
			"auth_url":     "http://example.com:5000/v3",
			"auth_methods": []string{"v3unsupported"},
		},
	}

	_, err := cloud.ToAuthOptions()
	th.AssertErr(t, err)

	unsupported, ok := err.(gophercloud.ErrUnsupportedAuthType)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "v3unsupported", unsupported.AuthType)
}

func canReauth(opts auth.Authenticator) bool {
	switch opts := opts.(type) {
	case auth.AuthOptionsV2:
		return opts.Auth.CanReauth()
	case auth.AuthOptionsV3:
		return opts.Auth.CanReauth()
	default:
		return false
	}
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

	opts, err := cloud.ToAuthOptions()
	th.AssertNoErr(t, err)

	v2Opts, ok := opts.(auth.AuthOptionsV2)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000", v2Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V2PasswordOpts{Username: "testuser", Password: "testpass", AllowReauth: true}, v2Opts.Auth)
}

func TestAuthOptionsFromCloudMissingAuthURL(t *testing.T) {
	cloud := clouds.Cloud{
		Auth: map[string]any{
			"username": "testuser",
			"password": "testpass",
		},
	}

	_, err := cloud.ToAuthOptions()
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrMissingInput)
	th.AssertEquals(t, true, ok)
}

func TestAuthOptionsFromCloudNilAuth(t *testing.T) {
	cloud := clouds.Cloud{}

	_, err := cloud.ToAuthOptions()
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrMissingInput)
	th.AssertEquals(t, true, ok)
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
			wantAuth: auth.V2PasswordOpts{Username: "testuser", Password: "testpass", AllowReauth: true},
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
			wantAuth: auth.V2TokenOpts{Token: "testtoken", AllowReauth: true},
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
			wantAuth: auth.V3PasswordOpts{Username: "testuser", Password: "testpass", Scope: &auth.Scope{}, AllowReauth: true},
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
			wantAuth: auth.V3ApplicationCredentialOpts{ApplicationCredentialID: "app-cred-id", ApplicationCredentialSecret: "app-cred-secret", AllowReauth: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := tt.cloud.ToAuthOptions()
			th.AssertNoErr(t, err)

			if tt.wantV2 {
				v2Opts, ok := opts.(auth.AuthOptionsV2)
				th.AssertEquals(t, true, ok)
				th.AssertDeepEquals(t, tt.wantAuth, v2Opts.Auth)
			} else {
				v3Opts, ok := opts.(auth.AuthOptionsV3)
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

	opts, err := cloud.ToAuthOptions(auth.WithPasscode("123456"))
	th.AssertNoErr(t, err)

	v3Opts, ok := opts.(auth.AuthOptionsV3)
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
		opts, err := cloud.ToAuthOptions(auth.WithPasscode("123456"))
		th.AssertNoErr(t, err)
		v3Opts := opts.(auth.AuthOptionsV3)
		_, ok := v3Opts.Auth.(auth.V3PasswordOpts)
		th.AssertEquals(t, true, ok)
	})

	t.Run("passcode beats token and appcred", func(t *testing.T) {
		m := maps.Clone(base)
		delete(m, "password")
		cloud := clouds.Cloud{Auth: m}
		opts, err := cloud.ToAuthOptions(auth.WithPasscode("123456"))
		th.AssertNoErr(t, err)
		v3Opts := opts.(auth.AuthOptionsV3)
		_, ok := v3Opts.Auth.(auth.V3TOTPOpts)
		th.AssertEquals(t, true, ok)
	})

	t.Run("token beats appcred", func(t *testing.T) {
		m := maps.Clone(base)
		delete(m, "password")
		cloud := clouds.Cloud{Auth: m}
		opts, err := cloud.ToAuthOptions()
		th.AssertNoErr(t, err)
		v3Opts := opts.(auth.AuthOptionsV3)
		_, ok := v3Opts.Auth.(auth.V3TokenOpts)
		th.AssertEquals(t, true, ok)
	})

	t.Run("appcred alone resolves", func(t *testing.T) {
		cloud := clouds.Cloud{Auth: map[string]any{
			"auth_url":                  "http://example.com:5000",
			"application_credential_id": "app-cred-id",
		}}
		opts, err := cloud.ToAuthOptions()
		th.AssertNoErr(t, err)
		v3Opts := opts.(auth.AuthOptionsV3)
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

	_, err := cloud.ToAuthOptions()
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

		opts, err := cloud.ToAuthOptions()
		th.AssertNoErr(t, err)

		v3Opts := opts.(auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3PasswordOpts{
			Username:       "testuser",
			Password:       "testpass",
			UserDomainName: "userdomain",
			AllowReauth:    true,
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

		opts, err := cloud.ToAuthOptions()
		th.AssertNoErr(t, err)

		v3Opts := opts.(auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3PasswordOpts{
			Username:     "testuser",
			Password:     "testpass",
			UserDomainID: "default",
			AllowReauth:  true,
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

		opts, err := cloud.ToAuthOptions()
		th.AssertNoErr(t, err)

		v3Opts := opts.(auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3PasswordOpts{
			Username:       "testuser",
			Password:       "testpass",
			UserDomainName: "shared",
			AllowReauth:    true,
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

		opts, err := cloud.ToAuthOptions()
		th.AssertNoErr(t, err)

		v3Opts := opts.(auth.AuthOptionsV3)
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

		opts, err := cloud.ToAuthOptions()
		th.AssertNoErr(t, err)

		v3Opts := opts.(auth.AuthOptionsV3)
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

		opts, err := cloud.ToAuthOptions()
		th.AssertNoErr(t, err)

		v3Opts := opts.(auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3ApplicationCredentialOpts{
			ApplicationCredentialID:     "app-cred-id",
			ApplicationCredentialSecret: "app-cred-secret",
			AllowReauth:                 true,
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
		opts, err := cloud.ToAuthOptions(auth.WithScope(explicitScope))
		th.AssertNoErr(t, err)

		v3Opts := opts.(auth.AuthOptionsV3)
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
		opts, err := baseCloud().ToAuthOptions(auth.WithUsername("override-user"), auth.WithPassword("override-pass"))
		th.AssertNoErr(t, err)
		v3Opts := opts.(auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3PasswordOpts{Username: "override-user", Password: "override-pass", Scope: &auth.Scope{}, AllowReauth: true}, v3Opts.Auth)
	})

	t.Run("WithUserID", func(t *testing.T) {
		opts, err := baseCloud().ToAuthOptions(auth.WithUserID("override-user-id"), auth.WithPassword("override-pass"))
		th.AssertNoErr(t, err)
		v3Opts := opts.(auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3PasswordOpts{UserID: "override-user-id", Password: "override-pass", Scope: &auth.Scope{}, AllowReauth: true}, v3Opts.Auth)
	})

	t.Run("WithToken", func(t *testing.T) {
		opts, err := baseCloud().ToAuthOptions(auth.WithToken("override-token"))
		th.AssertNoErr(t, err)
		v3Opts := opts.(auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3TokenOpts{Token: "override-token", Scope: &auth.Scope{}}, v3Opts.Auth)
	})

	t.Run("WithDomainID", func(t *testing.T) {
		opts, err := baseCloud().ToAuthOptions(auth.WithUsername("u"), auth.WithPassword("p"), auth.WithDomainID("override-domain-id"))
		th.AssertNoErr(t, err)
		v3Opts := opts.(auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3PasswordOpts{
			Username:     "u",
			Password:     "p",
			UserDomainID: "override-domain-id",
			Scope:        &auth.Scope{ProjectDomainID: "override-domain-id"},
			AllowReauth:  true,
		}, v3Opts.Auth)
	})

	t.Run("WithDomainName", func(t *testing.T) {
		opts, err := baseCloud().ToAuthOptions(auth.WithUsername("u"), auth.WithPassword("p"), auth.WithDomainName("override-domain"))
		th.AssertNoErr(t, err)
		v3Opts := opts.(auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3PasswordOpts{
			Username:       "u",
			Password:       "p",
			UserDomainName: "override-domain",
			Scope:          &auth.Scope{ProjectDomainName: "override-domain"},
			AllowReauth:    true,
		}, v3Opts.Auth)
	})

	t.Run("WithProjectID", func(t *testing.T) {
		opts, err := baseCloud().ToAuthOptions(auth.WithUsername("u"), auth.WithPassword("p"), auth.WithProjectID("override-project-id"))
		th.AssertNoErr(t, err)
		v3Opts := opts.(auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3PasswordOpts{
			Username:    "u",
			Password:    "p",
			Scope:       &auth.Scope{ProjectID: "override-project-id"},
			AllowReauth: true,
		}, v3Opts.Auth)
	})

	t.Run("WithProjectName", func(t *testing.T) {
		opts, err := baseCloud().ToAuthOptions(auth.WithUsername("u"), auth.WithPassword("p"), auth.WithProjectName("override-project-name"))
		th.AssertNoErr(t, err)
		v3Opts := opts.(auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3PasswordOpts{
			Username:    "u",
			Password:    "p",
			Scope:       &auth.Scope{ProjectName: "override-project-name"},
			AllowReauth: true,
		}, v3Opts.Auth)
	})

	t.Run("WithApplicationCredentialID and Secret", func(t *testing.T) {
		opts, err := baseCloud().ToAuthOptions(
			auth.WithApplicationCredentialID("override-appcred-id"),
			auth.WithApplicationCredentialSecret("override-appcred-secret"),
		)
		th.AssertNoErr(t, err)
		v3Opts := opts.(auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3ApplicationCredentialOpts{
			ApplicationCredentialID:     "override-appcred-id",
			ApplicationCredentialSecret: "override-appcred-secret",
			AllowReauth:                 true,
		}, v3Opts.Auth)
	})

	t.Run("WithApplicationCredentialName and Secret", func(t *testing.T) {
		opts, err := baseCloud().ToAuthOptions(
			auth.WithUsername("u"),
			auth.WithApplicationCredentialName("override-appcred-name"),
			auth.WithApplicationCredentialSecret("override-appcred-secret"),
		)
		th.AssertNoErr(t, err)
		v3Opts := opts.(auth.AuthOptionsV3)
		th.AssertDeepEquals(t, auth.V3ApplicationCredentialOpts{
			Username:                    "u",
			ApplicationCredentialName:   "override-appcred-name",
			ApplicationCredentialSecret: "override-appcred-secret",
			AllowReauth:                 true,
		}, v3Opts.Auth)
	})
}

func TestAuthOptionsFromCloudVersionAgnosticAuthType(t *testing.T) {
	t.Run("AuthPassword with IdentityAPIVersion 2.0 resolves to V2", func(t *testing.T) {
		fakeServer := setupIdentityVersion(t, "v2.0", "v2.0/")
		cloud := clouds.Cloud{
			AuthType:           auth.AuthPassword,
			IdentityAPIVersion: "2.0",
			Auth: map[string]any{
				"auth_url": fakeServer.Endpoint(),
				"username": "testuser",
				"password": "testpass",
			},
		}

		opts, err := cloud.ToAuthOptions()
		th.AssertNoErr(t, err)

		v2Opts, ok := opts.(auth.AuthOptionsV2)
		th.AssertEquals(t, true, ok)
		th.AssertDeepEquals(t, auth.V2PasswordOpts{Username: "testuser", Password: "testpass", AllowReauth: true}, v2Opts.Auth)
	})

	t.Run("AuthPassword without IdentityAPIVersion resolves to V3", func(t *testing.T) {
		fakeServer := setupIdentityVersion(t, "v3.0", "v3/")
		cloud := clouds.Cloud{
			AuthType: auth.AuthPassword,
			Auth: map[string]any{
				"auth_url": fakeServer.Endpoint(),
				"username": "testuser",
				"password": "testpass",
			},
		}

		opts, err := cloud.ToAuthOptions()
		th.AssertNoErr(t, err)

		v3Opts, ok := opts.(auth.AuthOptionsV3)
		th.AssertEquals(t, true, ok)
		th.AssertDeepEquals(t, auth.V3PasswordOpts{Username: "testuser", Password: "testpass", Scope: &auth.Scope{}, AllowReauth: true}, v3Opts.Auth)
	})

	t.Run("AuthToken with IdentityAPIVersion 2.0 resolves to V2", func(t *testing.T) {
		fakeServer := setupIdentityVersion(t, "v2.0", "v2.0/")
		cloud := clouds.Cloud{
			AuthType:           auth.AuthToken,
			IdentityAPIVersion: "2.0",
			Auth: map[string]any{
				"auth_url": fakeServer.Endpoint(),
				"token":    "testtoken",
			},
		}

		opts, err := cloud.ToAuthOptions()
		th.AssertNoErr(t, err)

		v2Opts, ok := opts.(auth.AuthOptionsV2)
		th.AssertEquals(t, true, ok)
		th.AssertDeepEquals(t, auth.V2TokenOpts{Token: "testtoken", AllowReauth: true}, v2Opts.Auth)
	})

	t.Run("AuthToken without IdentityAPIVersion resolves to V3", func(t *testing.T) {
		fakeServer := setupIdentityVersion(t, "v3.0", "v3/")
		cloud := clouds.Cloud{
			AuthType: auth.AuthToken,
			Auth: map[string]any{
				"auth_url": fakeServer.Endpoint(),
				"token":    "testtoken",
			},
		}

		opts, err := cloud.ToAuthOptions()
		th.AssertNoErr(t, err)

		v3Opts, ok := opts.(auth.AuthOptionsV3)
		th.AssertEquals(t, true, ok)
		th.AssertDeepEquals(t, auth.V3TokenOpts{Token: "testtoken", Scope: &auth.Scope{}}, v3Opts.Auth)
	})
}

func setupIdentityVersion(t *testing.T, versionID, suffix string) th.FakeServer {
	t.Helper()

	fakeServer := th.SetupHTTP()
	t.Cleanup(fakeServer.Teardown)
	fakeServer.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodGet)
		fmt.Fprintf(w, `{"versions":{"values":[{"id":%q,"status":"stable","links":[{"href":%q,"rel":"self"}]}]}}`, versionID, fakeServer.Endpoint()+suffix)
	})
	return fakeServer
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

		opts, err := cloud.ToAuthOptions()
		th.AssertNoErr(t, err)

		v3Opts, ok := opts.(auth.AuthOptionsV3)
		th.AssertEquals(t, true, ok)
		th.AssertDeepEquals(t, auth.V3PasswordOpts{Username: "testuser", Password: "testpass", Scope: &auth.Scope{}, AllowReauth: true}, v3Opts.Auth)
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

		opts, err := cloud.ToAuthOptions()
		th.AssertNoErr(t, err)

		v2Opts, ok := opts.(auth.AuthOptionsV2)
		th.AssertEquals(t, true, ok)
		th.AssertDeepEquals(t, auth.V2PasswordOpts{Username: "testuser", Password: "testpass", AllowReauth: true}, v2Opts.Auth)
	})
}
