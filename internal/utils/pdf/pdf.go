package pdf

import (
	"bytes"
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

// GenerateSimplePDF renders a basic multi-language PDF using gofpdf.
func GenerateSimplePDF(titleEn, titleAr string, sections []string) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8Font("ArialUnicode", "", "") // no-op if not provided
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, titleEn)
	pdf.Ln(12)

	if titleAr != "" {
		pdf.SetFont("Arial", "", 14)
		pdf.Cell(0, 10, titleAr)
		pdf.Ln(10)
	}

	pdf.SetFont("Arial", "", 12)
	for _, s := range sections {
		pdf.MultiCell(0, 7, s, "", "L", false)
		pdf.Ln(2)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("pdf output: %w", err)
	}
	return buf.Bytes(), nil
}

