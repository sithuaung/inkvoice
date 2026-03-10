package cli

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export invoices/clients to CSV/JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			format, _ := cmd.Flags().GetString("format")
			entity, _ := cmd.Flags().GetString("entity")
			output, _ := cmd.Flags().GetString("output")

			svc, db, err := openService()
			if err != nil {
				return err
			}
			defer db.Close()

			ctx := context.Background()

			var w *os.File
			if output != "" {
				w, err = os.Create(output)
				if err != nil {
					return fmt.Errorf("create file: %w", err)
				}
				defer w.Close()
			} else {
				w = os.Stdout
			}

			switch entity {
			case "clients":
				clients, err := svc.ListClients(ctx)
				if err != nil {
					return err
				}
				if format == "json" {
					enc := json.NewEncoder(w)
					enc.SetIndent("", "  ")
					return enc.Encode(clients)
				}
				cw := csv.NewWriter(w)
				cw.Write([]string{"id", "name", "email", "phone", "company"})
				for _, c := range clients {
					cw.Write([]string{c.ID, c.Name, c.Email, c.Phone, c.Company})
				}
				cw.Flush()
				return cw.Error()

			case "invoices":
				invoices, err := svc.ListInvoices(ctx, "", "")
				if err != nil {
					return err
				}
				if format == "json" {
					enc := json.NewEncoder(w)
					enc.SetIndent("", "  ")
					return enc.Encode(invoices)
				}
				cw := csv.NewWriter(w)
				cw.Write([]string{"id", "number", "client_id", "status", "total", "due_date"})
				for _, inv := range invoices {
					cw.Write([]string{inv.ID, inv.InvoiceNumber, inv.ClientID, inv.Status, fmt.Sprintf("%d", inv.Total), inv.DueDate})
				}
				cw.Flush()
				return cw.Error()

			default:
				return fmt.Errorf("unknown entity: %s (use 'clients' or 'invoices')", entity)
			}
		},
	}

	cmd.Flags().String("format", "csv", "output format: csv or json")
	cmd.Flags().String("entity", "invoices", "what to export: clients or invoices")
	cmd.Flags().StringP("output", "o", "", "output file (default: stdout)")
	return cmd
}
