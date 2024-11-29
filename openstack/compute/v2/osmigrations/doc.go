package osmigrations

/*
Package osmigrations provides the ability to list data on migrations.

Example to List os-migrations:

	pages, err := List(client, nil).AllPages(context.TODO())
	if err != nil {
		panic("fail to get os migration pages")
	}

	osMigrations, err := ExtractOsMigrations(pages)
	if err != nil {
		panic("fail to list os migrations")
	}

	for _, osMigration := range osMigrations {
		fmt.Println(osMigration)
	}

*/
