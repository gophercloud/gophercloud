package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/taas/tapmirrors"
	"github.com/gophercloud/gophercloud/v2/pagination"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/taas/tap_mirrors", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "tap_mirror": {
        "description": "description",
        "directions": {
            "IN": "1",
            "OUT": "2"
        },
        "mirror_type": "erspanv1",
        "name": "test",
        "port_id": "a25290e9-1a54-4c26-a5b3-34458d122acc",
        "remote_ip": "192.168.54.217"
    }
}
		`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, `
{
    "tap_mirror": {
        "id": "bd64a6e3-12b8-4092-a348-6fc7e27c298a",
        "project_id": "6776f022d64443a898ee3fab89dc8c05",
        "name": "test",
        "description": "description",
        "port_id": "a25290e9-1a54-4c26-a5b3-34458d122acc",
        "directions": {
            "IN": "1",
            "OUT": "2"
        },
        "remote_ip": "192.168.54.217",
        "mirror_type": "erspanv1",
        "tenant_id": "6776f022d64443a898ee3fab89dc8c05"
    }
}
    `)
	})

	options := tapmirrors.CreateOpts{
		Name:        "test",
		Description: "description",
		PortID:      "a25290e9-1a54-4c26-a5b3-34458d122acc",
		MirrorType:  tapmirrors.MirrorTypeErspanv1,
		RemoteIP:    "192.168.54.217",
		Directions: tapmirrors.Directions{
			In:  "1",
			Out: "2",
		},
	}
	actual, err := tapmirrors.Create(context.TODO(), fake.ServiceClient(fakeServer), options).Extract()
	th.AssertNoErr(t, err)
	expected := tapmirrors.TapMirror{
		ID:          "bd64a6e3-12b8-4092-a348-6fc7e27c298a",
		TenantID:    "6776f022d64443a898ee3fab89dc8c05",
		ProjectID:   "6776f022d64443a898ee3fab89dc8c05",
		Name:        "test",
		Description: "description",
		PortID:      "a25290e9-1a54-4c26-a5b3-34458d122acc",
		MirrorType:  "erspanv1",
		RemoteIP:    "192.168.54.217",
		Directions: tapmirrors.Directions{
			In:  "1",
			Out: "2",
		},
	}
	th.AssertDeepEquals(t, expected, *actual)
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/taas/tap_mirrors/0837b488-f0e2-4689-99b3-e3ed531f9b10", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `
{
    "tap_mirror": {
        "id": "0837b488-f0e2-4689-99b3-e3ed531f9b10",
        "project_id": "6776f022d64443a898ee3fab89dc8c05",
        "name": "test",
        "description": "description",
        "port_id": "a25290e9-1a54-4c26-a5b3-34458d122acc",
        "directions": {
            "IN": "1",
            "OUT": "2"
        },
        "remote_ip": "192.168.54.217",
        "mirror_type": "erspanv1",
        "tenant_id": "6776f022d64443a898ee3fab89dc8c05"
    }
}
        `)
	})

	actual, err := tapmirrors.Get(context.TODO(), fake.ServiceClient(fakeServer), "0837b488-f0e2-4689-99b3-e3ed531f9b10").Extract()
	th.AssertNoErr(t, err)
	expected := tapmirrors.TapMirror{
		ID:          "0837b488-f0e2-4689-99b3-e3ed531f9b10",
		TenantID:    "6776f022d64443a898ee3fab89dc8c05",
		ProjectID:   "6776f022d64443a898ee3fab89dc8c05",
		Name:        "test",
		Description: "description",
		PortID:      "a25290e9-1a54-4c26-a5b3-34458d122acc",
		MirrorType:  "erspanv1",
		RemoteIP:    "192.168.54.217",
		Directions: tapmirrors.Directions{
			In:  "1",
			Out: "2",
		},
	}
	th.AssertDeepEquals(t, expected, *actual)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/taas/tap_mirrors/0837b488-f0e2-4689-99b3-e3ed531f9b10", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := tapmirrors.Delete(context.TODO(), fake.ServiceClient(fakeServer), "0837b488-f0e2-4689-99b3-e3ed531f9b10")
	th.AssertNoErr(t, res.Err)
}

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/taas/tap_mirrors", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `
{
    "tap_mirrors": [
        {
            "id": "0837b488-f0e2-4689-99b3-e3ed531f9b10",
            "project_id": "6776f022d64443a898ee3fab89dc8c05",
            "name": "test",
            "description": "description",
            "port_id": "a25290e9-1a54-4c26-a5b3-34458d122acc",
            "directions": {
                "IN": "1",
                "OUT": "2"
            },
            "remote_ip": "192.168.54.217",
            "mirror_type": "erspanv1",
            "tenant_id": "6776f022d64443a898ee3fab89dc8c05"
        }
    ]
}
	  `)
	})

	count := 0

	err := tapmirrors.List(fake.ServiceClient(fakeServer), tapmirrors.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := tapmirrors.ExtractTapMirrors(page)
		if err != nil {
			t.Errorf("Failed to extract members: %v", err)
			return false, err
		}
		expected := []tapmirrors.TapMirror{
			{
				ID:          "0837b488-f0e2-4689-99b3-e3ed531f9b10",
				TenantID:    "6776f022d64443a898ee3fab89dc8c05",
				ProjectID:   "6776f022d64443a898ee3fab89dc8c05",
				Name:        "test",
				Description: "description",
				PortID:      "a25290e9-1a54-4c26-a5b3-34458d122acc",
				MirrorType:  "erspanv1",
				RemoteIP:    "192.168.54.217",
				Directions: tapmirrors.Directions{
					In:  "1",
					Out: "2",
				},
			},
		}
		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/taas/tap_mirrors/d031da31-fb9b-4bd9-8d37-aaf04a12d45f", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "tap_mirror": {
        "name": "new name",
        "description": "new description"
    }
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `
{
    "tap_mirror": {
        "id": "d031da31-fb9b-4bd9-8d37-aaf04a12d45f",
        "project_id": "6776f022d64443a898ee3fab89dc8c05",
        "name": "new name",
        "description": "new description",
        "port_id": "a25290e9-1a54-4c26-a5b3-34458d122acc",
        "directions": {
            "IN": "1",
            "OUT": "2"
        },
        "remote_ip": "192.168.54.217",
        "mirror_type": "erspanv1",
        "tenant_id": "6776f022d64443a898ee3fab89dc8c05"
    }
}
`)
	})

	updatedName := "new name"
	updatedDescription := "new description"
	options := tapmirrors.UpdateOpts{
		Name:        &updatedName,
		Description: &updatedDescription,
	}

	actual, err := tapmirrors.Update(context.TODO(), fake.ServiceClient(fakeServer), "d031da31-fb9b-4bd9-8d37-aaf04a12d45f", options).Extract()
	th.AssertNoErr(t, err)
	expected := tapmirrors.TapMirror{
		ID:          "d031da31-fb9b-4bd9-8d37-aaf04a12d45f",
		TenantID:    "6776f022d64443a898ee3fab89dc8c05",
		ProjectID:   "6776f022d64443a898ee3fab89dc8c05",
		Name:        "new name",
		Description: "new description",
		PortID:      "a25290e9-1a54-4c26-a5b3-34458d122acc",
		MirrorType:  "erspanv1",
		RemoteIP:    "192.168.54.217",
		Directions: tapmirrors.Directions{
			In:  "1",
			Out: "2",
		},
	}
	th.AssertDeepEquals(t, expected, *actual)
}
