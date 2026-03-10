package pdf

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Generator generates PDFs using Typst CLI.
type Generator struct {
	TemplatesDir string // directory containing .typ templates
}

// NewGenerator creates a new PDF generator.
func NewGenerator(templatesDir string) *Generator {
	return &Generator{TemplatesDir: templatesDir}
}

// Generate compiles a Typst template with the given input variables and returns the PDF bytes.
func (g *Generator) Generate(ctx context.Context, templatePath string, inputs map[string]string) ([]byte, error) {
	// Create a temp file for the output
	tmpFile, err := os.CreateTemp("", "inkvoice-*.pdf")
	if err != nil {
		return nil, fmt.Errorf("create temp file: %w", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// Build the typst compile command
	args := []string{"compile"}
	for k, v := range inputs {
		args = append(args, "--input", k+"="+v)
	}
	args = append(args, templatePath, tmpFile.Name())

	cmd := exec.CommandContext(ctx, "typst", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("typst compile: %w\n%s", err, string(output))
	}

	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("read pdf: %w", err)
	}
	return data, nil
}

// TypstAvailable checks if the typst CLI is available.
func TypstAvailable() bool {
	_, err := exec.LookPath("typst")
	return err == nil
}

// FindTemplates scans a directory for .typ files.
func FindTemplates(dir string) ([]string, error) {
	var templates []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".typ" {
			templates = append(templates, e.Name())
		}
	}
	return templates, nil
}
