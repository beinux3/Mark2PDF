package mark2pdf

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"time"
)

// PDFWriter gestisce la creazione di documenti PDF
type PDFWriter struct {
	objects      [][]byte
	pages        []int
	currentBuf   *bytes.Buffer
	yPosition    float64
	pageWidth    float64
	pageHeight   float64
	margin       float64
	currentPage  int
	fontSizes    map[string]float64
	pageContents []*bytes.Buffer
}

// NewPDFWriter crea un nuovo writer PDF
func NewPDFWriter() *PDFWriter {
	return &PDFWriter{
		objects:      make([][]byte, 0),
		pages:        make([]int, 0),
		currentBuf:   nil,
		pageWidth:    595.28, // A4 width in points
		pageHeight:   841.89, // A4 height in points
		margin:       50.0,
		yPosition:    0,
		currentPage:  -1,
		pageContents: make([]*bytes.Buffer, 0),
		fontSizes: map[string]float64{
			"h1":     24,
			"h2":     20,
			"h3":     16,
			"h4":     14,
			"h5":     12,
			"h6":     11,
			"normal": 10,
			"code":   9,
		},
	}
}

// addObject aggiunge un oggetto al PDF e restituisce il suo numero
func (p *PDFWriter) addObject(content []byte) int {
	p.objects = append(p.objects, content)
	return len(p.objects)
}

// newPage crea una nuova pagina
func (p *PDFWriter) newPage() {
	p.currentPage++
	p.yPosition = p.pageHeight - p.margin
	p.currentBuf = &bytes.Buffer{}
	p.pageContents = append(p.pageContents, p.currentBuf)
}

// writeText scrive testo alla posizione corrente
func (p *PDFWriter) writeText(text string, fontSize float64, isBold bool) {
	p.writeTextWithFont(text, fontSize, "F1") // Default font
}

// writeTextWithFont scrive testo con un font specifico
func (p *PDFWriter) writeTextWithFont(text string, fontSize float64, fontName string) {
	if p.currentBuf == nil {
		p.newPage()
	}

	// Check if we need a new page
	if p.yPosition < p.margin+20 {
		p.newPage()
	}

	p.currentBuf.WriteString("BT\n")
	p.currentBuf.WriteString(fmt.Sprintf("/%s %.2f Tf\n", fontName, fontSize))
	p.currentBuf.WriteString(fmt.Sprintf("%.2f %.2f Td\n", p.margin, p.yPosition))

	// Escape special characters in text
	escapedText := escapeString(text)
	p.currentBuf.WriteString(fmt.Sprintf("(%s) Tj\n", escapedText))
	p.currentBuf.WriteString("ET\n")

	p.yPosition -= fontSize * 1.5
}

// writeInlineText scrive testo inline senza andare a capo
func (p *PDFWriter) writeInlineText(text string, fontSize float64, fontName string, xOffset float64) {
	if p.currentBuf == nil {
		p.newPage()
	}

	// Check if we need a new page
	if p.yPosition < p.margin+20 {
		p.newPage()
	}

	p.currentBuf.WriteString("BT\n")
	p.currentBuf.WriteString(fmt.Sprintf("/%s %.2f Tf\n", fontName, fontSize))
	p.currentBuf.WriteString(fmt.Sprintf("%.2f %.2f Td\n", p.margin+xOffset, p.yPosition))

	// Escape special characters in text
	escapedText := escapeString(text)
	p.currentBuf.WriteString(fmt.Sprintf("(%s) Tj\n", escapedText))
	p.currentBuf.WriteString("ET\n")
}

// writeLine scrive una linea orizzontale
func (p *PDFWriter) writeLine(width float64) {
	if p.currentBuf == nil {
		p.newPage()
	}

	// Check if we need a new page
	if p.yPosition < p.margin+20 {
		p.newPage()
	}

	p.currentBuf.WriteString(fmt.Sprintf("%.2f %.2f m\n", p.margin, p.yPosition))
	p.currentBuf.WriteString(fmt.Sprintf("%.2f %.2f l\n", p.margin+width, p.yPosition))
	p.currentBuf.WriteString("S\n")

	p.yPosition -= 10
}

// addSpace aggiunge spazio verticale
func (p *PDFWriter) addSpace(points float64) {
	p.yPosition -= points
	if p.yPosition < p.margin+20 {
		p.newPage()
	}
}

// escapeString escapa caratteri speciali per PDF
func escapeString(s string) string {
	result := ""
	for _, r := range s {
		switch r {
		case '(', ')', '\\':
			result += "\\" + string(r)
		case '\n':
			continue
		case '\r':
			continue
		case '\t':
			result += "    "
		default:
			if r < 32 || r > 126 {
				// Replace non-ASCII with space for simplicity
				result += " "
			} else {
				result += string(r)
			}
		}
	}
	return result
}

// Build costruisce il PDF finale
func (p *PDFWriter) Build() ([]byte, error) {
	// If no pages were created, create an empty one
	if len(p.pageContents) == 0 {
		p.newPage()
	}

	output := &bytes.Buffer{}

	// PDF Header
	output.WriteString("%PDF-1.4\n")
	output.WriteString("%âăĎÓ\n") // Binary marker

	xrefPositions := make([]int, 0)
	xrefPositions = append(xrefPositions, 0) // Object 0 is always free

	// Catalog (Object 1)
	objNum := 1
	xrefPositions = append(xrefPositions, output.Len())
	output.WriteString(fmt.Sprintf("%d 0 obj\n", objNum))
	output.WriteString("<< /Type /Catalog /Pages 2 0 R >>\n")
	output.WriteString("endobj\n")
	objNum++

	// Pages object (Object 2) - must be referenced by Catalog
	pagesObjNum := objNum
	numPages := len(p.pageContents)

	// Reserve space for Pages object - we'll write it here
	pagesPos := output.Len()
	xrefPositions = append(xrefPositions, pagesPos)

	// Calculate object numbers
	fontObjNum := objNum + 1
	pageObjStart := fontObjNum + 4 // Now we have 4 fonts (F1, F2, F3, F4)
	contentObjStart := pageObjStart + numPages

	// Write Pages object
	output.WriteString(fmt.Sprintf("%d 0 obj\n", objNum))
	output.WriteString("<< /Type /Pages ")
	output.WriteString("/Kids [")
	for i := 0; i < numPages; i++ {
		output.WriteString(fmt.Sprintf("%d 0 R ", pageObjStart+i))
	}
	output.WriteString("] ")
	output.WriteString(fmt.Sprintf("/Count %d ", numPages))
	output.WriteString(">>\n")
	output.WriteString("endobj\n")
	objNum++

	// Font objects (F1=Helvetica, F2=Helvetica-Bold, F3=Helvetica-Oblique, F4=Courier)
	// F1 - Regular (Object 3)
	xrefPositions = append(xrefPositions, output.Len())
	output.WriteString(fmt.Sprintf("%d 0 obj\n", fontObjNum))
	output.WriteString("<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>\n")
	output.WriteString("endobj\n")

	// F2 - Bold (Object 4)
	xrefPositions = append(xrefPositions, output.Len())
	output.WriteString(fmt.Sprintf("%d 0 obj\n", fontObjNum+1))
	output.WriteString("<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica-Bold >>\n")
	output.WriteString("endobj\n")

	// F3 - Italic (Object 5)
	xrefPositions = append(xrefPositions, output.Len())
	output.WriteString(fmt.Sprintf("%d 0 obj\n", fontObjNum+2))
	output.WriteString("<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica-Oblique >>\n")
	output.WriteString("endobj\n")

	// F4 - Code/Monospace (Object 6)
	xrefPositions = append(xrefPositions, output.Len())
	output.WriteString(fmt.Sprintf("%d 0 obj\n", fontObjNum+3))
	output.WriteString("<< /Type /Font /Subtype /Type1 /BaseFont /Courier >>\n")
	output.WriteString("endobj\n")

	objNum = fontObjNum + 4

	// Page objects
	for i := range p.pageContents {
		xrefPositions = append(xrefPositions, output.Len())
		pageObjNum := pageObjStart + i
		contentObjNum := contentObjStart + i

		output.WriteString(fmt.Sprintf("%d 0 obj\n", pageObjNum))
		output.WriteString("<< /Type /Page ")
		output.WriteString(fmt.Sprintf("/Parent %d 0 R ", pagesObjNum))
		output.WriteString(fmt.Sprintf("/MediaBox [0 0 %.2f %.2f] ", p.pageWidth, p.pageHeight))
		output.WriteString(fmt.Sprintf("/Contents %d 0 R ", contentObjNum))
		// Include all 4 fonts in resources
		output.WriteString(fmt.Sprintf("/Resources << /Font << /F1 %d 0 R /F2 %d 0 R /F3 %d 0 R /F4 %d 0 R >> >> ",
			fontObjNum, fontObjNum+1, fontObjNum+2, fontObjNum+3))
		output.WriteString(">>\n")
		output.WriteString("endobj\n")
	}

	// Content streams
	for i, pageBuf := range p.pageContents {
		content := pageBuf.Bytes()

		// Compress content
		var compressed bytes.Buffer
		w := zlib.NewWriter(&compressed)
		w.Write(content)
		w.Close()

		contentObjNum := contentObjStart + i
		xrefPositions = append(xrefPositions, output.Len())
		output.WriteString(fmt.Sprintf("%d 0 obj\n", contentObjNum))
		output.WriteString(fmt.Sprintf("<< /Length %d /Filter /FlateDecode >>\n", compressed.Len()))
		output.WriteString("stream\n")
		output.Write(compressed.Bytes())
		output.WriteString("\nendstream\n")
		output.WriteString("endobj\n")
	}

	// xref table
	xrefPos := output.Len()
	output.WriteString("xref\n")
	output.WriteString(fmt.Sprintf("0 %d\n", len(xrefPositions)))
	output.WriteString("0000000000 65535 f \n")
	for i := 1; i < len(xrefPositions); i++ {
		output.WriteString(fmt.Sprintf("%010d 00000 n \n", xrefPositions[i]))
	}

	// Trailer
	output.WriteString("trailer\n")
	output.WriteString(fmt.Sprintf("<< /Size %d /Root 1 0 R >>\n", len(xrefPositions)))
	output.WriteString("startxref\n")
	output.WriteString(fmt.Sprintf("%d\n", xrefPos))
	output.WriteString("%%EOF\n")

	return output.Bytes(), nil
}

// WriteTo implementa io.WriterTo
func (p *PDFWriter) WriteTo(w io.Writer) (int64, error) {
	data, err := p.Build()
	if err != nil {
		return 0, err
	}
	n, err := w.Write(data)
	return int64(n), err
}

// GetCurrentY restituisce la posizione Y corrente
func (p *PDFWriter) GetCurrentY() float64 {
	return p.yPosition
}

// SetY imposta la posizione Y
func (p *PDFWriter) SetY(y float64) {
	p.yPosition = y
}

// GetFontSize restituisce la dimensione del font per un tipo di testo
func (p *PDFWriter) GetFontSize(textType string) float64 {
	if size, ok := p.fontSizes[textType]; ok {
		return size
	}
	return p.fontSizes["normal"]
}

// WriteMetadata aggiunge metadati al PDF (semplificato)
func (p *PDFWriter) WriteMetadata(title, author string) {
	// In una implementazione completa, questo aggiungerebbe un oggetto Info
	// Per semplicità, lo omettiamo in questa versione base
	_ = title
	_ = author
	_ = time.Now()
}
