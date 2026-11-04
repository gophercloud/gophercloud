// Package model defines the shared data model used to represent both the
// OpenAPI specification ("spec side") and the Gophercloud implementation
// ("impl side") of an OpenStack API, so the two can be diffed.
package model

// Field is a single request property, response property, or query parameter.
type Field struct {
	// Name is the wire name (JSON property or query-string key).
	Name string `json:"name"`
	// Required reports whether the schema marks this field required.
	Required bool `json:"required,omitempty"`
	// MinVer is the microversion that introduced the field, if known.
	MinVer string `json:"min_ver,omitempty"`
	// MaxVer is the microversion after which the field is unavailable, if any.
	MaxVer string `json:"max_ver,omitempty"`
	// Deprecated is true if the field is marked deprecated.
	Deprecated bool `json:"deprecated,omitempty"`
	// ManualHandled is set on the impl side for fields carried as json:"-"
	// (serialised by hand); such fields exist but cannot be verified by tag.
	ManualHandled bool `json:"manual_handled,omitempty"`
}

// Operation is a single API operation: an HTTP method against a path, with its
// request/response fields and query parameters. Actions (POST .../action) are
// represented as distinct operations distinguished by Action.
type Operation struct {
	// Service is the logical service key (e.g. "compute", "network").
	Service string `json:"service"`
	// Method is the upper-case HTTP verb (GET, POST, PUT, PATCH, DELETE, HEAD).
	Method string `json:"method"`
	// Path is the normalised path relative to the service root, with all path
	// parameters replaced by the "{}" wildcard (e.g. "/servers/{}/action").
	Path string `json:"path"`
	// Action is the action name for .../action endpoints (e.g. "reboot"),
	// otherwise empty.
	Action string `json:"action,omitempty"`
	// MinVer is the microversion that introduced the operation, if known.
	MinVer string `json:"min_ver,omitempty"`

	// OperationID is the OpenAPI operationId (spec side only), for reference.
	OperationID string `json:"operation_id,omitempty"`
	// Source is a human-friendly origin (impl side: package + function).
	Source string `json:"source,omitempty"`

	RequestFields  []Field `json:"request_fields,omitempty"`
	ResponseFields []Field `json:"response_fields,omitempty"`
	QueryParams    []Field `json:"query_params,omitempty"`
}

// Key returns the identity used to match a spec operation to an impl operation:
// method, normalised path, and action.
func (o Operation) Key() string {
	k := o.Method + " " + o.Path
	if o.Action != "" {
		k += ":" + o.Action
	}
	return k
}

// API is the full set of operations discovered for one service, from one side.
type API struct {
	Service    string      `json:"service"`
	Version    string      `json:"version,omitempty"`
	Operations []Operation `json:"operations"`
}
