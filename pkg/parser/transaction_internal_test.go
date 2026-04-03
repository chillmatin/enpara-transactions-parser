package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractType2NFCHintsFromPDFAutomatic(t *testing.T) {
	pdfPath := filepath.Join("..", "..", "tmp", "automatic.pdf")
	if _, err := os.Stat(pdfPath); err != nil {
		t.Skip("tmp/automatic.pdf is not available")
	}

	hints, err := extractType2NFCHintsFromPDF(pdfPath)
	if err != nil {
		t.Fatalf("extractType2NFCHintsFromPDF returned error: %v", err)
	}
	if len(hints) == 0 {
		t.Fatal("expected non-empty NFC hints")
	}

	trueCount := 0
	for _, value := range hints {
		if value {
			trueCount++
		}
	}
	if trueCount == 0 {
		t.Fatal("expected at least one true NFC hint")
	}
}

func TestParseStatementFromPDFAppliesExtractedHints(t *testing.T) {
	pdfPath := filepath.Join("..", "..", "tmp", "automatic.pdf")
	if _, err := os.Stat(pdfPath); err != nil {
		t.Skip("tmp/automatic.pdf is not available")
	}

	hints, err := extractType2NFCHintsFromPDF(pdfPath)
	if err != nil {
		t.Fatalf("extractType2NFCHintsFromPDF returned error: %v", err)
	}

	withHints, err := ParseStatementFromPDF(pdfPath, ParseOptions{PDFType: PDFType2, Type2NFCHints: hints})
	if err != nil {
		t.Fatalf("ParseStatementFromPDF with hints returned error: %v", err)
	}

	autoHints, err := ParseStatementFromPDF(pdfPath, ParseOptions{PDFType: PDFType2})
	if err != nil {
		t.Fatalf("ParseStatementFromPDF auto hints returned error: %v", err)
	}

	if len(withHints.Transactions) < 2 || len(autoHints.Transactions) < 2 {
		t.Fatalf("expected at least 2 transactions, got withHints=%d autoHints=%d", len(withHints.Transactions), len(autoHints.Transactions))
	}

	if withHints.Transactions[1].NFC && !autoHints.Transactions[1].NFC {
		t.Fatalf("expected auto-hint path to match explicit hints at row 2")
	}
}
