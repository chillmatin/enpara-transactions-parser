package parser_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/chillmatin/enpara-transactions-parser/pkg/converter"
	"github.com/chillmatin/enpara-transactions-parser/pkg/parser"
)

func TestParseRealPDFIntegration(t *testing.T) {
	pdfPath := os.Getenv("ENPARA_TEST_PDF_PATH")
	if pdfPath == "" {
		t.Skip("ENPARA_TEST_PDF_PATH is not set")
	}

	text, err := parser.ExtractTextFromPDF(pdfPath)
	if err != nil {
		t.Fatalf("ExtractTextFromPDF returned error: %v", err)
	}

	statement, err := parser.ParseStatement(text)
	if err != nil {
		t.Fatalf("ParseStatement returned error: %v", err)
	}

	if len(statement.Transactions) == 0 {
		t.Fatal("expected at least one transaction")
	}

	if _, err := converter.ToJSON(statement); err != nil {
		t.Fatalf("ToJSON returned error: %v", err)
	}

	if _, err := json.Marshal(statement.Transactions[0]); err != nil {
		t.Fatalf("marshal first transaction: %v", err)
	}
	if _, err := json.Marshal(statement.Transactions[len(statement.Transactions)-1]); err != nil {
		t.Fatalf("marshal last transaction: %v", err)
	}
}

func TestParseProvidedTmpPDFAccountHolder(t *testing.T) {
	pdfPath := "../../tmp/transaction.pdf"
	if _, err := os.Stat(pdfPath); err != nil {
		t.Skip("tmp/transaction.pdf is not available")
	}

	text, err := parser.ExtractTextFromPDF(pdfPath)
	if err != nil {
		t.Fatalf("ExtractTextFromPDF returned error: %v", err)
	}

	statement, err := parser.ParseStatement(text)
	if err != nil {
		t.Fatalf("ParseStatement returned error: %v", err)
	}

	if statement.AccountHolder == "" {
		t.Fatal("expected non-empty account holder")
	}

	normalized := strings.ToUpper(strings.TrimSpace(statement.AccountHolder))
	if normalized == "AD SOYAD" {
		t.Fatalf("account holder should be extracted value, got label text: %q", statement.AccountHolder)
	}

	if len(statement.Transactions) == 0 {
		t.Fatal("expected at least one transaction")
	}
}
