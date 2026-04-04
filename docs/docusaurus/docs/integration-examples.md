---
sidebar_position: 8
---

# Integration Examples

This page shows practical ways to use the converter in real workflows.

## Table of Contents

- [Web App Upload Flow](#web-app-upload-flow)
- [Accounting Workflow](#accounting-workflow)
- [Batch Automation Script](#batch-automation-script)
- [Backend Service Pattern](#backend-service-pattern)
- [Next Steps](#next-steps)

## Web App Upload Flow

Use enpara-api behind your frontend.

High-level flow:

1. User uploads PDF in browser.
2. Frontend sends multipart request to your backend.
3. Backend forwards to enpara-api /api/v1/convert.
4. Converted file is returned to user or stored.

Example Node backend route:

```javascript
app.post("/convert", upload.single("file"), async (req, res) => {
  const form = new FormData();
  form.append("file", new Blob([req.file.buffer]), req.file.originalname);
  form.append("format", req.body.format || "json");
  form.append("type", req.body.type || "auto");

  const apiRes = await fetch("http://localhost:8080/api/v1/convert", {
    method: "POST",
    body: form,
  });

  if (!apiRes.ok) {
    return res.status(422).json({error: "Conversion failed"});
  }

  const fileBuffer = Buffer.from(await apiRes.arrayBuffer());
  res.setHeader("Content-Type", apiRes.headers.get("content-type") || "application/octet-stream");
  res.send(fileBuffer);
});
```

## Accounting Workflow

Common simple flow:

1. Convert monthly PDF to OFX.
2. Import OFX into accounting tool.
3. Reconcile transactions.

```bash
./enpara-cli "March-2026.pdf" --format ofx --output "March-2026.ofx"
```

## Batch Automation Script

If you export statements monthly, automate conversion in a folder.

```bash
#!/usr/bin/env bash
set -euo pipefail

INPUT_DIR="./statements"
OUTPUT_DIR="./converted"
mkdir -p "$OUTPUT_DIR"

for pdf in "$INPUT_DIR"/*.pdf; do
  name="$(basename "$pdf" .pdf)"
  ./enpara-cli "$pdf" --format csv --type auto --output "$OUTPUT_DIR/$name.csv"
done
```

## Backend Service Pattern

For production integration:

- Keep conversion API internal (not public internet if possible)
- Add request size limits
- Add retries for transient failures
- Log conversion errors with request IDs

:::tip
For long-term integrations, use JSON output as your canonical data exchange format.
:::

## Next Steps

- Endpoint details: [API Usage](./api-usage.md)
- Data format decisions: [Output Formats](./output-formats.md)
