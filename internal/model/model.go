package model

import "time"

type Address struct {
	Street  string `json:"street,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Zip     string `json:"zip,omitempty"`
	Country string `json:"country,omitempty"`
}

type Settings struct {
	ID                string  `json:"id"`
	CompanyName       string  `json:"company_name"`
	CompanyEmail      string  `json:"company_email"`
	CompanyPhone      string  `json:"company_phone"`
	CompanyAddress    string  `json:"company_address"` // JSON
	InvoicePrefix     string  `json:"invoice_prefix"`
	NextInvoiceNumber int64   `json:"next_invoice_number"`
	DefaultDueDays    int64   `json:"default_due_days"`
	DefaultCurrency   string  `json:"default_currency"`
	DefaultTemplateID string  `json:"default_template_id"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

type Client struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Company   string `json:"company"`
	Address   string `json:"address"` // JSON
	Notes     string `json:"notes"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UnitPrice   int64  `json:"unit_price"` // cents
	Currency    string `json:"currency"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type Tax struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Rate      float64 `json:"rate"`
	IsDefault bool    `json:"is_default"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type Invoice struct {
	ID            string `json:"id"`
	InvoiceNumber string `json:"invoice_number"`
	ClientID      string `json:"client_id"`
	Status        string `json:"status"`
	IssueDate     string `json:"issue_date"`
	DueDate       string `json:"due_date"`
	Subtotal      int64  `json:"subtotal"`
	TaxTotal      int64  `json:"tax_total"`
	Total         int64  `json:"total"`
	AmountPaid    int64  `json:"amount_paid"`
	Currency      string `json:"currency"`
	Notes         string `json:"notes"`
	TemplateID    string `json:"template_id"`
	PDFPath       string `json:"pdf_path"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type InvoiceItem struct {
	ID          string  `json:"id"`
	InvoiceID   string  `json:"invoice_id"`
	ProductID   string  `json:"product_id"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   int64   `json:"unit_price"` // cents
	TaxID       string  `json:"tax_id"`
	TaxRate     float64 `json:"tax_rate"`
	Amount      int64   `json:"amount"` // cents
	SortOrder   int64   `json:"sort_order"`
	CreatedAt   string  `json:"created_at"`
}

type Payment struct {
	ID        string `json:"id"`
	InvoiceID string `json:"invoice_id"`
	Amount    int64  `json:"amount"` // cents
	Method    string `json:"method"`
	Reference string `json:"reference"`
	PaidAt    string `json:"paid_at"`
	Notes     string `json:"notes"`
	CreatedAt string `json:"created_at"`
}

type RecurringInvoice struct {
	ID         string `json:"id"`
	ClientID   string `json:"client_id"`
	Schedule   string `json:"schedule"`
	Status     string `json:"status"`
	NextRun    string `json:"next_run"`
	LastRun    string `json:"last_run"`
	Currency   string `json:"currency"`
	TemplateID string `json:"template_id"`
	Notes      string `json:"notes"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type RecurringInvoiceItem struct {
	ID                 string  `json:"id"`
	RecurringInvoiceID string  `json:"recurring_invoice_id"`
	ProductID          string  `json:"product_id"`
	Description        string  `json:"description"`
	Quantity           float64 `json:"quantity"`
	UnitPrice          int64   `json:"unit_price"`
	TaxID              string  `json:"tax_id"`
	TaxRate            float64 `json:"tax_rate"`
	SortOrder          int64   `json:"sort_order"`
	CreatedAt          string  `json:"created_at"`
}

type InvoiceTemplate struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	IsDefault bool   `json:"is_default"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// FormatMoney converts cents to a display string like "$12.50"
func FormatMoney(cents int64, currency string) string {
	dollars := float64(cents) / 100.0
	return currency + " " + formatFloat(dollars)
}

func formatFloat(f float64) string {
	s := ""
	if f < 0 {
		s = "-"
		f = -f
	}
	whole := int64(f)
	frac := int64((f - float64(whole) + 0.005) * 100)
	s += intToStr(whole) + "." + padTwo(frac)
	return s
}

func intToStr(n int64) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}

func padTwo(n int64) string {
	if n < 10 {
		return "0" + intToStr(n)
	}
	return intToStr(n)
}

// NowRFC3339 returns the current time in RFC3339 format
func NowRFC3339() string {
	return time.Now().UTC().Format(time.RFC3339)
}
