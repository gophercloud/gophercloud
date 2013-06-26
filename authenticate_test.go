package gophercloud

import (
	"testing"
	"net/http"
	"io/ioutil"
	"strings"
)


type testTransport struct {
	called int
	response string
}

func (t *testTransport) RoundTrip(req *http.Request) (rsp *http.Response, err error) {
	t.called++;

	headers := make(http.Header)
	headers.Add("Content-Type", "application/xml; charset=UTF-8")

	body := ioutil.NopCloser(strings.NewReader(t.response))

	rsp = &http.Response{
		Status: "200 OK",
		StatusCode: 200,
		Proto: "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header: headers,
		Body: body,
		ContentLength: -1,
		TransferEncoding: nil,
		Close: true,
		Trailer: nil,
		Request: req,
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

func TestUserNameAndPassword(t *testing.T) {
	c := TestContext()
	c.RegisterProvider("provider", &Provider{AuthEndpoint: "/"})

	auths := []AuthOptions{
		AuthOptions{},
		AuthOptions{Username: "u"},
		AuthOptions{Password: "p"},
	}
	for i, auth := range auths {
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
