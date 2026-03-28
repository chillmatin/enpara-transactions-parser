package converter

import (
	"encoding/json"

	"github.com/chillmatin/enpara-transactions-parser/pkg/models"
)

func ToJSON(statement *models.AccountStatement) ([]byte, error) {
	return json.MarshalIndent(statement, "", "  ")
}
