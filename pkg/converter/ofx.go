package converter

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/chillmatin/enpara-transactions-parser/pkg/models"
)

func ToOFX(statement *models.AccountStatement) ([]byte, error) {
	now := time.Now().UTC()
	startDate := statement.StartDate
	endDate := statement.EndDate

	if startDate.IsZero() {
		startDate = now
		if len(statement.Transactions) > 0 {
			startDate = statement.Transactions[0].Date
		}
	}
	if endDate.IsZero() {
		endDate = startDate
		if len(statement.Transactions) > 0 {
			endDate = statement.Transactions[0].Date
			for _, tx := range statement.Transactions[1:] {
				if tx.Date.After(endDate) {
					endDate = tx.Date
				}
			}
		}
	}

	transactions := make([]ofxStmtTrn, 0, len(statement.Transactions))
	fallbackDaySequence := make(map[string]int, len(statement.Transactions))
	for _, tx := range statement.Transactions {
		name := tx.Type
		if name == "" {
			name = tx.Description
		}

		daySequence := tx.DailySequence
		if daySequence <= 0 {
			dateKey := tx.Date.Format("2006-01-02")
			fallbackDaySequence[dateKey]++
			daySequence = fallbackDaySequence[dateKey]
		}

		transactions = append(transactions, ofxStmtTrn{
			TRNTYPE:  toOFXTrnType(tx.Amount),
			DTPOSTED: toOFXDate(tx.Date),
			TRNAMT:   fmt.Sprintf("%.2f", tx.Amount),
			FITID:    buildFITID(tx, daySequence),
			NAME:     name,
			MEMO:     tx.Description,
		})
	}

	balance := 0.0
	if len(statement.Transactions) > 0 {
		balance = statement.Transactions[len(statement.Transactions)-1].Balance
	}

	acctID := statement.AccountNumber
	if acctID == "" {
		acctID = statement.IBAN
	}
	if acctID == "" {
		acctID = "UNKNOWN"
	}

	doc := ofxDocument{
		XMLNS: "http://ofx.net/types/2003/04",
		SIGNONMSGSRSV1: ofxSignOnMsgs{
			SONRS: ofxSONRS{
				STATUS: ofxStatus{
					CODE:     "0",
					SEVERITY: "INFO",
				},
				DTSERVER: toOFXDate(now),
				LANGUAGE: "TUR",
			},
		},
		BANKMSGSRSV1: ofxBankMsgs{
			STMTTRNRS: ofxStmtTrnRs{
				TRNUID: "0",
				STATUS: ofxStatus{
					CODE:     "0",
					SEVERITY: "INFO",
				},
				STMTRS: ofxStmtRs{
					CURDEF: "TRY",
					BANKACCTFROM: ofxBankAcctFrom{
						BANKID:   "ENPARA",
						ACCTID:   acctID,
						ACCTTYPE: "CHECKING",
					},
					BANKTRANLIST: ofxBankTranList{
						DTSTART: toOFXDate(startDate),
						DTEND:   toOFXDate(endDate),
						STMTTRN: transactions,
					},
					LEDGERBAL: ofxLedgerBal{
						BALAMT: fmt.Sprintf("%.2f", balance),
						DTASOF: toOFXDate(endDate),
					},
				},
			},
		},
	}

	content, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal ofx xml: %w", err)
	}

	return append([]byte(xml.Header), content...), nil
}

func toOFXDate(value time.Time) string {
	if value.IsZero() {
		value = time.Now().UTC()
	}
	return value.UTC().Format("20060102150405")
}

func toOFXTrnType(amount float64) string {
	if amount < 0 {
		return "DEBIT"
	}
	return "CREDIT"
}

func buildFITID(tx models.Transaction, daySequence int) string {
	raw := fmt.Sprintf("%s|%d|%.2f|%.2f|%s", tx.Date.Format("2006-01-02"), daySequence, tx.Amount, tx.Balance, tx.RawText)
	hash := sha1.Sum([]byte(raw))
	return hex.EncodeToString(hash[:])[:24]
}

type ofxDocument struct {
	XMLName        xml.Name      `xml:"OFX"`
	XMLNS          string        `xml:"xmlns,attr"`
	SIGNONMSGSRSV1 ofxSignOnMsgs `xml:"SIGNONMSGSRSV1"`
	BANKMSGSRSV1   ofxBankMsgs   `xml:"BANKMSGSRSV1"`
}

type ofxSignOnMsgs struct {
	SONRS ofxSONRS `xml:"SONRS"`
}

type ofxSONRS struct {
	STATUS   ofxStatus `xml:"STATUS"`
	DTSERVER string    `xml:"DTSERVER"`
	LANGUAGE string    `xml:"LANGUAGE"`
}

type ofxBankMsgs struct {
	STMTTRNRS ofxStmtTrnRs `xml:"STMTTRNRS"`
}

type ofxStmtTrnRs struct {
	TRNUID string    `xml:"TRNUID"`
	STATUS ofxStatus `xml:"STATUS"`
	STMTRS ofxStmtRs `xml:"STMTRS"`
}

type ofxStatus struct {
	CODE     string `xml:"CODE"`
	SEVERITY string `xml:"SEVERITY"`
}

type ofxStmtRs struct {
	CURDEF       string          `xml:"CURDEF"`
	BANKACCTFROM ofxBankAcctFrom `xml:"BANKACCTFROM"`
	BANKTRANLIST ofxBankTranList `xml:"BANKTRANLIST"`
	LEDGERBAL    ofxLedgerBal    `xml:"LEDGERBAL"`
}

type ofxBankAcctFrom struct {
	BANKID   string `xml:"BANKID"`
	ACCTID   string `xml:"ACCTID"`
	ACCTTYPE string `xml:"ACCTTYPE"`
}

type ofxBankTranList struct {
	DTSTART string       `xml:"DTSTART"`
	DTEND   string       `xml:"DTEND"`
	STMTTRN []ofxStmtTrn `xml:"STMTTRN"`
}

type ofxStmtTrn struct {
	TRNTYPE  string `xml:"TRNTYPE"`
	DTPOSTED string `xml:"DTPOSTED"`
	TRNAMT   string `xml:"TRNAMT"`
	FITID    string `xml:"FITID"`
	NAME     string `xml:"NAME,omitempty"`
	MEMO     string `xml:"MEMO,omitempty"`
}

type ofxLedgerBal struct {
	BALAMT string `xml:"BALAMT"`
	DTASOF string `xml:"DTASOF"`
}
