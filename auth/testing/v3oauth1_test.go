package testing

import (
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/auth"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestV3OAuth1OptsRequiresCredentials(t *testing.T) {
	tests := []struct {
		name string
		opts auth.V3OAuth1Opts
	}{
		{name: "consumer key", opts: auth.V3OAuth1Opts{ConsumerSecret: "secret", Token: "token", TokenSecret: "secret", SignatureMethod: auth.OAuth1HMACSHA1}},
		{name: "consumer secret", opts: auth.V3OAuth1Opts{ConsumerKey: "key", Token: "token", TokenSecret: "secret", SignatureMethod: auth.OAuth1HMACSHA1}},
		{name: "token", opts: auth.V3OAuth1Opts{ConsumerKey: "key", ConsumerSecret: "secret", TokenSecret: "secret", SignatureMethod: auth.OAuth1HMACSHA1}},
		{name: "token secret", opts: auth.V3OAuth1Opts{ConsumerKey: "key", ConsumerSecret: "secret", Token: "token", SignatureMethod: auth.OAuth1HMACSHA1}},
		{name: "signature method", opts: auth.V3OAuth1Opts{ConsumerKey: "key", ConsumerSecret: "secret", Token: "token", TokenSecret: "secret"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.opts.ToAuthBody()
			th.AssertErr(t, err)
		})
	}
}

func TestV3OAuth1OptsBuildRequest(t *testing.T) {
	timestamp := time.Unix(0, 0)
	opts := auth.V3OAuth1Opts{
		ConsumerKey:     "7fea2d",
		ConsumerSecret:  "secretsecret",
		Token:           "accd36",
		TokenSecret:     "aa47da",
		SignatureMethod: auth.OAuth1HMACSHA1,
		Timestamp:       &timestamp,
		Nonce:           "66148873158553341551586804894",
		AllowReauth:     true,
	}

	body, headers, err := v3AuthRequest(t, opts)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]any{
		"auth": map[string]any{
			"identity": map[string]any{
				"methods": []string{"oauth1"},
				"oauth1":  map[string]any{},
			},
		},
	}, body)
	th.AssertEquals(t, `OAuth oauth_consumer_key="7fea2d", oauth_nonce="66148873158553341551586804894", oauth_signature_method="HMAC-SHA1", oauth_timestamp="0", oauth_token="accd36", oauth_version="1.0", oauth_signature="b4fZNXFnxnmjPa9lGMUrLOYlznM%3D"`, headers.Get("Authorization"))
	th.AssertEquals(t, auth.AuthV3OAuth1, opts.ToAuthType())
	th.AssertEquals(t, true, opts.CanReauth())
}

func TestV3OAuth1OptsBuildsHeadersWithRequestOptions(t *testing.T) {
	timestamp := time.Unix(0, 0)
	opts := auth.V3OAuth1Opts{
		ConsumerKey:     "7fea2d",
		ConsumerSecret:  "secretsecret",
		Token:           "accd36",
		TokenSecret:     "aa47da",
		SignatureMethod: auth.OAuth1HMACSHA1,
		Timestamp:       &timestamp,
		Nonce:           "66148873158553341551586804894",
	}

	_, headers, err := v3AuthRequest(t, opts)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, `OAuth oauth_consumer_key="7fea2d", oauth_nonce="66148873158553341551586804894", oauth_signature_method="HMAC-SHA1", oauth_timestamp="0", oauth_token="accd36", oauth_version="1.0", oauth_signature="b4fZNXFnxnmjPa9lGMUrLOYlznM%3D"`, headers.Get("Authorization"))
}
