package postgres

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	// Necessary to have pg driver to connect to the database
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/crounch-me/back/storage"
	log "github.com/sirupsen/logrus"
)

const packageName = "postgres"

// PostgresStorage handle postgres database connection
type PostgresStorage struct {
	session *sql.DB
	schema  string
}

// NewStorage create a new connection to the postgres database
func NewStorage(connectionURI, schema string) storage.Storage {
	db, err := sql.Open("postgres", connectionURI)
	if err != nil {
		log.WithError(err).Fatal("Unable to connect to postgres database")
	}
	err = db.Ping()
	if err != nil {
		log.WithError(err).Fatal("Unable to ping postgres database")
	}
	return &PostgresStorage{
		session: db,
		schema:  schema,
	}
}

// InitDB initialize data in database
func InitDB(connectionURI string) {
	log.Debug("Initializing database")
	m, err := migrate.New(
		"file://init-db",
		connectionURI,
	)
	if err != nil {
		log.WithError(err).Fatal("Unable to connect to database for initialization")
	}

	if err := m.Up(); err != nil && err.Error() != migrate.ErrNoChange.Error() {
		log.WithError(err).Fatal("Unable to initialize database")
	}
}
