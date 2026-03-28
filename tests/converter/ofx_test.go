package converter_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/chillmatin/enpara-transactions-parser/pkg/converter"
	"github.com/chillmatin/enpara-transactions-parser/pkg/models"
)

func TestToOFXFITIDDeterministicAcrossRuns(t *testing.T) {
	statement := &models.AccountStatement{
		StartDate: mustDate(t, "2024-07-22"),
		EndDate:   mustDate(t, "2024-07-22"),
		Transactions: []models.Transaction{
			{
				Date:          mustDate(t, "2024-07-22"),
				Amount:        -250.00,
				Balance:       1660.21,
				DailySequence: 1,
				Merchant:      "Hafiz Huseynzade",
				Description:   "Hafiz Huseynzade, mone",
				RawText:       "Hafiz Huseynzade, mone",
			},
			{
				Date:          mustDate(t, "2024-07-22"),
				Amount:        -250.00,
				Balance:       1660.21,
				DailySequence: 2,
				Merchant:      "Hafiz Huseynzade",
				Description:   "Hafiz Huseynzade, mone",
				RawText:       "Hafiz Huseynzade, mone",
			},
		},
	}

	firstOut, err := converter.ToOFX(statement)
	if err != nil {
		t.Fatalf("ToOFX first run error: %v", err)
	}

	secondOut, err := converter.ToOFX(statement)
	if err != nil {
		t.Fatalf("ToOFX second run error: %v", err)
	}

	firstFITIDs := extractFITIDs(t, string(firstOut))
	secondFITIDs := extractFITIDs(t, string(secondOut))

	if len(firstFITIDs) != 2 || len(secondFITIDs) != 2 {
		t.Fatalf("expected 2 FITIDs in both runs")
	}

	if firstFITIDs[0] == firstFITIDs[1] {
		t.Fatalf("expected same-day duplicates to have distinct FITIDs")
	}

	if firstFITIDs[0] != secondFITIDs[0] || firstFITIDs[1] != secondFITIDs[1] {
		t.Fatalf("expected FITIDs to be deterministic across identical runs")
	}
}

func TestToOFXUsesModelDailySequence(t *testing.T) {
	statement := &models.AccountStatement{
		StartDate: mustDate(t, "2024-07-22"),
		EndDate:   mustDate(t, "2024-07-22"),
		Transactions: []models.Transaction{
			{
				Date:          mustDate(t, "2024-07-22"),
				Amount:        -250.00,
				Balance:       1660.21,
				DailySequence: 5,
				Description:   "Hafiz Huseynzade, mone",
				RawText:       "Hafiz Huseynzade, mone",
			},
			{
				Date:          mustDate(t, "2024-07-22"),
				Amount:        -250.00,
				Balance:       1660.21,
				DailySequence: 7,
				Description:   "Hafiz Huseynzade, mone",
				RawText:       "Hafiz Huseynzade, mone",
			},
		},
	}

	out, err := converter.ToOFX(statement)
	if err != nil {
		t.Fatalf("ToOFX error: %v", err)
	}

	fitids := extractFITIDs(t, string(out))
	if len(fitids) != 2 {
		t.Fatalf("expected 2 FITIDs")
	}
	if fitids[0] == fitids[1] {
		t.Fatalf("expected distinct FITIDs for sequence 5 and 7")
	}
}

func extractFITIDs(t *testing.T, content string) []string {
	t.Helper()

	re := regexp.MustCompile(`<FITID>([^<]+)</FITID>`)
	matches := re.FindAllStringSubmatch(content, -1)
	fitids := make([]string, 0, len(matches))
	for _, m := range matches {
		if len(m) == 2 {
			fitids = append(fitids, m[1])
		}
	}

	return fitids
}

func mustDate(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		t.Fatalf("parse date %q: %v", value, err)
	}
	return parsed
}
