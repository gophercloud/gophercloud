package testing

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/extensions/ec2tokens"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	tokens_testing "github.com/gophercloud/gophercloud/openstack/identity/v3/tokens/testing"
	"github.com/gophercloud/gophercloud/testhelper"
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

	actual, err := ec2tokens.Create(&client, &options).Extract()
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
            "Host": "localhost"
         },
        "params": {
            "Action": "Test",
            "X-Amz-Algorithm": "AWS4-HMAC-SHA256",
            "X-Amz-Credential": "a7f1e798b7c2417cba4a02de97dc3cdc/00010101/region1/ec2/aws4_request",
            "X-Amz-Date": "00010101T000000Z",
            "X-Amz-SignedHeaders": ""
        },
        "path": "/",
        "signature": "f8b67f8bc0c8c7a6b84fcb37dd9007f1c728625fe783fe8616b6ba150e1d0655",
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
        "headers": {},
        "params": {
            "X-Amz-Algorithm": "AWS4-HMAC-SHA256",
            "X-Amz-Credential": "a7f1e798b7c2417cba4a02de97dc3cdc/00010101///aws4_request",
            "X-Amz-Date": "00010101T000000Z",
            "X-Amz-SignedHeaders": ""
        },
        "path": "",
        "signature": "a7f7a34f93ccd460e67079743d05d1b4488c6e9d708869b6b687d51244c68857",
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
            "Host": "localhost"
        },
        "params": {
            "Action": "Test",
            "X-Amz-Algorithm": "AWS4-HMAC-SHA256",
            "X-Amz-Credential": "a7f1e798b7c2417cba4a02de97dc3cdc/00010101/region1/ec2/aws4_request",
            "X-Amz-Date": "00010101T000000Z",
            "X-Amz-SignedHeaders": ""
        },
        "path": "/",
        "signature": "becc7d835c1d96835b772158198aed5583cb3fee4e7c7459a15fe710e3f533c6",
        "verb": "GET"
    }
}`)
}

func TestCreateV4WithSignature(t *testing.T) {
	credentials := ec2tokens.AuthOptions{
		Access:    "a7f1e798b7c2417cba4a02de97dc3cdc",
		BodyHash:  new(string),
		Path:      "/",
		Signature: "becc7d835c1d96835b772158198aed5583cb3fee4e7c7459a15fe710e3f533c6",
		Verb:      "GET",
		Headers: map[string]string{
			"Foo":  "Bar",
			"Host": "localhost",
		},
		Params: map[string]string{
			"Action":              "Test",
			"X-Amz-Algorithm":     "AWS4-HMAC-SHA256",
			"X-Amz-Credential":    "a7f1e798b7c2417cba4a02de97dc3cdc/00010101/region1/ec2/aws4_request",
			"X-Amz-Date":          "00010101T000000Z",
			"X-Amz-SignedHeaders": "",
		},
	}
	authTokenPost(t, credentials, `{
    "credentials": {
        "access": "a7f1e798b7c2417cba4a02de97dc3cdc",
        "body_hash": "",
        "host": "",
        "headers": {
            "Foo": "Bar",
            "Host": "localhost"
        },
        "params": {
            "Action": "Test",
            "X-Amz-Algorithm": "AWS4-HMAC-SHA256",
            "X-Amz-Credential": "a7f1e798b7c2417cba4a02de97dc3cdc/00010101/region1/ec2/aws4_request",
            "X-Amz-Date": "00010101T000000Z",
            "X-Amz-SignedHeaders": ""
        },
        "path": "/",
        "signature": "becc7d835c1d96835b772158198aed5583cb3fee4e7c7459a15fe710e3f533c6",
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
	}
	params := map[string]string{
		"Action": "foo",
		"Value":  "bar",
	}
	expected := "GET\nlocalhost\n/\nAction=foo&Value=bar"
	testhelper.CheckEquals(t, expected, ec2tokens.EC2CredentialsBuildStringToSignV2(opts, params))
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
	expected := "5f06633a42b1324477cb006b24b1266722703b5dcf9186481be6592b9554c38f"
	testhelper.CheckEquals(t, expected, hex.EncodeToString((ec2tokens.EC2CredentialsBuildSignatureKeyV4("foo", "bar", "baz", "qux"))))
}

func TestEC2CredentialsBuildSignatureV4(t *testing.T) {
	opts := ec2tokens.AuthOptions{
		Verb: "GET",
		Path: "/",
		Headers: map[string]string{
			"Host": "localhost",
		},
	}
	params := map[string]string{
		"Action":              "foo",
		"Value":               "bar",
		"X-Amz-SignedHeaders": "host",
	}
	expected := "76ecaab5616b005227ad98bcdf7e802120a16a4e840faf5193e49db7fe7db178"
	testhelper.CheckEquals(t, expected, ec2tokens.EC2CredentialsBuildSignatureV4(opts, params, time.Time{}, "foo"))
}
