package instancelogs

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// List retrieves all the logs for an instance.
func List(client *gophercloud.ServiceClient, id string) pagination.Pager {
	return pagination.NewPager(client, baseURL(client, id), func(r pagination.PageResult) pagination.Page {
		return LogPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Show details for a instance log.
func Show(client *gophercloud.ServiceClient, id, name string) (r ActionResult) {
	b := map[string]interface{}{"name": name}
	resp, err := client.Post(baseURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{OkCodes: []int{200}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Enable a log type for a instance.
func Enable(client *gophercloud.ServiceClient, id, name string) (r ActionResult) {
	b := map[string]interface{}{"name": name, "enable": 1}
	resp, err := client.Post(baseURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{OkCodes: []int{200}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Disable a log type for a instance.
func Disable(client *gophercloud.ServiceClient, id, name string) (r ActionResult) {
	b := map[string]interface{}{"name": name, "disable": 1}
	resp, err := client.Post(baseURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{OkCodes: []int{200}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Publish a log type for a instance.
// Publish will automatically enable a log.
func Publish(client *gophercloud.ServiceClient, id, name string) (r ActionResult) {
	b := map[string]interface{}{"name": name, "publish": 1}
	resp, err := client.Post(baseURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{OkCodes: []int{200}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Discard all previously published logs for a instance.
func Discard(client *gophercloud.ServiceClient, id, name string) (r ActionResult) {
	b := map[string]interface{}{"name": name, "discard": 1}
	resp, err := client.Post(baseURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{OkCodes: []int{200}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
