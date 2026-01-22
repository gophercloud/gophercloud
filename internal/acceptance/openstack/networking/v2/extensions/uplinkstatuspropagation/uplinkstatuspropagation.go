package uplinkstatuspropagation

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/uplinkstatuspropagation"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// PortWithUplinkStatusPropagationExt represents a port with the uplink status
// propagation extension fields.
type PortWithUplinkStatusPropagationExt struct {
	ports.Port
	uplinkstatuspropagation.PortPropagateUplinkStatusExt
}

// CreatePortWithUplinkStatusPropagation will create a port on the specified
// subnet with uplink status propagation enabled. An error will be returned if
// the port could not be created.
func CreatePortWithUplinkStatusPropagation(t *testing.T, client *gophercloud.ServiceClient, networkID, subnetID string, propagate bool) (PortWithUplinkStatusPropagationExt, error) {
	portName := tools.RandomString("TESTACC-", 8)
	portDescription := tools.RandomString("TESTACC-PORT-DESC-", 8)
	iFalse := false

	t.Logf("Attempting to create port with uplink status propagation: %s", portName)

	portCreateOpts := ports.CreateOpts{
		NetworkID:    networkID,
		Name:         portName,
		Description:  portDescription,
		AdminStateUp: &iFalse,
		FixedIPs:     []ports.IP{{SubnetID: subnetID}},
	}

	createOpts := uplinkstatuspropagation.PortPropagateUplinkStatusCreateOptsExt{
		CreateOptsBuilder:     portCreateOpts,
		PropagateUplinkStatus: &propagate,
	}

	var s PortWithUplinkStatusPropagationExt

	err := ports.Create(context.TODO(), client, createOpts).ExtractInto(&s)
	if err != nil {
		return s, err
	}

	t.Logf("Successfully created port: %s", portName)

	th.AssertEquals(t, s.Name, portName)
	th.AssertEquals(t, s.Description, portDescription)
	if s.PropagateUplinkStatus != nil {
		th.AssertEquals(t, *s.PropagateUplinkStatus, propagate)
	}

	return s, nil
}
