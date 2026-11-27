package testing

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	v2tokens "github.com/gophercloud/gophercloud/v2/openstack/identity/v2/tokens"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/ec2tokens"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/oauth1"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

type tokenBodyBuilder struct {
	err error
}

func (opts tokenBodyBuilder) ToTokenCreateMap() (map[string]any, error) {
	return map[string]any{"auth": map[string]any{"custom": true}}, opts.err
}

type tokenHeadersBuilder struct {
	tokenBodyBuilder
	headers map[string]string
	err     error
	method  string
	url     string
}

func (opts *tokenHeadersBuilder) ToTokenCreateHeaders(method, url string) (map[string]string, error) {
	opts.method, opts.url = method, url
	return opts.headers, opts.err
}

func TestCreateCustomBuilders(t *testing.T) {
	for _, withHeaders := range []bool{false, true} {
		name := "body only"
		if withHeaders {
			name = "body and headers"
		}
		t.Run(name, func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()
			c := client.ServiceClient(fakeServer)
			c.TokenID = "existing-token"
			fakeServer.Mux.HandleFunc("/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
				th.TestMethod(t, r, http.MethodPost)
				th.TestJSONRequest(t, r, `{"auth":{"custom":true}}`)
				th.AssertEquals(t, "", r.Header.Get("X-Auth-Token"))
				if withHeaders {
					th.TestHeader(t, r, "Authorization", "custom-signature")
					th.TestHeader(t, r, "X-Custom", "custom-value")
				} else {
					th.AssertEquals(t, "", r.Header.Get("Authorization"))
				}
				w.WriteHeader(http.StatusCreated)
				_, err := w.Write([]byte(`{"token":{}}`))
				th.CheckNoErr(t, err)
			})
			headers := &tokenHeadersBuilder{headers: map[string]string{
				"authorization": "custom-signature", "X-Custom": "custom-value",
			}}
			var opts tokens.CreateOptsBuilder = tokenBodyBuilder{}
			if withHeaders {
				opts = headers
			}
			result := tokens.Create(context.Background(), c, opts)
			th.AssertNoErr(t, result.Err)
			if withHeaders {
				th.AssertEquals(t, http.MethodPost, headers.method)
				th.AssertEquals(t, c.ServiceURL("auth", "tokens"), headers.url)
			}
		})
	}
}

func TestCreateBuilderErrors(t *testing.T) {
	expected := errors.New("builder failure")
	c := &gophercloud.ServiceClient{Endpoint: "http://example.com/"}
	ctx := context.Background()
	operations := map[string]func(tokens.CreateOptsBuilder) error{
		"v2":     func(opts tokens.CreateOptsBuilder) error { return v2tokens.Create(ctx, c, opts).Err },
		"v3":     func(opts tokens.CreateOptsBuilder) error { return tokens.Create(ctx, c, opts).Err },
		"ec2":    func(opts tokens.CreateOptsBuilder) error { return ec2tokens.Create(ctx, c, opts).Err },
		"oauth1": func(opts tokens.CreateOptsBuilder) error { return oauth1.Create(ctx, c, opts).Err },
	}
	for name, create := range operations {
		t.Run(name, func(t *testing.T) {
			for _, opts := range []tokens.CreateOptsBuilder{
				tokenBodyBuilder{err: expected},
				&tokenHeadersBuilder{err: expected},
			} {
				th.AssertEquals(t, expected, create(opts))
			}
		})
	}
}
