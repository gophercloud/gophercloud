package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	"github.com/gophercloud/gophercloud/v2/openstack/config/clouds"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// These tests call auth.AuthOptionsFromCloud directly, using a real
// clouds.Cloud as the auth.CloudSource, to prove the restored top-level
// function works standalone and not merely through Cloud.AuthOptions's
// wrapper (covered exhaustively in openstack/config/clouds/auth_test.go).

func TestAuthOptionsFromCloudDirect(t *testing.T) {
	cloud := clouds.Cloud{
		Auth: map[string]any{
			"auth_url": "http://example.com:5000",
			"username": "testuser",
			"password": "testpass",
		},
	}

	opts, err := auth.AuthOptionsFromCloud(cloud)
	th.AssertNoErr(t, err)

	v3Opts, ok := opts.(auth.AuthOptionsV3)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000", v3Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V3PasswordOpts{Username: "testuser", Password: "testpass", Scope: &auth.Scope{}, AllowReauth: true}, v3Opts.Auth)
}

func TestAuthOptionsFromCloudNoUsableCredentials(t *testing.T) {
	cloud := clouds.Cloud{
		Auth: map[string]any{
			"auth_url": "http://example.com:5000",
		},
	}

	opts, err := cloud.ToAuthOptions()
	th.AssertNoErr(t, err)

	_, err = opts.Authenticate(context.TODO(), nil)
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrMissingPassword)
	th.AssertEquals(t, true, ok)
}

// implementation of former "fails system_scope when value is not all" test in the clouds package
func TestAuthOptionsFromCloudFailsSystemScopeWhenValueIsNotAll(t *testing.T) {
	cloud := clouds.Cloud{
		Auth: map[string]any{
			"auth_url":            "https://example.com:5000/v3",
			"username":            "myuser",
			"password":            "mypassword",
			"user_domain_name":    "Default",
			"system_scope":        "something-else",
			"project_name":        "myproject",
			"project_domain_name": "Default",
		},
	}

	_, err := auth.AuthOptionsFromCloud(cloud)
	if err == nil {
		t.Fatalf("an error expected, got nil")
	}
}

func TestAuthOptionsFromCloudV2Direct(t *testing.T) {
	cloud := clouds.Cloud{
		AuthType: auth.AuthV2Password,
		Auth: map[string]any{
			"auth_url": "http://example.com:5000",
			"username": "testuser",
			"password": "testpass",
		},
	}

	v2Opts, err := auth.AuthOptionsFromCloudV2(cloud)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, auth.V2PasswordOpts{Username: "testuser", Password: "testpass", AllowReauth: true}, v2Opts.Auth)
}

func TestAuthOptionsFromCloudV3Direct(t *testing.T) {
	cloud := clouds.Cloud{
		Auth: map[string]any{
			"auth_url": "http://example.com:5000",
			"password": "secret",
		},
	}

	v3Opts, err := auth.AuthOptionsFromCloudV3(cloud, auth.WithUsername("Kris"))
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "Kris", v3Opts.Auth.(auth.V3PasswordOpts).Username)
}

type cloudSource struct {
	authType           auth.AuthType
	identityAPIVersion string
	authData           map[string]any
}

func (c cloudSource) GetAuthType() auth.AuthType    { return c.authType }
func (c cloudSource) GetIdentityAPIVersion() string { return c.identityAPIVersion }
func (c cloudSource) GetAuthData() map[string]any   { return c.authData }

func TestResolveCloudOptions(t *testing.T) {
	scope := &auth.Scope{ProjectID: "project-id"}
	got := auth.ResolveCloudOptions(
		auth.WithAuthURL("https://identity.example.com"),
		auth.WithUsername("testuser"),
		auth.WithUserID("user-id"),
		auth.WithPassword("testpass"),
		auth.WithToken("testtoken"),
		auth.WithPasscode("123456"),
		auth.WithUserDomainID("user-domain-id"),
		auth.WithUserDomainName("user-domain-name"),
		auth.WithDomainID("domain-id"),
		auth.WithDomainName("domain-name"),
		auth.WithProjectID("project-id"),
		auth.WithProjectName("project-name"),
		auth.WithProjectDomainID("project-domain-id"),
		auth.WithProjectDomainName("project-domain-name"),
		auth.WithSystemScope("all"),
		auth.WithTrustID("trust-id"),
		auth.WithApplicationCredentialID("credential-id"),
		auth.WithApplicationCredentialName("credential-name"),
		auth.WithApplicationCredentialSecret("secret"),
		auth.WithConsumerKey("consumer-key"),
		auth.WithConsumerSecret("consumer-secret"),
		auth.WithAccessKey("access-key"),
		auth.WithAccessSecret("access-secret"),
		auth.WithAllowReauth(false),
		auth.WithAuthMethods(auth.AuthV3Password, auth.AuthV3Totp),
		auth.WithScope(scope),
	)

	allowReauth := false

	th.AssertDeepEquals(t, auth.CloudOptions{
		AuthURL:                     "https://identity.example.com",
		Username:                    "testuser",
		UserID:                      "user-id",
		Password:                    "testpass",
		Token:                       "testtoken",
		Passcode:                    "123456",
		UserDomainID:                "user-domain-id",
		UserDomainName:              "user-domain-name",
		DomainID:                    "domain-id",
		DomainName:                  "domain-name",
		ProjectID:                   "project-id",
		ProjectName:                 "project-name",
		ProjectDomainID:             "project-domain-id",
		ProjectDomainName:           "project-domain-name",
		SystemScope:                 "all",
		TrustID:                     "trust-id",
		ApplicationCredentialID:     "credential-id",
		ApplicationCredentialName:   "credential-name",
		ApplicationCredentialSecret: "secret",
		ConsumerKey:                 "consumer-key",
		ConsumerSecret:              "consumer-secret",
		AccessKey:                   "access-key",
		AccessSecret:                "access-secret",
		AllowReauth:                 &allowReauth,
		AuthMethods:                 []auth.AuthType{auth.AuthV3Password, auth.AuthV3Totp},
		Scope:                       scope,
	}, got)
}

// TestAuthOptionsFromCloudV3Overrides checks that each override replaces the
// cloud's value in the resulting auth options.
func TestAuthOptionsFromCloudV3Overrides(t *testing.T) {
	newCloud := func(data map[string]any) cloudSource {
		data["auth_url"] = "https://identity.example.com"
		data["username"] = "testuser"
		data["password"] = "testpass"
		return cloudSource{authType: auth.AuthV3Password, authData: data}
	}

	t.Run("WithUserDomainID", func(t *testing.T) {
		got, err := auth.AuthOptionsFromCloudV3(newCloud(map[string]any{"user_domain_id": "cloud"}), auth.WithUserDomainID("override"))
		th.AssertNoErr(t, err)
		th.AssertEquals(t, "override", got.Auth.(auth.V3PasswordOpts).UserDomainID)
	})

	t.Run("WithUserDomainName", func(t *testing.T) {
		got, err := auth.AuthOptionsFromCloudV3(newCloud(map[string]any{"user_domain_name": "cloud"}), auth.WithUserDomainName("override"))
		th.AssertNoErr(t, err)
		th.AssertEquals(t, "override", got.Auth.(auth.V3PasswordOpts).UserDomainName)
	})

	t.Run("WithProjectDomainID", func(t *testing.T) {
		got, err := auth.AuthOptionsFromCloudV3(newCloud(map[string]any{"project_name": "project", "project_domain_id": "cloud"}), auth.WithProjectDomainID("override"))
		th.AssertNoErr(t, err)
		th.AssertDeepEquals(t, &auth.Scope{ProjectName: "project", ProjectDomainID: "override"}, got.Auth.(auth.V3PasswordOpts).Scope)
	})

	t.Run("WithProjectDomainName", func(t *testing.T) {
		got, err := auth.AuthOptionsFromCloudV3(newCloud(map[string]any{"project_name": "project", "project_domain_name": "cloud"}), auth.WithProjectDomainName("override"))
		th.AssertNoErr(t, err)
		th.AssertDeepEquals(t, &auth.Scope{ProjectName: "project", ProjectDomainName: "override"}, got.Auth.(auth.V3PasswordOpts).Scope)
	})

	t.Run("WithSystemScope", func(t *testing.T) {
		got, err := auth.AuthOptionsFromCloudV3(newCloud(map[string]any{}), auth.WithSystemScope("all"))
		th.AssertNoErr(t, err)
		th.AssertDeepEquals(t, &auth.Scope{System: true}, got.Auth.(auth.V3PasswordOpts).Scope)
	})

	t.Run("WithSystemScope invalid", func(t *testing.T) {
		_, err := auth.AuthOptionsFromCloudV3(newCloud(map[string]any{"system_scope": "all"}), auth.WithSystemScope("domain"))
		_, ok := err.(gophercloud.ErrInvalidSystemScope)
		th.AssertEquals(t, true, ok)
	})

	t.Run("WithTrustID", func(t *testing.T) {
		got, err := auth.AuthOptionsFromCloudV3(newCloud(map[string]any{"trust_id": "cloud"}), auth.WithTrustID("override"))
		th.AssertNoErr(t, err)
		th.AssertDeepEquals(t, &auth.Scope{TrustID: "override"}, got.Auth.(auth.V3PasswordOpts).Scope)
	})

	t.Run("OAuth1 overrides", func(t *testing.T) {
		cloud := cloudSource{authType: auth.AuthV3OAuth1, authData: map[string]any{
			"auth_url": "https://identity.example.com", "consumer_key": "cloud", "consumer_secret": "cloud", "access_key": "cloud", "access_secret": "cloud",
		}}
		got, err := auth.AuthOptionsFromCloudV3(cloud,
			auth.WithConsumerKey("consumer-key"), auth.WithConsumerSecret("consumer-secret"),
			auth.WithAccessKey("access-key"), auth.WithAccessSecret("access-secret"))
		th.AssertNoErr(t, err)
		oauth1 := got.Auth.(auth.V3OAuth1Opts)
		th.AssertEquals(t, "consumer-key", oauth1.ConsumerKey)
		th.AssertEquals(t, "consumer-secret", oauth1.ConsumerSecret)
		th.AssertEquals(t, "access-key", oauth1.Token)
		th.AssertEquals(t, "access-secret", oauth1.TokenSecret)
	})

	t.Run("WithAllowReauth", func(t *testing.T) {
		got, err := auth.AuthOptionsFromCloudV3(newCloud(map[string]any{"allow_reauth": true}), auth.WithAllowReauth(false))
		th.AssertNoErr(t, err)
		th.AssertEquals(t, false, got.Auth.CanReauth())
	})

	t.Run("WithAuthMethods", func(t *testing.T) {
		cloud := newCloud(map[string]any{"auth_methods": []any{"v3password"}, "passcode": "123456"})
		cloud.authType = auth.AuthV3MultiFactor
		got, err := auth.AuthOptionsFromCloudV3(cloud, auth.WithAuthMethods(auth.AuthV3Totp))
		th.AssertNoErr(t, err)
		methods := got.Auth.(auth.V3MultifactorOpts).AuthMethods
		th.AssertEquals(t, 1, len(methods))
		_, ok := methods[0].(auth.V3TOTPOpts)
		th.AssertEquals(t, true, ok)
	})
}

func TestAuthOptionsFromCloudDispatchesExplicitAuthTypes(t *testing.T) {
	tests := []struct {
		name     string
		authType auth.AuthType
		data     map[string]any
		wantAuth any
		wantV2   bool
	}{
		{name: "v2 password", authType: auth.AuthV2Password, data: map[string]any{"username": "testuser", "password": "testpass"}, wantAuth: auth.V2PasswordOpts{Username: "testuser", Password: "testpass", AllowReauth: true}, wantV2: true},
		{name: "v2 token", authType: auth.AuthV2Token, data: map[string]any{"token": "testtoken"}, wantAuth: auth.V2TokenOpts{Token: "testtoken", AllowReauth: true}, wantV2: true},
		{name: "v3 password", authType: auth.AuthV3Password, data: map[string]any{"username": "testuser", "password": "testpass"}, wantAuth: auth.V3PasswordOpts{Username: "testuser", Password: "testpass", Scope: &auth.Scope{}, AllowReauth: true}},
		{name: "v3 TOTP", authType: auth.AuthV3Totp, data: map[string]any{"username": "testuser", "passcode": "123456"}, wantAuth: auth.V3TOTPOpts{Username: "testuser", Passcode: "123456", Scope: &auth.Scope{}}},
		{name: "v3 token", authType: auth.AuthV3Token, data: map[string]any{"token": "testtoken"}, wantAuth: auth.V3TokenOpts{Token: "testtoken", Scope: &auth.Scope{}}},
		{name: "v3 application credential", authType: auth.AuthV3ApplicationCredential, data: map[string]any{"application_credential_id": "credential-id", "application_credential_secret": "secret"}, wantAuth: auth.V3ApplicationCredentialOpts{ApplicationCredentialID: "credential-id", ApplicationCredentialSecret: "secret", AllowReauth: true}},
		{
			name:     "v3 OAuth1",
			authType: auth.AuthV3OAuth1,
			data:     map[string]any{"consumer_key": "consumer-key", "consumer_secret": "consumer-secret", "access_key": "access-key", "access_secret": "access-secret"},
			wantAuth: auth.V3OAuth1Opts{
				ConsumerKey: "consumer-key", ConsumerSecret: "consumer-secret", Token: "access-key", TokenSecret: "access-secret",
				SignatureMethod: auth.OAuth1HMACSHA1, AllowReauth: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.data["auth_url"] = "https://identity.example.com"
			got, err := auth.AuthOptionsFromCloud(cloudSource{authType: test.authType, authData: test.data})
			th.AssertNoErr(t, err)
			if test.wantV2 {
				th.AssertDeepEquals(t, test.wantAuth, got.(auth.AuthOptionsV2).Auth)
			} else {
				th.AssertDeepEquals(t, test.wantAuth, got.(auth.AuthOptionsV3).Auth)
			}
		})
	}
}

func TestAuthOptionsFromCloudInfersV2AuthType(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]any
		wantAuth any
	}{
		{name: "password", data: map[string]any{"auth_url": "https://identity.example.com", "username": "testuser", "password": "testpass"}, wantAuth: auth.V2PasswordOpts{Username: "testuser", Password: "testpass", AllowReauth: true}},
		{name: "token", data: map[string]any{"auth_url": "https://identity.example.com", "token": "testtoken"}, wantAuth: auth.V2TokenOpts{Token: "testtoken", AllowReauth: true}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := auth.AuthOptionsFromCloudV2(cloudSource{authData: test.data})
			th.AssertNoErr(t, err)
			th.AssertDeepEquals(t, test.wantAuth, got.Auth)
		})
	}
}

func TestAuthOptionsFromCloudInfersV3AuthType(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]any
		wantAuth any
	}{
		{name: "password", data: map[string]any{"username": "testuser", "password": "testpass"}, wantAuth: auth.V3PasswordOpts{Username: "testuser", Password: "testpass", Scope: &auth.Scope{}, AllowReauth: true}},
		{name: "TOTP", data: map[string]any{"username": "testuser", "passcode": "123456"}, wantAuth: auth.V3TOTPOpts{Username: "testuser", Passcode: "123456", Scope: &auth.Scope{}}},
		{name: "token", data: map[string]any{"token": "testtoken"}, wantAuth: auth.V3TokenOpts{Token: "testtoken", Scope: &auth.Scope{}}},
		{name: "application credential ID", data: map[string]any{"application_credential_id": "credential-id", "application_credential_secret": "secret"}, wantAuth: auth.V3ApplicationCredentialOpts{ApplicationCredentialID: "credential-id", ApplicationCredentialSecret: "secret", AllowReauth: true}},
		{name: "application credential name", data: map[string]any{"application_credential_name": "credential-name", "application_credential_secret": "secret"}, wantAuth: auth.V3ApplicationCredentialOpts{ApplicationCredentialName: "credential-name", ApplicationCredentialSecret: "secret", AllowReauth: true}},
		{name: "fallback password", data: map[string]any{}, wantAuth: auth.V3PasswordOpts{Scope: &auth.Scope{}, AllowReauth: true}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.data["auth_url"] = "https://identity.example.com"
			got, err := auth.AuthOptionsFromCloudV3(cloudSource{authData: test.data})
			th.AssertNoErr(t, err)
			th.AssertDeepEquals(t, test.wantAuth, got.Auth)
		})
	}
}

func TestAuthOptionsFromCloudV3DomainsAndScope(t *testing.T) {
	tests := []struct {
		name         string
		data         map[string]any
		wantUserID   string
		wantUserName string
		wantScope    *auth.Scope
	}{
		{
			name: "domains map to their own fields",
			data: map[string]any{
				"user_domain_id": "user-domain-id", "project_domain_name": "project-domain",
				"project_name": "project-name", "default_domain": "ignored",
			},
			wantUserID: "user-domain-id",
			wantScope:  &auth.Scope{ProjectDomainName: "project-domain", ProjectName: "project-name"},
		},
		{
			name:      "system scope",
			data:      map[string]any{"system_scope": "all"},
			wantScope: &auth.Scope{System: true},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.data["auth_url"] = "https://identity.example.com"
			test.data["username"] = "testuser"
			test.data["password"] = "testpass"
			got, err := auth.AuthOptionsFromCloudV3(cloudSource{authData: test.data})
			th.AssertNoErr(t, err)
			password := got.Auth.(auth.V3PasswordOpts)
			th.AssertEquals(t, test.wantUserID, password.UserDomainID)
			th.AssertEquals(t, test.wantUserName, password.UserDomainName)
			th.AssertDeepEquals(t, test.wantScope, password.Scope)
		})
	}
}

func TestAuthOptionsFromCloudV3UsesExplicitScope(t *testing.T) {
	scope := &auth.Scope{DomainID: "explicit-domain"}
	got, err := auth.AuthOptionsFromCloudV3(cloudSource{authData: map[string]any{
		"auth_url": "https://identity.example.com", "username": "testuser", "password": "testpass", "system_scope": "invalid",
	}}, auth.WithScope(scope))
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, scope, got.Auth.(auth.V3PasswordOpts).Scope)
}

func TestAuthOptionsFromCloudV3Multifactor(t *testing.T) {
	got, err := auth.AuthOptionsFromCloudV3(cloudSource{
		authType: auth.AuthV3MultiFactor,
		authData: map[string]any{
			"auth_url": "https://identity.example.com", "auth_methods": []string{"v3password", "v3totp", "v3token", "v3applicationcredential"},
			"username": "testuser", "password": "testpass", "passcode": "123456", "token": "testtoken",
			"application_credential_id": "credential-id", "application_credential_secret": "secret",
		},
	})
	th.AssertNoErr(t, err)

	multifactor := got.Auth.(auth.V3MultifactorOpts)
	th.AssertEquals(t, 4, len(multifactor.AuthMethods))
	_, passwordOK := multifactor.AuthMethods[0].(auth.V3PasswordOpts)
	_, totpOK := multifactor.AuthMethods[1].(auth.V3TOTPOpts)
	_, tokenOK := multifactor.AuthMethods[2].(auth.V3TokenOpts)
	_, appCredOK := multifactor.AuthMethods[3].(auth.V3ApplicationCredentialOpts)
	th.AssertEquals(t, true, passwordOK && totpOK && tokenOK && appCredOK)
}

func TestAuthOptionsFromCloudErrors(t *testing.T) {
	tests := []struct {
		name        string
		call        func() error
		kind        any
		wantMessage string
	}{
		{name: "top-level missing URL", call: func() error { _, err := auth.AuthOptionsFromCloud(cloudSource{}); return err }, kind: gophercloud.ErrMissingInput{}},
		{name: "v2 missing URL", call: func() error { _, err := auth.AuthOptionsFromCloudV2(cloudSource{}); return err }, kind: gophercloud.ErrMissingInput{}},
		{name: "v3 missing URL", call: func() error { _, err := auth.AuthOptionsFromCloudV3(cloudSource{}); return err }, kind: gophercloud.ErrMissingInput{}},
		{name: "v2 unsupported type", call: func() error {
			_, err := auth.AuthOptionsFromCloudV2(cloudSource{authType: "unsupported", authData: map[string]any{"auth_url": "https://identity.example.com"}})
			return err
		}, kind: gophercloud.ErrUnsupportedAuthType{}},
		{name: "v3 unsupported type", call: func() error {
			_, err := auth.AuthOptionsFromCloudV3(cloudSource{authType: "unsupported", authData: map[string]any{"auth_url": "https://identity.example.com"}})
			return err
		}, kind: gophercloud.ErrUnsupportedAuthType{}},
		{name: "invalid system scope", call: func() error {
			_, err := auth.AuthOptionsFromCloudV3(cloudSource{authData: map[string]any{"auth_url": "https://identity.example.com", "system_scope": "domain"}})
			return err
		}, kind: gophercloud.ErrInvalidSystemScope{}, wantMessage: "Only system scope of all is supported"},
		{name: "multifactor without methods", call: func() error {
			_, err := auth.AuthOptionsFromCloudV3(cloudSource{authType: auth.AuthV3MultiFactor, authData: map[string]any{"auth_url": "https://identity.example.com"}})
			return err
		}, kind: gophercloud.ErrMissingInput{}},
		{name: "multifactor unsupported method", call: func() error {
			_, err := auth.AuthOptionsFromCloudV3(cloudSource{authType: auth.AuthV3MultiFactor, authData: map[string]any{"auth_url": "https://identity.example.com", "auth_methods": []string{"unsupported"}}})
			return err
		}, kind: gophercloud.ErrUnsupportedAuthType{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.call()
			th.AssertErr(t, err)
			if test.kind != nil {
				th.AssertEquals(t, fmt.Sprintf("%T", test.kind), fmt.Sprintf("%T", err))
			}
			if test.wantMessage != "" {
				th.AssertEquals(t, test.wantMessage, err.Error())
			}
		})
	}
}

func TestAuthOptionsFromCloudDiscoversIdentityVersion(t *testing.T) {
	tests := []struct {
		name      string
		versionID string
		suffix    string
		wantV2    bool
	}{
		{name: "v2", versionID: "v2.0", suffix: "v2.0/", wantV2: true},
		{name: "v3", versionID: "v3.0", suffix: "v3/"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()
			fakeServer.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				th.TestMethod(t, r, http.MethodGet)
				fmt.Fprintf(w, `{"versions":{"values":[{"id":%q,"status":"stable","links":[{"href":%q,"rel":"self"}]}]}}`, test.versionID, fakeServer.Endpoint()+test.suffix)
			})
			path := "/v3/auth/tokens"
			body := `{"token":{}}`
			if test.wantV2 {
				path = "/v2.0/tokens"
				body = `{"access":{"token":{"id":"discovered"}}}`
			}
			fakeServer.Mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
				th.TestMethod(t, r, http.MethodPost)
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-Subject-Token", "discovered")
				if !test.wantV2 {
					w.WriteHeader(http.StatusCreated)
				}
				fmt.Fprint(w, body)
			})

			got, err := auth.AuthOptionsFromCloud(cloudSource{authType: auth.AuthPassword, authData: map[string]any{
				"auth_url": fakeServer.Endpoint(), "username": "testuser", "password": "testpass", "user_domain_id": "default",
			}})
			th.AssertNoErr(t, err)
			result, err := got.Authenticate(context.Background(), nil)
			th.AssertNoErr(t, err)
			th.AssertEquals(t, "discovered", result.TokenID)
		})
	}
}
