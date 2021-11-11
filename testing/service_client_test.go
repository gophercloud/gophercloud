package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestServiceURL(t *testing.T) {
	c := &gophercloud.ServiceClient{Endpoint: "http://123.45.67.8/"}
	expected := "http://123.45.67.8/more/parts/here"
	actual := c.ServiceURL("more", "parts", "here")
	th.CheckEquals(t, expected, actual)
}

func TestMoreHeaders(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := new(gophercloud.ServiceClient)
	c.MoreHeaders = map[string]string{
		"custom": "header",
	}
	c.ProviderClient = new(gophercloud.ProviderClient)
	resp, err := c.Get(fmt.Sprintf("%s/route", th.Endpoint()), nil, nil)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, resp.Request.Header.Get("custom"), "header")
}

func TestServiceContext(t *testing.T) {
	t.Run("copied ServiceClient has original MoreHeaders", func(t *testing.T) {
		th.SetupHTTP()
		defer th.TeardownHTTP()
		th.Mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		var c *gophercloud.ServiceClient
		{
			original := new(gophercloud.ServiceClient)
			original.MoreHeaders = map[string]string{
				"custom": "header",
			}

			c = original.WithContext(context.Background())
		}
		c.ProviderClient = new(gophercloud.ProviderClient)

		resp, err := c.Get(fmt.Sprintf("%s/route", th.Endpoint()), nil, nil)
		th.AssertNoErr(t, err)
		th.AssertEquals(t, resp.Request.Header.Get("custom"), "header")
	})

	t.Run("original ServiceClient is untouched", func(t *testing.T) {
		t.Run("context", func(t *testing.T) {
			th.SetupHTTP()
			defer th.TeardownHTTP()
			th.Mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			original := new(gophercloud.ServiceClient)
			original.MoreHeaders = map[string]string{
				"custom": "header",
			}
			original.ProviderClient = new(gophercloud.ProviderClient)

			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			c := original.WithContext(ctx)

			{
				_, err := c.Get(fmt.Sprintf("%s/route", th.Endpoint()), nil, nil)
				th.AssertErr(t, err)
			}

			{
				_, err := original.Get(fmt.Sprintf("%s/route", th.Endpoint()), nil, nil)
				th.AssertNoErr(t, err)
			}
		})

		t.Run("headers", func(t *testing.T) {
			th.SetupHTTP()
			defer th.TeardownHTTP()
			th.Mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			original := new(gophercloud.ServiceClient)
			original.MoreHeaders = map[string]string{
				"custom": "header",
			}
			original.ProviderClient = new(gophercloud.ProviderClient)

			c := original.WithContext(context.Background())
			delete(c.MoreHeaders, "custom")

			{
				resp, err := c.Get(fmt.Sprintf("%s/route", th.Endpoint()), nil, nil)
				th.AssertNoErr(t, err)
				th.AssertEquals(t, resp.Request.Header.Get("custom"), "")
			}

			{
				resp, err := original.Get(fmt.Sprintf("%s/route", th.Endpoint()), nil, nil)
				th.AssertNoErr(t, err)
				th.AssertEquals(t, resp.Request.Header.Get("custom"), "header")
			}
		})
	})
}
