package cli

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"

	"github.com/sithuaung/inkvoice/internal/config"
)

var (
	version = "dev"
	dbPath  string
	appCfg  config.Config
)

// SetVersion sets the version string (called from main with ldflags).
func SetVersion(v string) {
	if v != "" {
		version = v
	}
}

func getVersion() string {
	if version != "dev" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, s := range info.Settings {
			if s.Key == "vcs.revision" && len(s.Value) >= 7 {
				return "dev-" + s.Value[:7]
			}
		}
	}
	return version
}

// NewRootCmd creates the root cobra command.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:     "inkvoice",
		Short:   "CLI invoicing tool for freelancers",
		Version: getVersion(),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Load .env file (silently skips if not found)
			config.LoadDotEnv(".env")

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			appCfg = cfg

			// Use env var for DB path if --db flag wasn't explicitly set
			if !cmd.Flags().Changed("db") && cfg.DBPath != "" {
				dbPath = cfg.DBPath
			}
			return nil
		},
	}

	root.PersistentFlags().StringVar(&dbPath, "db", "inkvoice.db", "path to SQLite database")

	root.AddCommand(
		newMigrateCmd(),
		newClientCmd(),
		newProductCmd(),
		newInvoiceCmd(),
		newRecurringCmd(),
		newSeedCmd(),
		newServeCmd(),
		newHealthCmd(),
		newBackupCmd(),
		newExportCmd(),
	)

	return root
}

// Execute runs the CLI.
func Execute() {
	root := NewRootCmd()
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
