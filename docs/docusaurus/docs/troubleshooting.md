---
sidebar_position: 7
---

# Troubleshooting

If conversion fails or output looks strange, start here.

## Table of Contents

- [No Output File Created](#no-output-file-created)
- [Unsupported Format or Type](#unsupported-format-or-type)
- [Rows Missing or Values Wrong](#rows-missing-or-values-wrong)
- [API Upload Issues](#api-upload-issues)
- [Path and Permissions Problems](#path-and-permissions-problems)
- [Next Steps](#next-steps)

## No Output File Created

Check that:

- Input PDF path is correct
- You are running in the folder that contains the binary
- The output path is writable

```bash
./enpara-cli "./tmp/statement.pdf" --format csv --output ./tmp/statement.csv
```

:::warning
If your filename has spaces, always wrap it in quotes.
:::

## Unsupported Format or Type

Use only supported values:

- format: csv, json, xlsx, ofx
- type: auto, type1, type2

Example:

```bash
./enpara-cli "statement.pdf" --format json --type auto
```

## Rows Missing or Values Wrong

Most common reason: parser type mismatch.

Try:

```bash
./enpara-cli "statement.pdf" --type type1 --format csv --output type1.csv
./enpara-cli "statement.pdf" --type type2 --format csv --output type2.csv
```

Compare both outputs and keep the better one.

## API Upload Issues

If API returns 400:

- Ensure file field exists and is named file
- Ensure request is multipart/form-data

If API returns 422:

- PDF may be unreadable or unsupported layout
- Retry with explicit type field

```bash
curl -X POST http://localhost:8080/api/v1/convert \
  -F "file=@statement.pdf" \
  -F "format=json" \
  -F "type=type2"
```

## Path and Permissions Problems

### Windows

Run terminal as a user that can access Downloads or Documents.

### macOS/Linux

Use absolute paths to avoid working-directory confusion.

```bash
./enpara-cli "/home/alice/Downloads/statement.pdf" --output "/home/alice/Downloads/statement.csv"
```

:::info
If this still fails, share the exact command and full error text when opening an issue.
:::

## Next Steps

- Review parser layouts in [PDF Types](./pdf-types.md)
- Build robust automation with [Integration Examples](./integration-examples.md)
