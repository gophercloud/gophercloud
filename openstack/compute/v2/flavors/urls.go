package flavors

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/rackspace/gophercloud"
)

func listURL(client *gophercloud.ServiceClient, lfo ListFilterOptions) string {
	v := url.Values{}
	if lfo.ChangesSince != "" {
		v.Set("changes-since", lfo.ChangesSince)
	}
	if lfo.MinDisk != 0 {
		v.Set("minDisk", strconv.Itoa(lfo.MinDisk))
	}
	if lfo.MinRAM != 0 {
		v.Set("minRam", strconv.Itoa(lfo.MinRAM))
	}
	if lfo.Marker != "" {
		v.Set("marker", lfo.Marker)
	}
	if lfo.Limit != 0 {
		v.Set("limit", strconv.Itoa(lfo.Limit))
	}
	tail := ""
	if len(v) > 0 {
		tail = fmt.Sprintf("?%s", v.Encode())
	}
	return client.ServiceURL("flavors", "detail") + tail
}

func flavorURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("flavors", id)
}
