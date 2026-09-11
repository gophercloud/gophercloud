package fakeservice

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Widget is the response resource. Its json-tagged fields feed the response
// field pool.
type Widget struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Size        int    `json:"size"`
	Status      string `json:"status"`
}

type commonResult struct {
	gophercloud.Result
}

// CreateResult is the result of a Create call.
type CreateResult struct {
	commonResult
}

// DeleteResult is the result of a Delete call.
type DeleteResult struct {
	gophercloud.ErrResult
}

// ActionResult is the result of an action call.
type ActionResult struct {
	gophercloud.ErrResult
}

// WidgetPage is a single page of Widget results.
type WidgetPage struct {
	pagination.LinkedPageBase
}

// IsEmpty satisfies pagination.Page.
func (p WidgetPage) IsEmpty() (bool, error) {
	return false, nil
}
