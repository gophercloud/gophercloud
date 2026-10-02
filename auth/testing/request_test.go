package testing

import (
	"errors"
	"testing"

	"github.com/gophercloud/gophercloud/v2/auth"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

type requestBuilder struct {
	body       map[string]map[string]any
	scope      map[string]any
	headers    map[string]any
	bodyErr    error
	scopeErr   error
	headersErr error
}

func (opts requestBuilder) ToAuthBody() (map[string]map[string]any, error) {
	return opts.body, opts.bodyErr
}

func (opts requestBuilder) ToAuthScope() (map[string]any, error) {
	return opts.scope, opts.scopeErr
}

func (opts requestBuilder) ToAuthHeaders(_ ...auth.RequestOption) (map[string]any, error) {
	return opts.headers, opts.headersErr
}

func (requestBuilder) ToAuthType() auth.AuthType {
	return auth.AuthV3MultiFactor
}

func (requestBuilder) CanReauth() bool {
	return false
}

func TestBuildV2Request(t *testing.T) {
	actual, err := auth.NewRequestV2(requestBuilder{
		body: map[string]map[string]any{
			"password": {
				"tenantId": "tenant",
				"passwordCredentials": map[string]any{
					"username": "user",
					"password": "secret",
				},
			},
		},
	})
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]any{
		"auth": map[string]any{
			"tenantId": "tenant",
			"passwordCredentials": map[string]any{
				"username": "user",
				"password": "secret",
			},
		},
	}, actual.JSONBody)
	th.AssertDeepEquals(t, []int{200, 203}, actual.OkCodes)
	th.AssertDeepEquals(t, []string{"X-Auth-Token", "Authorization"}, actual.OmitHeaders)
}

func TestBuildV3Request(t *testing.T) {
	opts := requestBuilder{
		body: map[string]map[string]any{
			"totp":     {"user": map[string]any{"id": "user", "passcode": "123456"}},
			"password": {"user": map[string]any{"id": "user", "password": "secret"}},
		},
		scope:   map[string]any{"project": map[string]any{"id": "project"}},
		headers: map[string]any{"OpenStack-Identity-Access-Rules": 1},
	}

	actual, err := auth.NewRequestV3(opts)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]any{
		"auth": map[string]any{
			"identity": map[string]any{
				"methods":  []string{"password", "totp"},
				"password": map[string]any{"user": map[string]any{"id": "user", "password": "secret"}},
				"totp":     map[string]any{"user": map[string]any{"id": "user", "passcode": "123456"}},
			},
			"scope": map[string]any{"project": map[string]any{"id": "project"}},
		},
	}, actual.JSONBody)
	th.AssertDeepEquals(t, map[string]string{"OpenStack-Identity-Access-Rules": "1"}, actual.MoreHeaders)
	th.AssertDeepEquals(t, []string{"X-Auth-Token", "Authorization"}, actual.OmitHeaders)
}

func TestBuildV3RequestPreservesBuilderAuthorization(t *testing.T) {
	actual, err := auth.NewRequestV3(requestBuilder{
		body:    map[string]map[string]any{"oauth1": {"access_token_id": "token"}},
		headers: map[string]any{"Authorization": "OAuth signed-request"},
	})
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]string{"Authorization": "OAuth signed-request"}, actual.MoreHeaders)
	th.AssertDeepEquals(t, []string{"X-Auth-Token"}, actual.OmitHeaders)
}

func TestBuildV3RequestWithoutScope(t *testing.T) {
	actual, err := auth.NewRequestV3(requestBuilder{
		body: map[string]map[string]any{
			"token": {"id": "token"},
		},
	})
	th.AssertNoErr(t, err)

	authBody := actual.JSONBody.(map[string]any)["auth"].(map[string]any)
	if _, ok := authBody["scope"]; ok {
		t.Fatal("request contains an unexpected scope")
	}
}

func TestBuildV3RequestPropagatesErrors(t *testing.T) {
	bodyErr := errors.New("body error")
	scopeErr := errors.New("scope error")
	headersErr := errors.New("headers error")

	tests := []struct {
		name string
		opts requestBuilder
		want error
	}{
		{
			name: "body",
			opts: requestBuilder{bodyErr: bodyErr},
			want: bodyErr,
		},
		{
			name: "scope",
			opts: requestBuilder{
				body:     map[string]map[string]any{"token": {"id": "token"}},
				scopeErr: scopeErr,
			},
			want: scopeErr,
		},
		{
			name: "headers",
			opts: requestBuilder{
				body:       map[string]map[string]any{"token": {"id": "token"}},
				headersErr: headersErr,
			},
			want: headersErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := auth.NewRequestV3(tt.opts)
			if !errors.Is(err, tt.want) {
				t.Fatalf("expected %v, got %v", tt.want, err)
			}
		})
	}
}
