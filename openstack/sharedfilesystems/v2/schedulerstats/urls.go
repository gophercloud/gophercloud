package schedulerstats

import "github.com/gophercloud/gophercloud/v2"

func poolsListURL(c gophercloud.Client) string {
	return c.ServiceURL("scheduler-stats", "pools")
}

func poolsListDetailURL(c gophercloud.Client) string {
	return c.ServiceURL("scheduler-stats", "pools", "detail")
}
