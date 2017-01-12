package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func MockCreateResponse(t *testing.T) {
	th.Mux.HandleFunc("/security-services", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
        {
            "security_service": {
                "description": "Creating my first Security Service",
                "dns_ip": "10.0.0.0/24",
                "user": "demo",
                "password": "***",
                "type": "kerberos",
                "name": "SecServ1"
            }
        }`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
        {
            "security_service": {
                "status": "new",
                "domain": null,
                "project_id": "16e1ab15c35a457e9c2b2aa189f544e1",
                "name": "SecServ1",
                "created_at": "2015-09-07T12:19:10.695211",
                "updated_at": null,
                "server": null,
                "dns_ip": "10.0.0.0/24",
                "user": "demo",
                "password": "supersecret",
                "type": "kerberos",
                "id": "3c829734-0679-4c17-9637-801da48c0d5f",
                "description": "Creating my first Security Service"
            }
        }`)
	})
}

func MockDeleteResponse(t *testing.T) {
	th.Mux.HandleFunc("/security-services/securityServiceID", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusAccepted)
	})
}
