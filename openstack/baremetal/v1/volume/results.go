package bmvolume

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type volumeResult struct {
	gophercloud.Result
}

type ListResult struct {
	volumeResult
}

func (r ListResult) Extract() (*Volume, error) {
	var s Volume
	err := r.ExtractInto(&s)
	return &s, err
}

type Volume struct {
	Connectors []interface{} `json:"connectors"`
	Targets    []interface{} `json:"targets"`
	Links      []interface{} `json:"links"`
}

type Connector struct {
	UUID          string                 `json:"uuid"`
	ConnectorType string                 `json:"type"`
	ConnectorId   string                 `json:"connector_id"`
	NodeUUID      string                 `json:"node_uuid"`
	Extra         map[string]interface{} `json:"extra"`
	Links         []interface{}          `json:"links"`
}

type connectorResult struct {
	gophercloud.Result
}

type ConnectorPage struct {
	pagination.LinkedPageBase
}

type Connectors struct {
	Connectors []Connector `json:"connectors"`
}

func ExtractConnectorsInto(r pagination.Page, v interface{}) error {
	return r.(ConnectorPage).Result.ExtractIntoSlicePtr(v, "connectors")
}
func ExtractConnectors(r pagination.Page) ([]Connector, error) {
	var s []Connector
	err := ExtractConnectorsInto(r, &s)
	return s, err
}

func (r ConnectorPage) IsEmpty() (bool, error) {
	s, err := ExtractConnectors(r)
	return len(s) == 0, err
}

func (r ConnectorPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"connector_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

type CreateConnectorResult struct {
	connectorResult
}

func (r CreateConnectorResult) Extract() (*Connector, error) {
	var c Connector
	err := r.ExtractInto(&c)
	return &c, err
}

type UpdateConnectorResult struct {
	connectorResult
}

func (r UpdateConnectorResult) Extract() (*Connector, error) {
	var c Connector
	err := r.ExtractInto(&c)
	return &c, err
}

type DeleteConnectorResult struct {
	gophercloud.ErrResult
}

func (r DeleteConnectorResult) ExtractErr() error {
	return r.ExtractErr()
}

type GetConnectorResult struct {
	connectorResult
}

func (r GetConnectorResult) Extract() (*Connector, error) {
	var c Connector
	err := r.ExtractInto(&c)
	return &c, err
}

type targetResult struct {
	gophercloud.Result
}

type Target struct {
	UUID       string                 `json:"uuid"`
	VolumeType string                 `json:"volume_type"`
	Properties map[string]interface{} `json:"properties"`
	BootIndex  string                 `json:"boot_index"`
	VolumeId   string                 `json:"volume_id"`
	Extra      map[string]interface{} `json:"extra"`
	NodeUUID   string                 `json:"node_uuid"`
	Links      []interface{}          `json:"links"`
}
type TargetPage struct {
	pagination.LinkedPageBase
}

func ExtractTargetsInto(r pagination.Page, v interface{}) error {
	return r.(TargetPage).Result.ExtractIntoSlicePtr(v, "targets")
}
func ExtractTargets(r pagination.Page) ([]Target, error) {
	var s []Target
	err := ExtractTargetsInto(r, &s)
	return s, err
}

func (r TargetPage) IsEmpty() (bool, error) {
	s, err := ExtractTargets(r)
	return len(s) == 0, err
}

func (r TargetPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"target_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

type CreateTargetResult struct {
	targetResult
}

func (r CreateTargetResult) Extract() (*Target, error) {
	var t Target
	err := r.ExtractInto(&t)
	return &t, err
}

type UpdateTargetResult struct {
	targetResult
}

func (r UpdateTargetResult) Extract() (*Target, error) {
	var t Target
	err := r.ExtractInto(&t)
	return &t, err
}

type DeleteTargetResult struct {
	gophercloud.ErrResult
}

func (r DeleteTargetResult) ExtractErr() error {
	return r.ExtractErr()
}

type GetTargetResult struct {
	targetResult
}

func (r GetTargetResult) Extract() (*Target, error) {
	var t Target
	err := r.ExtractInto(&t)
	return &t, err
}
