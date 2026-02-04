package uplinkstatuspropagation

type PortPropagateUplinkStatusExt struct {
	// PropagateUplinkStatus specifies whether port uplink status
	// propagation is enabled or disabled.
	PropagateUplinkStatus *bool `json:"propagate_uplink_status"`
}
