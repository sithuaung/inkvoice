package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/sithuaung/inkvoice/internal/database"
	"github.com/sithuaung/inkvoice/internal/database/dbsqlc"
	"github.com/sithuaung/inkvoice/internal/model"
)

// Service contains business logic.
type Service struct {
	DB *database.DB
}

// New creates a new Service.
func New(db *database.DB) *Service {
	return &Service{DB: db}
}

func newID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String()
}

func now() string {
	return model.NowRFC3339()
}

// --- Clients ---

func (s *Service) CreateClient(ctx context.Context, name, email, phone, company, address, notes string) (string, error) {
	id := newID()
	ts := now()
	err := s.DB.Queries.CreateClient(ctx, dbsqlc.CreateClientParams{
		ID: id, Name: name, Email: email, Phone: phone,
		Company: company, Address: address, Notes: notes,
		CreatedAt: ts, UpdatedAt: ts,
	})
	if err != nil {
		return "", fmt.Errorf("create client: %w", err)
	}
	return id, nil
}

func (s *Service) GetClient(ctx context.Context, id string) (dbsqlc.Client, error) {
	return s.DB.Queries.GetClient(ctx, id)
}

func (s *Service) ListClients(ctx context.Context) ([]dbsqlc.Client, error) {
	return s.DB.Queries.ListClients(ctx)
}

func (s *Service) UpdateClient(ctx context.Context, id, name, email, phone, company, address, notes string) error {
	return s.DB.Queries.UpdateClient(ctx, dbsqlc.UpdateClientParams{
		Name: name, Email: email, Phone: phone,
		Company: company, Address: address, Notes: notes,
		UpdatedAt: now(), ID: id,
	})
}

func (s *Service) DeleteClient(ctx context.Context, id string) error {
	return s.DB.Queries.DeleteClient(ctx, id)
}

// --- Products ---

func (s *Service) CreateProduct(ctx context.Context, name, description string, unitPrice int64, currency string) (string, error) {
	id := newID()
	ts := now()
	err := s.DB.Queries.CreateProduct(ctx, dbsqlc.CreateProductParams{
		ID: id, Name: name, Description: description,
		UnitPrice: unitPrice, Currency: currency,
		CreatedAt: ts, UpdatedAt: ts,
	})
	if err != nil {
		return "", fmt.Errorf("create product: %w", err)
	}
	return id, nil
}

func (s *Service) GetProduct(ctx context.Context, id string) (dbsqlc.Product, error) {
	return s.DB.Queries.GetProduct(ctx, id)
}

func (s *Service) ListProducts(ctx context.Context) ([]dbsqlc.Product, error) {
	return s.DB.Queries.ListProducts(ctx)
}

func (s *Service) UpdateProduct(ctx context.Context, id, name, description string, unitPrice int64, currency string) error {
	return s.DB.Queries.UpdateProduct(ctx, dbsqlc.UpdateProductParams{
		Name: name, Description: description,
		UnitPrice: unitPrice, Currency: currency,
		UpdatedAt: now(), ID: id,
	})
}

func (s *Service) DeleteProduct(ctx context.Context, id string) error {
	return s.DB.Queries.DeleteProduct(ctx, id)
}

// --- Invoices ---

func (s *Service) CreateInvoice(ctx context.Context, clientID, notes string) (string, error) {
	settings, err := s.DB.Queries.GetSettings(ctx)
	if err != nil {
		return "", fmt.Errorf("get settings: %w", err)
	}

	// Get next invoice number atomically
	num, err := s.DB.Queries.IncrementInvoiceNumber(ctx, now())
	if err != nil {
		return "", fmt.Errorf("increment invoice number: %w", err)
	}

	id := newID()
	ts := now()
	invoiceNumber := fmt.Sprintf("%s-%04d", settings.InvoicePrefix, num)
	issueDate := time.Now().UTC().Format("2006-01-02")
	dueDate := time.Now().UTC().AddDate(0, 0, int(settings.DefaultDueDays)).Format("2006-01-02")

	err = s.DB.Queries.CreateInvoice(ctx, dbsqlc.CreateInvoiceParams{
		ID:            id,
		InvoiceNumber: invoiceNumber,
		ClientID:      clientID,
		Status:        "draft",
		IssueDate:     issueDate,
		DueDate:       dueDate,
		Subtotal:      0,
		TaxTotal:      0,
		Total:         0,
		AmountPaid:    0,
		Currency:      settings.DefaultCurrency,
		Notes:         notes,
		TemplateID:    settings.DefaultTemplateID,
		PdfPath:       "",
		CreatedAt:     ts,
		UpdatedAt:     ts,
	})
	if err != nil {
		return "", fmt.Errorf("create invoice: %w", err)
	}
	return id, nil
}

func (s *Service) GetInvoice(ctx context.Context, id string) (dbsqlc.Invoice, error) {
	return s.DB.Queries.GetInvoice(ctx, id)
}

func (s *Service) GetInvoiceByNumber(ctx context.Context, number string) (dbsqlc.Invoice, error) {
	return s.DB.Queries.GetInvoiceByNumber(ctx, number)
}

func (s *Service) ListInvoices(ctx context.Context, status, clientID string) ([]dbsqlc.Invoice, error) {
	if status != "" {
		return s.DB.Queries.ListInvoicesByStatus(ctx, status)
	}
	if clientID != "" {
		return s.DB.Queries.ListInvoicesByClient(ctx, clientID)
	}
	return s.DB.Queries.ListInvoices(ctx)
}

func (s *Service) AddInvoiceItem(ctx context.Context, invoiceID, productID, description string, quantity float64, unitPrice int64, taxID string, taxRate float64) (string, error) {
	id := newID()
	amount := int64(math.Round(float64(unitPrice) * quantity))
	ts := now()

	// Get current max sort order
	items, err := s.DB.Queries.ListInvoiceItems(ctx, invoiceID)
	if err != nil {
		return "", fmt.Errorf("list invoice items: %w", err)
	}
	sortOrder := int64(len(items))

	err = s.DB.Queries.CreateInvoiceItem(ctx, dbsqlc.CreateInvoiceItemParams{
		ID: id, InvoiceID: invoiceID,
		ProductID:   sql.NullString{String: productID, Valid: productID != ""},
		Description: description, Quantity: quantity,
		UnitPrice: unitPrice,
		TaxID:     sql.NullString{String: taxID, Valid: taxID != ""},
		TaxRate:   taxRate,
		Amount:    amount, SortOrder: sortOrder, CreatedAt: ts,
	})
	if err != nil {
		return "", fmt.Errorf("create invoice item: %w", err)
	}

	// Recalculate totals
	if err := s.recalcInvoiceTotals(ctx, invoiceID); err != nil {
		return "", err
	}

	return id, nil
}

func (s *Service) RemoveInvoiceItem(ctx context.Context, itemID string) error {
	item, err := s.DB.Queries.GetInvoiceItem(ctx, itemID)
	if err != nil {
		return fmt.Errorf("get invoice item: %w", err)
	}
	if err := s.DB.Queries.DeleteInvoiceItem(ctx, itemID); err != nil {
		return fmt.Errorf("delete invoice item: %w", err)
	}
	return s.recalcInvoiceTotals(ctx, item.InvoiceID)
}

func (s *Service) recalcInvoiceTotals(ctx context.Context, invoiceID string) error {
	items, err := s.DB.Queries.ListInvoiceItems(ctx, invoiceID)
	if err != nil {
		return fmt.Errorf("list invoice items: %w", err)
	}

	var subtotal, taxTotal int64
	for _, item := range items {
		subtotal += item.Amount
		taxAmount := int64(math.Round(float64(item.Amount) * item.TaxRate / 100))
		taxTotal += taxAmount
	}

	return s.DB.Queries.UpdateInvoiceTotals(ctx, dbsqlc.UpdateInvoiceTotalsParams{
		Subtotal:  subtotal,
		TaxTotal:  taxTotal,
		Total:     subtotal + taxTotal,
		UpdatedAt: now(),
		ID:        invoiceID,
	})
}

func (s *Service) ListInvoiceItems(ctx context.Context, invoiceID string) ([]dbsqlc.InvoiceItem, error) {
	return s.DB.Queries.ListInvoiceItems(ctx, invoiceID)
}

func (s *Service) UpdateInvoiceStatus(ctx context.Context, id, status string) error {
	return s.DB.Queries.UpdateInvoiceStatus(ctx, dbsqlc.UpdateInvoiceStatusParams{
		Status: status, UpdatedAt: now(), ID: id,
	})
}

func (s *Service) UpdateInvoicePDFPath(ctx context.Context, id, pdfPath string) error {
	return s.DB.Queries.UpdateInvoicePDFPath(ctx, dbsqlc.UpdateInvoicePDFPathParams{
		PdfPath: pdfPath, UpdatedAt: now(), ID: id,
	})
}

// --- Settings ---

func (s *Service) GetSettings(ctx context.Context) (dbsqlc.Setting, error) {
	return s.DB.Queries.GetSettings(ctx)
}

// --- Templates ---

func (s *Service) CreateTemplate(ctx context.Context, name, path string, isDefault bool) (string, error) {
	id := newID()
	ts := now()
	def := int64(0)
	if isDefault {
		def = 1
	}
	err := s.DB.Queries.CreateTemplate(ctx, dbsqlc.CreateTemplateParams{
		ID: id, Name: name, Path: path, IsDefault: def,
		CreatedAt: ts, UpdatedAt: ts,
	})
	if err != nil {
		return "", fmt.Errorf("create template: %w", err)
	}
	return id, nil
}

func (s *Service) GetTemplateByPath(ctx context.Context, path string) (dbsqlc.InvoiceTemplate, error) {
	return s.DB.Queries.GetTemplateByPath(ctx, path)
}

func (s *Service) ListTemplates(ctx context.Context) ([]dbsqlc.InvoiceTemplate, error) {
	return s.DB.Queries.ListTemplates(ctx)
}

// --- Recurring ---

func (s *Service) ListDueRecurringInvoices(ctx context.Context) ([]dbsqlc.RecurringInvoice, error) {
	return s.DB.Queries.ListDueRecurringInvoices(ctx, now())
}

func (s *Service) ListRecurringInvoiceItems(ctx context.Context, recurringID string) ([]dbsqlc.RecurringInvoiceItem, error) {
	return s.DB.Queries.ListRecurringInvoiceItems(ctx, recurringID)
}

func (s *Service) UpdateRecurringInvoiceNextRun(ctx context.Context, id, nextRun string) error {
	return s.DB.Queries.UpdateRecurringInvoiceNextRun(ctx, dbsqlc.UpdateRecurringInvoiceNextRunParams{
		NextRun: nextRun, LastRun: now(), UpdatedAt: now(), ID: id,
	})
}

func (s *Service) ListRecurringInvoices(ctx context.Context) ([]dbsqlc.RecurringInvoice, error) {
	return s.DB.Queries.ListRecurringInvoices(ctx)
}

func (s *Service) GetRecurringInvoice(ctx context.Context, id string) (dbsqlc.RecurringInvoice, error) {
	return s.DB.Queries.GetRecurringInvoice(ctx, id)
}

func (s *Service) UpdateRecurringInvoiceStatus(ctx context.Context, id, status string) error {
	return s.DB.Queries.UpdateRecurringInvoiceStatus(ctx, dbsqlc.UpdateRecurringInvoiceStatusParams{
		Status: status, UpdatedAt: now(), ID: id,
	})
}

// --- Payments ---

func (s *Service) RecordPayment(ctx context.Context, invoiceID string, amount int64, method, reference, notes string) (string, error) {
	id := newID()
	ts := now()
	err := s.DB.Queries.CreatePayment(ctx, dbsqlc.CreatePaymentParams{
		ID: id, InvoiceID: invoiceID, Amount: amount,
		Method: method, Reference: reference, PaidAt: ts,
		Notes: notes, CreatedAt: ts,
	})
	if err != nil {
		return "", fmt.Errorf("create payment: %w", err)
	}

	// Update invoice amount_paid
	totalPaidVal, err := s.DB.Queries.SumPaymentsByInvoice(ctx, invoiceID)
	if err != nil {
		return "", fmt.Errorf("sum payments: %w", err)
	}
	if err := s.DB.Queries.UpdateInvoiceAmountPaid(ctx, dbsqlc.UpdateInvoiceAmountPaidParams{
		AmountPaid: totalPaidVal, UpdatedAt: now(), ID: invoiceID,
	}); err != nil {
		return "", fmt.Errorf("update amount paid: %w", err)
	}

	// Auto-mark as paid if fully paid
	inv, err := s.DB.Queries.GetInvoice(ctx, invoiceID)
	if err != nil {
		return "", fmt.Errorf("get invoice: %w", err)
	}
	if inv.AmountPaid >= inv.Total && inv.Status != "paid" {
		_ = s.DB.Queries.UpdateInvoiceStatus(ctx, dbsqlc.UpdateInvoiceStatusParams{
			Status: "paid", UpdatedAt: now(), ID: invoiceID,
		})
	}

	return id, nil
}
