package scheduler

import (
	"context"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/sithuaung/inkvoice/internal/service"
)

// Scheduler runs recurring invoice processing on a cron schedule.
type Scheduler struct {
	svc  *service.Service
	cron *cron.Cron
}

// New creates a new Scheduler.
func New(svc *service.Service) *Scheduler {
	return &Scheduler{
		svc:  svc,
		cron: cron.New(),
	}
}

// Start begins the cron scheduler. It checks for due recurring invoices every minute.
func (s *Scheduler) Start() {
	s.cron.AddFunc("* * * * *", func() {
		s.processDue()
	})
	s.cron.Start()
	slog.Info("scheduler started")
}

// Stop stops the cron scheduler.
func (s *Scheduler) Stop() {
	s.cron.Stop()
	slog.Info("scheduler stopped")
}

func (s *Scheduler) processDue() {
	ctx := context.Background()
	dues, err := s.svc.ListDueRecurringInvoices(ctx)
	if err != nil {
		slog.Error("list due recurring invoices", "error", err)
		return
	}

	for _, ri := range dues {
		invoiceID, err := s.svc.CreateInvoice(ctx, ri.ClientID, ri.Notes)
		if err != nil {
			slog.Error("create invoice from recurring", "recurring_id", ri.ID, "error", err)
			continue
		}

		// Copy recurring items to the new invoice
		items, err := s.svc.ListRecurringInvoiceItems(ctx, ri.ID)
		if err != nil {
			slog.Error("list recurring items", "recurring_id", ri.ID, "error", err)
			continue
		}
		for _, item := range items {
			_, err := s.svc.AddInvoiceItem(ctx, invoiceID,
				item.ProductID.String, item.Description,
				item.Quantity, item.UnitPrice,
				item.TaxID.String, item.TaxRate,
			)
			if err != nil {
				slog.Error("add invoice item from recurring", "error", err)
			}
		}

		// Parse the cron schedule to calculate next run
		parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		sched, err := parser.Parse(ri.Schedule)
		if err != nil {
			slog.Error("parse cron schedule", "schedule", ri.Schedule, "error", err)
			continue
		}
		nextRun := sched.Next(time.Now().UTC()).Format(time.RFC3339)
		if err := s.svc.UpdateRecurringInvoiceNextRun(ctx, ri.ID, nextRun); err != nil {
			slog.Error("update next run", "recurring_id", ri.ID, "error", err)
		}

		slog.Info("generated invoice from recurring", "invoice_id", invoiceID, "recurring_id", ri.ID)
	}
}

// ProcessDue is exported for manual triggering via CLI.
func (s *Scheduler) ProcessDue() {
	s.processDue()
}
