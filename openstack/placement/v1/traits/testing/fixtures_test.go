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
const CustomTraitToCreate = "CUSTOM_TRAIT_TO_CREATE"
const CustomTraitToDelete = CustomTraitToCreate
const StandardHardwareTrait = "HW_CPU_X86_AVX"

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

func HandleCreateTraitSuccess(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/traits/"+CustomTraitToCreate,
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PUT")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.WriteHeader(http.StatusCreated)
		})
}

func HandleCreateTraitThatAlreadyExists(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/traits/"+PresentTrait,
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PUT")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.WriteHeader(http.StatusNoContent)
		})
}

// Trait names created via the API must be prefixed with CUSTOM_.
func HandleCreateTraitInvalidName(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/traits/"+AbsentTrait,
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PUT")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.WriteHeader(http.StatusBadRequest)
		})
}

func HandleDeleteTraitSuccess(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/traits/"+CustomTraitToDelete,
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "DELETE")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.WriteHeader(http.StatusNoContent)
		})
}

func HandleDeleteTraitNotFound(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/traits/"+AbsentTrait,
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "DELETE")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.WriteHeader(http.StatusNotFound)
		})
}

func HandleDeleteStandardTraitFailure(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/traits/"+StandardHardwareTrait,
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "DELETE")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.WriteHeader(http.StatusBadRequest)
		})
}

func HandleDeleteTraitInUseFailure(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/traits/"+PresentTrait,
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "DELETE")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.WriteHeader(http.StatusConflict)
		})
}
