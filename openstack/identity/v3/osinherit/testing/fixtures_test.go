package testing

import (
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func HandleAssignSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/OS-INHERIT/projects/{project_id}/users/{user_id}/roles/{role_id}/inherited_to_projects", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	th.Mux.HandleFunc("/OS-INHERIT/projects/{project_id}/groups/{group_id}/roles/{role_id}/inherited_to_projects", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	th.Mux.HandleFunc("/OS-INHERIT/domains/{domain_id}/users/{user_id}/roles/{role_id}/inherited_to_projects", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	th.Mux.HandleFunc("/OS-INHERIT/domains/{domain_id}/groups/{group_id}/roles/{role_id}/inherited_to_projects", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})
}

func HandleValidateSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/OS-INHERIT/projects/{project_id}/users/{user_id}/roles/{role_id}/inherited_to_projects", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "HEAD")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	th.Mux.HandleFunc("/OS-INHERIT/projects/{project_id}/groups/{group_id}/roles/{role_id}/inherited_to_projects", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "HEAD")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	th.Mux.HandleFunc("/OS-INHERIT/domains/{domain_id}/users/{user_id}/roles/{role_id}/inherited_to_projects", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "HEAD")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	th.Mux.HandleFunc("/OS-INHERIT/domains/{domain_id}/groups/{group_id}/roles/{role_id}/inherited_to_projects", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "HEAD")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})
}

func HandleUnassignSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/OS-INHERIT/projects/{project_id}/users/{user_id}/roles/{role_id}/inherited_to_projects", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	th.Mux.HandleFunc("/OS-INHERIT/projects/{project_id}/groups/{group_id}/roles/{role_id}/inherited_to_projects", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	th.Mux.HandleFunc("/OS-INHERIT/domains/{domain_id}/users/{user_id}/roles/{role_id}/inherited_to_projects", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	th.Mux.HandleFunc("/OS-INHERIT/domains/{domain_id}/groups/{group_id}/roles/{role_id}/inherited_to_projects", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})
}
