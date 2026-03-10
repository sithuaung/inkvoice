package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"text/tabwriter"

	"github.com/sithuaung/inkvoice/internal/email"
	"github.com/sithuaung/inkvoice/internal/model"
	"github.com/sithuaung/inkvoice/internal/pdf"
	"github.com/sithuaung/inkvoice/internal/storage"
	"github.com/spf13/cobra"
)

func newInvoiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invoice",
		Short: "Manage invoices",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a draft invoice for a client",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientID, _ := cmd.Flags().GetString("client")
			notes, _ := cmd.Flags().GetString("notes")

			if clientID == "" {
				return fmt.Errorf("--client is required")
			}

			svc, db, err := openService()
			if err != nil {
				return err
			}
			defer db.Close()

			id, err := svc.CreateInvoice(context.Background(), clientID, notes)
			if err != nil {
				return err
			}

			inv, _ := svc.GetInvoice(context.Background(), id)
			fmt.Printf("Invoice created: %s (%s)\n", inv.InvoiceNumber, id)
			return nil
		},
	}
	createCmd.Flags().String("client", "", "client ID (required)")
	createCmd.Flags().String("notes", "", "invoice notes")

	addItemCmd := &cobra.Command{
		Use:   "add-item",
		Short: "Add a line item to an invoice",
		RunE: func(cmd *cobra.Command, args []string) error {
			invoiceID, _ := cmd.Flags().GetString("invoice")
			productID, _ := cmd.Flags().GetString("product")
			desc, _ := cmd.Flags().GetString("description")
			qty, _ := cmd.Flags().GetFloat64("quantity")
			price, _ := cmd.Flags().GetInt64("price")
			taxRate, _ := cmd.Flags().GetFloat64("tax-rate")

			if invoiceID == "" {
				return fmt.Errorf("--invoice is required")
			}
			if desc == "" && productID == "" {
				return fmt.Errorf("--description or --product is required")
			}

			svc, db, err := openService()
			if err != nil {
				return err
			}
			defer db.Close()

			ctx := context.Background()

			// If product specified, use its details as defaults
			if productID != "" {
				p, err := svc.GetProduct(ctx, productID)
				if err != nil {
					return fmt.Errorf("get product: %w", err)
				}
				if desc == "" {
					desc = p.Name
				}
				if !cmd.Flags().Changed("price") {
					price = p.UnitPrice
				}
			}

			if qty == 0 {
				qty = 1
			}

			itemID, err := svc.AddInvoiceItem(ctx, invoiceID, productID, desc, qty, price, "", taxRate)
			if err != nil {
				return err
			}
			fmt.Printf("Item added: %s\n", itemID)
			return nil
		},
	}
	addItemCmd.Flags().String("invoice", "", "invoice ID (required)")
	addItemCmd.Flags().String("product", "", "product ID")
	addItemCmd.Flags().String("description", "", "item description")
	addItemCmd.Flags().Float64("quantity", 1, "quantity")
	addItemCmd.Flags().Int64("price", 0, "unit price in cents")
	addItemCmd.Flags().Float64("tax-rate", 0, "tax rate percentage")

	removeItemCmd := &cobra.Command{
		Use:   "remove-item [item-id]",
		Short: "Remove a line item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, db, err := openService()
			if err != nil {
				return err
			}
			defer db.Close()

			if err := svc.RemoveInvoiceItem(context.Background(), args[0]); err != nil {
				return err
			}
			fmt.Println("Item removed.")
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List invoices",
		RunE: func(cmd *cobra.Command, args []string) error {
			status, _ := cmd.Flags().GetString("status")
			clientID, _ := cmd.Flags().GetString("client")

			svc, db, err := openService()
			if err != nil {
				return err
			}
			defer db.Close()

			invoices, err := svc.ListInvoices(context.Background(), status, clientID)
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NUMBER\tSTATUS\tTOTAL\tDUE DATE\tID")
			for _, inv := range invoices {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					inv.InvoiceNumber, inv.Status,
					model.FormatMoney(inv.Total, inv.Currency),
					inv.DueDate, inv.ID)
			}
			w.Flush()
			return nil
		},
	}
	listCmd.Flags().String("status", "", "filter by status")
	listCmd.Flags().String("client", "", "filter by client ID")

	showCmd := &cobra.Command{
		Use:   "show [id-or-number]",
		Short: "Show invoice details + line items",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, db, err := openService()
			if err != nil {
				return err
			}
			defer db.Close()

			ctx := context.Background()

			// Try by ID first, then by number
			inv, err := svc.GetInvoice(ctx, args[0])
			if err != nil {
				inv, err = svc.GetInvoiceByNumber(ctx, args[0])
				if err != nil {
					return fmt.Errorf("invoice not found: %s", args[0])
				}
			}

			client, _ := svc.GetClient(ctx, inv.ClientID)

			fmt.Printf("Invoice:  %s\n", inv.InvoiceNumber)
			fmt.Printf("ID:       %s\n", inv.ID)
			fmt.Printf("Client:   %s (%s)\n", client.Name, client.ID)
			fmt.Printf("Status:   %s\n", inv.Status)
			fmt.Printf("Issued:   %s\n", inv.IssueDate)
			fmt.Printf("Due:      %s\n", inv.DueDate)
			fmt.Printf("Currency: %s\n", inv.Currency)
			fmt.Println()

			items, err := svc.ListInvoiceItems(ctx, inv.ID)
			if err != nil {
				return err
			}

			if len(items) > 0 {
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				fmt.Fprintln(w, "  #\tDESCRIPTION\tQTY\tPRICE\tAMOUNT")
				for i, item := range items {
					fmt.Fprintf(w, "  %d\t%s\t%.2f\t%s\t%s\n",
						i+1, item.Description, item.Quantity,
						model.FormatMoney(item.UnitPrice, inv.Currency),
						model.FormatMoney(item.Amount, inv.Currency))
				}
				w.Flush()
				fmt.Println()
			}

			fmt.Printf("Subtotal:    %s\n", model.FormatMoney(inv.Subtotal, inv.Currency))
			fmt.Printf("Tax:         %s\n", model.FormatMoney(inv.TaxTotal, inv.Currency))
			fmt.Printf("Total:       %s\n", model.FormatMoney(inv.Total, inv.Currency))
			fmt.Printf("Paid:        %s\n", model.FormatMoney(inv.AmountPaid, inv.Currency))
			fmt.Printf("Balance:     %s\n", model.FormatMoney(inv.Total-inv.AmountPaid, inv.Currency))
			return nil
		},
	}

	pdfCmd := &cobra.Command{
		Use:   "pdf [id-or-number]",
		Short: "Generate PDF to file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			output, _ := cmd.Flags().GetString("output")

			svc, db, err := openService()
			if err != nil {
				return err
			}
			defer db.Close()

			ctx := context.Background()
			inv, err := svc.GetInvoice(ctx, args[0])
			if err != nil {
				inv, err = svc.GetInvoiceByNumber(ctx, args[0])
				if err != nil {
					return fmt.Errorf("invoice not found: %s", args[0])
				}
			}

			client, err := svc.GetClient(ctx, inv.ClientID)
			if err != nil {
				return fmt.Errorf("get client: %w", err)
			}

			settings, err := svc.GetSettings(ctx)
			if err != nil {
				return fmt.Errorf("get settings: %w", err)
			}

			gen := pdf.NewGenerator("invoice-templates")
			inputs := map[string]string{
				"invoice-number": inv.InvoiceNumber,
				"issue-date":     inv.IssueDate,
				"due-date":       inv.DueDate,
				"client-name":    client.Name,
				"client-email":   client.Email,
				"client-company": client.Company,
				"company-name":   settings.CompanyName,
				"company-email":  settings.CompanyEmail,
				"subtotal":       fmt.Sprintf("%d", inv.Subtotal),
				"tax-total":      fmt.Sprintf("%d", inv.TaxTotal),
				"total":          fmt.Sprintf("%d", inv.Total),
				"currency":       inv.Currency,
				"notes":          inv.Notes,
			}

			// Build items JSON
			items, _ := svc.ListInvoiceItems(ctx, inv.ID)
			itemsJSON := "["
			for i, item := range items {
				if i > 0 {
					itemsJSON += ","
				}
				itemsJSON += fmt.Sprintf(`{"description":"%s","quantity":%.2f,"unit_price":%d,"amount":%d}`,
					item.Description, item.Quantity, item.UnitPrice, item.Amount)
			}
			itemsJSON += "]"
			inputs["items"] = itemsJSON

			templatePath := filepath.Join("invoice-templates", "template-a.typ")
			pdfData, err := gen.Generate(ctx, templatePath, inputs)
			if err != nil {
				return fmt.Errorf("generate PDF: %w", err)
			}

			if output == "" {
				// Save to storage
				store := storage.NewLocalStore("data/invoices")
				relPath, err := store.Save(inv.ID, pdfData)
				if err != nil {
					return fmt.Errorf("save PDF: %w", err)
				}
				if err := svc.UpdateInvoicePDFPath(ctx, inv.ID, relPath); err != nil {
					return fmt.Errorf("update pdf path: %w", err)
				}
				fmt.Printf("PDF saved: %s\n", store.Path(relPath))
			} else {
				if err := os.WriteFile(output, pdfData, 0o644); err != nil {
					return fmt.Errorf("write file: %w", err)
				}
				fmt.Printf("PDF saved: %s\n", output)
			}
			return nil
		},
	}
	pdfCmd.Flags().StringP("output", "o", "", "output file path (default: save to storage)")

	previewCmd := &cobra.Command{
		Use:   "preview [id-or-number]",
		Short: "Generate PDF and open in default viewer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tmpFile, err := os.CreateTemp("", "inkvoice-preview-*.pdf")
			if err != nil {
				return err
			}
			tmpPath := tmpFile.Name()
			tmpFile.Close()

			// Set output flag on the pdfCmd itself, then invoke with pdfCmd as receiver
			pdfCmd.Flags().Set("output", tmpPath)
			if err := pdfCmd.RunE(pdfCmd, args); err != nil {
				return err
			}

			var openCmd string
			switch runtime.GOOS {
			case "darwin":
				openCmd = "open"
			case "linux":
				openCmd = "xdg-open"
			default:
				openCmd = "start"
			}
			return exec.Command(openCmd, tmpPath).Start()
		},
	}

	sendCmd := &cobra.Command{
		Use:   "send [id-or-number]",
		Short: "Send invoice via email / send reminder",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if appCfg.SMTPHost == "" {
				return fmt.Errorf("SMTP not configured. Set SMTP_HOST, SMTP_PORT, SMTP_USER, SMTP_PASS, SMTP_FROM in .env")
			}

			svc, db, err := openService()
			if err != nil {
				return err
			}
			defer db.Close()

			ctx := context.Background()
			inv, err := svc.GetInvoice(ctx, args[0])
			if err != nil {
				inv, err = svc.GetInvoiceByNumber(ctx, args[0])
				if err != nil {
					return fmt.Errorf("invoice not found: %s", args[0])
				}
			}

			client, err := svc.GetClient(ctx, inv.ClientID)
			if err != nil {
				return fmt.Errorf("get client: %w", err)
			}

			if client.Email == "" {
				return fmt.Errorf("client %s has no email address", client.Name)
			}

			// Generate PDF if not already generated
			pdfFilePath := ""
			if inv.PdfPath != "" {
				store := storage.NewLocalStore(appCfg.StorageDir)
				pdfFilePath = store.Path(inv.PdfPath)
			} else {
				// Generate PDF first
				pdfCmd.Flags().Set("output", "")
				if err := pdfCmd.RunE(pdfCmd, args); err != nil {
					return fmt.Errorf("generate PDF: %w", err)
				}
				// Re-fetch invoice to get the updated pdf_path
				inv, _ = svc.GetInvoice(ctx, inv.ID)
				if inv.PdfPath != "" {
					store := storage.NewLocalStore(appCfg.StorageDir)
					pdfFilePath = store.Path(inv.PdfPath)
				}
			}

			subject := fmt.Sprintf("Invoice %s", inv.InvoiceNumber)
			body := fmt.Sprintf("Hi %s,\n\nPlease find attached invoice %s for %s.\n\nDue date: %s\n\nThank you!",
				client.Name, inv.InvoiceNumber,
				model.FormatMoney(inv.Total, inv.Currency), inv.DueDate)

			sender := email.NewSMTPSender(email.SMTPConfig{
				Host:     appCfg.SMTPHost,
				Port:     appCfg.SMTPPort,
				Username: appCfg.SMTPUser,
				Password: appCfg.SMTPPass,
				From:     appCfg.SMTPFrom,
			})

			fmt.Printf("Sending invoice %s to %s...\n", inv.InvoiceNumber, client.Email)
			if err := sender.SendInvoice(client.Email, subject, body, pdfFilePath); err != nil {
				return fmt.Errorf("send email: %w", err)
			}

			// Update status to sent if still draft
			if inv.Status == "draft" {
				svc.UpdateInvoiceStatus(ctx, inv.ID, "sent")
			}

			fmt.Println("Invoice sent successfully.")
			return nil
		},
	}

	cmd.AddCommand(createCmd, addItemCmd, removeItemCmd, listCmd, showCmd, pdfCmd, previewCmd, sendCmd)
	return cmd
}
