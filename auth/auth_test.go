package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestAuthOptionsV3AuthenticatePropagatesScopeError(t *testing.T) {
	expected := errors.New("scope failure")
	opts := AuthOptionsV3{
		AuthURL: "https://identity.example.com",
		Auth: testV3Auth{
			body:     map[string]map[string]any{"token": {"id": "testtoken"}},
			scopeErr: expected,
		},
	}

	_, err := opts.Authenticate(context.TODO(), nil)
	th.AssertEquals(t, expected, err)
}

func TestAuthOptionsV3AuthenticatePropagatesHeaderError(t *testing.T) {
	expected := errors.New("header failure")
	opts := AuthOptionsV3{
		AuthURL: "https://identity.example.com",
		Auth: testV3Auth{
			body:      map[string]map[string]any{"token": {"id": "testtoken"}},
			headerErr: expected,
		},
	}

	_, err := opts.Authenticate(context.TODO(), nil)
	th.AssertEquals(t, expected, err)
}

func TestAuthOptionsV3AuthenticateAddsBuilderHeaders(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestHeader(t, r, "X-Test-Number", "123")
		w.Header().Set("X-Subject-Token", "testtoken")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"token": {}}`)
	})

	opts := AuthOptionsV3{
		AuthURL: fakeServer.Endpoint(),
		Auth: testV3Auth{
			body:    map[string]map[string]any{"token": {"id": "testtoken"}},
			headers: map[string]any{"X-Test-Number": 123},
		},
	}

	result, err := opts.Authenticate(context.TODO(), nil)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "testtoken", result.TokenID)
}

type testV3Auth struct {
	body      map[string]map[string]any
	scope     map[string]any
	headers   map[string]any
	scopeErr  error
	headerErr error
}

func (a testV3Auth) ToAuthBody() (map[string]map[string]any, error) {
	return a.body, nil
}

func (a testV3Auth) CanReauth() bool {
	return false
}

func (a testV3Auth) ToAuthScope() (map[string]any, error) {
	return a.scope, a.scopeErr
}

func (a testV3Auth) ToAuthHeaders(_ ...requestOption) (map[string]any, error) {
	return a.headers, a.headerErr
}

func (a testV3Auth) ToAuthType() AuthType { return AuthV3Token }

// ec2V4Opts are EC2 credentials signed with a deterministic AWS signature V4.
func ec2V4Opts() EC2TokenOpts {
	bodyHash := "foo"
	return EC2TokenOpts{
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

// TestNewRequestEC2 verifies that the token, which is only used to validate S3
// tokens, is left out of an ec2tokens request.
func TestNewRequestEC2(t *testing.T) {
	request, err := newRequestEC2(ec2V4Opts(), "token")
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
	request, err := newRequestEC2(ec2V4Opts(), "body_hash", "headers", "host", "params", "path", "verb")
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
