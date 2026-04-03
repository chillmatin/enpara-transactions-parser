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
	pdfPath := "../../tmp/manual.pdf"
	if _, err := os.Stat(pdfPath); err != nil {
		t.Skip("tmp/manual.pdf is not available")
	}

	text, err := parser.ExtractTextFromPDF(pdfPath)
	if err != nil {
		t.Fatalf("ExtractTextFromPDF returned error: %v", err)
	}

	statement, err := parser.ParseStatementWithOptions(text, parser.ParseOptions{PDFType: parser.PDFType1})
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

func TestParseProvidedAutomaticPDFAutoDetect(t *testing.T) {
	pdfPath := "../../tmp/automatic.pdf"
	if _, err := os.Stat(pdfPath); err != nil {
		t.Skip("tmp/automatic.pdf is not available")
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
}

func TestParseProvidedAutomaticPDFType2(t *testing.T) {
	pdfPath := "../../tmp/automatic.pdf"
	if _, err := os.Stat(pdfPath); err != nil {
		t.Skip("tmp/automatic.pdf is not available")
	}

	text, err := parser.ExtractTextFromPDF(pdfPath)
	if err != nil {
		t.Fatalf("ExtractTextFromPDF returned error: %v", err)
	}

	statement, err := parser.ParseStatementWithOptions(text, parser.ParseOptions{PDFType: parser.PDFType2})
	if err != nil {
		t.Fatalf("ParseStatementWithOptions(type2) returned error: %v", err)
	}

	if len(statement.Transactions) == 0 {
		t.Fatal("expected at least one transaction")
	}
}

func TestParseProvidedAutomaticPDFNFCByIcon(t *testing.T) {
	pdfPath := "../../tmp/automatic.pdf"
	if _, err := os.Stat(pdfPath); err != nil {
		t.Skip("tmp/automatic.pdf is not available")
	}

	statement, err := parser.ParseStatementFromPDF(pdfPath, parser.ParseOptions{PDFType: parser.PDFType2})
	if err != nil {
		t.Fatalf("ParseStatementFromPDF(type2) returned error: %v", err)
	}

	if len(statement.Transactions) == 0 {
		t.Fatal("expected at least one transaction")
	}

	var odealFound bool
	var transferFound bool

	for _, tx := range statement.Transactions {
		if strings.Contains(strings.ToUpper(tx.Description), "ODEAL//YASEMIN YILDI") {
			odealFound = true
			if !tx.NFC {
				t.Fatalf("expected NFC=true for ODEAL row, got false: %q", tx.Description)
			}
		}

		normalizedDescription := strings.ToUpper(tx.Description)
		if strings.Contains(normalizedDescription, "DEBIT KARTTAN ENPARA DEBIT KARTA") {
			transferFound = true
			if tx.NFC {
				t.Fatalf("expected NFC=false for transfer row, got true: %q", tx.Description)
			}
		}
	}

	if !odealFound {
		t.Fatal("could not find ODEAL sample row for NFC assertion")
	}
	if !transferFound {
		t.Fatal("could not find transfer sample row for NFC assertion")
	}
}
