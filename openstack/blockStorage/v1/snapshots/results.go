package snapshots

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type Snapshot struct {
	CreatedAt   string
	Description string
	ID          string
	Metadata    map[string]interface{}
	Name        string
	Size        int
	Status      string
	VolumeID    string
}

type GetResult struct {
	err error
	r   map[string]interface{}
}

func (gr GetResult) ExtractSnapshot() (*Snapshot, error) {
	if gr.err != nil {
		return nil, gr.err
	}

	var response struct {
		Snapshot *Snapshot `json:"snapshot"`
	}

	err := mapstructure.Decode(gr.r, &response)
	if err != nil {
		return nil, fmt.Errorf("snapshots: Error decoding snapshot.GetResult: %v", err)
	}
	return response.Snapshot, nil
}
