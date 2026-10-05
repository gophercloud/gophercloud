package testing

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// ec2V4Opts are EC2 credentials signed with a deterministic AWS signature V4.
func ec2V4Opts() auth.EC2TokenOpts {
	bodyHash := "foo"
	return auth.EC2TokenOpts{
		Access:    "a7f1e798b7c2417cba4a02de97dc3cdc",
		Secret:    "18f4f6761ada4e3795fa5273c30349b9",
		BodyHash:  &bodyHash,
		Timestamp: new(time.Time),
		Region:    "region1",
		Service:   "ec2",
		Path:      "/",
		Verb:      "GET",
		Headers:   map[string]string{"Host": "localhost"},
		Params:    map[string]string{"Action": "Test"},
	}
}

const ec2V4Signature = "f36f79118f75d7d6ec86ead9a61679cbdcf94c0cbfe5e9cf2407e8406aa82028"

const ec2V4Token = "QVdTNC1ITUFDLVNIQTI1NgowMDAxMDEwMVQwMDAwMDBaCjAwMDEwMTAxL3JlZ2lvbjEvZWMyL2F3czRfcmVxdWVzdAoyNTg2ODVjMTcyNTI5YTM4NzhiZDliMDk1NTE5ZWQxODE5YTJiMTQ4Y2VkODkyOTYwMmY1YmYxMDQ4NWE0MjJj"

func TestEC2TokenOptsToAuthBodyV2(t *testing.T) {
	opts := auth.EC2TokenOpts{
		Access: "a7f1e798b7c2417cba4a02de97dc3cdc",
		Secret: "18f4f6761ada4e3795fa5273c30349b9",
		Host:   "localhost",
		Path:   "/",
		Verb:   "GET",
		// body_hash and headers are unused by a signature V2 request
		BodyHash: new(string),
		Headers:  map[string]string{"Foo": "Bar"},
		Params: map[string]string{
			"Action":           "Test",
			"SignatureMethod":  "HmacSHA256",
			"SignatureVersion": "2",
		},
	}

	body, err := opts.ToAuthBody()
	th.AssertNoErr(t, err)
	th.AssertJSONEquals(t, `
		{
			"credentials": {
				"access": "a7f1e798b7c2417cba4a02de97dc3cdc",
				"host": "localhost",
				"path": "/",
				"verb": "GET",
				"params": {
					"Action": "Test",
					"SignatureMethod": "HmacSHA256",
					"SignatureVersion": "2"
				},
				"signature": "Up+MbVbbrvdR5FRkUz+n3nc+VW6xieuN50wh6ONEJ4w="
			}
		}
	`, body)
}

func TestEC2TokenOptsToAuthBodyV2HmacSha1(t *testing.T) {
	opts := auth.EC2TokenOpts{
		Access: "a7f1e798b7c2417cba4a02de97dc3cdc",
		Secret: "18f4f6761ada4e3795fa5273c30349b9",
		Host:   "localhost",
		Path:   "/",
		Verb:   "GET",
		Params: map[string]string{
			"Action":           "Test",
			"SignatureMethod":  "HmacSHA1",
			"SignatureVersion": "2",
		},
	}

	body, err := opts.ToAuthBody()
	th.AssertNoErr(t, err)
	th.AssertJSONEquals(t, `
		{
			"credentials": {
				"access": "a7f1e798b7c2417cba4a02de97dc3cdc",
				"host": "localhost",
				"path": "/",
				"verb": "GET",
				"params": {
					"Action": "Test",
					"SignatureMethod": "HmacSHA1",
					"SignatureVersion": "2"
				},
				"signature": "pCRaSBxB487K0ua5q6ud2iW0Y4I="
			}
		}
	`, body)
}

// TestEC2TokenOptsToAuthBodyV4 verifies that a signature V4 body carries the
// calculated signature, the AWS authorization headers and the signed string as
// a token. The token is only consumed by the s3tokens endpoint.
func TestEC2TokenOptsToAuthBodyV4(t *testing.T) {
	body, err := ec2V4Opts().ToAuthBody()
	th.AssertNoErr(t, err)
	th.AssertJSONEquals(t, fmt.Sprintf(`
		{
			"credentials": {
				"access": "a7f1e798b7c2417cba4a02de97dc3cdc",
				"body_hash": "foo",
				"host": "",
				"path": "/",
				"verb": "GET",
				"headers": {
					"Host": "localhost",
					"Authorization": "AWS4-HMAC-SHA256 Credential=a7f1e798b7c2417cba4a02de97dc3cdc/00010101/region1/ec2/aws4_request, SignedHeaders=, Signature=%s",
					"X-Amz-Date": "00010101T000000Z"
				},
				"params": {"Action": "Test"},
				"signature": "%s",
				"token": "%s"
			}
		}
	`, ec2V4Signature, ec2V4Signature, ec2V4Token), body)
}

// TestEC2TokenOptsToAuthBodyWithSignature verifies that a preset signature is
// used as is, without calculating one.
func TestEC2TokenOptsToAuthBodyWithSignature(t *testing.T) {
	opts := auth.EC2TokenOpts{
		Access:    "a7f1e798b7c2417cba4a02de97dc3cdc",
		Path:      "/",
		Verb:      "GET",
		BodyHash:  new(string),
		Signature: "presetsignature",
		Token:     []byte("presettoken"),
	}

	body, err := opts.ToAuthBody()
	th.AssertNoErr(t, err)
	th.AssertJSONEquals(t, `
		{
			"credentials": {
				"access": "a7f1e798b7c2417cba4a02de97dc3cdc",
				"body_hash": "",
				"host": "",
				"path": "/",
				"verb": "GET",
				"headers": null,
				"params": null,
				"signature": "presetsignature",
				"token": "cHJlc2V0dG9rZW4="
			}
		}
	`, body)
}

func TestEC2TokenOptsRequiresAccess(t *testing.T) {
	_, err := auth.EC2TokenOpts{Secret: "18f4f6761ada4e3795fa5273c30349b9"}.ToAuthBody()
	th.AssertErr(t, err)

	missing, ok := err.(gophercloud.ErrMissingInput)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "Access", missing.Argument)
}

func TestEC2TokenOptsSignatureErrors(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]string
		want   string
	}{
		{
			name:   "unsupported signature version",
			params: map[string]string{"SignatureVersion": "3"},
			want:   "unsupported signature version: 3",
		},
		{
			name:   "missing signature method",
			params: map[string]string{"SignatureVersion": "2"},
			want:   "signature method must be provided",
		},
		{
			name:   "unsupported signature method",
			params: map[string]string{"SignatureVersion": "2", "SignatureMethod": "HmacSHA512"},
			want:   "unsupported signature method: HmacSHA512",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			opts := auth.EC2TokenOpts{
				Access: "a7f1e798b7c2417cba4a02de97dc3cdc",
				Secret: "18f4f6761ada4e3795fa5273c30349b9",
				Params: test.params,
			}

			_, err := opts.ToAuthBody()
			th.AssertErr(t, err)
			th.AssertEquals(t, test.want, err.Error())
		})
	}
}

func TestEC2TokenOptsCanReauth(t *testing.T) {
	th.AssertEquals(t, false, auth.EC2TokenOpts{}.CanReauth())
	th.AssertEquals(t, true, auth.EC2TokenOpts{AllowReauth: true}.CanReauth())
}

// TestNewRequestEC2 verifies that the token, which is only used to validate S3
// tokens, is left out of an ec2tokens request.
func TestNewRequestEC2(t *testing.T) {
	request, err := auth.NewRequestEC2(ec2V4Opts())
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, []int{200}, request.OkCodes)
	th.AssertDeepEquals(t, []string(nil), request.OmitHeaders)
	th.AssertJSONEquals(t, fmt.Sprintf(`
		{
			"credentials": {
				"access": "a7f1e798b7c2417cba4a02de97dc3cdc",
				"body_hash": "foo",
				"host": "",
				"path": "/",
				"verb": "GET",
				"headers": {
					"Host": "localhost",
					"Authorization": "AWS4-HMAC-SHA256 Credential=a7f1e798b7c2417cba4a02de97dc3cdc/00010101/region1/ec2/aws4_request, SignedHeaders=, Signature=%s",
					"X-Amz-Date": "00010101T000000Z"
				},
				"params": {"Action": "Test"},
				"signature": "%s"
			}
		}
	`, ec2V4Signature, ec2V4Signature), request.JSONBody)
}

// TestNewRequestS3 verifies that the elements which are only used by the
// ec2tokens endpoint are left out of an s3tokens request.
func TestNewRequestS3(t *testing.T) {
	request, err := auth.NewRequestS3(ec2V4Opts())
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, []int{200}, request.OkCodes)
	th.AssertDeepEquals(t, []string(nil), request.OmitHeaders)
	th.AssertJSONEquals(t, fmt.Sprintf(`
		{
			"credentials": {
				"access": "a7f1e798b7c2417cba4a02de97dc3cdc",
				"signature": "%s",
				"token": "%s"
			}
		}
	`, ec2V4Signature, ec2V4Token), request.JSONBody)
}

func TestAuthOptionsEC2Authenticate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v3/ec2tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", "service-token")
		th.TestJSONRequest(t, r, fmt.Sprintf(`
			{
				"credentials": {
					"access": "a7f1e798b7c2417cba4a02de97dc3cdc",
					"body_hash": "foo",
					"host": "",
					"path": "/",
					"verb": "GET",
					"headers": {
						"Host": "localhost",
						"Authorization": "AWS4-HMAC-SHA256 Credential=a7f1e798b7c2417cba4a02de97dc3cdc/00010101/region1/ec2/aws4_request, SignedHeaders=, Signature=%s",
						"X-Amz-Date": "00010101T000000Z"
					},
					"params": {"Action": "Test"},
					"signature": "%s"
				}
			}
		`, ec2V4Signature, ec2V4Signature))

		w.Header().Set("X-Subject-Token", "the-token-id")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `
			{
				"token": {
					"methods": ["ec2credential"],
					"expires_at": "2026-09-01T12:00:00.000000Z",
					"issued_at": "2026-09-01T11:00:00.000000Z",
					"user": {
						"id": "user-id",
						"name": "testuser",
						"domain": {"id": "default", "name": "Default"}
					},
					"project": {"id": "project-id", "name": "project-name"},
					"roles": [{"id": "role-id", "name": "member"}],
					"catalog": [
						{
							"id": "service-id",
							"name": "nova",
							"type": "compute",
							"endpoints": [
								{
									"id": "endpoint-id",
									"region": "RegionOne",
									"region_id": "RegionOne",
									"interface": "public",
									"url": "http://public.example.com/compute"
								}
							]
						}
					]
				}
			}
		`)
	})

	provider := &gophercloud.ProviderClient{TokenID: "service-token"}
	opts := auth.AuthOptionsEC2{
		ServiceProvider: provider,
		AuthURL:         fakeServer.Endpoint(),
		Auth:            ec2V4Opts(),
	}

	result, err := opts.Authenticate(context.TODO(), nil)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "the-token-id", result.TokenID)
	th.AssertEquals(t, "user-id", result.User.ID)
	th.AssertEquals(t, "project-id", result.Project.ID)
	th.AssertEquals(t, 1, len(result.Roles))
	th.AssertEquals(t, "member", result.Roles[0].Name)
	th.AssertEquals(t, 1, len(result.Catalog.Entries))
	th.AssertEquals(t, false, result.CanReauth)

	expectedExpires := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	th.AssertEquals(t, true, result.ExpiresAt.Equal(expectedExpires))
}

func TestAuthOptionsEC2AuthenticateAllowReauth(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v3/ec2tokens", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Subject-Token", "the-token-id")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"token": {"user": {"id": "user-id"}}}`)
	})

	credentials := ec2V4Opts()
	credentials.AllowReauth = true

	opts := auth.AuthOptionsEC2{
		AuthURL: fakeServer.Endpoint(),
		Auth:    credentials,
	}

	result, err := opts.Authenticate(context.TODO(), nil)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, true, result.CanReauth)
}

func TestAuthOptionsEC2GetAuthURL(t *testing.T) {
	tests := []struct {
		authURL string
		want    string
	}{
		{authURL: "http://example.com:5000/", want: "http://example.com:5000/v3/"},
		{authURL: "http://example.com:5000/v3", want: "http://example.com:5000/v3/"},
		{authURL: "http://example.com:5000/v2.0", want: "http://example.com:5000/v3/"},
	}

	for _, test := range tests {
		t.Run(test.authURL, func(t *testing.T) {
			th.AssertEquals(t, test.want, auth.AuthOptionsEC2{AuthURL: test.authURL}.GetAuthURL())
		})
	}
}

func TestAuthOptionsEC2AuthenticateNilAuth(t *testing.T) {
	opts := auth.AuthOptionsEC2{AuthURL: "http://example.com:5000/v3"}

	_, err := opts.Authenticate(context.TODO(), nil)
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrMissingInput)
	th.AssertEquals(t, true, ok)
}

func TestAuthOptionsEC2AuthenticateToAuthBodyError(t *testing.T) {
	opts := auth.AuthOptionsEC2{
		AuthURL: "http://example.com:5000/v3",
		Auth:    auth.EC2TokenOpts{}, // no access set
	}

	_, err := opts.Authenticate(context.TODO(), nil)
	th.AssertErr(t, err)
}

func TestAuthOptionsEC2AuthenticateHTTPFailure(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v3/ec2tokens", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"error": {"message": "bad credentials"}}`)
	})

	opts := auth.AuthOptionsEC2{
		AuthURL: fakeServer.Endpoint(),
		Auth:    ec2V4Opts(),
	}

	_, err := opts.Authenticate(context.TODO(), nil)
	th.AssertErr(t, err)
}

func TestEC2CredentialsBuildCanonicalQueryStringV2(t *testing.T) {
	params := map[string]string{
		"Action": "foo",
		"Value":  "bar",
	}
	expected := "Action=foo&Value=bar"
	th.CheckEquals(t, expected, auth.EC2CredentialsBuildCanonicalQueryStringV2(params))
}

func TestEC2CredentialsBuildStringToSignV2(t *testing.T) {
	opts := auth.EC2TokenOpts{
		Verb: "GET",
		Host: "localhost",
		Path: "/",
		Params: map[string]string{
			"Action": "foo",
			"Value":  "bar",
		},
	}
	expected := []byte("GET\nlocalhost\n/\nAction=foo&Value=bar")
	th.CheckDeepEquals(t, expected, auth.EC2CredentialsBuildStringToSignV2(opts))
}

func TestEC2CredentialsBuildCanonicalQueryStringV4(t *testing.T) {
	params := map[string]string{
		"Action": "foo",
		"Value":  "bar",
	}
	expected := "Action=foo&Value=bar"
	th.CheckEquals(t, expected, auth.EC2CredentialsBuildCanonicalQueryStringV4("foo", params))
	th.CheckEquals(t, "", auth.EC2CredentialsBuildCanonicalQueryStringV4("POST", params))
}

func TestEC2CredentialsBuildCanonicalHeadersV4(t *testing.T) {
	headers := map[string]string{
		"Foo": "bar",
		"Baz": "qux",
	}
	signedHeaders := "foo;baz"
	expected := "foo:bar\nbaz:qux\n"
	th.CheckEquals(t, expected, auth.EC2CredentialsBuildCanonicalHeadersV4(headers, signedHeaders))
}

func TestEC2CredentialsBuildSignatureKeyV4(t *testing.T) {
	expected := "246626bd815b0a0cae4bedc3f4e124ca25e208cd75fd812d836aeae184de038a"
	th.CheckEquals(t, expected, hex.EncodeToString(auth.EC2CredentialsBuildSignatureKeyV4("foo", "bar", "baz", time.Time{})))
}

func TestEC2CredentialsBuildSignatureV4(t *testing.T) {
	opts := auth.EC2TokenOpts{
		Verb: "GET",
		Path: "/",
		Headers: map[string]string{
			"Host": "localhost",
		},
		Params: map[string]string{
			"Action": "foo",
			"Value":  "bar",
		},
	}
	expected := "6a5febe41427bf601f0ae7c34dbb0fd67094776138b03fb8e65783d733d302a5"

	date := time.Time{}
	stringToSign := auth.EC2CredentialsBuildStringToSignV4(opts, "host", "foo", date)
	key := auth.EC2CredentialsBuildSignatureKeyV4("", "", "", date)

	th.CheckEquals(t, expected, auth.EC2CredentialsBuildSignatureV4(key, stringToSign))
}

func TestEC2CredentialsBuildAuthorizationHeaderV4(t *testing.T) {
	opts := auth.EC2TokenOpts{
		Access:  "a7f1e798b7c2417cba4a02de97dc3cdc",
		Region:  "region1",
		Service: "ec2",
	}
	expected := "AWS4-HMAC-SHA256 Credential=a7f1e798b7c2417cba4a02de97dc3cdc/00010101/region1/ec2/aws4_request, SignedHeaders=host, Signature=foo"

	th.CheckEquals(t, expected, auth.EC2CredentialsBuildAuthorizationHeaderV4(opts, "host", "foo", time.Time{}))
}
