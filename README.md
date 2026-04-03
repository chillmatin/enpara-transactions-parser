<div align="center">

[![License](https://img.shields.io/badge/license-MIT-blue)](#)
[![Issues - enpara-transactions-parser](https://img.shields.io/github/issues/chillmatin/enpara-transactions-parser)](https://github.com/chillmatin/enpara-transactions-parser/issues)
[![GitHub Release](https://img.shields.io/github/v/release/chillmatin/enpara-transactions-parser)](#)
[![Go Version](https://img.shields.io/github/go-mod/go-version/chillmatin/enpara-transactions-parser)](#)

</div>

&nbsp;

<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://github.com/chillmatin/enpara-transactions-parser/blob/main/docs/assets/logo.png?raw=true">
    <source media="(prefers-color-scheme: light)" srcset="https://github.com/chillmatin/enpara-transactions-parser/blob/main/docs/assets/logo.png?raw=true">
    <img alt="matinhuseynzade.com logo" src="https://github.com/chillmatin/enpara-transactions-parser/blob/main/docs/assets/logo.png?raw=true" width=512>
  </picture>
</div>

<h3 align="center">
  <a>Enpara Transactions Parser</a>
  <br/>
  Convert Enpara PDF account statements to CSV, JSON, XLSX, or OFX.
</h3>

## Quick Start


Download a binary for your OS from the [Releases](https://github.com/chillmatin/enpara-transactions-parser/releases) page, or build from source:
   ```sh
   git clone https://github.com/chillmatin/enpara-transactions-parser.git
   cd enpara-transactions-parser
   make build
   cd bin
   ```
   Binaries are in the `bin/` directory.

After obtaining binaries, there are two ways you can use this tool: `enpara-api` and  `enpara-cli`

### `enpara-api` tool
1. Start the API server:
   ```sh
   ./enpara-api --swagger
   ```
   Visit http://localhost:8080/swagger for interactive API docs.

### `enpara-cli` tool
2. Place your Enpara statement PDF (e.g. `1- Enpara Hesap Hareketleri.pdf`) in the current directory. 

3. Convert to CSV (default):
   ```sh
   ./enpara-cli "1- Enpara Hesap Hareketleri.pdf"
   ```
   This creates `1- Enpara Hesap Hareketleri.csv` in the same folder.

4. Convert to JSON, XLSX, or OFX:
   ```sh
   ./enpara-cli "1- Enpara Hesap Hareketleri.pdf" --format json
   ./enpara-cli "1- Enpara Hesap Hareketleri.pdf" --format xlsx
   ./enpara-cli "1- Enpara Hesap Hareketleri.pdf" --format ofx
   ```

5. Choose PDF parser type (default: auto):
   ```sh
   ./enpara-cli "1- Enpara Hesap Hareketleri.pdf" --type auto
   ./enpara-cli "1- Enpara Hesap Hareketleri.pdf" --type type1
   ./enpara-cli "1- Enpara Hesap Hareketleri.pdf" --type type2
   ```

6. Set output file name:
   ```sh
   ./enpara-cli "1- Enpara Hesap Hareketleri.pdf" --format csv --output my.csv
   ```

## PDF Types

- `type1`: Manual statement layout (existing parser behavior).
- `type2`: Automatic monthly statement layout with columns `Tarih`, `Açıklama`, `Tutar`, `Bakiye`.
- `auto`: Detects layout from PDF text and chooses `type1` or `type2`.

For `type2`, an `NFC` field is included in JSON/CSV/XLSX outputs as `1` or `0`.

## API Parameters

`POST /api/v1/convert` multipart form fields:

- `file` (required): PDF file.
- `format` (optional): `json|csv|xlsx|ofx` (default `json`).
- `type` (optional): `auto|type1|type2` (default `auto`).


## Usage and Help

Show all Makefile targets:
```sh
make help
```

Show CLI help:
```sh
./enpara-cli --help
```

Show API help:
```sh
./enpara-api --help
```


## Build for All Platforms

Generate release artifacts:
```sh
make release
```

Artifacts are written to `dist/`.

## Verify Release Integrity and Authenticity

Assume you downloaded the release files into one folder.

1. Create a local keyring from the bundled public key and verify the checksum signature:
   ```sh
   gpg --no-default-keyring --keyring ./release-public-key.gpg --verify CHECKSUMS.sha256.asc CHECKSUMS.sha256
   ```

2. Verify file integrity:
   ```sh
   sha256sum -c CHECKSUMS.sha256
   ```

The first check confirms the checksums were signed by the release key without adding it to your normal keyring. The second check confirms the zip contents match the published checksums.


## Need Help?
Open an issue on [GitHub](https://github.com/chillmatin/enpara-transactions-parser/issues).
