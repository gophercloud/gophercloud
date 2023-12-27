package snapshotextensions

// SnapshotExtension is an extension to the base Snapshot object
type SnapshotExtension struct {
	// A percentage value for the build progress.
	Progress string `json:"os-extended-snapshot-attributes:progress"`
	// The UUID of the owning project.
	ProjectId string `json:"os-extended-snapshot-attributes:project_id"`
}
