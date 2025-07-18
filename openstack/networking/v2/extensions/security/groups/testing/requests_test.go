package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/security-groups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, SecurityGroupListResponse)
	})

	count := 0

	err := groups.List(fake.ServiceClient(fakeServer), groups.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := groups.ExtractGroups(page)
		if err != nil {
			t.Errorf("Failed to extract secgroups: %v", err)
			return false, err
		}

		expected := []groups.SecGroup{SecurityGroup1}
		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/security-groups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, SecurityGroupCreateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, SecurityGroupCreateResponse)
	})

	opts := groups.CreateOpts{Name: "new-webservers", Description: "security group for webservers"}
	_, err := groups.Create(context.TODO(), fake.ServiceClient(fakeServer), opts).Extract()
	th.AssertNoErr(t, err)
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/security-groups/2076db17-a522-4506-91de-c6dd8e837028",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PUT")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, SecurityGroupUpdateRequest)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, SecurityGroupUpdateResponse)
		})

	name := "newer-webservers"
	opts := groups.UpdateOpts{Name: &name}
	sg, err := groups.Update(context.TODO(), fake.ServiceClient(fakeServer), "2076db17-a522-4506-91de-c6dd8e837028", opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "newer-webservers", sg.Name)
	th.AssertEquals(t, "security group for webservers", sg.Description)
	th.AssertEquals(t, "2076db17-a522-4506-91de-c6dd8e837028", sg.ID)
	th.AssertEquals(t, "2019-06-30T04:15:37Z", sg.CreatedAt.Format(time.RFC3339))
	th.AssertEquals(t, "2019-06-30T05:18:49Z", sg.UpdatedAt.Format(time.RFC3339))
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/security-groups/85cc3048-abc3-43cc-89b3-377341426ac5", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, SecurityGroupGetResponse)
	})

	sg, err := groups.Get(context.TODO(), fake.ServiceClient(fakeServer), "85cc3048-abc3-43cc-89b3-377341426ac5").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "default", sg.Description)
	th.AssertEquals(t, "85cc3048-abc3-43cc-89b3-377341426ac5", sg.ID)
	th.AssertEquals(t, "default", sg.Name)
	th.AssertEquals(t, 2, len(sg.Rules))
	th.AssertEquals(t, "e4f50856753b4dc6afee5fa6b9b6c550", sg.TenantID)
	th.AssertEquals(t, "2019-06-30T04:15:37Z", sg.CreatedAt.Format(time.RFC3339))
	th.AssertEquals(t, "2019-06-30T05:18:49Z", sg.UpdatedAt.Format(time.RFC3339))
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/security-groups/4ec89087-d057-4e2c-911f-60a3b47ee304", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := groups.Delete(context.TODO(), fake.ServiceClient(fakeServer), "4ec89087-d057-4e2c-911f-60a3b47ee304")
	th.AssertNoErr(t, res.Err)
}
