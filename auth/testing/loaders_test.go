package testing

import (
	"strings"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// TestAuthOptionsLoadersGenerateSameRequest checks that equivalent
// environment and cloud configurations produce the same auth body and scope.
func TestAuthOptionsLoadersGenerateSameRequest(t *testing.T) {
	credentials := map[string]string{
		"username":       "testuser",
		"password":       "testpass",
		"user_domain_id": "user-domain",
	}

	testCases := []struct {
		name      string
		settings  map[string]string
		wantScope string
	}{
		{
			name:      "domain scope",
			settings:  map[string]string{"domain_id": "scope-domain"},
			wantScope: `{"domain":{"id":"scope-domain"}}`,
		},
		{
			name:      "domain name scope",
			settings:  map[string]string{"domain_name": "scope-domain"},
			wantScope: `{"domain":{"name":"scope-domain"}}`,
		},
		{
			name:      "project scope",
			settings:  map[string]string{"project_id": "project-id"},
			wantScope: `{"project":{"id":"project-id"}}`,
		},
		{
			name:      "project name scope",
			settings:  map[string]string{"project_name": "project-name", "project_domain_id": "project-domain"},
			wantScope: `{"project":{"name":"project-name","domain":{"id":"project-domain"}}}`,
		},
		{
			name:      "system scope",
			settings:  map[string]string{"system_scope": "all"},
			wantScope: `{"system":{"all":true}}`,
		},
		{
			name:      "trust scope",
			settings:  map[string]string{"trust_id": "trust-id"},
			wantScope: `{"OS-TRUST:trust":{"id":"trust-id"}}`,
		},
		{
			name:     "unscoped",
			settings: map[string]string{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			defer CleanupEnv(t)
			CleanupEnv(t)

			cloudData := map[string]any{"auth_url": "http://example.com:5000/v3"}
			t.Setenv("OS_AUTH_URL", "http://example.com:5000/v3")
			for _, settings := range []map[string]string{credentials, testCase.settings} {
				for k, v := range settings {
					cloudData[k] = v
					t.Setenv("OS_"+strings.ToUpper(k), v)
				}
			}

			envOpts, err := auth.AuthOptionsFromEnvV3()
			th.AssertNoErr(t, err)
			cloudOpts, err := auth.AuthOptionsFromCloudV3(cloudSource{authType: auth.AuthV3Password, authData: cloudData})
			th.AssertNoErr(t, err)

			envBody, err := envOpts.Auth.ToAuthBody()
			th.AssertNoErr(t, err)
			cloudBody, err := cloudOpts.Auth.ToAuthBody()
			th.AssertNoErr(t, err)
			th.AssertJSONEquals(t, `{"password":{"user":{"name":"testuser","password":"testpass","domain":{"id":"user-domain"}}}}`, envBody)
			th.AssertJSONEquals(t, `{"password":{"user":{"name":"testuser","password":"testpass","domain":{"id":"user-domain"}}}}`, cloudBody)

			envScope, err := envOpts.Auth.ToAuthScope()
			th.AssertNoErr(t, err)
			cloudScope, err := cloudOpts.Auth.ToAuthScope()
			th.AssertNoErr(t, err)
			if testCase.wantScope == "" {
				th.AssertEquals(t, 0, len(envScope))
				th.AssertEquals(t, 0, len(cloudScope))
				return
			}
			th.AssertJSONEquals(t, testCase.wantScope, envScope)
			th.AssertJSONEquals(t, testCase.wantScope, cloudScope)
		})
	}
}

// TestAuthOptionsFromCloudV3DomainWithOtherScope checks that the generic
// domain is a domain scope, so combining it with another scope is an error.
func TestAuthOptionsFromCloudV3DomainWithOtherScope(t *testing.T) {
	testCases := []struct {
		name    string
		data    map[string]any
		checkFn func(error) bool
	}{
		{
			name: "project scope",
			data: map[string]any{"domain_name": "Default", "project_name": "project-name"},
			checkFn: func(err error) bool {
				_, ok := err.(gophercloud.ErrProjectScopeAndDomainScope)
				return ok
			},
		},
		{
			name: "system scope",
			data: map[string]any{"domain_name": "Default", "system_scope": "all"},
			checkFn: func(err error) bool {
				_, ok := err.(gophercloud.ErrSystemScopeAndProjectOrDomainScope)
				return ok
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.data["auth_url"] = "http://example.com:5000/v3"
			testCase.data["username"] = "testuser"
			testCase.data["password"] = "testpass"

			opts, err := auth.AuthOptionsFromCloudV3(cloudSource{authData: testCase.data})
			th.AssertNoErr(t, err)

			_, err = opts.Auth.ToAuthScope()
			th.AssertErr(t, err)
			th.AssertEquals(t, true, testCase.checkFn(err))
		})
	}
}

// TestAuthOptionsLoadersV3OAuth1 checks that the environment and cloud
// loaders build OAuth1 options that sign requests the same way, using the
// keystoneauth option names.
func TestAuthOptionsLoadersV3OAuth1(t *testing.T) {
	defer CleanupEnv(t)
	CleanupEnv(t)
	t.Setenv("OS_AUTH_URL", "http://127.0.0.1:33199/v3")
	t.Setenv("OS_AUTH_TYPE", "v3oauth1")
	t.Setenv("OS_CONSUMER_KEY", "7fea2d")
	t.Setenv("OS_CONSUMER_SECRET", "secretsecret")
	t.Setenv("OS_ACCESS_KEY", "accd36")
	t.Setenv("OS_ACCESS_SECRET", "aa47da")

	envOpts, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)
	cloudOpts, err := auth.AuthOptionsFromCloud(cloudSource{authType: auth.AuthV3OAuth1, authData: map[string]any{
		"auth_url":        "http://127.0.0.1:33199/v3",
		"consumer_key":    "7fea2d",
		"consumer_secret": "secretsecret",
		"access_key":      "accd36",
		"access_secret":   "aa47da",
	}})
	th.AssertNoErr(t, err)

	// Same request and signature as TestV3OAuth1OptsBuildRequest.
	want := `OAuth oauth_consumer_key="7fea2d", oauth_nonce="66148873158553341551586804894", oauth_signature_method="HMAC-SHA1", oauth_timestamp="0", oauth_token="accd36", oauth_version="1.0", oauth_signature="b4fZNXFnxnmjPa9lGMUrLOYlznM%3D"`
	timestamp := time.Unix(0, 0)
	for name, opts := range map[string]auth.Authenticator{"env": envOpts, "cloud": cloudOpts} {
		t.Run(name, func(t *testing.T) {
			v3Opts, ok := opts.(auth.AuthOptionsV3)
			th.AssertEquals(t, true, ok)
			oauth1 := v3Opts.Auth.(auth.V3OAuth1Opts)
			th.AssertEquals(t, true, oauth1.CanReauth())
			oauth1.Timestamp = &timestamp
			oauth1.Nonce = "66148873158553341551586804894"

			request, err := auth.NewRequestV3(oauth1, auth.WithTokenURL("http://127.0.0.1:33199/v3/auth/tokens"))
			th.AssertNoErr(t, err)
			th.AssertEquals(t, want, request.MoreHeaders["Authorization"])
		})
	}
}

// TestAuthOptionsLoadersV3OAuth1InMultifactor checks that, as in
// keystoneauth, OAuth1 can be combined with other methods, and that the
// request carries both the method bodies and the OAuth1 signature.
func TestAuthOptionsLoadersV3OAuth1InMultifactor(t *testing.T) {
	defer CleanupEnv(t)
	CleanupEnv(t)
	t.Setenv("OS_AUTH_URL", "http://127.0.0.1:33199/v3")
	t.Setenv("OS_AUTH_TYPE", "v3multifactor")
	t.Setenv("OS_AUTH_METHODS", "v3password,v3oauth1")
	t.Setenv("OS_USERID", "user-id")
	t.Setenv("OS_PASSWORD", "testpass")
	t.Setenv("OS_CONSUMER_KEY", "7fea2d")
	t.Setenv("OS_CONSUMER_SECRET", "secretsecret")
	t.Setenv("OS_ACCESS_KEY", "accd36")
	t.Setenv("OS_ACCESS_SECRET", "aa47da")

	envOpts, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)
	cloudOpts, err := auth.AuthOptionsFromCloud(cloudSource{authType: auth.AuthV3MultiFactor, authData: map[string]any{
		"auth_url":        "http://127.0.0.1:33199/v3",
		"auth_methods":    []any{"v3password", "v3oauth1"},
		"user_id":         "user-id",
		"password":        "testpass",
		"consumer_key":    "7fea2d",
		"consumer_secret": "secretsecret",
		"access_key":      "accd36",
		"access_secret":   "aa47da",
	}})
	th.AssertNoErr(t, err)

	// Same signature as TestV3OAuth1OptsBuildRequest.
	wantHeader := `OAuth oauth_consumer_key="7fea2d", oauth_nonce="66148873158553341551586804894", oauth_signature_method="HMAC-SHA1", oauth_timestamp="0", oauth_token="accd36", oauth_version="1.0", oauth_signature="b4fZNXFnxnmjPa9lGMUrLOYlznM%3D"`
	timestamp := time.Unix(0, 0)
	for name, opts := range map[string]auth.Authenticator{"env": envOpts, "cloud": cloudOpts} {
		t.Run(name, func(t *testing.T) {
			multifactor := opts.(auth.AuthOptionsV3).Auth.(auth.V3MultifactorOpts)
			th.AssertEquals(t, 2, len(multifactor.AuthMethods))
			oauth1 := multifactor.AuthMethods[1].(auth.V3OAuth1Opts)
			oauth1.Timestamp = &timestamp
			oauth1.Nonce = "66148873158553341551586804894"
			multifactor.AuthMethods[1] = oauth1

			request, err := auth.NewRequestV3(multifactor, auth.WithTokenURL("http://127.0.0.1:33199/v3/auth/tokens"))
			th.AssertNoErr(t, err)
			th.AssertJSONEquals(t, `{"auth":{"identity":{"methods":["oauth1","password"],"oauth1":{},"password":{"user":{"id":"user-id","password":"testpass"}}}}}`, request.JSONBody)
			th.AssertEquals(t, wantHeader, request.MoreHeaders["Authorization"])
		})
	}
}

// TestAuthOptionsLoadersIdentityAPIVersionFallback checks which identity
// version is used when the auth type does not select one.
func TestAuthOptionsLoadersIdentityAPIVersionFallback(t *testing.T) {
	testCases := []struct {
		version string
		wantV2  bool
	}{
		{version: "2", wantV2: true},
		{version: "2.0", wantV2: true},
		{version: "3", wantV2: false},
		{version: "", wantV2: false},
	}

	for _, testCase := range testCases {
		t.Run("version "+testCase.version, func(t *testing.T) {
			defer CleanupEnv(t)
			CleanupEnv(t)
			t.Setenv("OS_AUTH_URL", "http://example.com:5000")
			t.Setenv("OS_IDENTITY_API_VERSION", testCase.version)
			t.Setenv("OS_USERNAME", "testuser")
			t.Setenv("OS_PASSWORD", "testpass")

			envOpts, err := auth.AuthOptionsFromEnv()
			th.AssertNoErr(t, err)
			cloudOpts, err := auth.AuthOptionsFromCloud(cloudSource{identityAPIVersion: testCase.version, authData: map[string]any{
				"auth_url": "http://example.com:5000", "username": "testuser", "password": "testpass",
			}})
			th.AssertNoErr(t, err)

			_, envV2 := envOpts.(auth.AuthOptionsV2)
			_, cloudV2 := cloudOpts.(auth.AuthOptionsV2)
			th.AssertEquals(t, testCase.wantV2, envV2)
			th.AssertEquals(t, testCase.wantV2, cloudV2)
		})
	}
}

// TestAuthOptionsLoadersExplicitV3AuthTypeIgnoresIdentityAPIVersion checks
// that an explicit v3 auth type is not sent to identity v2.
func TestAuthOptionsLoadersExplicitV3AuthTypeIgnoresIdentityAPIVersion(t *testing.T) {
	for _, authType := range []auth.AuthType{
		auth.AuthV3Password, auth.AuthV3Totp, auth.AuthV3Token,
		auth.AuthV3ApplicationCredential, auth.AuthV3MultiFactor, auth.AuthV3OAuth1,
	} {
		t.Run(string(authType), func(t *testing.T) {
			defer CleanupEnv(t)
			CleanupEnv(t)
			t.Setenv("OS_AUTH_URL", "http://example.com:5000/v3")
			t.Setenv("OS_AUTH_TYPE", string(authType))
			t.Setenv("OS_AUTH_METHODS", "v3password")
			t.Setenv("OS_IDENTITY_API_VERSION", "2")

			envOpts, err := auth.AuthOptionsFromEnv()
			th.AssertNoErr(t, err)
			cloudOpts, err := auth.AuthOptionsFromCloud(cloudSource{authType: authType, identityAPIVersion: "2", authData: map[string]any{
				"auth_url": "http://example.com:5000/v3", "auth_methods": []any{"v3password"},
			}})
			th.AssertNoErr(t, err)

			_, envV3 := envOpts.(auth.AuthOptionsV3)
			_, cloudV3 := cloudOpts.(auth.AuthOptionsV3)
			th.AssertEquals(t, true, envV3)
			th.AssertEquals(t, true, cloudV3)
		})
	}
}
