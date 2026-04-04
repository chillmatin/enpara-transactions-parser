---
sidebar_position: 5
---

# Output Formats

Pick the format that fits your workflow.

## Table of Contents

- [Quick Comparison](#quick-comparison)
- [Fields Included](#fields-included)
- [Sample Outputs](#sample-outputs)
- [When to Use Which](#when-to-use-which)
- [Next Steps](#next-steps)

## Quick Comparison

| Format | Best For | Human-Readable | App-Friendly | Typical Next Step |
| --- | --- | --- | --- | --- |
| CSV | Excel, Google Sheets, quick filtering | Yes | Medium | Open in spreadsheet |
| JSON | APIs, scripts, backend systems | Yes | High | Parse in code |
| XLSX | Finance teams using Excel features | Yes | Medium | Pivot tables and formulas |
| OFX | Accounting and personal finance software | Limited | High (for finance tools) | Import into bookkeeping app |

## Fields Included

All conversion formats represent these transaction-level fields:

- Date
- Transaction Type
- Description
- NFC (0 or 1)
- Amount
- Balance

Model-level metadata may also be present or used internally (account holder, account number, IBAN, statement date range).

:::tip
For automation, JSON is usually the easiest format to process reliably.
:::

## Sample Outputs

### CSV snippet

```csv
Tarih;Hareket tipi;Açıklama;NFC;İşlem Tutarı;Bakiye
02.03.2026;Diğer;Diğer Banka Debit Karttan ENPARA Debit Karta Para Transferi;0;14.954,09;15.048,78
03.03.2026;Diğer;000000003620060-ODEAL//YASEMIN YILDI MANISA TR;1;-145,00;14.903,78
```

### JSON snippet

```json
[
  {
    "Tarih": "28.03.2026",
    "Hareket tipi": "Diğer",
    "Açıklama": "000000004228140-DED COFFEE IZMIR TR",
    "NFC": 0,
    "İşlem Tutarı": -215,
    "Bakiye": 2142.27
  }
]
```

### XLSX

XLSX contains the same columns as CSV, with formatting suitable for Excel.

### OFX

OFX exports transactions in a banking interchange format with fields such as TRNTYPE, DTPOSTED, TRNAMT, FITID, NAME, and MEMO.

```xml
<STMTTRN>
  <TRNTYPE>CREDIT</TRNTYPE>
  <DTPOSTED>20260302000000</DTPOSTED>
  <TRNAMT>14954.09</TRNAMT>
  <FITID>8f6f64b78db7e905e9fa806d</FITID>
  <NAME>Diğer</NAME>
  <MEMO>Diğer Banka Debit Karttan ENPARA Debit Karta Para Transferi</MEMO>
</STMTTRN>
```

## When to Use Which

- Choose CSV for quick manual checks.
- Choose JSON for service-to-service integration.
- Choose XLSX when you need rich spreadsheet formatting.
- Choose OFX when your accounting tool supports direct OFX import.

## Next Steps

- Learn parser layout matching in [PDF Types](./pdf-types.md)
- Fix common format issues in [Troubleshooting](./troubleshooting.md)
