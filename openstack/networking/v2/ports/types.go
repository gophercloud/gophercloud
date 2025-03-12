package ports

type PortStatus string

const (
	PortStatusActive PortStatus = "ACTIVE"
	PortStatusBuild  PortStatus = "BUILD"
	PortStatusDown   PortStatus = "DOWN"
	PortStatusError  PortStatus = "ERROR"
)
