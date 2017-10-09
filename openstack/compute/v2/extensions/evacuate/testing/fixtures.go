package testing

import (
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func mockEvacuateResponse(t *testing.T, id string) {
	th.Mux.HandleFunc("/servers/"+id+"/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `
		{
		    "evacuate": {
		    	"OnSharedStorage": false,
			    "adminPass": "true",
			    "host": "derp"
		    }
		}
		      `)
		w.WriteHeader(http.StatusAccepted)
	})
}

func mockEvacuateResponseWithHost(t *testing.T, id string) {
	th.Mux.HandleFunc("/servers/"+id+"/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `
		{
		    "evacuate": {
		    	"OnSharedStorage": false,
			    "host": "derp"
		    }
		}
		      `)
		w.WriteHeader(http.StatusAccepted)
	})
}

func mockEvacuateResponseWithNoOpts(t *testing.T, id string) {
	th.Mux.HandleFunc("/servers/"+id+"/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `
		{
		    "evacuate": {
		    	"OnSharedStorage": false
		  }
		}
		      `)
		w.WriteHeader(http.StatusAccepted)
	})
}
