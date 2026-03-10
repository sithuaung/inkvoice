package cli

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/sithuaung/inkvoice/internal/scheduler"
	"github.com/spf13/cobra"
)

func newRecurringCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recurring",
		Short: "Manage recurring invoice schedules",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List recurring schedules",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, db, err := openService()
			if err != nil {
				return err
			}
			defer db.Close()

			recurring, err := svc.ListRecurringInvoices(context.Background())
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tCLIENT\tSCHEDULE\tSTATUS\tNEXT RUN")
			for _, r := range recurring {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", r.ID, r.ClientID, r.Schedule, r.Status, r.NextRun)
			}
			w.Flush()
			return nil
		},
	}

	showCmd := &cobra.Command{
		Use:   "show [id]",
		Short: "Show schedule details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, db, err := openService()
			if err != nil {
				return err
			}
			defer db.Close()

			ctx := context.Background()
			r, err := svc.GetRecurringInvoice(ctx, args[0])
			if err != nil {
				return err
			}

			fmt.Printf("ID:       %s\n", r.ID)
			fmt.Printf("Client:   %s\n", r.ClientID)
			fmt.Printf("Schedule: %s\n", r.Schedule)
			fmt.Printf("Status:   %s\n", r.Status)
			fmt.Printf("Next Run: %s\n", r.NextRun)
			fmt.Printf("Last Run: %s\n", r.LastRun)
			fmt.Printf("Created:  %s\n", r.CreatedAt)

			items, err := svc.ListRecurringInvoiceItems(ctx, r.ID)
			if err != nil {
				return err
			}
			if len(items) > 0 {
				fmt.Println("\nItems:")
				for i, item := range items {
					fmt.Printf("  %d. %s (qty: %.2f, price: %d)\n", i+1, item.Description, item.Quantity, item.UnitPrice)
				}
			}
			return nil
		},
	}

	pauseCmd := &cobra.Command{
		Use:   "pause [id]",
		Short: "Pause a schedule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, db, err := openService()
			if err != nil {
				return err
			}
			defer db.Close()

			if err := svc.UpdateRecurringInvoiceStatus(context.Background(), args[0], "paused"); err != nil {
				return err
			}
			fmt.Println("Schedule paused.")
			return nil
		},
	}

	resumeCmd := &cobra.Command{
		Use:   "resume [id]",
		Short: "Resume a paused schedule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, db, err := openService()
			if err != nil {
				return err
			}
			defer db.Close()

			if err := svc.UpdateRecurringInvoiceStatus(context.Background(), args[0], "active"); err != nil {
				return err
			}
			fmt.Println("Schedule resumed.")
			return nil
		},
	}

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Process all due recurring invoices",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, db, err := openService()
			if err != nil {
				return err
			}
			defer db.Close()

			s := scheduler.New(svc)
			s.ProcessDue()
			fmt.Println("Done processing due recurring invoices.")
			return nil
		},
	}

	triggerCmd := &cobra.Command{
		Use:   "trigger [id]",
		Short: "Manually generate invoice now",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Manual trigger will be implemented with full recurring create support.")
			return nil
		},
	}

	cmd.AddCommand(listCmd, showCmd, pauseCmd, resumeCmd, runCmd, triggerCmd)
	return cmd
}
