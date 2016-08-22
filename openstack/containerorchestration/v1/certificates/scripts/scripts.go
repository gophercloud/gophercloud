package scripts

import (
	"fmt"

	"github.com/gophercloud/gophercloud/openstack/containerorchestration/v1/bays"
)

func Generate(coe string, bay *bays.Bay) map[string][]byte {
	w, err := getScriptWriter(coe)
	if err != nil {
		return nil
	}
	return w.Generate(bay)
}

type scriptWriter interface {
	Generate(bay *bays.Bay) map[string][]byte
}

func getScriptWriter(coe string) (scriptWriter, error) {
	switch coe {
	case "swarm":
		return &swarmWriter{}, nil
	case "kubernetes":
		return &kubernetesWriter{}, nil
	default:
		return nil, fmt.Errorf("Unsupported COE: %s", coe)
	}
}
