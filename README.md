# Enpara Transactions Parser

Convert Enpara PDF account statements to CSV, JSON, XLSX, or OFX.


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

5. Set output file name:
   ```sh
   ./enpara-cli "1- Enpara Hesap Hareketleri.pdf" --format csv --output my.csv
   ```


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

To build binaries for all supported OS/architectures:
```sh
make release
ls dist/
```
Binaries are in the `dist/` folder.


## Need Help?
Open an issue on [GitHub](https://github.com/chillmatin/enpara-transactions-parser/issues).
