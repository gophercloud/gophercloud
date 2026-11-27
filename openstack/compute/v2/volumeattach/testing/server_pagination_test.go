package testing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/volumeattach"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// Nova slices all block-device mappings before excluding non-volume devices.
// Keep the single-page contract: returned attachment counts cannot determine
// the next offset, and an empty response does not prove there are no more BDMs.
// https://github.com/openstack/nova/blob/stable/2026.1/nova/api/openstack/compute/volume_attachments.py
func TestVolumeAttachmentListSinglePageMixedBlockDevices(t *testing.T) {
	for _, tc := range []struct {
		name     string
		devices  []string
		pageSize int
		expected []string
	}{
		{name: "LocalBeforeVolume", devices: []string{"", "1", "2"}, pageSize: 3, expected: []string{"1", "2"}},
		{name: "ServerCap", devices: []string{"", "1", "2"}, pageSize: 2, expected: []string{"1"}},
		{name: "EmptyFirstPage", devices: []string{"", "", "1"}, pageSize: 2, expected: []string{}},
		{name: "EmptyIntermediatePage", devices: []string{"1", "2", "", "", "3"}, pageSize: 2, expected: []string{"1", "2"}},
	} {
		for _, allPages := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/AllPages=%t", tc.name, allPages), func(t *testing.T) {
				fakeServer := th.SetupHTTP()
				defer fakeServer.Teardown()
				requests := 0
				fakeServer.Mux.HandleFunc("/servers/server/os-volume_attachments", func(w http.ResponseWriter, r *http.Request) {
					requests++
					th.TestMethod(t, r, "GET")
					th.CheckEquals(t, "", r.URL.RawQuery)
					if requests > 1 {
						http.Error(w, "unexpected pagination request", http.StatusInternalServerError)
						return
					}
					items := []map[string]string{}
					for _, id := range tc.devices[:min(tc.pageSize, len(tc.devices))] {
						if id != "" {
							items = append(items, map[string]string{"id": id, "volumeId": id, "serverId": "server"})
						}
					}
					w.Header().Set("Content-Type", "application/json")
					th.AssertNoErr(t, json.NewEncoder(w).Encode(map[string]any{"volumeAttachments": items}))
				})

				pager := volumeattach.List(client.ServiceClient(fakeServer), "server")
				actual := []string{}
				collect := func(page pagination.Page) error {
					items, err := volumeattach.ExtractVolumeAttachments(page)
					if err != nil {
						return err
					}
					for _, item := range items {
						actual = append(actual, item.VolumeID)
					}
					return nil
				}
				if allPages {
					page, err := pager.AllPages(context.Background())
					th.AssertNoErr(t, err)
					th.AssertNoErr(t, collect(page))
				} else {
					err := pager.EachPage(context.Background(), func(_ context.Context, page pagination.Page) (bool, error) {
						return true, collect(page)
					})
					th.AssertNoErr(t, err)
				}
				th.CheckDeepEquals(t, tc.expected, actual)
				th.CheckEquals(t, 1, requests)
			})
		}
	}
}
