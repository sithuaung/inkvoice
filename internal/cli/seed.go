package cli

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/sithuaung/inkvoice/internal/pdf"
	"github.com/spf13/cobra"
)

func newSeedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "seed",
		Short: "Seed demo data and templates",
	}

	dataCmd := &cobra.Command{
		Use:   "data",
		Short: "Insert demo clients, products, invoices",
		RunE: func(cmd *cobra.Command, args []string) error {
			skipIfExists, _ := cmd.Flags().GetBool("skip-if-exists")

			svc, db, err := openService()
			if err != nil {
				return err
			}
			defer db.Close()

			ctx := context.Background()

			// Check if data already exists
			if skipIfExists {
				count, err := db.Queries.CountClients(ctx)
				if err == nil && count > 0 {
					fmt.Println("Data already exists, skipping.")
					return nil
				}
			}

			// Create demo clients
			c1, err := svc.CreateClient(ctx, "Acme Corp", "billing@acme.com", "+1-555-0100", "Acme Corporation", "{}", "Demo client")
			if err != nil {
				return fmt.Errorf("create client: %w", err)
			}
			c2, err := svc.CreateClient(ctx, "Jane Smith", "jane@example.com", "+1-555-0200", "Smith Consulting", "{}", "")
			if err != nil {
				return fmt.Errorf("create client: %w", err)
			}

			// Create demo products
			p1, err := svc.CreateProduct(ctx, "Web Development", "Full-stack web development services", 15000, "USD")
			if err != nil {
				return fmt.Errorf("create product: %w", err)
			}
			p2, err := svc.CreateProduct(ctx, "UI/UX Design", "User interface and experience design", 12000, "USD")
			if err != nil {
				return fmt.Errorf("create product: %w", err)
			}
			_, err = svc.CreateProduct(ctx, "Consulting", "Technical consulting per hour", 20000, "USD")
			if err != nil {
				return fmt.Errorf("create product: %w", err)
			}

			// Create demo invoices
			inv1, err := svc.CreateInvoice(ctx, c1, "Thank you for your business!")
			if err != nil {
				return fmt.Errorf("create invoice: %w", err)
			}
			_, err = svc.AddInvoiceItem(ctx, inv1, p1, "Web Development", 40, 15000, "", 0)
			if err != nil {
				return fmt.Errorf("add item: %w", err)
			}
			_, err = svc.AddInvoiceItem(ctx, inv1, p2, "UI/UX Design", 20, 12000, "", 0)
			if err != nil {
				return fmt.Errorf("add item: %w", err)
			}

			inv2, err := svc.CreateInvoice(ctx, c2, "")
			if err != nil {
				return fmt.Errorf("create invoice: %w", err)
			}
			_, err = svc.AddInvoiceItem(ctx, inv2, "", "Custom project work", 1, 500000, "", 0)
			if err != nil {
				return fmt.Errorf("add item: %w", err)
			}

			fmt.Println("Demo data seeded successfully.")
			fmt.Printf("  Clients:  %s, %s\n", c1[:8], c2[:8])
			fmt.Printf("  Products: %s, %s\n", p1[:8], p2[:8])
			fmt.Printf("  Invoices: %s, %s\n", inv1[:8], inv2[:8])
			return nil
		},
	}
	dataCmd.Flags().Bool("skip-if-exists", false, "skip if data already exists")

	templateCmd := &cobra.Command{
		Use:   "template",
		Short: "Scan invoice-templates/ folder, register .typ files in DB",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, db, err := openService()
			if err != nil {
				return err
			}
			defer db.Close()

			ctx := context.Background()
			templatesDir := "invoice-templates"
			templates, err := pdf.FindTemplates(templatesDir)
			if err != nil {
				slog.Warn("no templates directory found", "dir", templatesDir)
				return nil
			}

			registered := 0
			for _, t := range templates {
				path := filepath.Join(templatesDir, t)
				// Check if already registered
				_, err := svc.GetTemplateByPath(ctx, path)
				if err == nil {
					continue // already exists
				}
				if err != sql.ErrNoRows {
					return fmt.Errorf("check template: %w", err)
				}

				name := t[:len(t)-len(filepath.Ext(t))] // strip .typ
				isDefault := registered == 0              // first one is default
				id, err := svc.CreateTemplate(ctx, name, path, isDefault)
				if err != nil {
					return fmt.Errorf("register template: %w", err)
				}
				registered++
				fmt.Printf("Registered template: %s (%s)\n", name, id[:8])
			}

			if registered == 0 {
				fmt.Println("No new templates found.")
			} else {
				fmt.Printf("Registered %d template(s).\n", registered)
			}
			return nil
		},
	}

	cmd.AddCommand(dataCmd, templateCmd)
	return cmd
}
