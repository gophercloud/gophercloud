package testing

import (
	"fmt"
	"net/http"
	"testing"

	"crypto/x509"
	"encoding/pem"
	"github.com/gophercloud/gophercloud/openstack/containerorchestration/v1/certificates"
	fake "github.com/gophercloud/gophercloud/openstack/containerorchestration/v1/common"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestCreateCertificate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/bays/mccluster", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "status": "CREATE_COMPLETE",
  "uuid": "1f085f23-5206-4192-99da-53c0558626ee",
  "links": [
    {
      "href": "http://65.61.151.130:9511/v1/bays/1f085f23-5206-4192-99da-53c0558626ee",
      "rel": "self"
    },
    {
      "href": "http://65.61.151.130:9511/bays/1f085f23-5206-4192-99da-53c0558626ee",
      "rel": "bookmark"
    }
  ],
  "stack_id": "7590db28-9d80-4e33-8aeb-57c09085828a",
  "created_at": "2016-08-16T14:11:31+00:00",
  "api_address": "https://172.29.248.97:6443",
  "discovery_url": "https://discovery.etcd.io/0d1e5076f84678db8865ab99debb5516",
  "updated_at": "2016-08-16T14:15:14+00:00",
  "master_count": 1,
  "baymodel_id": "472807c2-f175-4946-9765-149701a5aba7",
  "master_addresses": [
    "172.29.248.97"
  ],
  "node_count": 1,
  "node_addresses": [
    "172.29.248.98"
  ],
  "status_reason": "Stack CREATE completed successfully",
  "bay_create_timeout": 60,
  "name": "mccluster"
}
			`)
	})

	th.Mux.HandleFunc("/v1/certificates", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
  "bay_uuid": "1f085f23-5206-4192-99da-53c0558626ee",
  "pem": "-----BEGIN CERTIFICATE-----\nMIIDxzCCAq+gAwIBAgIQDoAuMAUcRw6eBgV53L7+bjANBgkqhkiG9w0BAQsFADAU\nMRIwEAYDVQQDDAltY2NsdXN0ZXIwHhcNMTYwODE1MDkyNzU0WhcNMjEwODE1MDky\nNzU0WjASMRAwDgYDVQQDDAdDYXJvbHluMIICIjANBgkqhkiG9w0BAQEFAAOCAg8A\nMIICCgKCAgEAxWUD6G1O6sYjo+6xQW9scFdrJP2MzBJ6tKumc6sGET9Jl+XWTI+d\ng+a6rA+wRer+GkBn91tWuTVKOq/p7SPjl49GtN9jwUvFH2K8JOdXMGAa7oX/oUvo\nmbxwfpq/k02Y7qYU+RMK7hfTLFS63ufZDXyoG89c0x0u+b8ZtB94XJyAbRZfo35W\nb8Q/8Ql1qxBwD7XSL8PPf5TxXWT9cbX34wdbJzdXOZUb5Tmp0/Zsii2HdUwL363+\nc4gWFb9fVBCZHQOZ68NWE5B3OkY0S9y9o/l4YkI5jidxlCGwVaSw/dqQfHW14aFw\nefGpnksvQCd6IW/aGfJb/tSvqsn8GE41hSBP9/gF+NoPICPtb+LJwvaaLoeDJPtL\nNG0C9+TBkDXPzRki3vx7R+v8Y0/54l71K1M/xoDL1WnnPc7tLS8Bckfh2ryCmPbz\nTa11HGVZwEg1yMAhtcb/OKgvsKnU/Y7HNZ4tFkKC32TMNmPKLAgpy9B179u8tkv2\nlQ73j4iyIXrdbFr8gGzWh8EXO9ykOMV8IYuw6DXW+awSxUW3v+3hAhoAPZxJMfqp\nmULKKN3JravknDmAYRbnqPc0LY6iPLmtVmmkFug/nmgWBqfiNE2YPn2C+5toJLul\n/QRZiGMaDkfx642PzOzbXKxgDhu7OmrIL7sM7lCtlNxFPB/Road8fBECAwEAAaMX\nMBUwEwYDVR0lBAwwCgYIKwYBBQUHAwIwDQYJKoZIhvcNAQELBQADggEBAId8LhM6\nbOSubyLDkRoJLbkk0ADgMcBWo9RGR+kpT2tJSbFSfOlGKyPqVjwmk8vTGv1SgMTS\n20tqKaWwtPAV+AiK8EG578/uoPSWiEKI49T3aNEb85Kk8t0YRWeh1giQEIe1SDzw\nqKBse7Pxu77xJ2h7VlYCrk3gfc/7YE7Qo+EBycTuDWOjjX2vDPuF1z+FQ3VWJPgb\nIDbyTVpozO+ZS1P9TLTszUTYdv9751RmLuvIzcnVcWrtNMXFc0JYLA55oUMckYXG\n2b8ZnriZ7L4CPlvmvOrjHuXsB4lEwZcdAM+3H1LbsNjo9/W4VrwMsXkXF+UvgDAj\nOwozelKfCQzHzhg=\n-----END CERTIFICATE-----\n"
}
		`)
	})

	pemContents := "-----BEGIN RSA PRIVATE KEY-----\n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\n-----END RSA PRIVATE KEY-----\n"

	block, _ := pem.Decode([]byte(pemContents))
	privateKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

	r := certificates.CreateCertificate(fake.ServiceClient(), "mccluster", privateKey)
	th.AssertNoErr(t, r.Err)

	c, err := r.Extract()

	th.AssertNoErr(t, err)
	th.AssertEquals(t, "1f085f23-5206-4192-99da-53c0558626ee", c.BayID)

	certContents := "-----BEGIN CERTIFICATE-----\nMIIDxzCCAq+gAwIBAgIQDoAuMAUcRw6eBgV53L7+bjANBgkqhkiG9w0BAQsFADAU\nMRIwEAYDVQQDDAltY2NsdXN0ZXIwHhcNMTYwODE1MDkyNzU0WhcNMjEwODE1MDky\nNzU0WjASMRAwDgYDVQQDDAdDYXJvbHluMIICIjANBgkqhkiG9w0BAQEFAAOCAg8A\nMIICCgKCAgEAxWUD6G1O6sYjo+6xQW9scFdrJP2MzBJ6tKumc6sGET9Jl+XWTI+d\ng+a6rA+wRer+GkBn91tWuTVKOq/p7SPjl49GtN9jwUvFH2K8JOdXMGAa7oX/oUvo\nmbxwfpq/k02Y7qYU+RMK7hfTLFS63ufZDXyoG89c0x0u+b8ZtB94XJyAbRZfo35W\nb8Q/8Ql1qxBwD7XSL8PPf5TxXWT9cbX34wdbJzdXOZUb5Tmp0/Zsii2HdUwL363+\nc4gWFb9fVBCZHQOZ68NWE5B3OkY0S9y9o/l4YkI5jidxlCGwVaSw/dqQfHW14aFw\nefGpnksvQCd6IW/aGfJb/tSvqsn8GE41hSBP9/gF+NoPICPtb+LJwvaaLoeDJPtL\nNG0C9+TBkDXPzRki3vx7R+v8Y0/54l71K1M/xoDL1WnnPc7tLS8Bckfh2ryCmPbz\nTa11HGVZwEg1yMAhtcb/OKgvsKnU/Y7HNZ4tFkKC32TMNmPKLAgpy9B179u8tkv2\nlQ73j4iyIXrdbFr8gGzWh8EXO9ykOMV8IYuw6DXW+awSxUW3v+3hAhoAPZxJMfqp\nmULKKN3JravknDmAYRbnqPc0LY6iPLmtVmmkFug/nmgWBqfiNE2YPn2C+5toJLul\n/QRZiGMaDkfx642PzOzbXKxgDhu7OmrIL7sM7lCtlNxFPB/Road8fBECAwEAAaMX\nMBUwEwYDVR0lBAwwCgYIKwYBBQUHAwIwDQYJKoZIhvcNAQELBQADggEBAId8LhM6\nbOSubyLDkRoJLbkk0ADgMcBWo9RGR+kpT2tJSbFSfOlGKyPqVjwmk8vTGv1SgMTS\n20tqKaWwtPAV+AiK8EG578/uoPSWiEKI49T3aNEb85Kk8t0YRWeh1giQEIe1SDzw\nqKBse7Pxu77xJ2h7VlYCrk3gfc/7YE7Qo+EBycTuDWOjjX2vDPuF1z+FQ3VWJPgb\nIDbyTVpozO+ZS1P9TLTszUTYdv9751RmLuvIzcnVcWrtNMXFc0JYLA55oUMckYXG\n2b8ZnriZ7L4CPlvmvOrjHuXsB4lEwZcdAM+3H1LbsNjo9/W4VrwMsXkXF+UvgDAj\nOwozelKfCQzHzhg=\n-----END CERTIFICATE-----\n"

	th.AssertEquals(t, certContents, c.String())
}

func TestCreateCertificateFailed(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/bays/mccluster", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `
{
  "errors": [
    {
      "status": 404,
      "code": "client",
      "links": [],
      "title": "Bay mccluster could not be found",
      "detail": "Bay mccluster could not be found.",
      "request_id": ""
    }
  ]
}
		`)
	})

	pemContents := "-----BEGIN RSA PRIVATE KEY-----\n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\n-----END RSA PRIVATE KEY-----\n"

	block, _ := pem.Decode([]byte(pemContents))
	privateKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

	resp := certificates.CreateCertificate(fake.ServiceClient(), "mccluster", privateKey)
	th.AssertEquals(t, "Bay mccluster could not be found.", resp.Err.Error())

	er, ok := resp.Err.(*fake.ErrorResponse)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, http.StatusNotFound, er.Actual)
}

func TestCreateCredentialsBundle(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/bays/mccluster", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "status": "CREATE_COMPLETE",
  "uuid": "1f085f23-5206-4192-99da-53c0558626ee",
  "links": [
    {
      "href": "http://65.61.151.130:9511/v1/bays/1f085f23-5206-4192-99da-53c0558626ee",
      "rel": "self"
    },
    {
      "href": "http://65.61.151.130:9511/bays/1f085f23-5206-4192-99da-53c0558626ee",
      "rel": "bookmark"
    }
  ],
  "stack_id": "7590db28-9d80-4e33-8aeb-57c09085828a",
  "created_at": "2016-08-16T14:11:31+00:00",
  "api_address": "https://172.29.248.97:6443",
  "discovery_url": "https://discovery.etcd.io/0d1e5076f84678db8865ab99debb5516",
  "updated_at": "2016-08-16T14:15:14+00:00",
  "master_count": 1,
  "baymodel_id": "472807c2-f175-4946-9765-149701a5aba7",
  "master_addresses": [
    "172.29.248.97"
  ],
  "node_count": 1,
  "node_addresses": [
    "172.29.248.98"
  ],
  "status_reason": "Stack CREATE completed successfully",
  "bay_create_timeout": 60,
  "name": "mccluster"
}
			`)
	})

	th.Mux.HandleFunc("/v1/bays/1f085f23-5206-4192-99da-53c0558626ee", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "status": "CREATE_COMPLETE",
  "uuid": "1f085f23-5206-4192-99da-53c0558626ee",
  "links": [
    {
      "href": "http://65.61.151.130:9511/v1/bays/1f085f23-5206-4192-99da-53c0558626ee",
      "rel": "self"
    },
    {
      "href": "http://65.61.151.130:9511/bays/1f085f23-5206-4192-99da-53c0558626ee",
      "rel": "bookmark"
    }
  ],
  "stack_id": "7590db28-9d80-4e33-8aeb-57c09085828a",
  "created_at": "2016-08-16T14:11:31+00:00",
  "api_address": "https://172.29.248.97:6443",
  "discovery_url": "https://discovery.etcd.io/0d1e5076f84678db8865ab99debb5516",
  "updated_at": "2016-08-16T14:15:14+00:00",
  "master_count": 1,
  "baymodel_id": "472807c2-f175-4946-9765-149701a5aba7",
  "master_addresses": [
    "172.29.248.97"
  ],
  "node_count": 1,
  "node_addresses": [
    "172.29.248.98"
  ],
  "status_reason": "Stack CREATE completed successfully",
  "bay_create_timeout": 60,
  "name": "mccluster"
}
			`)
	})

	th.Mux.HandleFunc("/v1/certificates", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
  "bay_uuid": "1f085f23-5206-4192-99da-53c0558626ee",
  "pem": "-----BEGIN CERTIFICATE-----\nMIIDxzCCAq+gAwIBAgIQDoAuMAUcRw6eBgV53L7+bjANBgkqhkiG9w0BAQsFADAU\nMRIwEAYDVQQDDAltY2NsdXN0ZXIwHhcNMTYwODE1MDkyNzU0WhcNMjEwODE1MDky\nNzU0WjASMRAwDgYDVQQDDAdDYXJvbHluMIICIjANBgkqhkiG9w0BAQEFAAOCAg8A\nMIICCgKCAgEAxWUD6G1O6sYjo+6xQW9scFdrJP2MzBJ6tKumc6sGET9Jl+XWTI+d\ng+a6rA+wRer+GkBn91tWuTVKOq/p7SPjl49GtN9jwUvFH2K8JOdXMGAa7oX/oUvo\nmbxwfpq/k02Y7qYU+RMK7hfTLFS63ufZDXyoG89c0x0u+b8ZtB94XJyAbRZfo35W\nb8Q/8Ql1qxBwD7XSL8PPf5TxXWT9cbX34wdbJzdXOZUb5Tmp0/Zsii2HdUwL363+\nc4gWFb9fVBCZHQOZ68NWE5B3OkY0S9y9o/l4YkI5jidxlCGwVaSw/dqQfHW14aFw\nefGpnksvQCd6IW/aGfJb/tSvqsn8GE41hSBP9/gF+NoPICPtb+LJwvaaLoeDJPtL\nNG0C9+TBkDXPzRki3vx7R+v8Y0/54l71K1M/xoDL1WnnPc7tLS8Bckfh2ryCmPbz\nTa11HGVZwEg1yMAhtcb/OKgvsKnU/Y7HNZ4tFkKC32TMNmPKLAgpy9B179u8tkv2\nlQ73j4iyIXrdbFr8gGzWh8EXO9ykOMV8IYuw6DXW+awSxUW3v+3hAhoAPZxJMfqp\nmULKKN3JravknDmAYRbnqPc0LY6iPLmtVmmkFug/nmgWBqfiNE2YPn2C+5toJLul\n/QRZiGMaDkfx642PzOzbXKxgDhu7OmrIL7sM7lCtlNxFPB/Road8fBECAwEAAaMX\nMBUwEwYDVR0lBAwwCgYIKwYBBQUHAwIwDQYJKoZIhvcNAQELBQADggEBAId8LhM6\nbOSubyLDkRoJLbkk0ADgMcBWo9RGR+kpT2tJSbFSfOlGKyPqVjwmk8vTGv1SgMTS\n20tqKaWwtPAV+AiK8EG578/uoPSWiEKI49T3aNEb85Kk8t0YRWeh1giQEIe1SDzw\nqKBse7Pxu77xJ2h7VlYCrk3gfc/7YE7Qo+EBycTuDWOjjX2vDPuF1z+FQ3VWJPgb\nIDbyTVpozO+ZS1P9TLTszUTYdv9751RmLuvIzcnVcWrtNMXFc0JYLA55oUMckYXG\n2b8ZnriZ7L4CPlvmvOrjHuXsB4lEwZcdAM+3H1LbsNjo9/W4VrwMsXkXF+UvgDAj\nOwozelKfCQzHzhg=\n-----END CERTIFICATE-----\n"
}
		`)
	})

	th.Mux.HandleFunc("/v1/certificates/1f085f23-5206-4192-99da-53c0558626ee", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "pem": "-----BEGIN CERTIFICATE-----\nMIIC0zCCAbugAwIBAgIRAIz6zUqixUcXvPZpDMcvwJMwDQYJKoZIhvcNAQELBQAw\nETEPMA0GA1UEAwwGazhzYmF5MB4XDTE2MDcxMzE4NTg0OFoXDTIxMDcxMzE4NTg0\nOFowETEPMA0GA1UEAwwGazhzYmF5MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB\nCgKCAQEA543UyfsnZvANjkA6OiEml2V9pQyvtdT300gCr4IWznfpayV43o4WIHnI\nzQy77IrASbeCtzZE/TuKq0pYC0lsS9HZ3GBhGgsmP432cJIR1K8ROhTCdZ/lCcsN\ni+rfsehuaUglavrnkgjm+9PK4jwadem/7HnO4T+ghaMLfH/KoMGPqpN7HVL1ItXf\n+0N9rYp16nD1TVVHrWFziIKRVWh787ivkw0AeS/iH6dQzCHpc/6a0PWsNGX2Ag9s\n3GCRAp2fXqGs+VMUtXOWZ19AUfUx3aipVJM1MUV+GtRw9r6gmDnO3GdK1g7gxNIT\nb1iigawn8BsEtk+WXs9LTIOgrDW7SQIDAQABoyYwJDASBgNVHRMBAf8ECDAGAQH/\nAgEAMA4GA1UdDwEB/wQEAwICBDANBgkqhkiG9w0BAQsFAAOCAQEAYm2bAuLJxWJE\nur5SFPlStwYI/3qNuzv5DoCyKpLvdJjA29COp/wYtxPPskx04G29JCrhlZAs++Ny\niBlGPCmPZy7XmKJGwvi3YwCWRaHyaXctDRJp2SMtgp0277gjDmGWWTc4WW1yGc4I\nZwWTX00xBLmoZDhxUffR8MrNv7fJmsv2/4CAd2FasNGI71j/+B3Mt6BBo+icu33D\ns+CWjTZMcad/y35r7LW6neUx48unqiDcdzgxxcsdiE8f1pei364bdD+wLSb33ugw\nTUphEFcVV0VWaI6JHPcGrO4O3OKNKk1XCbb+UtN75nAIVc7dfcdigi/skRBI7fyc\nGCPKTZC97Q==\n-----END CERTIFICATE-----\n",
  "bay_uuid": "1f085f23-5206-4192-99da-53c0558626ee",
  "links": [
    {
      "href": "http://65.61.151.130:9511/v1/certificates/1f085f23-5206-4192-99da-53c0558626ee",
      "rel": "self"
    },
    {
      "href": "http://65.61.151.130:9511/certificates/1f085f23-5206-4192-99da-53c0558626ee",
      "rel": "bookmark"
    }
  ]
}
		`)
	})

	b, err := certificates.CreateCredentialsBundle(fake.ServiceClient(), "mccluster")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "1f085f23-5206-4192-99da-53c0558626ee", b.BayID)
	th.AssertEquals(t, "https://172.29.248.97:6443", b.COEEndpoint)
	th.AssertEquals(t, true, b.PrivateKey.Bytes != nil)
	th.AssertEquals(t, true, b.Certificate.Bytes != nil)
	th.AssertEquals(t, true, b.CACertificate.Bytes != nil)
}

func TestCreateCredentialsBundleFailed(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/bays/mccluster", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `
{
  "errors": [
    {
      "status": 404,
      "code": "client",
      "links": [],
      "title": "Bay mccluster could not be found",
      "detail": "Bay mccluster could not be found.",
      "request_id": ""
    }
  ]
}
		`)
	})

	_, err := certificates.CreateCredentialsBundle(fake.ServiceClient(), "mccluster")
	th.AssertEquals(t, "Bay mccluster could not be found.", err.Error())

	er, ok := err.(*fake.ErrorResponse)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, http.StatusNotFound, er.Actual)
}
