---
sidebar_position: 3
---

# CLI Usage

Use enpara-cli when you want fast local conversion from terminal.

## Table of Contents

- [Command Syntax](#command-syntax)
- [Flags](#flags)
- [Practical Examples](#practical-examples)
- [PDF Type Selection](#pdf-type-selection)
- [Next Steps](#next-steps)

## Command Syntax

```bash
enpara-cli <input.pdf> [flags]
```

Example:

```bash
./bin/enpara-cli ./tmp/manual.pdf --type type1 --format json --output ./tmp/manual/statement.json
```

## Flags

| Flag | Short | Default | Description |
| --- | --- | --- | --- |
| --format | -f | csv | Output format: csv, json, xlsx, ofx |
| --output | -o | auto name | Output file path |
| --type | -t | auto | PDF parser type: auto, type1, type2 |

:::tip
If you skip --output, the tool writes input-base.format in your current folder.
:::

## Practical Examples

### Default conversion (CSV)

```bash
./bin/enpara-cli ./tmp/manual.pdf --type type1 --format csv --output ./tmp/manual/statement.csv
```

### Convert to JSON

```bash
./bin/enpara-cli ./tmp/automatic.pdf --type type2 --format json --output ./tmp/auto/statement.json
```

### Convert to XLSX

```bash
./bin/enpara-cli ./tmp/automatic.pdf --type type2 --format xlsx --output ./tmp/auto/statement.xlsx
```

### Convert to OFX for accounting import

```bash
./bin/enpara-cli ./tmp/automatic.pdf --format ofx --output ./tmp/auto/statement.ofx
```

### Choose a custom output filename

```bash
./bin/enpara-cli ./tmp/manual.pdf --type type1 --format csv --output ./tmp/manual/my-transactions.csv
```

### Parse as type2 explicitly

```bash
./bin/enpara-cli ./tmp/automatic.pdf --type type2 --format json --output ./tmp/auto/type2.json
```

### Verified output preview

```text
Tarih;Hareket tipi;Açıklama;NFC;İşlem Tutarı;Bakiye
28.03.2026;Diğer;000000004228140-DED COFFEE IZMIR TR;0;-215,00;2.142,27
```

## PDF Type Selection

- auto: recommended first choice. Detects type1 or type2 automatically.
- type1: manual statement layout parser.
- type2: monthly automatic statement layout parser.

If conversion looks wrong, rerun with an explicit type.

:::warning
Wrong parser type can cause missing rows or malformed values. If auto fails, force type1 or type2 and compare outputs.
:::

## Next Steps

- Learn request-based usage in [API Usage](./api-usage.md)
- Compare formats in [Output Formats](./output-formats.md)
