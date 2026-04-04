---
sidebar_position: 6
---

import type1Example from './type1-steps/10.png';
import type2Example from './type2-steps/2.png';

# PDF Types Explained

This tool supports two ways Enpara can provide transaction data in PDF format.

## Table of Contents

- [Why PDF Types Matter](#why-pdf-types-matter)
- [Type1](#type1)
- [Type2](#type2)
- [Auto Detection](#auto-detection)
- [How to Pick the Right Type](#how-to-pick-the-right-type)
- [Next Steps](#next-steps)

## Why PDF Types Matter

PDF parsing depends on layout patterns. If layout and parser type do not match, output quality can drop.

## Type1

Type1 is the manual-style statement layout.

Typical signs:

- Date appears like 28.03.2026
- Rows are often sentence-like and may wrap lines
- Transaction text may include merchant and city in one description block

### Screenshot Placeholder

<img className="screenshot-single" src={type1Example} alt="Type1 PDF example with dd.mm.yyyy dates and dense transaction rows" loading="lazy" />

### Where to add your screenshot notes

- Manual statement layout with sentence-like rows.
- Date format appears as dd.mm.yyyy.

## Type2

Type2 is the automatic monthly statement layout.

Typical signs:

- Columns are clearly separated (Date, Description, Amount, Balance)
- Date often appears like 03/03/26
- Usually easier to parse row-by-row

### Screenshot

<img className="screenshot-single" src={type2Example} alt="Type2 PDF example with clear columns and separate date format" loading="lazy" />

## Auto Detection

Default type is auto.

In auto mode, the parser inspects the PDF text and chooses type1 or type2.

:::tip
Start with auto first. Force type1 or type2 only if output looks wrong.
:::

## How to Pick the Right Type

1. Run once with auto.
2. If data looks incomplete, rerun with --type type1.
3. If still wrong, rerun with --type type2.
4. Compare row counts and values.

## Next Steps

- See practical rerun examples in [CLI Usage](./cli-usage.md)
- Check repair steps in [Troubleshooting](./troubleshooting.md)
