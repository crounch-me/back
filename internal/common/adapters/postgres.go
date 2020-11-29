package adapters

import "database/sql"

func GetDatabaseConnection(connectionURI string) *sql.DB {
	db, err := sql.Open("postgres", connectionURI)
	if err != nil {
		panic("Unable to connect to postgres database")
	}
	err = db.Ping()

	if err != nil {
		panic("Unable to ping postgres database")
	}

	return db
}
