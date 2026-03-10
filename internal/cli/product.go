package cli

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/sithuaung/inkvoice/internal/model"
	"github.com/spf13/cobra"
)

func newProductCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "product",
		Short: "Manage products",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Add a product/service",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			desc, _ := cmd.Flags().GetString("description")
			price, _ := cmd.Flags().GetInt64("price")
			currency, _ := cmd.Flags().GetString("currency")

			if name == "" {
				return fmt.Errorf("--name is required")
			}

			svc, db, err := openService()
			if err != nil {
				return err
			}
			defer db.Close()

			id, err := svc.CreateProduct(context.Background(), name, desc, price, currency)
			if err != nil {
				return err
			}
			fmt.Printf("Product created: %s\n", id)
			return nil
		},
	}
	createCmd.Flags().String("name", "", "product name (required)")
	createCmd.Flags().String("description", "", "product description")
	createCmd.Flags().Int64("price", 0, "unit price in cents")
	createCmd.Flags().String("currency", "USD", "currency")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List products",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, db, err := openService()
			if err != nil {
				return err
			}
			defer db.Close()

			products, err := svc.ListProducts(context.Background())
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tPRICE\tCURRENCY")
			for _, p := range products {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", p.ID, p.Name, model.FormatMoney(p.UnitPrice, p.Currency), p.Currency)
			}
			w.Flush()
			return nil
		},
	}

	updateCmd := &cobra.Command{
		Use:   "update [id]",
		Short: "Update product",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, db, err := openService()
			if err != nil {
				return err
			}
			defer db.Close()

			ctx := context.Background()
			p, err := svc.GetProduct(ctx, args[0])
			if err != nil {
				return err
			}

			if v, _ := cmd.Flags().GetString("name"); v != "" {
				p.Name = v
			}
			if v, _ := cmd.Flags().GetString("description"); v != "" {
				p.Description = v
			}
			if cmd.Flags().Changed("price") {
				v, _ := cmd.Flags().GetInt64("price")
				p.UnitPrice = v
			}
			if v, _ := cmd.Flags().GetString("currency"); v != "" {
				p.Currency = v
			}

			if err := svc.UpdateProduct(ctx, p.ID, p.Name, p.Description, p.UnitPrice, p.Currency); err != nil {
				return err
			}
			fmt.Println("Product updated.")
			return nil
		},
	}
	updateCmd.Flags().String("name", "", "product name")
	updateCmd.Flags().String("description", "", "description")
	updateCmd.Flags().Int64("price", 0, "unit price in cents")
	updateCmd.Flags().String("currency", "", "currency")

	deleteCmd := &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete product",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, db, err := openService()
			if err != nil {
				return err
			}
			defer db.Close()

			if err := svc.DeleteProduct(context.Background(), args[0]); err != nil {
				return err
			}
			fmt.Println("Product deleted.")
			return nil
		},
	}

	cmd.AddCommand(createCmd, listCmd, updateCmd, deleteCmd)
	return cmd
}
