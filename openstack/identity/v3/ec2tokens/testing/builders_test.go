package testing

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/ec2tokens"
	tokens_testing "github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens/testing"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestValidateS3Token(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	c := client.ServiceClient(fakeServer)
	c.TokenID = "service-token"
	fakeServer.Mux.HandleFunc("/s3tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodPost)
		th.TestHeader(t, r, "X-Auth-Token", "service-token")
		th.TestJSONRequest(t, r, `{"credentials":{"access":"access-id","signature":"signature","token":"dG9rZW4="}}`)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, tokens_testing.TokenOutput)
	})
	opts := s3BodyBuilder{body: map[string]any{
		"credentials": map[string]any{
			"access": "access-id", "signature": "signature", "token": "dG9rZW4=",
		},
	}}

	_, err := ec2tokens.ValidateS3Token(context.Background(), c, opts).Extract()
	th.AssertNoErr(t, err)
}

type s3BodyBuilder struct {
	body map[string]any
	err  error
}

func (opts s3BodyBuilder) ToS3TokenValidateMap() (map[string]any, error) {
	if opts.body != nil {
		return opts.body, opts.err
	}
	return map[string]any{"credentials": map[string]any{}}, opts.err
}

type s3HeadersBuilder struct {
	s3BodyBuilder
	err    error
	method string
	url    string
}

func (opts *s3HeadersBuilder) ToS3TokenValidateHeaders(method, url string) (map[string]string, error) {
	opts.method, opts.url = method, url
	return map[string]string{"X-Custom": "custom-value"}, opts.err
}

func TestValidateS3TokenHeaders(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	c := client.ServiceClient(fakeServer)
	fakeServer.Mux.HandleFunc("/s3tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestHeader(t, r, "X-Custom", "custom-value")
		th.TestJSONRequest(t, r, `{"credentials":{}}`)
		fmt.Fprint(w, tokens_testing.TokenOutput)
	})
	opts := &s3HeadersBuilder{}
	_, err := ec2tokens.ValidateS3Token(context.Background(), c, opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, http.MethodPost, opts.method)
	th.AssertEquals(t, c.ServiceURL("s3tokens"), opts.url)
}

func TestValidateS3TokenBuilderErrors(t *testing.T) {
	expected := errors.New("builder failure")
	c := &gophercloud.ServiceClient{Endpoint: "http://example.com/"}
	for _, opts := range []ec2tokens.ValidateS3TokenOptsBuilder{
		s3BodyBuilder{err: expected}, &s3HeadersBuilder{err: expected},
	} {
		result := ec2tokens.ValidateS3Token(context.Background(), c, opts)
		th.AssertEquals(t, expected, result.Err)
	}
}
