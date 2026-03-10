package cli

import (
	"fmt"
	"os"

	"github.com/sithuaung/inkvoice/internal/pdf"
	"github.com/spf13/cobra"
)

func newHealthCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Check DB, Typst, storage status",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Health Check")
			fmt.Println("============")

			// Database
			db, err := openDB()
			if err != nil {
				fmt.Printf("Database:   FAIL (%v)\n", err)
			} else {
				fmt.Printf("Database:   OK (%s)\n", dbPath)
				db.Close()
			}

			// Typst
			if pdf.TypstAvailable() {
				fmt.Println("Typst CLI:  OK")
			} else {
				fmt.Println("Typst CLI:  NOT FOUND (PDF generation unavailable)")
			}

			// Storage directory
			storageDir := "data/invoices"
			if info, err := os.Stat(storageDir); err == nil && info.IsDir() {
				fmt.Printf("Storage:    OK (%s)\n", storageDir)
			} else {
				fmt.Printf("Storage:    NOT FOUND (%s) — will be created on first PDF\n", storageDir)
			}

			// Templates directory
			templatesDir := "invoice-templates"
			templates, err := pdf.FindTemplates(templatesDir)
			if err != nil {
				fmt.Printf("Templates:  NOT FOUND (%s)\n", templatesDir)
			} else {
				fmt.Printf("Templates:  %d found in %s\n", len(templates), templatesDir)
			}

			return nil
		},
	}
}
