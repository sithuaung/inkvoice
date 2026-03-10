package cli

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/sithuaung/inkvoice/internal/scheduler"
	"github.com/spf13/cobra"
)

func newServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start cron scheduler (recurring invoices)",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, db, err := openService()
			if err != nil {
				return err
			}
			defer db.Close()

			s := scheduler.New(svc)
			s.Start()

			fmt.Println("Inkvoice scheduler running. Press Ctrl+C to stop.")
			slog.Info("serve started", "db", dbPath)

			// Wait for shutdown signal
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
			<-sig

			s.Stop()
			fmt.Println("\nShutdown complete.")
			return nil
		},
	}
}
