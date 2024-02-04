/*
Package snapshotextensions provides the ability to extend a snapshot
result with tenant/project information. Example:

	type SnapshotWithExt struct {
		snapshots.Snapshot
		snapshotextensions.SnapshotExtension
	}

	var allSnapshot []SnapshotWithExt

	allPages, err := snapshots.List(client, nil).AllPages()
	if err != nil {
		panic("Unable to retrieve snapshots: %s", err)
	}

	err = snapshots.ExtractVolumesInto(allPages, &allSnapshots)
	if err != nil {
		panic("Unable to extract snapshots: %s", err)
	}

	for _, snapshot := range allSnapshots {
		fmt.Println(snapshot.ProjectId)
	}
*/
package snapshotextensions
