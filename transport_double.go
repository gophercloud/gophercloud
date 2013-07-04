package gophercloud

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type transport struct {
	called         int
	response       string
	expectTenantId bool
	tenantIdFound  bool
}

func (t *transport) RoundTrip(req *http.Request) (rsp *http.Response, err error) {
	var authContainer *AuthContainer

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

func newTransport() *transport {
	return &transport{}
}

func (t *transport) IgnoreTenantId() *transport {
	t.expectTenantId = false
	return t
}

func (t *transport) ExpectTenantId() *transport {
	t.expectTenantId = true
	return t
}

func (t *transport) WithResponse(r string) *transport {
	t.response = r
	return t
}
