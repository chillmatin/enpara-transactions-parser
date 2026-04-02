package converter

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/chillmatin/enpara-transactions-parser/pkg/models"
)

func ToCSV(statement *models.AccountStatement) ([]byte, error) {
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)
	writer.Comma = ';'

	headers := []string{"Tarih", "Hareket tipi", "Açıklama", "İşlem Tutarı", "Bakiye"}
	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("write csv headers: %w", err)
	}

	for _, transaction := range statement.Transactions {
		row := []string{
			transaction.Date.Format("02.01.2006"),
			transaction.Type,
			transaction.Description,
			formatTurkishDecimal(transaction.Amount),
			formatTurkishDecimal(transaction.Balance),
		}

		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("write csv row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("flush csv writer: %w", err)
	}

	return buffer.Bytes(), nil
}

func formatTurkishDecimal(value float64) string {
	raw := fmt.Sprintf("%.2f", value)
	sign := ""
	if strings.HasPrefix(raw, "-") {
		sign = "-"
		raw = strings.TrimPrefix(raw, "-")
	}

	parts := strings.Split(raw, ".")
	integerPart := parts[0]
	decimalPart := parts[1]

	if len(integerPart) > 3 {
		var grouped strings.Builder
		for i, digit := range integerPart {
			if i > 0 && (len(integerPart)-i)%3 == 0 {
				grouped.WriteByte('.')
			}
			grouped.WriteRune(digit)
		}
		integerPart = grouped.String()
	}

	return sign + integerPart + "," + decimalPart
}
