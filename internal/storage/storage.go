package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// InvoiceStore defines the interface for storing generated PDFs.
type InvoiceStore interface {
	// Save stores a PDF and returns the storage path.
	Save(invoiceID string, data []byte) (string, error)
	// Path returns the full filesystem path for a stored invoice.
	Path(relativePath string) string
}

// LocalStore stores invoices on the local filesystem.
type LocalStore struct {
	BaseDir string // e.g. "./data/invoices"
}

// NewLocalStore creates a new local file store.
func NewLocalStore(baseDir string) *LocalStore {
	return &LocalStore{BaseDir: baseDir}
}

// Save writes PDF data to disk with year/month directory structure.
func (s *LocalStore) Save(invoiceID string, data []byte) (string, error) {
	now := time.Now().UTC()
	dir := filepath.Join(s.BaseDir, fmt.Sprintf("%d", now.Year()), fmt.Sprintf("%02d", now.Month()))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create directory: %w", err)
	}

	relPath := filepath.Join(fmt.Sprintf("%d", now.Year()), fmt.Sprintf("%02d", now.Month()), invoiceID+".pdf")
	fullPath := filepath.Join(s.BaseDir, relPath)
	if err := os.WriteFile(fullPath, data, 0o644); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}
	return relPath, nil
}

// Path returns the full path for a relative storage path.
func (s *LocalStore) Path(relativePath string) string {
	return filepath.Join(s.BaseDir, relativePath)
}
