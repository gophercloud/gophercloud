package bmvolume

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type ListOpts struct {
}

func List(client *gophercloud.ServiceClient, opts ListOpts) (r ListResult) {
	resp, err := client.Get(listURL(client), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type ListConnectorsOptsBuilder interface {
	ToConnectorsListQuery() (string, error)
}
type ListConnectorsOpts struct {
	// node uuid
	Node string `q:"node"`
	// One or more fields to be returned in the response.
	Fields []string `q:"fields"`
	// Requests a page size of items.
	Limit int `q:"limit"`
	// The ID of the last-seen item
	Marker string `q:"marker"`
	// Sorts the response by the requested sort direction.
	// Valid value is asc (ascending) or desc (descending). Default is asc.
	SortDir string `q:"sort_dir"`
	// Sorts the response by the this attribute value. Default is id.
	SortKey string `q:"sort_key"`
}

func (opts ListConnectorsOpts) ToConnectorsListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

func ListConnectors(client *gophercloud.ServiceClient, opts ListConnectorsOptsBuilder) pagination.Pager {
	url := listConnectorsURL(client)
	if opts != nil {
		query, err := opts.ToConnectorsListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ConnectorPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

type CreateConnectorOptsBuilder interface {
	ToConnectorCreateMap() (map[string]interface{}, error)
}

type CreateConnectorOpts struct {
	NodeUUID      string                 `json:"node_uuid,omitempty"`
	ConnectorType string                 `json:"type,omitempty"`
	ConnectorId   string                 `json:"connector_id,omitempty"`
	Extra         map[string]interface{} `json:"Extra,omitempty"`
	UUID          string                 `json:"uuid,omitempty"`
}

func (opts CreateConnectorOpts) ToConnectorCreateMap() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return body, nil
}

func CreateConnector(client *gophercloud.ServiceClient, opts CreateConnectorOptsBuilder) (r CreateConnectorResult) {
	reqBody, err := opts.ToConnectorCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(createConnectorsURL(client), reqBody, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type Patch interface {
	ToUpdateMap() map[string]interface{}
}

type UpdateOpts []Patch

type UpdateOp string

const (
	ReplaceOp UpdateOp = "replace"
	AddOp     UpdateOp = "add"
	RemoveOp  UpdateOp = "remove"
)

type UpdateOperation struct {
	Op    UpdateOp    `json:"op" required:"true"`
	Path  string      `json:"path" required:"true"`
	Value interface{} `json:"value,omitempty"`
}

func (opts UpdateOperation) ToUpdateMap() map[string]interface{} {
	return map[string]interface{}{
		"op":    opts.Op,
		"path":  opts.Path,
		"value": opts.Value,
	}
}

func UpdateConnector(client *gophercloud.ServiceClient, id string, opts UpdateOpts) (r UpdateConnectorResult) {
	body := make([]map[string]interface{}, len(opts))
	for i, patch := range opts {
		body[i] = patch.ToUpdateMap()
	}
	resp, err := client.Patch(patchConnectorsURL(client, id), body, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func DeleteConnector(client *gophercloud.ServiceClient, id string) (r DeleteConnectorResult) {
	resp, err := client.Delete(deleteConnectorsURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func GetConnector(client *gophercloud.ServiceClient, id string) (r GetConnectorResult) {
	resp, err := client.Get(getConnectorsURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type ListTargetsOptsBuilder interface {
	ToTargetsListQuery() (string, error)
}
type ListTargetsOpts struct {
	// node uuid
	Node string `q:"node"`
	// One or more fields to be returned in the response.
	Fields []string `q:"fields"`
	// Requests a page size of items.
	Limit int `q:"limit"`
	// The ID of the last-seen item
	Marker string `q:"marker"`
	// Sorts the response by the requested sort direction.
	// Valid value is asc (ascending) or desc (descending). Default is asc.
	SortDir string `q:"sort_dir"`
	// Sorts the response by the this attribute value. Default is id.
	SortKey string `q:"sort_key"`
}

func (opts ListTargetsOpts) ToTargetsListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

func ListTargets(client *gophercloud.ServiceClient, opts ListTargetsOptsBuilder) pagination.Pager {
	url := listTargetsURL(client)
	if opts != nil {
		query, err := opts.ToTargetsListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return TargetPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

type CreateTargetOptsBuilder interface {
	ToTargetCreateMap() (map[string]interface{}, error)
}

type CreateTargetOpts struct {
	NodeUUID   string                 `json:"node_uuid,omitempty"`
	VolumeType string                 `json:"volume_type,omitempty"` //iscsi, fibre_channel
	Properties map[string]interface{} `json:"properties,omitempty"`
	BootIndex  string                 `json:"boot_index,omitempty"`
	VolumeId   string                 `json:"volume_id,omitempty"`
	Extra      map[string]interface{} `json:"Extra,omitempty"`
	UUID       string                 `json:"uuid,omitempty"`
}

func (opts CreateTargetOpts) ToTargetCreateMap() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return body, nil
}

func CreateTarget(client *gophercloud.ServiceClient, opts CreateTargetOptsBuilder) (r CreateTargetResult) {
	reqBody, err := opts.ToTargetCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(createTargetsURL(client), reqBody, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func UpdateTarget(client *gophercloud.ServiceClient, id string, opts UpdateOpts) (r UpdateTargetResult) {
	body := make([]map[string]interface{}, len(opts))
	for i, patch := range opts {
		body[i] = patch.ToUpdateMap()
	}
	resp, err := client.Patch(patchTargetsURL(client, id), body, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func DeleteTarget(client *gophercloud.ServiceClient, id string) (r DeleteTargetResult) {
	resp, err := client.Delete(deleteTargetsURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func GetTarget(client *gophercloud.ServiceClient, id string) (r GetTargetResult) {
	resp, err := client.Get(getTargetsURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
