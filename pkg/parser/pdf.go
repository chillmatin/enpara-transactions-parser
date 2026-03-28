package parser

import (
	"fmt"
	"strings"

	"github.com/ledongthuc/pdf"
)

func ExtractTextFromPDF(filepath string) (string, error) {
	file, reader, err := pdf.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("open pdf %q: %w", filepath, err)
	}
	defer file.Close()

	var content strings.Builder
	for pageIndex := 1; pageIndex <= reader.NumPage(); pageIndex++ {
		page := reader.Page(pageIndex)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			return "", fmt.Errorf("extract text from page %d: %w", pageIndex, err)
		}
		content.WriteString(text)
		if !strings.HasSuffix(text, "\n") {
			content.WriteByte('\n')
		}
	}

	return content.String(), nil
}
