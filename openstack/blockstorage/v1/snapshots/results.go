package snapshots

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

type GetResult map[string]interface{}
