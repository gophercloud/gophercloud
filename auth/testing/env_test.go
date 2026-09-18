package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

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
			expected: auth.V2PasswordOpts{Username: "testuser", Password: "testpass", AllowReauth: true},
		},
		{
			name:     "token",
			authType: "v2token",
			setCreds: func(t *testing.T) {
				t.Setenv("OS_TOKEN", "testtoken")
			},
			expected: auth.V2TokenOpts{Token: "testtoken", AllowReauth: true},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			defer CleanupEnv(t)
			CleanupEnv(t)
			t.Setenv("OS_AUTH_URL", "http://example.com:5000/v2.0")
			t.Setenv("OS_AUTH_TYPE", testCase.authType)
			t.Setenv("OS_IDENTITY_API_VERSION", "3")
			testCase.setCreds(t)

			opts, err := auth.AuthOptionsFromEnv()
			th.AssertNoErr(t, err)

			v2Opts, ok := opts.(auth.AuthOptionsV2)
			th.AssertEquals(t, true, ok)
			th.AssertDeepEquals(t, testCase.expected, v2Opts.Auth)
		})
	}
}

func TestAuthOptionsFromEnvGenericPasswordResolvesV2(t *testing.T) {
	defer CleanupEnv(t)
	CleanupEnv(t)
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fakeServer.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		fmt.Fprintf(w, `{"versions":{"values":[{"id":"v2.0","status":"stable","links":[{"href":%q,"rel":"self"}]}]}}`, fakeServer.Endpoint()+"v2.0/")
	})
	t.Setenv("OS_AUTH_URL", fakeServer.Endpoint())
	t.Setenv("OS_AUTH_TYPE", "password")
	t.Setenv("OS_IDENTITY_API_VERSION", "3")
	t.Setenv("OS_USERNAME", "testuser")
	t.Setenv("OS_PASSWORD", "testpass")

	opts, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	v2Opts, ok := opts.(auth.AuthOptionsV2)
	th.AssertEquals(t, true, ok)
	th.AssertDeepEquals(t, auth.V2PasswordOpts{Username: "testuser", Password: "testpass", AllowReauth: true}, v2Opts.Auth)
}

func TestAuthOptionsFromEnvGenericTokenResolvesV3(t *testing.T) {
	defer CleanupEnv(t)
	CleanupEnv(t)
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fakeServer.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		fmt.Fprintf(w, `{"versions":{"values":[{"id":"v3.0","status":"stable","links":[{"href":%q,"rel":"self"}]}]}}`, fakeServer.Endpoint()+"v3/")
	})
	t.Setenv("OS_AUTH_URL", fakeServer.Endpoint())
	t.Setenv("OS_AUTH_TYPE", "token")
	t.Setenv("OS_IDENTITY_API_VERSION", "2.0")
	t.Setenv("OS_TOKEN", "testtoken")

	opts, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	v3Opts, ok := opts.(auth.AuthOptionsV3)
	th.AssertEquals(t, true, ok)
	th.AssertDeepEquals(t, auth.V3TokenOpts{Token: "testtoken", Scope: &auth.Scope{}}, v3Opts.Auth)
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

func TestAuthOptionsFromEnvV3Direct(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV3Password(t)

	opts, err := auth.AuthOptionsFromEnvV3()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "http://example.com:5000/v3", opts.AuthURL)
	th.AssertDeepEquals(t, auth.V3PasswordOpts{Username: "testuser", Password: "testpass", Scope: &auth.Scope{}, AllowReauth: true}, opts.Auth)
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
