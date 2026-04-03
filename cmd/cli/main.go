package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chillmatin/enpara-transactions-parser/pkg/converter"
	"github.com/chillmatin/enpara-transactions-parser/pkg/parser"
	"github.com/spf13/cobra"
)

func main() {
	var formatFlag string
	var outputFlag string
	var pdfTypeFlag string

	rootCmd := &cobra.Command{
		Use:   "enpara-cli <input.pdf>",
		Short: "Convert Enpara PDF statement to multiple output formats",
		Long: "Convert an Enpara PDF statement into JSON, CSV, XLSX, or OFX. " +
			"If output is omitted, a file is written next to the current directory with the same base name.",
		Example: strings.Join([]string{
			"  enpara-cli ./tmp/transaction.pdf",
			"  enpara-cli ./tmp/transaction.pdf --format csv",
			"  enpara-cli ./tmp/transaction.pdf -f ofx -o ./tmp/statement.ofx",
			"  enpara-cli ./tmp/automatic.pdf --type type2 -f json",
		}, "\n"),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputPath := args[0]
			return runConversion(inputPath, formatFlag, outputFlag, pdfTypeFlag)
		},
	}

	rootCmd.Flags().StringVarP(&formatFlag, "format", "f", "csv", "Output format (csv|json|xlsx|ofx)")
	rootCmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Output file path (default: <input-base>.<format>)")
	rootCmd.Flags().StringVarP(&pdfTypeFlag, "type", "t", parser.PDFTypeAuto, "Parser type (auto|type1|type2)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runConversion(inputPath string, format string, outputPath string, pdfType string) error {
	format = strings.ToLower(strings.TrimSpace(format))
	if format == "" {
		format = "csv"
	}

	pdfType = strings.ToLower(strings.TrimSpace(pdfType))
	if pdfType == "" {
		pdfType = parser.PDFTypeAuto
	}
	if pdfType != parser.PDFTypeAuto && pdfType != parser.PDFType1 && pdfType != parser.PDFType2 {
		return fmt.Errorf("unsupported pdf type %q (supported: auto, type1, type2)", pdfType)
	}

	fmt.Printf("PDF Type: %s\n", pdfType)
	statement, err := parser.ParseStatementFromPDF(inputPath, parser.ParseOptions{PDFType: pdfType})
	if err != nil {
		return fmt.Errorf("parse statement: %w", err)
	}

	var data []byte
	switch format {
	case "json":
		data, err = converter.ToJSON(statement)
	case "csv":
		data, err = converter.ToCSV(statement)
	case "xlsx":
		data, err = converter.ToXLSX(statement)
	case "ofx":
		data, err = converter.ToOFX(statement)
	default:
		return fmt.Errorf("unsupported format %q (supported: json, csv, xlsx, ofx)", format)
	}
	if err != nil {
		return fmt.Errorf("convert statement to %s: %w", format, err)
	}

	if outputPath == "" {
		outputPath = defaultOutputPath(inputPath, format)
	}

	if err := os.WriteFile(outputPath, data, 0o644); err != nil {
		return fmt.Errorf("write output file %q: %w", outputPath, err)
	}

	fmt.Printf("Conversion successful: %s\n", outputPath)
	return nil
}

func defaultOutputPath(inputPath string, format string) string {
	ext := "." + format
	if format == "ofx" {
		ext = ".ofx"
	}

	baseName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	if baseName == "" {
		baseName = "statement"
	}

	return baseName + ext
}
