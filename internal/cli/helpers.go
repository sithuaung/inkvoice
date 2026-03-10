package cli

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/sithuaung/inkvoice/internal/database"
	"github.com/sithuaung/inkvoice/internal/service"
)

// MigrationsFS is set by the main package to provide the embedded migrations.
var MigrationsFS fs.FS

func openDB() (*database.DB, error) {
	return database.Open(dbPath)
}

func openService() (*service.Service, *database.DB, error) {
	db, err := openDB()
	if err != nil {
		return nil, nil, err
	}
	return service.New(db), db, nil
}

func exitErr(msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
	os.Exit(1)
}
