package testing

import (
	"io"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// HandlePutImageDataSuccessfully setup
func HandlePutImageDataSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/images/da3b75d9-3f4a-40e7-8a2c-bfab23927dea/file", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		b, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Unable to read request body: %v", err)
		}

		th.AssertByteArrayEquals(t, []byte{5, 3, 7, 24}, b)

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleStageImageDataSuccessfully setup
func HandleStageImageDataSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/images/da3b75d9-3f4a-40e7-8a2c-bfab23927dea/stage", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		b, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Unable to read request body: %v", err)
		}

		th.AssertByteArrayEquals(t, []byte{5, 3, 7, 24}, b)

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleGetImageDataSuccessfully setup
func HandleGetImageDataSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/images/da3b75d9-3f4a-40e7-8a2c-bfab23927dea/file", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte{34, 87, 0, 23, 23, 23, 56, 255, 254, 0})
		th.AssertNoErr(t, err)
	})
}
