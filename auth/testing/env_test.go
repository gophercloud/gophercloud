package testing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestAuthOptionsFromEnvDefaultsVersionlessURLToV3(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/identity/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Subject-Token", "testtoken")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"token":{}}`)
	})

	t.Setenv("OS_AUTH_URL", fakeServer.Endpoint()+"identity")
	t.Setenv("OS_USERID", "user-id")
	t.Setenv("OS_PASSWORD", "password")

	opts, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)
	result, err := opts.Authenticate(context.Background(), nil)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "testtoken", result.TokenID)
}

func TestAuthOptionsFromEnvV2Password(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV2Password(t)

	opts, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	v2Opts, ok := opts.(auth.AuthOptionsV2)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000/v2.0", v2Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V2PasswordOpts{Username: "testuser", Password: "testpass", AllowReauth: true}, v2Opts.Auth)
}

func TestAuthOptionsFromEnvV2Token(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV2Token(t)

	opts, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	v2Opts, ok := opts.(auth.AuthOptionsV2)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000/v2.0", v2Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V2TokenOpts{Token: "testtoken", AllowReauth: true}, v2Opts.Auth)
}

func TestAuthOptionsFromEnvExplicitV2AuthType(t *testing.T) {
	testCases := []struct {
		name     string
		authType string
		setCreds func(*testing.T)
		expected auth.AuthOptionsBuilderV2
	}{
		{
			name:     "password",
			authType: "v2password",
			setCreds: func(t *testing.T) {
				t.Setenv("OS_USERNAME", "testuser")
				t.Setenv("OS_PASSWORD", "testpass")
			},
			expected: auth.V2PasswordOpts{Username: "testuser", Password: "testpass", TenantID: "tenant-id", TenantName: "tenant-name", AllowReauth: true},
		},
		{
			name:     "token",
			authType: "v2token",
			setCreds: func(t *testing.T) {
				t.Setenv("OS_TOKEN", "testtoken")
			},
			expected: auth.V2TokenOpts{Token: "testtoken", TenantID: "tenant-id", TenantName: "tenant-name", AllowReauth: true},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			defer CleanupEnv(t)
			CleanupEnv(t)
			t.Setenv("OS_AUTH_URL", "http://example.com:5000/v2.0")
			t.Setenv("OS_AUTH_TYPE", testCase.authType)
			t.Setenv("OS_IDENTITY_API_VERSION", "3")
			t.Setenv("OS_PROJECT_ID", "project-id")
			t.Setenv("OS_PROJECT_NAME", "project-name")
			t.Setenv("OS_TENANT_ID", "tenant-id")
			t.Setenv("OS_TENANT_NAME", "tenant-name")
			testCase.setCreds(t)

			opts, err := auth.AuthOptionsFromEnv()
			th.AssertNoErr(t, err)

			v2Opts, ok := opts.(auth.AuthOptionsV2)
			th.AssertEquals(t, true, ok)
			th.AssertDeepEquals(t, testCase.expected, v2Opts.Auth)
		})
	}
}

func TestAuthOptionsFromEnvGenericAuthResolvesV2ScopeAliases(t *testing.T) {
	testCases := []struct {
		name         string
		authType     string
		env          map[string]string
		expected     auth.AuthOptionsBuilderV2
		expectedBody map[string]map[string]any
	}{
		{
			name:     "password project scope",
			authType: "password",
			env:      map[string]string{"OS_USERNAME": "testuser", "OS_PASSWORD": "testpass", "OS_PROJECT_ID": "project-id", "OS_PROJECT_NAME": "project-name"},
			expected: auth.V2PasswordOpts{Username: "testuser", Password: "testpass", TenantID: "project-id", TenantName: "project-name", AllowReauth: true},
			expectedBody: map[string]map[string]any{"password": {
				"tenantId": "project-id", "tenantName": "project-name",
				"passwordCredentials": map[string]any{"username": "testuser", "password": "testpass"},
			}},
		},
		{
			name:     "password tenant scope",
			authType: "password",
			env:      map[string]string{"OS_USERNAME": "testuser", "OS_PASSWORD": "testpass", "OS_TENANT_ID": "tenant-id", "OS_TENANT_NAME": "tenant-name"},
			expected: auth.V2PasswordOpts{Username: "testuser", Password: "testpass", TenantID: "tenant-id", TenantName: "tenant-name", AllowReauth: true},
			expectedBody: map[string]map[string]any{"password": {
				"tenantId": "tenant-id", "tenantName": "tenant-name",
				"passwordCredentials": map[string]any{"username": "testuser", "password": "testpass"},
			}},
		},
		{
			name:     "password project scope takes precedence",
			authType: "password",
			env: map[string]string{
				"OS_USERNAME": "testuser", "OS_PASSWORD": "testpass",
				"OS_PROJECT_ID": "project-id", "OS_PROJECT_NAME": "project-name",
				"OS_TENANT_ID": "tenant-id", "OS_TENANT_NAME": "tenant-name",
			},
			expected: auth.V2PasswordOpts{Username: "testuser", Password: "testpass", TenantID: "project-id", TenantName: "project-name", AllowReauth: true},
			expectedBody: map[string]map[string]any{"password": {
				"tenantId": "project-id", "tenantName": "project-name",
				"passwordCredentials": map[string]any{"username": "testuser", "password": "testpass"},
			}},
		},
		{
			name:     "token project scope",
			authType: "token",
			env:      map[string]string{"OS_TOKEN": "testtoken", "OS_PROJECT_ID": "project-id", "OS_PROJECT_NAME": "project-name"},
			expected: auth.V2TokenOpts{Token: "testtoken", TenantID: "project-id", TenantName: "project-name", AllowReauth: true},
			expectedBody: map[string]map[string]any{"token": {
				"tenantId": "project-id", "tenantName": "project-name", "token": map[string]any{"id": "testtoken"},
			}},
		},
		{
			name:     "token tenant scope",
			authType: "token",
			env:      map[string]string{"OS_TOKEN": "testtoken", "OS_TENANT_ID": "tenant-id", "OS_TENANT_NAME": "tenant-name"},
			expected: auth.V2TokenOpts{Token: "testtoken", TenantID: "tenant-id", TenantName: "tenant-name", AllowReauth: true},
			expectedBody: map[string]map[string]any{"token": {
				"tenantId": "tenant-id", "tenantName": "tenant-name", "token": map[string]any{"id": "testtoken"},
			}},
		},
		{
			name:     "token project scope takes precedence",
			authType: "token",
			env: map[string]string{
				"OS_TOKEN":      "testtoken",
				"OS_PROJECT_ID": "project-id", "OS_PROJECT_NAME": "project-name",
				"OS_TENANT_ID": "tenant-id", "OS_TENANT_NAME": "tenant-name",
			},
			expected: auth.V2TokenOpts{Token: "testtoken", TenantID: "project-id", TenantName: "project-name", AllowReauth: true},
			expectedBody: map[string]map[string]any{"token": {
				"tenantId": "project-id", "tenantName": "project-name", "token": map[string]any{"id": "testtoken"},
			}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			defer CleanupEnv(t)
			CleanupEnv(t)
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()
			fakeServer.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				th.TestMethod(t, r, http.MethodGet)
				fmt.Fprintf(w, `{"versions":{"values":[{"id":"v2.0","status":"stable","links":[{"href":%q,"rel":"self"}]}]}}`, fakeServer.Endpoint()+"v2.0/")
			})
			var authBody map[string]any
			for _, authBody = range testCase.expectedBody {
			}
			expectedJSON, err := json.Marshal(map[string]any{"auth": authBody})
			th.AssertNoErr(t, err)
			fakeServer.Mux.HandleFunc("/v2.0/tokens", func(w http.ResponseWriter, r *http.Request) {
				th.TestMethod(t, r, http.MethodPost)
				th.TestJSONRequest(t, r, string(expectedJSON))
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{"access":{"token":{"id":"discovered"}}}`)
			})
			t.Setenv("OS_AUTH_URL", fakeServer.Endpoint())
			t.Setenv("OS_AUTH_TYPE", testCase.authType)
			for key, value := range testCase.env {
				t.Setenv(key, value)
			}

			opts, err := auth.AuthOptionsFromEnv()
			th.AssertNoErr(t, err)
			result, err := opts.Authenticate(context.Background(), nil)
			th.AssertNoErr(t, err)
			th.AssertEquals(t, "discovered", result.TokenID)
			th.AssertEquals(t, testCase.expected.CanReauth(), result.CanReauth)
		})
	}
}

func TestAuthOptionsFromEnvGenericTokenResolvesV3(t *testing.T) {
	defer CleanupEnv(t)
	CleanupEnv(t)
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fakeServer.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodGet)
		fmt.Fprintf(w, `{"versions":{"values":[{"id":"v3.0","status":"stable","links":[{"href":%q,"rel":"self"}]}]}}`, fakeServer.Endpoint()+"v3/")
	})
	fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodPost)
		th.TestJSONRequest(t, r, `{"auth":{"identity":{"methods":["token"],"token":{"id":"testtoken"}}}}`)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Subject-Token", "discovered")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"token":{}}`)
	})
	t.Setenv("OS_AUTH_URL", fakeServer.Endpoint())
	t.Setenv("OS_AUTH_TYPE", "token")
	t.Setenv("OS_IDENTITY_API_VERSION", "2.0")
	t.Setenv("OS_TOKEN", "testtoken")

	opts, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	result, err := opts.Authenticate(context.Background(), nil)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "discovered", result.TokenID)
}

func TestAuthOptionsFromEnvV3Password(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV3Password(t)

	opts, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	v3Opts, ok := opts.(auth.AuthOptionsV3)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000/v3", v3Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V3PasswordOpts{Username: "testuser", Password: "testpass", Scope: &auth.Scope{}, AllowReauth: true}, v3Opts.Auth)
}

func TestAuthOptionsFromEnvV3Token(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV3Token(t)

	opts, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	v3Opts, ok := opts.(auth.AuthOptionsV3)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000/v3", v3Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V3TokenOpts{Token: "testtoken", Scope: &auth.Scope{}}, v3Opts.Auth)
}

func TestAuthOptionsFromEnvV3TOTP(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV3TOTP(t)

	opts, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	v3Opts, ok := opts.(auth.AuthOptionsV3)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000/v3", v3Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V3TOTPOpts{Passcode: "123456", Scope: &auth.Scope{}}, v3Opts.Auth)
}

func TestAuthOptionsFromEnvV3AppCred(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV3AppCred(t)

	opts, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	v3Opts, ok := opts.(auth.AuthOptionsV3)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000/v3", v3Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V3ApplicationCredentialOpts{ApplicationCredentialID: "app-cred-id", AllowReauth: true}, v3Opts.Auth)
}

func TestAuthOptionsFromEnvV3AppCredName(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV3AppCredName(t)

	opts, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	v3Opts, ok := opts.(auth.AuthOptionsV3)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000/v3", v3Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V3ApplicationCredentialOpts{ApplicationCredentialName: "app-cred-name", AllowReauth: true}, v3Opts.Auth)
}

func TestAuthOptionsFromEnvV3ExplicitAuthType(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV3ExplicitAuthType(t)

	opts, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	v3Opts, ok := opts.(auth.AuthOptionsV3)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000/v3", v3Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V3PasswordOpts{Username: "testuser", Password: "testpass", Scope: &auth.Scope{}, AllowReauth: true}, v3Opts.Auth)
}

func TestAuthOptionsFromEnvV3WithAuthMethods(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV3WithAuthMethods(t)

	opts, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	v3Opts, ok := opts.(auth.AuthOptionsV3)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000/v3", v3Opts.AuthURL)
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
}

func TestAuthOptionsFromEnvMissingAuthURL(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvMissingAuthURL(t)

	_, err := auth.AuthOptionsFromEnv()
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrMissingEnvironmentVariable)
	th.AssertEquals(t, true, ok)
}

func TestAuthOptionsFromEnvUnsupportedAuthType(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvUnsupportedAuthType(t)

	_, err := auth.AuthOptionsFromEnv()
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrUnsupportedAuthType)
	th.AssertEquals(t, true, ok)
}

func TestAuthOptionsFromEnvNoCredentials(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvNoCredentials(t)

	_, err := auth.AuthOptionsFromEnv()
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrUnsupportedAuthType)
	th.AssertEquals(t, true, ok)
}

func TestAuthOptionsFromEnvV2Direct(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV2Password(t)

	opts, err := auth.AuthOptionsFromEnvV2()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "http://example.com:5000/v2.0", opts.AuthURL)
	th.AssertDeepEquals(t, auth.V2PasswordOpts{Username: "testuser", Password: "testpass", AllowReauth: true}, opts.Auth)
}

func TestAuthOptionsFromEnvIdentityV2UsesTenantScope(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV2Password(t)
	t.Setenv("OS_PROJECT_ID", "project-id")
	t.Setenv("OS_PROJECT_NAME", "project-name")
	t.Setenv("OS_TENANT_ID", "tenant-id")
	t.Setenv("OS_TENANT_NAME", "tenant-name")

	opts, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)
	v2Opts, ok := opts.(auth.AuthOptionsV2)
	th.AssertEquals(t, true, ok)
	th.AssertDeepEquals(t, auth.V2PasswordOpts{
		Username: "testuser", Password: "testpass",
		TenantID: "tenant-id", TenantName: "tenant-name", AllowReauth: true,
	}, v2Opts.Auth)
}

func TestAuthOptionsFromEnvV2UsesTenantScope(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV2Password(t)
	t.Setenv("OS_PROJECT_ID", "project-id")
	t.Setenv("OS_PROJECT_NAME", "project-name")
	t.Setenv("OS_TENANT_ID", "tenant-id")
	t.Setenv("OS_TENANT_NAME", "tenant-name")

	opts, err := auth.AuthOptionsFromEnvV2()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, auth.V2PasswordOpts{
		Username: "testuser", Password: "testpass",
		TenantID: "tenant-id", TenantName: "tenant-name", AllowReauth: true,
	}, opts.Auth)

	body, err := opts.Auth.ToAuthBody()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]map[string]any{"password": {
		"tenantId": "tenant-id", "tenantName": "tenant-name",
		"passwordCredentials": map[string]any{"username": "testuser", "password": "testpass"},
	}}, body)
}

func TestAuthOptionsFromEnvV3Direct(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV3Password(t)

	opts, err := auth.AuthOptionsFromEnvV3()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "http://example.com:5000/v3", opts.AuthURL)
	th.AssertDeepEquals(t, auth.V3PasswordOpts{Username: "testuser", Password: "testpass", Scope: &auth.Scope{}, AllowReauth: true}, opts.Auth)
}

func TestAuthOptionsFromEnvV3SystemScope(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV3Password(t)
	t.Setenv("OS_SYSTEM_SCOPE", "all")

	opts, err := auth.AuthOptionsFromEnvV3()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, auth.V3PasswordOpts{
		Username:    "testuser",
		Password:    "testpass",
		Scope:       &auth.Scope{System: true},
		AllowReauth: true,
	}, opts.Auth)
}

func TestAuthOptionsFromEnvV3InvalidSystemScope(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV3Password(t)
	t.Setenv("OS_SYSTEM_SCOPE", "domain")

	_, err := auth.AuthOptionsFromEnvV3()
	th.AssertErr(t, err)
	_, ok := err.(gophercloud.ErrInvalidSystemScope)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "Only system scope of all is supported", err.Error())
}

func TestAuthOptionsFromEnvV3SystemScopeConflict(t *testing.T) {
	testCases := map[string]string{
		"OS_PROJECT_ID":   "project-id",
		"OS_PROJECT_NAME": "project-name",
		"OS_DOMAIN_ID":    "domain-id",
		"OS_DOMAIN_NAME":  "domain-name",
	}

	for envVar, value := range testCases {
		t.Run(envVar, func(t *testing.T) {
			defer CleanupEnv(t)
			SetupEnvV3Password(t)
			t.Setenv("OS_SYSTEM_SCOPE", "all")
			t.Setenv(envVar, value)

			opts, err := auth.AuthOptionsFromEnvV3()
			th.AssertNoErr(t, err)

			_, err = opts.Auth.ToAuthScope()
			th.AssertErr(t, err)
			_, ok := err.(gophercloud.ErrSystemScopeAndProjectOrDomainScope)
			th.AssertEquals(t, true, ok)
		})
	}
}

func TestAuthOptionsFromEnvV2PasswordCanReauth(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV2Password(t)

	opts, err := auth.AuthOptionsFromEnvV2()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, true, opts.Auth.CanReauth())
}

func TestAuthOptionsFromEnvV2TokenCanReauth(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV2Token(t)

	opts, err := auth.AuthOptionsFromEnvV2()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, true, opts.Auth.CanReauth())
}

func TestAuthOptionsFromEnvV3PasswordCanReauth(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV3Password(t)

	opts, err := auth.AuthOptionsFromEnvV3()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, true, opts.Auth.CanReauth())
}

func TestAuthOptionsFromEnvV3AppCredCanReauth(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV3AppCred(t)

	opts, err := auth.AuthOptionsFromEnvV3()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, true, opts.Auth.CanReauth())
}
