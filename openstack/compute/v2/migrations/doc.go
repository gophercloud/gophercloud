package migrations

/*
Package migrations provides the ability to list data on migrations.

Example to List os-migrations:

	pages, err := List(client, nil).AllPages(context.TODO())
	if err != nil {
		panic("fail to get migration pages")
	}

	migrations, err := ExtractMigrations(pages)
	if err != nil {
		panic("fail to list migrations")
	}

	for _, migration := range migrations {
		fmt.Println(migration)
	}

*/
