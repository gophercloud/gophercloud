package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const PresentTrait = "CUSTOM_HW_FPGA_CLASS1"
const AbsentTrait = "NON_EXISTENT_TRAIT"

const TraitsListResultAll = `
{
    "traits": [
        "CUSTOM_HW_FPGA_CLASS1",
        "CUSTOM_HW_FPGA_CLASS2",
        "HW_CPU_X86_AVX"
    ]
}`

const TraitsListFilteredCustomResult = `
{
    "traits": [
        "CUSTOM_HW_FPGA_CLASS1",
        "CUSTOM_HW_FPGA_CLASS2"
    ]
}`

const TraitsListFilteredAssociatedResult = TraitsListResultAll

var ExpectedTraitsListResultAll = []string{
	"CUSTOM_HW_FPGA_CLASS1",
	"CUSTOM_HW_FPGA_CLASS2",
	"HW_CPU_X86_AVX",
}

var ExpectedTraitsListFilteredNameResult = []string{
	"CUSTOM_HW_FPGA_CLASS1",
	"CUSTOM_HW_FPGA_CLASS2",
}

var ExpectedTraitsListFilteredAssociatedResult = ExpectedTraitsListResultAll

func HandleListTraitsAll(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/traits",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, TraitsListResultAll)
		})
}

func HandleListTraitsFilteredName(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/traits",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			th.TestFormValues(t, r, map[string]string{"name": "startswith:CUSTOM"})

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, TraitsListFilteredCustomResult)
		})
}

func HandleListTraitsFilteredAssociated(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/traits",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			th.TestFormValues(t, r, map[string]string{"associated": "true"})

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, TraitsListFilteredAssociatedResult)
		})
}

func HandleGetTraitSuccess(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/traits/"+PresentTrait,
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.WriteHeader(http.StatusNoContent)
		})
}

func HandleGetTraitNotFound(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/traits/"+AbsentTrait,
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.WriteHeader(http.StatusNotFound)
		})
}
