/*
Package backups provides information and interaction with backups in the
OpenStack Block Storage service. A backup is a point in time copy of the
data contained in an external storage volume, and can be controlled
programmatically.

Example to List Backups

	listOpts := backups.ListOpts{
		VolumeID: "uuid",
	}

	allPages, err := backups.List(client, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allBackups, err := backups.ExtractBackups(allPages)
	if err != nil {
		panic(err)
	}

	for _, backup := range allBackups {
		fmt.Println(backup)
	}

Example to Create a Backup

	createOpts := backups.CreateOpts{
		VolumeID: "uuid",
		Name:     "my-backup",
	}

	backup, err := backups.Create(client, createOpts).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Println(backup)

Example to Update a Backup

	updateOpts := backups.UpdateOpts{
		Name: "new-name",
	}

	backup, err := backups.Update(client, "uuid", updateOpts).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Println(backup)

Example to Restore a Backup to a Volume

	options := backups.RestoreOpts{
		VolumeID: "1234",
		Name:     "vol-001",
	}

	restore, err := backups.RestoreFromBackup(client, "uuid", options).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Println(restore)

Example to Delete a Backup

	err := backups.Delete(client, "uuid").ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Export a Backup

	export, err := backups.Export(client, "uuid").Extract()
	if err != nil {
		panic(err)
	}

	fmt.Println(export)

Example to Import a Backup

	availabilityZone := "region1b"
	backupURL, _ := json.Marshal(backups.ImportBackup{
		ID:               "d32019d3-bc6e-4319-9c1d-6722fc136a22",
		Status:           "available",
		AvailabilityZone: &availabilityZone,
		VolumeID:         "cf9bc6fa-c5bc-41f6-bc4e-6e76c0bea959",
		UpdatedAt:        time.Date(2020, 3, 11, 19, 29, 8, 0, time.UTC),
		Host:             "cinder-backup-host1",
		UserID:           "93514e04-a026-4f60-8176-395c859501dd",
		ServiceMetadata:  "volume_cf9bc6fa-c5bc-41f6-bc4e-6e76c0bea959/20200311192855/az_regionb_backup_b87bb1e5-0d4e-445e-a548-5ae742562bac",
		Size:             1,
		ObjectCount:      2,
		Container:        "my-test-backup",
		Service:          "cinder.backup.drivers.swift.SwiftBackupDriver",
		CreatedAt:        time.Date(2020, 3, 11, 19, 25, 24, 0, time.UTC),
		DataTimestamp:    time.Date(2020, 3, 11, 19, 25, 24, 0, time.UTC),
		ProjectID:        "14f1c1f5d12b4755b94edef78ff8b325",
	})

	options := backups.ImportOpts{
		BackupService: "cinder.backup.drivers.swift.SwiftBackupDriver",
		BackupURL:     backupURL,
	}

	backup, err := backups.Import(client, options).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Println(backup)
*/
package backups
