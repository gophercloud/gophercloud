package resetstate

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/extensions"
)

// ServerState refers to the states usable in ResetState Action
type ServerState string

const (
	// StateActive returns the state of the server as active
	StateActive ServerState = "active"

	// StateError returns the state of the server as error
	StateError ServerState = "error"
)

// ResetState will reset the state of a server
func ResetState(ctx context.Context, client *gophercloud.ServiceClient, id string, state ServerState) (r ResetResult) {
	stateMap := map[string]interface{}{"state": state}
	resp, err := client.PostWithContext(ctx, extensions.ActionURL(client, id), map[string]interface{}{"os-resetState": stateMap}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
