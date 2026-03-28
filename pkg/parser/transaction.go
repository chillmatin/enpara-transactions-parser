package parser

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/chillmatin/enpara-transactions-parser/pkg/models"
)

const (
	turkishDateLayout = "02.01.2006"
	currencySuffixTL  = "TL"

	transactionTypeIncomingTransfer = "Gelen Transfer"
	transactionTypeOutgoingTransfer = "Giden Transfer"
	transactionTypeOther            = "Diğer"
	transactionTypeCardSpend        = "Encard Harcaması"
	transactionTypeCashWithdraw     = "Para Çekme"
	transactionTypeFee              = "Masraf/Ücret"
	transactionTypeRefund           = "İptal/İade"
)

var (
	dateLinePattern      = regexp.MustCompile(`^\d{2}\.\d{2}\.\d{4}\b`)
	trAmountPattern      = regexp.MustCompile(`([+-]?\d{1,3}(?:\.\d{3})*,\d{2})\s*` + currencySuffixTL)
	ibanPattern          = regexp.MustCompile(`\bTR\d{2}(?:\s?\d{4}){5}\s?\d{2}\b`)
	dateRangePattern     = regexp.MustCompile(`(\d{2}\.\d{2}\.\d{4}).*?(\d{2}\.\d{2}\.\d{4})`)
	numericPrefixPattern = regexp.MustCompile(`^\d+\s*-\s*`)
	compactPrefixPattern = regexp.MustCompile(`^\d+-`)
	extraSpacePattern    = regexp.MustCompile(`\s+`)
	merchantCityTRSuffix = regexp.MustCompile(`\s+[\p{L}]+\s+TR$`)

	turkishToASCIIReplacer = strings.NewReplacer(
		"İ", "I",
		"I", "I",
		"Ş", "S",
		"Ğ", "G",
		"Ü", "U",
		"Ö", "O",
		"Ç", "C",
	)
)

var knownTransactionTypes = []string{
	transactionTypeIncomingTransfer,
	transactionTypeOutgoingTransfer,
	transactionTypeOther,
	transactionTypeCardSpend,
	transactionTypeCashWithdraw,
	transactionTypeFee,
	transactionTypeRefund,
}

var metadataAccountHolderKeywords = []string{
	"HESAP SAHIBI",
	"AD SOYAD",
}

var metadataAccountNumberKeywords = []string{
	"HESAP NO",
	"HESAP NUMARASI",
}

var continuationSkipKeywords = []string{
	"ISLEM TARIHI",
	"ACIKLAMA",
	"SAYFA",
}

func init() {
	sort.SliceStable(knownTransactionTypes, func(i, j int) bool {
		return len(knownTransactionTypes[i]) > len(knownTransactionTypes[j])
	})
}

func ParseStatement(pdfText string) (*models.AccountStatement, error) {
	lines := normalizeLines(pdfText)
	statement := &models.AccountStatement{}

	extractMetadata(lines, statement)
	rows := collectTransactionRows(lines)
	if len(rows) == 0 {
		return nil, fmt.Errorf("no transaction rows detected")
	}

	transactions := make([]models.Transaction, 0, len(rows))
	errorCount := 0

	for _, row := range rows {
		tx, err := parseTransactionRow(row)
		if err != nil {
			errorCount++
			continue
		}
		transactions = append(transactions, tx)
	}

	if len(transactions) == 0 {
		return nil, fmt.Errorf("failed to parse all %d transaction rows", len(rows))
	}

	if errorCount > len(rows)/2 {
		return nil, fmt.Errorf("too many malformed rows: %d of %d", errorCount, len(rows))
	}

	assignDailySequence(transactions)

	statement.Transactions = transactions
	if statement.StartDate.IsZero() || statement.EndDate.IsZero() {
		statement.StartDate, statement.EndDate = deriveDateRange(transactions)
	}

	return statement, nil
}

func normalizeLines(pdfText string) []string {
	pdfText = strings.ReplaceAll(pdfText, "\r\n", "\n")
	pdfText = strings.ReplaceAll(pdfText, "\r", "\n")

	rawLines := strings.Split(pdfText, "\n")
	lines := make([]string, 0, len(rawLines))
	for _, line := range rawLines {
		cleaned := strings.TrimSpace(extraSpacePattern.ReplaceAllString(line, " "))
		if cleaned == "" {
			continue
		}
		lines = append(lines, cleaned)
	}

	return lines
}

func extractMetadata(lines []string, statement *models.AccountStatement) {
	for i, line := range lines {
		normalizedLine := normalizeForKeywordMatch(line)

		if statement.IBAN == "" {
			if iban := ibanPattern.FindString(line); iban != "" {
				statement.IBAN = strings.ReplaceAll(iban, " ", "")
			}
		}

		if statement.AccountHolder == "" && containsAnyKeyword(normalizedLine, metadataAccountHolderKeywords) {
			statement.AccountHolder = extractAccountHolderValue(lines, i)
		}

		if statement.AccountNumber == "" && containsAnyKeyword(normalizedLine, metadataAccountNumberKeywords) {
			statement.AccountNumber = valueAfterColon(line)
		}

		if statement.StartDate.IsZero() || statement.EndDate.IsZero() {
			matches := dateRangePattern.FindStringSubmatch(line)
			if len(matches) == 3 {
				startDate, startErr := parseTurkishDate(matches[1])
				endDate, endErr := parseTurkishDate(matches[2])
				if startErr == nil && endErr == nil {
					statement.StartDate = startDate
					statement.EndDate = endDate
				}
			}
		}
	}
}

func extractAccountHolderValue(lines []string, currentIndex int) string {
	current := strings.TrimSpace(valueAfterColon(lines[currentIndex]))
	if isLikelyAccountHolderValue(current) {
		return current
	}

	for i := currentIndex + 1; i < len(lines) && i <= currentIndex+3; i++ {
		candidate := strings.TrimSpace(lines[i])
		if isLikelyAccountHolderValue(candidate) {
			return candidate
		}
	}

	return ""
}

func isLikelyAccountHolderValue(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" || value == ":" {
		return false
	}

	normalized := normalizeForKeywordMatch(value)
	if containsAnyKeyword(normalized, metadataAccountHolderKeywords) {
		return false
	}

	if strings.Contains(normalized, "TC KIMLIK") {
		return false
	}

	if containsAnyKeyword(normalized, continuationSkipKeywords) {
		return false
	}

	// Avoid table/header tokens that can appear right after label-only lines.
	rejectedHeaders := []string{"TARIH", "HAREKET TIPI", "TUTAR", "BAKIYE", "SERI/SIRA NO"}
	if containsAnyKeyword(normalized, rejectedHeaders) {
		return false
	}

	return strings.Count(value, " ") >= 1
}

func collectTransactionRows(lines []string) []string {
	rows := make([]string, 0)
	var current strings.Builder

	flushCurrent := func() {
		if current.Len() == 0 {
			return
		}
		rows = append(rows, strings.TrimSpace(current.String()))
		current.Reset()
	}

	for _, line := range lines {
		if dateLinePattern.MatchString(line) {
			flushCurrent()
			current.WriteString(line)
			continue
		}

		if current.Len() == 0 {
			continue
		}
		if isSkippableContinuation(line) {
			continue
		}

		current.WriteByte(' ')
		current.WriteString(line)
	}

	flushCurrent()
	return rows
}

func isSkippableContinuation(line string) bool {
	return containsAnyKeyword(normalizeForKeywordMatch(line), continuationSkipKeywords)
}

func parseTransactionRow(row string) (models.Transaction, error) {
	matches := trAmountPattern.FindAllStringSubmatchIndex(row, -1)
	if len(matches) < 2 {
		return models.Transaction{}, fmt.Errorf("missing amount/balance fields")
	}

	amountBounds := matches[len(matches)-2]
	balanceBounds := matches[len(matches)-1]

	amountRaw := row[amountBounds[2]:amountBounds[3]]
	balanceRaw := row[balanceBounds[2]:balanceBounds[3]]
	prefix := strings.TrimSpace(row[:amountBounds[0]])

	parts := strings.Fields(prefix)
	if len(parts) < 2 {
		return models.Transaction{}, fmt.Errorf("missing date/type section")
	}

	dateStr := parts[0]
	date, err := parseTurkishDate(dateStr)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("parse date %q: %w", dateStr, err)
	}

	typeAndDescription := strings.TrimSpace(strings.TrimPrefix(prefix, dateStr))
	transactionType, description := splitTypeAndDescription(typeAndDescription)

	amount, err := parseTurkishAmount(amountRaw)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("parse amount %q: %w", amountRaw, err)
	}

	balance, err := parseTurkishAmount(balanceRaw)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("parse balance %q: %w", balanceRaw, err)
	}

	amount = normalizeAmountSign(amount, transactionType)
	merchant := extractMerchant(description)

	return models.Transaction{
		Date:        date,
		Type:        transactionType,
		Description: description,
		Merchant:    merchant,
		Amount:      amount,
		Balance:     balance,
		RawText:     row,
	}, nil
}

func splitTypeAndDescription(value string) (string, string) {
	for _, transactionType := range knownTransactionTypes {
		if strings.HasPrefix(value, transactionType+" ") {
			return transactionType, strings.TrimSpace(strings.TrimPrefix(value, transactionType))
		}
		if value == transactionType {
			return transactionType, ""
		}
	}

	upperValue := strings.ToUpper(value)
	for _, transactionType := range knownTransactionTypes {
		if strings.HasPrefix(upperValue, strings.ToUpper(transactionType)+" ") {
			return transactionType, strings.TrimSpace(value[len(transactionType):])
		}
	}

	return transactionTypeOther, value
}

func parseTurkishDate(dateStr string) (time.Time, error) {
	return time.Parse(turkishDateLayout, dateStr)
}

func parseTurkishAmount(raw string) (float64, error) {
	clean := strings.TrimSpace(raw)
	clean = strings.TrimSuffix(clean, currencySuffixTL)
	clean = strings.TrimSpace(clean)
	clean = strings.ReplaceAll(clean, ".", "")
	clean = strings.ReplaceAll(clean, ",", ".")

	return strconv.ParseFloat(clean, 64)
}

func normalizeAmountSign(amount float64, transactionType string) float64 {
	if amount < 0 {
		return amount
	}

	switch transactionType {
	case transactionTypeIncomingTransfer, transactionTypeRefund:
		return amount
	case transactionTypeOutgoingTransfer, transactionTypeCardSpend, transactionTypeCashWithdraw, transactionTypeFee:
		return -amount
	default:
		return amount
	}
}

func extractMerchant(description string) string {
	merchant := strings.TrimSpace(description)
	merchant = numericPrefixPattern.ReplaceAllString(merchant, "")
	merchant = compactPrefixPattern.ReplaceAllString(merchant, "")
	merchant = merchantCityTRSuffix.ReplaceAllString(merchant, "")

	merchant = strings.TrimSpace(strings.Trim(merchant, "-"))
	if merchant == "" {
		return strings.TrimSpace(description)
	}

	return merchant
}

func deriveDateRange(transactions []models.Transaction) (time.Time, time.Time) {
	if len(transactions) == 0 {
		return time.Time{}, time.Time{}
	}

	startDate := transactions[0].Date
	endDate := transactions[0].Date
	for _, transaction := range transactions[1:] {
		if transaction.Date.Before(startDate) {
			startDate = transaction.Date
		}
		if transaction.Date.After(endDate) {
			endDate = transaction.Date
		}
	}

	return startDate, endDate
}

func valueAfterColon(line string) string {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[1])
	}
	return strings.TrimSpace(line)
}

func normalizeForKeywordMatch(value string) string {
	normalized := strings.ToUpper(value)
	normalized = turkishToASCIIReplacer.Replace(normalized)
	normalized = extraSpacePattern.ReplaceAllString(normalized, " ")
	return strings.TrimSpace(normalized)
}

func containsAnyKeyword(normalizedValue string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(normalizedValue, keyword) {
			return true
		}
	}

	return false
}

func assignDailySequence(transactions []models.Transaction) {
	daySequence := make(map[string]int, len(transactions))
	for i := range transactions {
		dateKey := transactions[i].Date.Format("2006-01-02")
		daySequence[dateKey]++
		transactions[i].DailySequence = daySequence[dateKey]
	}
}
