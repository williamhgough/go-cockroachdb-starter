package cockroach

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/cockroachdb"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"

	"github.com/williamhgough/go-cockroachdb-starter/internal/repository/cockroach/migrations"
)

//go:generate go-bindata -pkg migrations -prefix migrations -nometadata -ignore bindata -ignore BUILD -o ./migrations/bindata.go ./migrations

// Repository is a data store backed by a SQL database
type Repository struct {
	db *sql.DB
	sb squirrel.StatementBuilderType
}

// NewRepository creates a data store backed by an SQL database.
func NewRepository(dbDSN string) (*Repository, error) {
	c, err := pgx.ParseConfig(dbDSN)
	if err != nil {
		return nil, err
	}

	db := stdlib.OpenDB(*c)
	err = migrateUp(db)
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}, nil
}

// migrateUp will run our migrations on the given DB connection.
// Taken from my friends post: https://jbrandhorst.com/post/postgres/
func migrateUp(db *sql.DB) error {
	sourceInstance, err := bindata.WithInstance(bindata.Resource(migrations.AssetNames(), migrations.Asset))
	if err != nil {
		return err
	}

	targetInstance, err := cockroachdb.WithInstance(db, new(cockroachdb.Config))
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("go-bindata", sourceInstance, "cockroachdb", targetInstance)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return sourceInstance.Close()
}
