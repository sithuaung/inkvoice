package database

import (
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "modernc.org/sqlite"

	"github.com/sithuaung/inkvoice/internal/database/dbsqlc"
)

// DB wraps the database connection and sqlc queries.
type DB struct {
	Conn    *sql.DB
	Queries *dbsqlc.Queries
}

// Open opens a SQLite database at the given path.
func Open(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(wal)&_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &DB{
		Conn:    conn,
		Queries: dbsqlc.New(conn),
	}, nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	return db.Conn.Close()
}

// MigrateUp runs all pending migrations.
func (db *DB) MigrateUp(migrationsFS fs.FS) error {
	m, err := db.newMigrate(migrationsFS)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil
}

// MigrateDown rolls back the last migration.
func (db *DB) MigrateDown(migrationsFS fs.FS) error {
	m, err := db.newMigrate(migrationsFS)
	if err != nil {
		return err
	}
	if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate down: %w", err)
	}
	return nil
}

// MigrateVersion returns the current migration version and dirty state.
func (db *DB) MigrateVersion() (uint, bool, error) {
	// We need a dummy migrate instance just for version check, but we can
	// query the schema_migrations table directly.
	var version uint
	var dirty bool
	err := db.Conn.QueryRow("SELECT version, dirty FROM schema_migrations LIMIT 1").Scan(&version, &dirty)
	if err != nil {
		return 0, false, fmt.Errorf("get migration version: %w", err)
	}
	return version, dirty, nil
}

func (db *DB) newMigrate(migrationsFS fs.FS) (*migrate.Migrate, error) {
	sourceDriver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("create migration source: %w", err)
	}
	dbDriver, err := sqlite.WithInstance(db.Conn, &sqlite.Config{
		NoTxWrap: true,
	})
	if err != nil {
		return nil, fmt.Errorf("create migration db driver: %w", err)
	}
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "sqlite", dbDriver)
	if err != nil {
		return nil, fmt.Errorf("create migrator: %w", err)
	}
	return m, nil
}
