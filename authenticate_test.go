package gophercloud

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

const SUCCESSFUL_RESPONSE = `{
	"access": {
		"serviceCatalog": [{
			"endpoints": [{
				"publicURL": "https://ord.servers.api.rackspacecloud.com/v2/12345",
				"region": "ORD",
				"tenantId": "12345",
				"versionId": "2",
				"versionInfo": "https://ord.servers.api.rackspacecloud.com/v2",
				"versionList": "https://ord.servers.api.rackspacecloud.com/"
			},{
				"publicURL": "https://dfw.servers.api.rackspacecloud.com/v2/12345",
				"region": "DFW",
				"tenantId": "12345",
				"versionId": "2",
				"versionInfo": "https://dfw.servers.api.rackspacecloud.com/v2",
				"versionList": "https://dfw.servers.api.rackspacecloud.com/"
			}],
			"name": "cloudServersOpenStack",
			"type": "compute"
		},{
			"endpoints": [{
				"publicURL": "https://ord.databases.api.rackspacecloud.com/v1.0/12345",
				"region": "ORD",
				"tenantId": "12345"
			}],
			"name": "cloudDatabases",
			"type": "rax:database"
		}],
		"token": {
			"expires": "2012-04-13T13:15:00.000-05:00",
			"id": "aaaaa-bbbbb-ccccc-dddd"
		},
		"user": {
			"RAX-AUTH:defaultRegion": "DFW",
			"id": "161418",
			"name": "demoauthor",
			"roles": [{
				"description": "User Admin Role.",
				"id": "3",
				"name": "identity:user-admin"
			}]
		}
	}
}
`

type testTransport struct {
	called   int
	response string
}

func (t *testTransport) RoundTrip(req *http.Request) (rsp *http.Response, err error) {
	t.called++

	headers := make(http.Header)
	headers.Add("Content-Type", "application/xml; charset=UTF-8")

	body := ioutil.NopCloser(strings.NewReader(t.response))

	rsp = &http.Response{
		Status:           "200 OK",
		StatusCode:       200,
		Proto:            "HTTP/1.1",
		ProtoMajor:       1,
		ProtoMinor:       1,
		Header:           headers,
		Body:             body,
		ContentLength:    -1,
		TransferEncoding: nil,
		Close:            true,
		Trailer:          nil,
		Request:          req,
	}
	return
}

type tenantIdCheckTransport struct {
	expectTenantId bool
	tenantIdFound  bool
}

func (t *tenantIdCheckTransport) RoundTrip(req *http.Request) (rsp *http.Response, err error) {
	var authContainer *AuthContainer

	headers := make(http.Header)
	headers.Add("Content-Type", "application/xml; charset=UTF-8")

	body := ioutil.NopCloser(strings.NewReader("t.response"))

	rsp = &http.Response{
		Status:           "200 OK",
		StatusCode:       200,
		Proto:            "HTTP/1.1",
		ProtoMajor:       1,
		ProtoMinor:       1,
		Header:           headers,
		Body:             body,
		ContentLength:    -1,
		TransferEncoding: nil,
		Close:            true,
		Trailer:          nil,
		Request:          req,
	}

	bytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &authContainer)
	if err != nil {
		return nil, err
	}
	t.tenantIdFound = (authContainer.Auth.TenantId != "")

	if t.tenantIdFound != t.expectTenantId {
		rsp.Status = "500 Internal Server Error"
		rsp.StatusCode = 500
	}
	return
}

func TestAuthProvider(t *testing.T) {
	c := TestContext()
	tt := &testTransport{}
	c.UseCustomClient(&http.Client{
		Transport: tt,
	})

	_, err := c.Authenticate("", AuthOptions{})
	if err == nil {
		t.Error("Expected error for empty provider string")
		return
	}
	_, err = c.Authenticate("unknown-provider", AuthOptions{Username: "u", Password: "p"})
	if err == nil {
		t.Error("Expected error for unknown service provider")
		return
	}

	err = c.RegisterProvider("provider", &Provider{AuthEndpoint: "/"})
	if err != nil {
		t.Error(err)
		return
	}
	_, err = c.Authenticate("provider", AuthOptions{Username: "u", Password: "p"})
	if err != nil {
		t.Error(err)
		return
	}
	if tt.called != 1 {
		t.Error("Expected transport to be called once.")
		return
	}
}

func TestTenantIdEncoding(t *testing.T) {
	c := TestContext()
	tt := &tenantIdCheckTransport{}
	c.UseCustomClient(&http.Client{
		Transport: tt,
	})
	c.RegisterProvider("provider", &Provider{AuthEndpoint: "/"})

	tt.expectTenantId = false
	_, err := c.Authenticate("provider", AuthOptions{
		Username: "u",
		Password: "p",
	})
	if err != nil {
		t.Error(err)
		return
	}
	if tt.tenantIdFound {
		t.Error("Tenant ID should not have been encoded")
		return
	}

	tt.expectTenantId = true
	_, err = c.Authenticate("provider", AuthOptions{
		Username: "u",
		Password: "p",
		TenantId: "t",
	})
	if err != nil {
		t.Error(err)
		return
	}
	if !tt.tenantIdFound {
		t.Error("Tenant ID should have been encoded")
		return
	}
}

func TestUserNameAndPassword(t *testing.T) {
	c := TestContext()
	c.UseCustomClient(&http.Client{Transport: &testTransport{}})
	c.RegisterProvider("provider", &Provider{AuthEndpoint: "http://localhost/"})

	credentials := []AuthOptions{
		AuthOptions{},
		AuthOptions{Username: "u"},
		AuthOptions{Password: "p"},
	}
	for i, auth := range credentials {
		_, err := c.Authenticate("provider", auth)
		if err == nil {
			t.Error("Expected error from missing credentials (%d)", i)
			return
		}
	}

	_, err := c.Authenticate("provider", AuthOptions{Username: "u", Password: "p"})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestTokenAcquisition(t *testing.T) {
	c := TestContext()
	tt := &testTransport{}
	tt.response = SUCCESSFUL_RESPONSE
	c.UseCustomClient(&http.Client{Transport: tt})
	c.RegisterProvider("provider", &Provider{AuthEndpoint: "http://localhost"})

	acc, err := c.Authenticate("provider", AuthOptions{Username: "u", Password: "p"})
	if err != nil {
		t.Error(err)
		return
	}

	tok := acc.Token
	if (tok.Id == "") || (tok.Expires == "") {
		t.Error("Expected a valid token for successful login; got %s, %s", tok.Id, tok.Expires)
		return
	}
}

func TestServiceCatalogAcquisition(t *testing.T) {
	c := TestContext()
	tt := &testTransport{}
	tt.response = SUCCESSFUL_RESPONSE
	c.UseCustomClient(&http.Client{Transport: tt})
	c.RegisterProvider("provider", &Provider{AuthEndpoint: "http://localhost"})

	acc, err := c.Authenticate("provider", AuthOptions{Username: "u", Password: "p"})
	if err != nil {
		t.Error(err)
		return
	}

	svcs := acc.ServiceCatalog
	if len(svcs) < 2 {
		t.Error("Expected 2 service catalog entries; got %d", len(svcs))
		return
	}

	types := map[string]bool {
		"compute": true,
		"rax:database": true,
	}
	for _, entry := range svcs {
		if !types[entry.Type] {
			t.Error("Expected to find type %s.", entry.Type)
			return
		}
	}
}

func TestUserAcquisition(t *testing.T) {
	c := TestContext()
	tt := &testTransport{}
	tt.response = SUCCESSFUL_RESPONSE
	c.UseCustomClient(&http.Client{Transport: tt})
	c.RegisterProvider("provider", &Provider{AuthEndpoint: "http://localhost"})

	acc, err := c.Authenticate("provider", AuthOptions{Username: "u", Password: "p"})
	if err != nil {
		t.Error(err)
		return
	}

	u := acc.User
	if u.Id != "161418" {
		t.Error("Expected user ID of 16148; got", u.Id)
		return
	}
}