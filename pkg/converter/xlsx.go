package converter

import (
	"fmt"

	"github.com/chillmatin/enpara-transactions-parser/pkg/models"
	"github.com/xuri/excelize/v2"
)

func ToXLSX(statement *models.AccountStatement) ([]byte, error) {
	file := excelize.NewFile()
	defaultSheet := file.GetSheetName(0)
	sheetName := "Transactions"
	file.SetSheetName(defaultSheet, sheetName)

	headers := []string{"Date", "Type", "Merchant", "Description", "Amount", "Balance"}
	for i, header := range headers {
		cell, err := excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			return nil, fmt.Errorf("header cell coordinates: %w", err)
		}
		if err := file.SetCellValue(sheetName, cell, header); err != nil {
			return nil, fmt.Errorf("set header value: %w", err)
		}
	}

	headerStyle, err := file.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	if err != nil {
		return nil, fmt.Errorf("create header style: %w", err)
	}
	if err := file.SetCellStyle(sheetName, "A1", "F1", headerStyle); err != nil {
		return nil, fmt.Errorf("apply header style: %w", err)
	}

	dateFmt := "dd.mm.yyyy"
	dateStyle, err := file.NewStyle(&excelize.Style{CustomNumFmt: &dateFmt})
	if err != nil {
		return nil, fmt.Errorf("create date style: %w", err)
	}

	amountFmt := "#,##0.00"
	amountStyle, err := file.NewStyle(&excelize.Style{CustomNumFmt: &amountFmt})
	if err != nil {
		return nil, fmt.Errorf("create amount style: %w", err)
	}

	columnWidths := []float64{12, 18, 22, 40, 14, 14}

	for rowIndex, transaction := range statement.Transactions {
		excelRow := rowIndex + 2
		if err := file.SetCellValue(sheetName, fmt.Sprintf("A%d", excelRow), transaction.Date); err != nil {
			return nil, fmt.Errorf("set date value: %w", err)
		}
		if err := file.SetCellValue(sheetName, fmt.Sprintf("B%d", excelRow), transaction.Type); err != nil {
			return nil, fmt.Errorf("set type value: %w", err)
		}
		if err := file.SetCellValue(sheetName, fmt.Sprintf("C%d", excelRow), transaction.Merchant); err != nil {
			return nil, fmt.Errorf("set merchant value: %w", err)
		}
		if err := file.SetCellValue(sheetName, fmt.Sprintf("D%d", excelRow), transaction.Description); err != nil {
			return nil, fmt.Errorf("set description value: %w", err)
		}
		if err := file.SetCellValue(sheetName, fmt.Sprintf("E%d", excelRow), transaction.Amount); err != nil {
			return nil, fmt.Errorf("set amount value: %w", err)
		}
		if err := file.SetCellValue(sheetName, fmt.Sprintf("F%d", excelRow), transaction.Balance); err != nil {
			return nil, fmt.Errorf("set balance value: %w", err)
		}
		columnWidths[1] = maxFloat(columnWidths[1], float64(len([]rune(transaction.Type))+2))
		columnWidths[2] = maxFloat(columnWidths[2], float64(len([]rune(transaction.Merchant))+2))
		columnWidths[3] = maxFloat(columnWidths[3], float64(len([]rune(transaction.Description))+2))
	}

	if len(statement.Transactions) > 0 {
		last := len(statement.Transactions) + 1
		if err := file.SetCellStyle(sheetName, "A2", fmt.Sprintf("A%d", last), dateStyle); err != nil {
			return nil, fmt.Errorf("apply date style: %w", err)
		}
		if err := file.SetCellStyle(sheetName, "E2", fmt.Sprintf("F%d", last), amountStyle); err != nil {
			return nil, fmt.Errorf("apply amount style: %w", err)
		}
	}

	for i, width := range columnWidths {
		columnName, err := excelize.ColumnNumberToName(i + 1)
		if err != nil {
			return nil, fmt.Errorf("column number to name: %w", err)
		}
		if err := file.SetColWidth(sheetName, columnName, columnName, width); err != nil {
			return nil, fmt.Errorf("set column width: %w", err)
		}
	}

	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("write xlsx buffer: %w", err)
	}

	return buffer.Bytes(), nil
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
