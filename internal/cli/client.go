package cli

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func newClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "client",
		Short: "Manage clients",
	}

	// create
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Add a client",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			email, _ := cmd.Flags().GetString("email")
			phone, _ := cmd.Flags().GetString("phone")
			company, _ := cmd.Flags().GetString("company")
			notes, _ := cmd.Flags().GetString("notes")

			if name == "" {
				return fmt.Errorf("--name is required")
			}

			svc, db, err := openService()
			if err != nil {
				return err
			}
			defer db.Close()

			id, err := svc.CreateClient(context.Background(), name, email, phone, company, "{}", notes)
			if err != nil {
				return err
			}
			fmt.Printf("Client created: %s\n", id)
			return nil
		},
	}
	createCmd.Flags().String("name", "", "client name (required)")
	createCmd.Flags().String("email", "", "client email")
	createCmd.Flags().String("phone", "", "client phone")
	createCmd.Flags().String("company", "", "client company")
	createCmd.Flags().String("notes", "", "notes")

	// list
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List clients",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, db, err := openService()
			if err != nil {
				return err
			}
			defer db.Close()

			clients, err := svc.ListClients(context.Background())
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tEMAIL\tCOMPANY")
			for _, c := range clients {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", c.ID, c.Name, c.Email, c.Company)
			}
			w.Flush()
			return nil
		},
	}

	// show
	showCmd := &cobra.Command{
		Use:   "show [id]",
		Short: "Show client details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, db, err := openService()
			if err != nil {
				return err
			}
			defer db.Close()

			c, err := svc.GetClient(context.Background(), args[0])
			if err != nil {
				return err
			}

			fmt.Printf("ID:      %s\n", c.ID)
			fmt.Printf("Name:    %s\n", c.Name)
			fmt.Printf("Email:   %s\n", c.Email)
			fmt.Printf("Phone:   %s\n", c.Phone)
			fmt.Printf("Company: %s\n", c.Company)
			fmt.Printf("Notes:   %s\n", c.Notes)
			fmt.Printf("Created: %s\n", c.CreatedAt)
			return nil
		},
	}

	// update
	updateCmd := &cobra.Command{
		Use:   "update [id]",
		Short: "Update client",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, db, err := openService()
			if err != nil {
				return err
			}
			defer db.Close()

			ctx := context.Background()
			c, err := svc.GetClient(ctx, args[0])
			if err != nil {
				return err
			}

			if v, _ := cmd.Flags().GetString("name"); v != "" {
				c.Name = v
			}
			if v, _ := cmd.Flags().GetString("email"); v != "" {
				c.Email = v
			}
			if v, _ := cmd.Flags().GetString("phone"); v != "" {
				c.Phone = v
			}
			if v, _ := cmd.Flags().GetString("company"); v != "" {
				c.Company = v
			}
			if v, _ := cmd.Flags().GetString("notes"); v != "" {
				c.Notes = v
			}

			if err := svc.UpdateClient(ctx, c.ID, c.Name, c.Email, c.Phone, c.Company, c.Address, c.Notes); err != nil {
				return err
			}
			fmt.Println("Client updated.")
			return nil
		},
	}
	updateCmd.Flags().String("name", "", "client name")
	updateCmd.Flags().String("email", "", "client email")
	updateCmd.Flags().String("phone", "", "client phone")
	updateCmd.Flags().String("company", "", "client company")
	updateCmd.Flags().String("notes", "", "notes")

	// delete
	deleteCmd := &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete client",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, db, err := openService()
			if err != nil {
				return err
			}
			defer db.Close()

			if err := svc.DeleteClient(context.Background(), args[0]); err != nil {
				return err
			}
			fmt.Println("Client deleted.")
			return nil
		},
	}

	cmd.AddCommand(createCmd, listCmd, showCmd, updateCmd, deleteCmd)
	return cmd
}
