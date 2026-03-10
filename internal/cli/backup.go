package cli

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func newBackupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Safe copy of SQLite database",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, _ := cmd.Flags().GetString("output")

			if output == "" {
				output = fmt.Sprintf("inkvoice-backup-%s.db", time.Now().Format("20060102-150405"))
			}

			src, err := os.Open(dbPath)
			if err != nil {
				return fmt.Errorf("open database: %w", err)
			}
			defer src.Close()

			dst, err := os.Create(output)
			if err != nil {
				return fmt.Errorf("create backup: %w", err)
			}
			defer dst.Close()

			if _, err := io.Copy(dst, src); err != nil {
				return fmt.Errorf("copy database: %w", err)
			}

			fmt.Printf("Backup saved to: %s\n", output)
			return nil
		},
	}

	cmd.Flags().StringP("output", "o", "", "output file path")
	return cmd
}
