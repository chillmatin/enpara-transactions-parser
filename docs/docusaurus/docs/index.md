---
sidebar_position: 1
slug: /
---

# Enpara Bank Statement Converter

Need to turn your Enpara statement PDF into something useful? You are in the right place.

This tool converts Enpara PDF statements into:

- CSV for spreadsheet work
- JSON for apps and scripts
- XLSX for Excel-based workflows
- OFX for accounting and finance software

It supports both quick one-off conversions and full API integration.

:::tip Best starting point
If this is your first time, go to Getting Started and finish your first conversion in under 2 minutes.
:::

## Table of Contents

- [What Is This?](#what-is-this)
- [Who Is This For?](#who-is-this-for)
- [Fast Path](#fast-path)
- [Documentation Map](#documentation-map)
- [Next Steps](#next-steps)

## What Is This?

Enpara Bank Statement Converter reads statement data from a PDF and exports clean, structured transaction rows.

You can filter, sort, import, and automate your bank data instead of manually copying from PDF screens.

## Who Is This For?

- Non-technical users who want to convert a statement quickly
- Developers who want to automate conversion in a backend, script, or web app

## Fast Path

1. Download your Enpara statement PDF.
2. Run:

```bash
./enpara-cli "1- Enpara Hesap Hareketleri.pdf"
```

3. Open the generated CSV file.

## Documentation Map

- [Getting Started](./getting-started.md): install and first success
- [CLI Usage](./cli-usage.md): all flags and command examples
- [API Usage](./api-usage.md): endpoints plus curl, JavaScript, Python
- [Output Formats](./output-formats.md): field mapping and format comparison
- [PDF Types](./pdf-types.md): type1, type2, and auto mode explained
- [Troubleshooting](./troubleshooting.md): common issues and fixes
- [Integration Examples](./integration-examples.md): real workflow patterns

## Next Steps

- Start with [Getting Started](./getting-started.md)
