---
sidebar_position: 2
---

import step1 from './type1-steps/1.PNG';
import step2 from './type1-steps/2.PNG';
import step3 from './type1-steps/3.PNG';
import step4 from './type1-steps/4.PNG';
import step5 from './type1-steps/5.PNG';
import step6 from './type1-steps/6.PNG';
import step7 from './type1-steps/7.PNG';
import step8 from './type1-steps/8.png';
import step9 from './type1-steps/9.png';
import step10 from './type1-steps/10.png';
import type2step1 from './type2-steps/1.png';
import type2step2 from './type2-steps/2.png';

# Getting Started

This guide gets you from PDF to converted file in under 2 minutes.

## Table of Contents

- [Install](#install)
- [Get PDF from Enpara (Screenshots)](#get-pdf-from-enpara-screenshots)
- [Converting Your First Statement](#converting-your-first-statement)
- [File Paths on Windows, macOS, Linux](#file-paths-on-windows-macos-linux)
- [Before and After](#before-and-after)
- [Next Steps](#next-steps)

## Install

1. Download a release binary from GitHub Releases.
2. Or build from source:

```bash
git clone https://github.com/chillmatin/enpara-transactions-parser.git
cd enpara-transactions-parser
make build
cd bin
```

You will get two binaries:

- enpara-cli
- enpara-api

:::info
Use enpara-cli for local command-line conversion. Use enpara-api for programmatic integration.
:::

## Get PDF from Enpara (Screenshots)

Use this Type1 retrieval flow to request and download the statement PDF.

### Type1 steps (in filename order)

<div className="screenshot-grid">
	<figure className="screenshot-card">
		<img src={step1} alt="Step 1 - Main Menu" loading="lazy" />
		<figcaption>1 - Open menu by clicking top left icon.</figcaption>
	</figure>
	<figure className="screenshot-card">
		<img src={step2} alt="Step 2 - Menu with Belgelerim expanded showing Belge talebi" loading="lazy" />
		<figcaption>2 - Open Belgelerim.</figcaption>
	</figure>
	<figure className="screenshot-card">
		<img src={step3} alt="Step 3 - Menu with Belgelerim expanded showing Belge talebi" loading="lazy" />
		<figcaption>3 - Select :Belge talebi".</figcaption>
	</figure>
	<figure className="screenshot-card">
		<img src={step4} alt="Step 4 - Belge tipi screen with Yeni Belge Talep Et option" loading="lazy" />
		<figcaption>4 - Select "Yeni Belge Talep Et".</figcaption>
	</figure>
	<figure className="screenshot-card">
		<img src={step5} alt="Step 5 - Belge Talebi Menu" loading="lazy" />
		<figcaption>5 - Select "Hesap hareketleri".</figcaption>
	</figure>
	<figure className="screenshot-card">
		<img src={step6} alt="Step 6 - Belge Talebi account selection" loading="lazy" />
		<figcaption>6 - First define period.Then, choose your account and select "Onayla".</figcaption>
	</figure>
	<figure className="screenshot-card">
		<img src={step7} alt="Step 7 - Mail sent confirmation menu" loading="lazy" />
		<figcaption>7 - Make sure you get the mail.</figcaption>
	</figure>
	<figure className="screenshot-card">
		<img src={step8} alt="Step 8 - Enpara mail message" loading="lazy" />
		<figcaption>8 - Open the mail.</figcaption>
	</figure>
	<figure className="screenshot-card">
		<img src={step9} alt="Step 9 - Download action for Enpara Hesap Hareketleri PDF" loading="lazy" />
		<figcaption>9 - Download the PDF.</figcaption>
	</figure>
	<figure className="screenshot-card">
		<img src={step10} alt="Step 10 - Opened Type1 PDF with transaction rows" loading="lazy" />
		<figcaption>10 - Opened Type1 PDF example.</figcaption>
	</figure>
</div>

:::tip
If you receive more than one PDF, pick the one named similar to Enpara Hesap Hareketleri.pdf.
:::

### Type2 steps (in filename order)

<div className="screenshot-grid">
	<figure className="screenshot-card">
		<img src={type2step1} alt="Step 1 - Type2 PDF retrieval" loading="lazy" />
		<figcaption>1 - Browser-based Type2 retrieval flow.</figcaption>
	</figure>
	<figure className="screenshot-card">
		<img src={type2step2} alt="Step 2 - Type2 PDF example" loading="lazy" />
		<figcaption>2 - Opened Type2 PDF example.</figcaption>
	</figure>
</div>

## Converting Your First Statement

1. Put your PDF in your current working folder.
2. Run:

```bash
./enpara-cli "1- Enpara Hesap Hareketleri.pdf"
```

3. Open the generated CSV file.

That is it. Your transactions are now in a sortable, importable format.

## File Paths on Windows, macOS, Linux

Use quotes around paths that include spaces.

### Windows (PowerShell)

```powershell
.\enpara-cli.exe "C:\Users\Alice\Downloads\1- Enpara Hesap Hareketleri.pdf" --format csv
```

### macOS

```bash
./enpara-cli "/Users/alice/Downloads/1- Enpara Hesap Hareketleri.pdf" --format csv
```

### Linux

```bash
./enpara-cli "/home/alice/Downloads/1- Enpara Hesap Hareketleri.pdf" --format csv
```

## Before and After

### Before (inside PDF)

You usually see rows like this in a statement page:

```text
28.03.2026 Diger 000000004228140-DED COFFEE IZMIR TR -52,50 TL 2.357,27 TL
29.03.2026 Gelen Transfer Maasi Odemesi 5.000,00 TL 7.357,27 TL
```

### After (CSV output)

```csv
Tarih;Hareket tipi;Aciklama;NFC;Islem Tutari;Bakiye
28.03.2026;Diger;DED COFFEE;0;-52,50;2.357,27
29.03.2026;Gelen Transfer;Maasi Odemesi;0;5.000,00;7.357,27
```

## Next Steps

- Learn all CLI options in [CLI Usage](./cli-usage.md)
- If you want integration, jump to [API Usage](./api-usage.md)
