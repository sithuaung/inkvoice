-- Settings: single-row app config
CREATE TABLE IF NOT EXISTS settings (
    id          TEXT PRIMARY KEY DEFAULT 'default',
    company_name    TEXT NOT NULL DEFAULT '',
    company_email   TEXT NOT NULL DEFAULT '',
    company_phone   TEXT NOT NULL DEFAULT '',
    company_address TEXT NOT NULL DEFAULT '{}', -- JSON
    invoice_prefix  TEXT NOT NULL DEFAULT 'INK',
    next_invoice_number INTEGER NOT NULL DEFAULT 1,
    default_due_days    INTEGER NOT NULL DEFAULT 30,
    default_currency    TEXT NOT NULL DEFAULT 'USD',
    default_template_id TEXT NOT NULL DEFAULT '',
    created_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);

-- Insert default settings row
INSERT OR IGNORE INTO settings (id) VALUES ('default');

-- Clients
CREATE TABLE IF NOT EXISTS clients (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    email       TEXT NOT NULL DEFAULT '',
    phone       TEXT NOT NULL DEFAULT '',
    company     TEXT NOT NULL DEFAULT '',
    address     TEXT NOT NULL DEFAULT '{}', -- JSON
    notes       TEXT NOT NULL DEFAULT '',
    created_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);

-- Products
CREATE TABLE IF NOT EXISTS products (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    unit_price  INTEGER NOT NULL DEFAULT 0, -- cents
    currency    TEXT NOT NULL DEFAULT 'USD',
    created_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);

-- Taxes
CREATE TABLE IF NOT EXISTS taxes (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    rate        REAL NOT NULL DEFAULT 0, -- percentage e.g. 10.0 for 10%
    is_default  INTEGER NOT NULL DEFAULT 0,
    created_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);

-- Invoices
CREATE TABLE IF NOT EXISTS invoices (
    id              TEXT PRIMARY KEY,
    invoice_number  TEXT NOT NULL UNIQUE,
    client_id       TEXT NOT NULL REFERENCES clients(id),
    status          TEXT NOT NULL DEFAULT 'draft', -- draft, sent, paid, overdue, cancelled
    issue_date      TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    due_date        TEXT NOT NULL,
    subtotal        INTEGER NOT NULL DEFAULT 0, -- cents
    tax_total       INTEGER NOT NULL DEFAULT 0, -- cents
    total           INTEGER NOT NULL DEFAULT 0, -- cents
    amount_paid     INTEGER NOT NULL DEFAULT 0, -- cents
    currency        TEXT NOT NULL DEFAULT 'USD',
    notes           TEXT NOT NULL DEFAULT '',
    template_id     TEXT NOT NULL DEFAULT '',
    pdf_path        TEXT NOT NULL DEFAULT '',
    created_at      TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at      TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);

CREATE INDEX idx_invoices_client_id ON invoices(client_id);
CREATE INDEX idx_invoices_status ON invoices(status);
CREATE INDEX idx_invoices_invoice_number ON invoices(invoice_number);

-- Invoice Items
CREATE TABLE IF NOT EXISTS invoice_items (
    id          TEXT PRIMARY KEY,
    invoice_id  TEXT NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    product_id  TEXT, -- nullable, can be ad-hoc item
    description TEXT NOT NULL,
    quantity    REAL NOT NULL DEFAULT 1,
    unit_price  INTEGER NOT NULL DEFAULT 0, -- cents
    tax_id      TEXT, -- nullable
    tax_rate    REAL NOT NULL DEFAULT 0,
    amount      INTEGER NOT NULL DEFAULT 0, -- cents (quantity * unit_price)
    sort_order  INTEGER NOT NULL DEFAULT 0,
    created_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);

CREATE INDEX idx_invoice_items_invoice_id ON invoice_items(invoice_id);

-- Payments
CREATE TABLE IF NOT EXISTS payments (
    id          TEXT PRIMARY KEY,
    invoice_id  TEXT NOT NULL REFERENCES invoices(id),
    amount      INTEGER NOT NULL DEFAULT 0, -- cents
    method      TEXT NOT NULL DEFAULT '', -- bank_transfer, cash, card, etc.
    reference   TEXT NOT NULL DEFAULT '',
    paid_at     TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    notes       TEXT NOT NULL DEFAULT '',
    created_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);

CREATE INDEX idx_payments_invoice_id ON payments(invoice_id);

-- Recurring Invoices
CREATE TABLE IF NOT EXISTS recurring_invoices (
    id              TEXT PRIMARY KEY,
    client_id       TEXT NOT NULL REFERENCES clients(id),
    schedule        TEXT NOT NULL, -- cron expression
    status          TEXT NOT NULL DEFAULT 'active', -- active, paused
    next_run        TEXT NOT NULL,
    last_run        TEXT NOT NULL DEFAULT '',
    currency        TEXT NOT NULL DEFAULT 'USD',
    template_id     TEXT NOT NULL DEFAULT '',
    notes           TEXT NOT NULL DEFAULT '',
    created_at      TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at      TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);

CREATE INDEX idx_recurring_invoices_status ON recurring_invoices(status);
CREATE INDEX idx_recurring_invoices_next_run ON recurring_invoices(next_run);

-- Recurring Invoice Items
CREATE TABLE IF NOT EXISTS recurring_invoice_items (
    id                  TEXT PRIMARY KEY,
    recurring_invoice_id TEXT NOT NULL REFERENCES recurring_invoices(id) ON DELETE CASCADE,
    product_id          TEXT,
    description         TEXT NOT NULL,
    quantity            REAL NOT NULL DEFAULT 1,
    unit_price          INTEGER NOT NULL DEFAULT 0, -- cents
    tax_id              TEXT,
    tax_rate            REAL NOT NULL DEFAULT 0,
    sort_order          INTEGER NOT NULL DEFAULT 0,
    created_at          TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);

CREATE INDEX idx_recurring_invoice_items_recurring_id ON recurring_invoice_items(recurring_invoice_id);

-- Invoice Templates
CREATE TABLE IF NOT EXISTS invoice_templates (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    path        TEXT NOT NULL UNIQUE, -- path to .typ file
    is_default  INTEGER NOT NULL DEFAULT 0,
    created_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);
