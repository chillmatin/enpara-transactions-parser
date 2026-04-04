---
sidebar_position: 4
---

# API Usage

Use enpara-api when you want your app or service to upload a PDF and receive converted output.

## Table of Contents

- [Start the API Server](#start-the-api-server)
- [Endpoints](#endpoints)
- [Request Parameters](#request-parameters)
- [Example Requests](#example-requests)
- [Response Behavior](#response-behavior)
- [Next Steps](#next-steps)

## Start the API Server

```bash
./enpara-api --swagger
```

Swagger UI will be available at http://localhost:8080/swagger.

You can also set host, port, and swagger defaults with environment variables.

```bash
ENPARA_API_HOST=0.0.0.0 ENPARA_API_PORT=8080 ENPARA_API_SWAGGER=true ./enpara-api
```

## Endpoints

- POST /api/v1/convert
- GET /api/v1/formats
- GET /api/v1/health

## Request Parameters

POST /api/v1/convert accepts multipart/form-data fields:

- file (required): PDF statement file
- format (optional): json, csv, xlsx, ofx. Default: json
- type (optional): auto, type1, type2. Default: auto

## Example Requests

### curl

```bash
curl -X POST http://localhost:8080/api/v1/convert \
  -F "file=@./tmp/automatic.pdf" \
  -F "format=csv" \
  -F "type=auto" \
  --output ./tmp/auto/statement.csv
```

Verified output first line:

```text
Tarih;Hareket tipi;Açıklama;NFC;İşlem Tutarı;Bakiye
```

### JavaScript (fetch)

```javascript
import fs from "node:fs";

const form = new FormData();
form.append("file", new Blob([fs.readFileSync("./tmp/statement.pdf")]), "statement.pdf");
form.append("format", "json");
form.append("type", "auto");

const res = await fetch("http://localhost:8080/api/v1/convert", {
  method: "POST",
  body: form,
});

if (!res.ok) {
  throw new Error(`Conversion failed: ${res.status}`);
}

const output = await res.text();
fs.writeFileSync("./statement.json", output);
```

### Python (requests)

```python
import requests

url = "http://localhost:8080/api/v1/convert"
files = {"file": open("./tmp/statement.pdf", "rb")}
data = {"format": "xlsx", "type": "auto"}

resp = requests.post(url, files=files, data=data, timeout=60)
resp.raise_for_status()

with open("statement.xlsx", "wb") as f:
    f.write(resp.content)
```

## Response Behavior

- Success: binary content of converted file
- Content-Disposition: attachment filename is generated from original PDF name
- Errors:
  - 400 for invalid request fields or missing file
  - 422 when parsing or conversion fails

:::info
If your filename includes non-English characters, the API sends proper filename headers for better client compatibility.
:::

### Endpoint quick checks

```bash
curl http://localhost:8080/api/v1/formats
curl http://localhost:8080/api/v1/health
```

Expected responses:

```json
{"formats":["json","csv","xlsx","ofx"]}
```

```json
{"status":"ok"}
```

## Next Steps

- Understand format tradeoffs in [Output Formats](./output-formats.md)
- Learn parser choices in [PDF Types](./pdf-types.md)
