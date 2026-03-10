# Inkvoice

A CLI-first invoicing tool for freelancers and small businesses. Written in Go.

Single-user, local-first. Your data stays on your machine in a single SQLite file.

## Features

- Create and manage invoices, clients, and products from the terminal
- Generate PDF invoices using [Typst](https://typst.app) templates
- Recurring invoices with automatic scheduling
- Send invoices and payment reminders via email
- Export data to CSV/JSON
- One-command database backup

## Requirements

- Go 1.25+
- [Typst](https://typst.app) CLI (`brew install typst`)

## Install

```bash
go install github.com/sithuaung/inkvoice/cmd/inkvoice@latest
```

Or build from source:

```bash
git clone https://github.com/sithuaung/inkvoice.git
cd inkvoice
go build -o inkvoice ./cmd/inkvoice
```

## Try it with Docker

```bash
git clone https://github.com/sithuaung/inkvoice.git
cd inkvoice
docker compose up -d
```

Done. Database is migrated, demo data and templates are loaded. Use the CLI:

```bash
docker compose exec inkvoice inkvoice invoice list
docker compose exec inkvoice inkvoice invoice preview INK-0001
```

## Quick Start (local install)

```bash
# Set up the database
inkvoice migrate up

# Load demo data to try it out
inkvoice seed data

# Load invoice templates
inkvoice seed template

# Add a client
inkvoice client create --name "Acme Corp" --email "billing@acme.com"

# Add a product
inkvoice product create --name "Web Development" --unit-price 15000 --unit hour

# Create an invoice
inkvoice invoice create --client "Acme Corp"

# Generate PDF
inkvoice invoice pdf INK-0001 --output invoice.pdf

# Send it
inkvoice invoice send INK-0001
```

## Commands

```
inkvoice
├── serve              Start cron scheduler (recurring invoices)
├── migrate            Database migrations (up, down, status)
├── invoice            Create, list, show, send, generate PDF
├── client             Create, list, show, update, delete clients
├── product            Create, list, update, delete products
├── recurring          Manage recurring invoice schedules
├── seed
│   ├── data           Insert demo data
│   └── template       Register invoice templates from invoice-templates/
├── export             Export invoices/clients to CSV/JSON
├── backup             Safe copy of SQLite database
├── health             Check system status
└── --version          Print version info
```

Use `inkvoice [command] --help` for details on any command.

## Configuration

Inkvoice uses environment variables for configuration:

```bash
# Database (default: ./inkvoice.db)
INKVOICE_DB_PATH=./inkvoice.db

# Email (SMTP)
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=you@example.com
SMTP_PASSWORD=your-password

# Invoice storage (default: ./data/invoices)
INVOICE_STORAGE_PATH=./data/invoices
```

Or pass `--db` and `--config` flags to any command.

## PDF Templates

Invoice templates are Typst files in the `invoice-templates/` directory. Register them with:

```bash
inkvoice seed template
```

Templates receive invoice data via `sys.inputs`:

```typst
#let client_name = sys.inputs.at("client_name", default: "")
#let invoice_number = sys.inputs.at("invoice_number", default: "")
```

## License

MIT
