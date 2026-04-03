package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/chillmatin/enpara-transactions-parser/pkg/converter"
	"github.com/chillmatin/enpara-transactions-parser/pkg/parser"
	"github.com/gin-gonic/gin"
)

var supportedFormats = map[string]struct{}{
	"json": {},
	"csv":  {},
	"xlsx": {},
	"ofx":  {},
}

var supportedPDFTypes = map[string]struct{}{
	parser.PDFTypeAuto: {},
	parser.PDFType1:    {},
	parser.PDFType2:    {},
}

func HandleConvert(c *gin.Context) {
	format := strings.ToLower(strings.TrimSpace(c.DefaultPostForm("format", "json")))
	if _, ok := supportedFormats[format]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported format"})
		return
	}

	pdfType := strings.ToLower(strings.TrimSpace(c.DefaultPostForm("type", parser.PDFTypeAuto)))
	if _, ok := supportedPDFTypes[pdfType]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported type"})
		return
	}

	uploadedFile, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	tempFile, err := os.CreateTemp("", "enpara-upload-*.pdf")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create temp file"})
		return
	}
	tempPath := tempFile.Name()
	defer os.Remove(tempPath)

	src, err := uploadedFile.Open()
	if err != nil {
		_ = tempFile.Close()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open uploaded file"})
		return
	}

	if _, err := io.Copy(tempFile, src); err != nil {
		_ = src.Close()
		_ = tempFile.Close()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save uploaded file"})
		return
	}

	_ = src.Close()
	if err := tempFile.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to finalize uploaded file"})
		return
	}

	data, filename, contentType, err := convertPDF(tempPath, format, pdfType, uploadedFile.Filename)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Use RFC 5987 encoding for internationalized filenames
	dispositionHeader := fmt.Sprintf("attachment; filename=%q", filename)
	// Add RFC 5987 filename* for proper Unicode support in filename
	// Use PathEscape for proper percent-encoding (spaces as %20, not +)
	dispositionHeader += fmt.Sprintf("; filename*=UTF-8''%s", url.PathEscape(filename))
	c.Header("Content-Disposition", dispositionHeader)
	c.Data(http.StatusOK, contentType, data)
}

func HandleFormats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"formats": []string{"json", "csv", "xlsx", "ofx"}})
}

func HandleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func convertPDF(pdfPath string, format string, pdfType string, originalFilename string) ([]byte, string, string, error) {
	statement, err := parser.ParseStatementFromPDF(pdfPath, parser.ParseOptions{PDFType: pdfType})
	if err != nil {
		return nil, "", "", fmt.Errorf("parse statement: %w", err)
	}

	base := strings.TrimSuffix(filepath.Base(originalFilename), filepath.Ext(originalFilename))
	if base == "" {
		base = "statement"
	}

	switch format {
	case "json":
		out, err := converter.ToJSON(statement)
		return out, base + ".json", "application/json; charset=utf-8", err
	case "csv":
		out, err := converter.ToCSV(statement)
		return out, base + ".csv", "text/csv; charset=utf-8", err
	case "xlsx":
		out, err := converter.ToXLSX(statement)
		return out, base + ".xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", err
	case "ofx":
		out, err := converter.ToOFX(statement)
		return out, base + ".ofx", "application/x-ofx; charset=utf-8", err
	default:
		return nil, "", "", fmt.Errorf("unsupported format %q", format)
	}
}
