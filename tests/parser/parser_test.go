package parser_test

import (
	"testing"

	"github.com/chillmatin/enpara-transactions-parser/pkg/parser"
)

func TestParseStatementBasic(t *testing.T) {
	pdfText := `
	Hesap Sahibi: Test User
	Hesap No: ****1234
	IBAN: TR12 3456 7890 1234 5678 9012 34
	Donem: 01.03.2026 - 31.03.2026
	Islem Tarihi Aciklama Tutar Bakiye
	28.03.2026 Diğer 000000004228140-DED COFFEE IZMIR TR -52,50 TL 2.357,27 TL
	29.03.2026 Gelen Transfer Maasi Odemesi 5.000,00 TL 7.357,27 TL
	`

	statement, err := parser.ParseStatement(pdfText)
	if err != nil {
		t.Fatalf("ParseStatement returned error: %v", err)
	}

	if statement.AccountHolder != "Test User" {
		t.Fatalf("unexpected account holder: %q", statement.AccountHolder)
	}

	if statement.IBAN != "TR123456789012345678901234" {
		t.Fatalf("unexpected iban: %q", statement.IBAN)
	}

	if len(statement.Transactions) != 2 {
		t.Fatalf("unexpected transaction count: %d", len(statement.Transactions))
	}

	first := statement.Transactions[0]
	if first.Type != "Diğer" {
		t.Fatalf("unexpected type: %q", first.Type)
	}
	if first.Merchant != "DED COFFEE" {
		t.Fatalf("unexpected merchant: %q", first.Merchant)
	}
	if first.Amount != -52.50 {
		t.Fatalf("unexpected amount: %v", first.Amount)
	}
	if first.Balance != 2357.27 {
		t.Fatalf("unexpected balance: %v", first.Balance)
	}

	last := statement.Transactions[len(statement.Transactions)-1]
	if last.Amount <= 0 {
		t.Fatalf("expected positive amount for income, got %v", last.Amount)
	}
}

func TestParseStatementWrappedDescription(t *testing.T) {
	pdfText := `
	28.03.2026 Encard Harcaması 1802326580 -MAVIS EV
	YEMEKLERI ŞANLIURFA TR -120,75 TL 2.236,52 TL
	`

	statement, err := parser.ParseStatement(pdfText)
	if err != nil {
		t.Fatalf("ParseStatement returned error: %v", err)
	}

	if len(statement.Transactions) != 1 {
		t.Fatalf("unexpected transaction count: %d", len(statement.Transactions))
	}

	tx := statement.Transactions[0]
	if tx.Type != "Encard Harcaması" {
		t.Fatalf("unexpected transaction type: %q", tx.Type)
	}
	if tx.Merchant != "MAVIS EV YEMEKLERI" {
		t.Fatalf("unexpected merchant: %q", tx.Merchant)
	}
}

func TestParseStatementAssignsDailySequence(t *testing.T) {
	pdfText := `
	28.03.2026 Diğer 000000004228140-DED COFFEE IZMIR TR -52,50 TL 2.357,27 TL
	28.03.2026 Diğer 000000004228140-DED COFFEE IZMIR TR -12,50 TL 2.344,77 TL
	29.03.2026 Gelen Transfer Maasi Odemesi 5.000,00 TL 7.344,77 TL
	`

	statement, err := parser.ParseStatement(pdfText)
	if err != nil {
		t.Fatalf("ParseStatement returned error: %v", err)
	}

	if len(statement.Transactions) != 3 {
		t.Fatalf("unexpected transaction count: %d", len(statement.Transactions))
	}

	if got := statement.Transactions[0].DailySequence; got != 1 {
		t.Fatalf("first transaction daily sequence = %d, expected 1", got)
	}
	if got := statement.Transactions[1].DailySequence; got != 2 {
		t.Fatalf("second transaction daily sequence = %d, expected 2", got)
	}
	if got := statement.Transactions[2].DailySequence; got != 1 {
		t.Fatalf("third transaction daily sequence = %d, expected 1", got)
	}
}

func TestParseStatementMetadataAndSkips(t *testing.T) {
	pdfText := `
	Hesap Sahibi : Test User
	Hesap Numarası: 1234
	İŞLEM TARİHİ AÇIKLAMA TUTAR BAKİYE
	28.03.2026 Diğer ORNEK MARKET IZMIR TR -10,00 TL 90,00 TL
	Sayfa 1 / 3
	`

	statement, err := parser.ParseStatement(pdfText)
	if err != nil {
		t.Fatalf("ParseStatement returned error: %v", err)
	}

	if statement.AccountHolder != "Test User" {
		t.Fatalf("unexpected account holder: %q", statement.AccountHolder)
	}
	if statement.AccountNumber != "1234" {
		t.Fatalf("unexpected account number: %q", statement.AccountNumber)
	}
	if len(statement.Transactions) != 1 {
		t.Fatalf("unexpected transaction count: %d", len(statement.Transactions))
	}
}

func TestParseStatementAccountHolderFromAdSoyad(t *testing.T) {
	pdfText := `
	Ad Soyad: Test User
	Hesap Numarası: 1234
	28.03.2026 Diğer ORNEK MARKET IZMIR TR -10,00 TL 90,00 TL
	`

	statement, err := parser.ParseStatement(pdfText)
	if err != nil {
		t.Fatalf("ParseStatement returned error: %v", err)
	}

	if statement.AccountHolder != "Test User" {
		t.Fatalf("unexpected account holder: %q", statement.AccountHolder)
	}
}
