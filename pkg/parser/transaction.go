package parser

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/chillmatin/enpara-transactions-parser/pkg/models"
)

const (
	turkishDateLayout = "02.01.2006"
	type2DateLayout   = "02/01/06"
	currencySuffixTL  = "TL"

	PDFTypeAuto = "auto"
	PDFType1    = "type1"
	PDFType2    = "type2"

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
	type2DateLinePattern = regexp.MustCompile(`^\d{2}/\d{2}/\d{2}\b`)
	type2AmountPattern   = regexp.MustCompile(`([+-]?\s?\d{1,3}(?:\.\d{3})*,\d{2})\s*` + currencySuffixTL)
	ibanPattern          = regexp.MustCompile(`\bTR\d{2}(?:\s?\d{4}){5}\s?\d{2}\b`)
	dateRangePattern     = regexp.MustCompile(`(\d{2}\.\d{2}\.\d{4}).*?(\d{2}\.\d{2}\.\d{4})`)
	numericPrefixPattern = regexp.MustCompile(`^\d+\s*-\s*`)
	compactPrefixPattern = regexp.MustCompile(`^\d+-`)
	extraSpacePattern    = regexp.MustCompile(`\s+`)
	merchantCityTRSuffix = regexp.MustCompile(`\s+[\p{L}]+\s+TR$`)
	type2NFCPattern      = regexp.MustCompile(`(?i)\b(nfc|temassiz|temasiz|contactless)\b`)
	xmlTagPattern        = regexp.MustCompile(`<[^>]+>`)

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

type ParseOptions struct {
	PDFType       string
	Type2NFCHints []bool
}

func ParseStatement(pdfText string) (*models.AccountStatement, error) {
	return ParseStatementWithOptions(pdfText, ParseOptions{PDFType: PDFTypeAuto})
}

func ParseStatementFromPDF(pdfPath string, opts ParseOptions) (*models.AccountStatement, error) {
	pdfText, err := ExtractTextFromPDF(pdfPath)
	if err != nil {
		return nil, fmt.Errorf("extract text from pdf: %w", err)
	}

	lines := normalizeLines(pdfText)
	pdfType, err := resolvePDFType(lines, opts)
	if err != nil {
		return nil, err
	}

	if pdfType == PDFType2 && len(opts.Type2NFCHints) == 0 {
		if hints, hintErr := extractType2NFCHintsFromPDF(pdfPath); hintErr == nil {
			opts.Type2NFCHints = hints
		}
	}

	return parseStatementByType(lines, pdfType, opts)
}

func ParseStatementWithOptions(pdfText string, opts ParseOptions) (*models.AccountStatement, error) {
	lines := normalizeLines(pdfText)
	pdfType, err := resolvePDFType(lines, opts)
	if err != nil {
		return nil, err
	}

	return parseStatementByType(lines, pdfType, opts)
}

func parseStatementByType(lines []string, pdfType string, opts ParseOptions) (*models.AccountStatement, error) {

	switch pdfType {
	case PDFType1:
		return parseStatementType1(lines)
	case PDFType2:
		return parseStatementType2(lines, opts)
	default:
		return nil, fmt.Errorf("unsupported pdf type %q", pdfType)
	}
}

func resolvePDFType(lines []string, opts ParseOptions) (string, error) {
	pdfType := strings.ToLower(strings.TrimSpace(opts.PDFType))
	if pdfType == "" {
		pdfType = PDFTypeAuto
	}

	switch pdfType {
	case PDFType1, PDFType2:
		return pdfType, nil
	case PDFTypeAuto:
		if detectType2(lines) {
			return PDFType2, nil
		}
		return PDFType1, nil
	default:
		return "", fmt.Errorf("unsupported pdf type %q (supported: auto, type1, type2)", pdfType)
	}
}

func detectType2(lines []string) bool {
	hasColumns := false
	hasSlashDate := false
	hasDetailSection := false

	for _, line := range lines {
		normalized := normalizeForKeywordMatch(line)
		if strings.Contains(normalized, "HAREKETLERININ DETAYI") {
			hasDetailSection = true
		}
		if normalized == "TARIH" || normalized == "ACIKLAMA" || normalized == "TUTAR" || normalized == "BAKIYE" {
			hasColumns = true
		}
		if type2DateLinePattern.MatchString(line) {
			hasSlashDate = true
		}
	}

	return hasDetailSection && hasColumns && hasSlashDate
}

func parseStatementType1(lines []string) (*models.AccountStatement, error) {
	statement := &models.AccountStatement{}

	extractMetadata(lines, statement)
	rows := collectTransactionRowsType1(lines)
	if len(rows) == 0 {
		return nil, fmt.Errorf("no transaction rows detected")
	}

	transactions := make([]models.Transaction, 0, len(rows))
	errorCount := 0

	for _, row := range rows {
		tx, err := parseTransactionRowType1(row)
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

func collectTransactionRowsType1(lines []string) []string {
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

func parseTransactionRowType1(row string) (models.Transaction, error) {
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
	transactionType, description := splitType1TypeAndDescription(typeAndDescription)

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

func parseStatementType2(lines []string, opts ParseOptions) (*models.AccountStatement, error) {
	statement := &models.AccountStatement{}
	extractMetadata(lines, statement)

	rows := collectTransactionRowsType2(lines)
	if len(rows) == 0 {
		return nil, fmt.Errorf("no transaction rows detected")
	}

	transactions := make([]models.Transaction, 0, len(rows))
	errorCount := 0
	for i, row := range rows {
		rowNFCHint := false
		if i < len(opts.Type2NFCHints) {
			rowNFCHint = opts.Type2NFCHints[i]
		}

		tx, err := parseTransactionRowType2(row, rowNFCHint)
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

func collectTransactionRowsType2(lines []string) []string {
	rows := make([]string, 0)
	var current strings.Builder

	inDetailSection := false
	headerSeen := false

	flushCurrent := func() {
		if current.Len() == 0 {
			return
		}
		rows = append(rows, strings.TrimSpace(current.String()))
		current.Reset()
	}

	for i, line := range lines {
		normalized := normalizeForKeywordMatch(line)
		if strings.Contains(normalized, "HAREKETLERININ DETAYI") {
			inDetailSection = true
			headerSeen = false
			flushCurrent()
			continue
		}

		if !inDetailSection {
			continue
		}

		if !headerSeen && isType2HeaderBlock(lines, i) {
			headerSeen = true
			continue
		}

		if !headerSeen {
			continue
		}

		if type2DateLinePattern.MatchString(line) {
			flushCurrent()
			current.WriteString(line)
			continue
		}

		if current.Len() == 0 {
			continue
		}

		if isSkippableType2Continuation(line) {
			continue
		}

		current.WriteByte(' ')
		current.WriteString(line)
	}

	flushCurrent()
	return rows
}

func isType2HeaderBlock(lines []string, start int) bool {
	windowEnd := start + 6
	if windowEnd > len(lines)-1 {
		windowEnd = len(lines) - 1
	}

	seen := map[string]bool{}
	for i := start; i <= windowEnd; i++ {
		normalized := normalizeForKeywordMatch(lines[i])
		if normalized == "TARIH" || normalized == "ACIKLAMA" || normalized == "TUTAR" || normalized == "BAKIYE" {
			seen[normalized] = true
		}
	}

	return seen["TARIH"] && seen["ACIKLAMA"] && seen["TUTAR"] && seen["BAKIYE"]
}

func isSkippableType2Continuation(line string) bool {
	normalized := normalizeForKeywordMatch(line)
	if normalized == "" {
		return true
	}
	if strings.Contains(normalized, "SAYFA") {
		return true
	}
	if strings.Contains(normalized, "ENPARA BANK") {
		return true
	}
	if strings.Contains(normalized, "ESENTEPE MAH") {
		return true
	}
	if normalized == "TARIH" || normalized == "ACIKLAMA" || normalized == "TUTAR" || normalized == "BAKIYE" {
		return true
	}
	return false
}

func parseTransactionRowType2(row string, nfcHint bool) (models.Transaction, error) {
	matches := type2AmountPattern.FindAllStringSubmatchIndex(row, -1)
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
		return models.Transaction{}, fmt.Errorf("missing date/description section")
	}

	dateStr := parts[0]
	date, err := parseType2Date(dateStr)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("parse date %q: %w", dateStr, err)
	}

	typeAndDescription := strings.TrimSpace(strings.TrimPrefix(prefix, dateStr))
	transactionType, description := splitType2TypeAndDescription(typeAndDescription)

	amount, err := parseTurkishAmount(amountRaw)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("parse amount %q: %w", amountRaw, err)
	}
	balance, err := parseTurkishAmount(balanceRaw)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("parse balance %q: %w", balanceRaw, err)
	}

	nfc := detectNFCFromRow(typeAndDescription)
	if !nfc && nfcHint {
		nfc = true
	}

	return models.Transaction{
		Date:        date,
		Type:        transactionType,
		Description: description,
		Merchant:    extractMerchant(description),
		NFC:         nfc,
		Amount:      amount,
		Balance:     balance,
		RawText:     row,
	}, nil
}

type pdf2XMLDocument struct {
	Pages []pdf2XMLPage `xml:"page"`
}

type pdf2XMLPage struct {
	Number int            `xml:"number,attr"`
	Images []pdf2XMLImage `xml:"image"`
	Texts  []pdf2XMLText  `xml:"text"`
}

type pdf2XMLImage struct {
	Top    int    `xml:"top,attr"`
	Left   int    `xml:"left,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
	Src    string `xml:"src,attr"`
}

type pdf2XMLText struct {
	Top   int    `xml:"top,attr"`
	Left  int    `xml:"left,attr"`
	Value string `xml:",innerxml"`
}

func (t pdf2XMLText) plainText() string {
	value := xmlTagPattern.ReplaceAllString(t.Value, "")
	return strings.TrimSpace(value)
}

func extractType2NFCHintsFromPDF(pdfPath string) ([]bool, error) {
	if _, err := exec.LookPath("pdftohtml"); err != nil {
		return nil, fmt.Errorf("pdftohtml not found: %w", err)
	}

	tempDir, err := os.MkdirTemp("", "enpara-pdftohtml-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	xmlOutputPath := filepath.Join(tempDir, "statement.xml")
	cmd := exec.Command("pdftohtml", "-xml", "-hidden", "-nodrm", pdfPath, xmlOutputPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("run pdftohtml: %w (%s)", err, strings.TrimSpace(string(output)))
	}

	xmlBytes, err := os.ReadFile(xmlOutputPath)
	if err != nil {
		return nil, fmt.Errorf("read pdftohtml xml: %w", err)
	}

	var doc pdf2XMLDocument
	if err := xml.Unmarshal(xmlBytes, &doc); err != nil {
		return nil, fmt.Errorf("decode pdftohtml xml: %w", err)
	}

	hints := make([]bool, 0)
	for _, page := range doc.Pages {
		detailHeaderTop := findType2DetailHeaderTop(page.Texts)
		if detailHeaderTop == 0 {
			continue
		}

		amountTops := make([]int, 0)
		for _, text := range page.Texts {
			value := text.plainText()
			if text.Top <= detailHeaderTop {
				continue
			}
			if text.Left < 590 || text.Left > 690 {
				continue
			}
			if type2AmountPattern.MatchString(value) {
				amountTops = append(amountTops, text.Top)
			}
		}

		iconTops := make([]int, 0)
		for _, image := range page.Images {
			if image.Left < 120 || image.Left > 180 {
				continue
			}
			if image.Width < 14 || image.Width > 26 || image.Height < 18 || image.Height > 30 {
				continue
			}
			iconTops = append(iconTops, image.Top)
		}

		sort.Ints(amountTops)
		sort.Ints(iconTops)

		for _, amountTop := range amountTops {
			hasIcon := false
			for _, iconTop := range iconTops {
				delta := amountTop - iconTop
				if delta >= 10 && delta <= 30 {
					hasIcon = true
					break
				}
			}
			hints = append(hints, hasIcon)
		}
	}

	return hints, nil
}

func findType2DetailHeaderTop(texts []pdf2XMLText) int {
	headerTop := 0
	for _, text := range texts {
		if normalizeForKeywordMatch(text.plainText()) != "ACIKLAMA" {
			continue
		}
		if text.Left < 140 || text.Left > 320 {
			continue
		}

		if headerTop == 0 || text.Top < headerTop {
			headerTop = text.Top
		}
	}

	return headerTop
}

func splitType2TypeAndDescription(value string) (string, string) {
	parts := strings.SplitN(value, ",", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}

	return transactionTypeOther, strings.TrimSpace(value)
}

func detectNFCFromRow(value string) bool {
	return type2NFCPattern.MatchString(value)
}

func splitType1TypeAndDescription(value string) (string, string) {
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

func parseType2Date(dateStr string) (time.Time, error) {
	return time.Parse(type2DateLayout, dateStr)
}

func parseTurkishAmount(raw string) (float64, error) {
	clean := strings.TrimSpace(raw)
	clean = strings.TrimSuffix(clean, currencySuffixTL)
	clean = strings.TrimSpace(clean)
	clean = strings.ReplaceAll(clean, " ", "")
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
