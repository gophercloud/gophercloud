package vlantransparent

import (
	"net/url"
	"strconv"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
)

// ListOptsExt adds the vlan-transparent network options to the base ListOpts.
type ListOptsExt struct {
	networks.ListOptsBuilder
	VLANTransparent *bool `q:"vlan_transparent"`
}

// ToNetworkListQuery adds the vlan_transparent option to the base network
// list options.
func (opts ListOptsExt) ToNetworkListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts.ListOptsBuilder)
	if err != nil {
		return "", err
	}

	params := q.Query()
	if opts.VLANTransparent != nil {
		v := strconv.FormatBool(*opts.VLANTransparent)
		params.Add("vlan_transparent", v)
	}

	q = &url.URL{RawQuery: params.Encode()}
	return q.String(), err
}
