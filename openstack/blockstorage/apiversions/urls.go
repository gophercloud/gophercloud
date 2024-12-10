package apiversions

import (
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/utils"
)

func listURL(c gophercloud.Client) string {
	baseEndpoint, _ := utils.BaseEndpoint(c.EndpointURL())
	endpoint := strings.TrimRight(baseEndpoint, "/") + "/"
	return endpoint
}
