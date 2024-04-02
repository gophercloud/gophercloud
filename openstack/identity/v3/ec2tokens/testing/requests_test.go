package testing

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/ec2tokens"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	tokens_testing "github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens/testing"
	"github.com/gophercloud/gophercloud/v2/testhelper"
)

// authTokenPost verifies that providing certain AuthOptions and Scope results in an expected JSON structure.
func authTokenPost(t *testing.T, options ec2tokens.AuthOptions, requestJSON string) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       testhelper.Endpoint(),
	}

	testhelper.Mux.HandleFunc("/ec2tokens", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "POST")
		testhelper.TestHeader(t, r, "Content-Type", "application/json")
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestJSONRequest(t, r, requestJSON)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, tokens_testing.TokenOutput)
	})

	expected := &tokens.Token{
		ExpiresAt: time.Date(2017, 6, 3, 2, 19, 49, 0, time.UTC),
	}

	actual, err := ec2tokens.Create(context.TODO(), &client, &options).Extract()
	testhelper.AssertNoErr(t, err)
	testhelper.CheckDeepEquals(t, expected, actual)
}

func TestCreateV2(t *testing.T) {
	credentials := ec2tokens.AuthOptions{
		Access: "a7f1e798b7c2417cba4a02de97dc3cdc",
		Host:   "localhost",
		Path:   "/",
		Secret: "18f4f6761ada4e3795fa5273c30349b9",
		Verb:   "GET",
		// this should be removed from JSON request
		BodyHash: new(string),
		// this should be removed from JSON request
		Headers: map[string]string{
			"Foo": "Bar",
		},
		Params: map[string]string{
			"Action":           "Test",
			"SignatureMethod":  "HmacSHA256",
			"SignatureVersion": "2",
		},
	}
	authTokenPost(t, credentials, `{
    "credentials": {
        "access": "a7f1e798b7c2417cba4a02de97dc3cdc",
        "host": "localhost",
        "params": {
            "Action": "Test",
            "SignatureMethod": "HmacSHA256",
            "SignatureVersion": "2"
        },
        "path": "/",
        "signature": "Up+MbVbbrvdR5FRkUz+n3nc+VW6xieuN50wh6ONEJ4w=",
        "verb": "GET"
    }
}`)
}

func TestCreateV4(t *testing.T) {
	bodyHash := "foo"
	credentials := ec2tokens.AuthOptions{
		Access:    "a7f1e798b7c2417cba4a02de97dc3cdc",
		BodyHash:  &bodyHash,
		Timestamp: new(time.Time),
		Region:    "region1",
		Service:   "ec2",
		Path:      "/",
		Secret:    "18f4f6761ada4e3795fa5273c30349b9",
		Verb:      "GET",
		Headers: map[string]string{
			"Host": "localhost",
		},
		Params: map[string]string{
			"Action": "Test",
		},
	}
	authTokenPost(t, credentials, `{
    "credentials": {
        "access": "a7f1e798b7c2417cba4a02de97dc3cdc",
        "body_hash": "foo",
        "host": "",
        "headers": {
            "Host": "localhost",
            "Authorization": "AWS4-HMAC-SHA256 Credential=a7f1e798b7c2417cba4a02de97dc3cdc/00010101/region1/ec2/aws4_request, SignedHeaders=, Signature=f36f79118f75d7d6ec86ead9a61679cbdcf94c0cbfe5e9cf2407e8406aa82028",
            "X-Amz-Date": "00010101T000000Z"
         },
        "params": {
            "Action": "Test"
        },
        "path": "/",
        "signature": "f36f79118f75d7d6ec86ead9a61679cbdcf94c0cbfe5e9cf2407e8406aa82028",
        "verb": "GET"
    }
}`)
}

func TestCreateV4Empty(t *testing.T) {
	credentials := ec2tokens.AuthOptions{
		Access:    "a7f1e798b7c2417cba4a02de97dc3cdc",
		Secret:    "18f4f6761ada4e3795fa5273c30349b9",
		BodyHash:  new(string),
		Timestamp: new(time.Time),
	}
	authTokenPost(t, credentials, `{
    "credentials": {
        "access": "a7f1e798b7c2417cba4a02de97dc3cdc",
        "body_hash": "",
        "host": "",
        "headers": {
            "Authorization": "AWS4-HMAC-SHA256 Credential=a7f1e798b7c2417cba4a02de97dc3cdc/00010101///aws4_request, SignedHeaders=, Signature=140a31abf1efe93a607dcac6cd8f66887b86d2bc8f712c290d9aa06edf428608",
            "X-Amz-Date": "00010101T000000Z"
        },
        "params": {},
        "path": "",
        "signature": "140a31abf1efe93a607dcac6cd8f66887b86d2bc8f712c290d9aa06edf428608",
        "verb": ""
    }
}`)
}

func TestCreateV4Headers(t *testing.T) {
	credentials := ec2tokens.AuthOptions{
		Access:    "a7f1e798b7c2417cba4a02de97dc3cdc",
		BodyHash:  new(string),
		Timestamp: new(time.Time),
		Region:    "region1",
		Service:   "ec2",
		Path:      "/",
		Secret:    "18f4f6761ada4e3795fa5273c30349b9",
		Verb:      "GET",
		Headers: map[string]string{
			"Foo":  "Bar",
			"Host": "localhost",
		},
		Params: map[string]string{
			"Action": "Test",
		},
	}
	authTokenPost(t, credentials, `{
    "credentials": {
        "access": "a7f1e798b7c2417cba4a02de97dc3cdc",
        "body_hash": "",
        "host": "",
        "headers": {
            "Foo": "Bar",
            "Host": "localhost",
            "Authorization": "AWS4-HMAC-SHA256 Credential=a7f1e798b7c2417cba4a02de97dc3cdc/00010101/region1/ec2/aws4_request, SignedHeaders=, Signature=f5cd6995be98e5576a130b30cca277375f10439217ea82169aa8386e83965611",
            "X-Amz-Date": "00010101T000000Z"
        },
        "params": {
            "Action": "Test"
        },
        "path": "/",
        "signature": "f5cd6995be98e5576a130b30cca277375f10439217ea82169aa8386e83965611",
        "verb": "GET"
    }
}`)
}

func TestCreateV4WithSignature(t *testing.T) {
	credentials := ec2tokens.AuthOptions{
		Access:    "a7f1e798b7c2417cba4a02de97dc3cdc",
		BodyHash:  new(string),
		Path:      "/",
		Signature: "f5cd6995be98e5576a130b30cca277375f10439217ea82169aa8386e83965611",
		Verb:      "GET",
		Headers: map[string]string{
			"Foo":           "Bar",
			"Host":          "localhost",
			"Authorization": "AWS4-HMAC-SHA256 Credential=a7f1e798b7c2417cba4a02de97dc3cdc/00010101/region1/ec2/aws4_request, SignedHeaders=, Signature=f5cd6995be98e5576a130b30cca277375f10439217ea82169aa8386e83965611",
			"X-Amz-Date":    "00010101T000000Z",
		},
		Params: map[string]string{
			"Action": "Test",
		},
	}
	authTokenPost(t, credentials, `{
    "credentials": {
        "access": "a7f1e798b7c2417cba4a02de97dc3cdc",
        "body_hash": "",
        "host": "",
        "headers": {
            "Foo": "Bar",
            "Host": "localhost",
            "Authorization": "AWS4-HMAC-SHA256 Credential=a7f1e798b7c2417cba4a02de97dc3cdc/00010101/region1/ec2/aws4_request, SignedHeaders=, Signature=f5cd6995be98e5576a130b30cca277375f10439217ea82169aa8386e83965611",
            "X-Amz-Date": "00010101T000000Z"
        },
        "params": {
            "Action": "Test"
        },
        "path": "/",
        "signature": "f5cd6995be98e5576a130b30cca277375f10439217ea82169aa8386e83965611",
        "verb": "GET"
    }
}`)
}

func TestEC2CredentialsBuildCanonicalQueryStringV2(t *testing.T) {
	params := map[string]string{
		"Action": "foo",
		"Value":  "bar",
	}
	expected := "Action=foo&Value=bar"
	testhelper.CheckEquals(t, expected, ec2tokens.EC2CredentialsBuildCanonicalQueryStringV2(params))
}

func TestEC2CredentialsBuildStringToSignV2(t *testing.T) {
	opts := ec2tokens.AuthOptions{
		Verb: "GET",
		Host: "localhost",
		Path: "/",
		Params: map[string]string{
			"Action": "foo",
			"Value":  "bar",
		},
	}
	expected := []byte("GET\nlocalhost\n/\nAction=foo&Value=bar")
	testhelper.CheckDeepEquals(t, expected, ec2tokens.EC2CredentialsBuildStringToSignV2(opts))
}

func TestEC2CredentialsBuildCanonicalQueryStringV4(t *testing.T) {
	params := map[string]string{
		"Action": "foo",
		"Value":  "bar",
	}
	expected := "Action=foo&Value=bar"
	testhelper.CheckEquals(t, expected, ec2tokens.EC2CredentialsBuildCanonicalQueryStringV4("foo", params))
	testhelper.CheckEquals(t, "", ec2tokens.EC2CredentialsBuildCanonicalQueryStringV4("POST", params))
}

func TestEC2CredentialsBuildCanonicalHeadersV4(t *testing.T) {
	headers := map[string]string{
		"Foo": "bar",
		"Baz": "qux",
	}
	signedHeaders := "foo;baz"
	expected := "foo:bar\nbaz:qux\n"
	testhelper.CheckEquals(t, expected, ec2tokens.EC2CredentialsBuildCanonicalHeadersV4(headers, signedHeaders))
}

func TestEC2CredentialsBuildSignatureKeyV4(t *testing.T) {
	expected := "246626bd815b0a0cae4bedc3f4e124ca25e208cd75fd812d836aeae184de038a"
	testhelper.CheckEquals(t, expected, hex.EncodeToString((ec2tokens.EC2CredentialsBuildSignatureKeyV4("foo", "bar", "baz", time.Time{}))))
}

func TestEC2CredentialsBuildSignatureV4(t *testing.T) {
	opts := ec2tokens.AuthOptions{
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
	stringToSign := ec2tokens.EC2CredentialsBuildStringToSignV4(opts, "host", "foo", date)
	key := ec2tokens.EC2CredentialsBuildSignatureKeyV4("", "", "", date)

	testhelper.CheckEquals(t, expected, ec2tokens.EC2CredentialsBuildSignatureV4(key, stringToSign))
}
