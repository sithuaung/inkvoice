package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Database migration commands",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "up",
			Short: "Apply pending migrations",
			RunE: func(cmd *cobra.Command, args []string) error {
				db, err := openDB()
				if err != nil {
					return err
				}
				defer db.Close()
				if err := db.MigrateUp(MigrationsFS); err != nil {
					return fmt.Errorf("migrate up: %w", err)
				}
				fmt.Println("Migrations applied successfully.")
				return nil
			},
		},
		&cobra.Command{
			Use:   "down",
			Short: "Rollback last migration",
			RunE: func(cmd *cobra.Command, args []string) error {
				db, err := openDB()
				if err != nil {
					return err
				}
				defer db.Close()
				if err := db.MigrateDown(MigrationsFS); err != nil {
					return fmt.Errorf("migrate down: %w", err)
				}
				fmt.Println("Migration rolled back successfully.")
				return nil
			},
		},
		&cobra.Command{
			Use:   "status",
			Short: "Show migration status",
			RunE: func(cmd *cobra.Command, args []string) error {
				db, err := openDB()
				if err != nil {
					return err
				}
				defer db.Close()
				version, dirty, err := db.MigrateVersion()
				if err != nil {
					return fmt.Errorf("migration status: %w", err)
				}
				fmt.Printf("Version: %d\nDirty:   %v\n", version, dirty)
				return nil
			},
		},
	)

	return cmd
}
