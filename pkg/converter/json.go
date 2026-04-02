package converter

import (
	"encoding/json"

	"github.com/chillmatin/enpara-transactions-parser/pkg/models"
)

type exportJSONTransaction struct {
	Date        string  `json:"Tarih"`
	Type        string  `json:"Hareket tipi"`
	Description string  `json:"Açıklama"`
	Amount      float64 `json:"İşlem Tutarı"`
	Balance     float64 `json:"Bakiye"`
}

func ToJSON(statement *models.AccountStatement) ([]byte, error) {
	rows := make([]exportJSONTransaction, 0, len(statement.Transactions))
	for _, tx := range statement.Transactions {
		rows = append(rows, exportJSONTransaction{
			Date:        tx.Date.Format("02.01.2006"),
			Type:        tx.Type,
			Description: tx.Description,
			Amount:      tx.Amount,
			Balance:     tx.Balance,
		})
	}

	return json.MarshalIndent(rows, "", "  ")
}
