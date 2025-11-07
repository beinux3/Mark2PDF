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
	objects     [][]byte
	pages       []int
	buffer      *bytes.Buffer
	yPosition   float64
	pageWidth   float64
	pageHeight  float64
	margin      float64
	currentPage int
	fontSizes   map[string]float64
}

// NewPDFWriter crea un nuovo writer PDF
func NewPDFWriter() *PDFWriter {
	return &PDFWriter{
		objects:     make([][]byte, 0),
		pages:       make([]int, 0),
		buffer:      &bytes.Buffer{},
		pageWidth:   595.28, // A4 width in points
		pageHeight:  841.89, // A4 height in points
		margin:      50.0,
		yPosition:   791.89, // Start position (top - margin)
		currentPage: 0,
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
func (p *PDFWriter) newPage() int {
	p.yPosition = p.pageHeight - p.margin
	p.currentPage++

	// Create page content stream
	contentBuf := &bytes.Buffer{}
	pageNum := p.addObject(contentBuf.Bytes())
	p.pages = append(p.pages, pageNum)

	return pageNum
}

// writeText scrive testo alla posizione corrente
func (p *PDFWriter) writeText(text string, fontSize float64, isBold bool) {
	if p.currentPage == 0 {
		p.newPage()
	}

	// Check if we need a new page
	if p.yPosition < p.margin+20 {
		p.newPage()
	}

	p.buffer.WriteString(fmt.Sprintf("BT\n"))
	p.buffer.WriteString(fmt.Sprintf("/F1 %.2f Tf\n", fontSize))
	p.buffer.WriteString(fmt.Sprintf("%.2f %.2f Td\n", p.margin, p.yPosition))

	// Escape special characters in text
	escapedText := escapeString(text)
	p.buffer.WriteString(fmt.Sprintf("(%s) Tj\n", escapedText))
	p.buffer.WriteString("ET\n")

	p.yPosition -= fontSize * 1.5
}

// writeLine scrive una linea orizzontale
func (p *PDFWriter) writeLine(width float64) {
	if p.currentPage == 0 {
		p.newPage()
	}

	p.buffer.WriteString(fmt.Sprintf("%.2f %.2f m\n", p.margin, p.yPosition))
	p.buffer.WriteString(fmt.Sprintf("%.2f %.2f l\n", p.margin+width, p.yPosition))
	p.buffer.WriteString("S\n")

	p.yPosition -= 10
}

// addSpace aggiunge spazio verticale
func (p *PDFWriter) addSpace(points float64) {
	p.yPosition -= points
	if p.yPosition < p.margin {
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
	// Finalize current page content
	if p.buffer.Len() > 0 {
		// Update the last page object with actual content
		if len(p.pages) > 0 {
			content := p.buffer.Bytes()

			// Compress content
			var compressed bytes.Buffer
			w := zlib.NewWriter(&compressed)
			w.Write(content)
			w.Close()

			streamObj := fmt.Sprintf("<< /Length %d /Filter /FlateDecode >>\nstream\n", compressed.Len())
			streamObj += compressed.String()
			streamObj += "\nendstream"

			p.objects[p.pages[len(p.pages)-1]-1] = []byte(streamObj)
		}
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

	// Pages object (Object 2)
	xrefPositions = append(xrefPositions, output.Len())
	output.WriteString(fmt.Sprintf("%d 0 obj\n", objNum))
	output.WriteString("<< /Type /Pages ")
	output.WriteString("/Kids [")
	for i := range p.pages {
		output.WriteString(fmt.Sprintf("%d 0 R ", objNum+1+i))
	}
	output.WriteString("] ")
	output.WriteString(fmt.Sprintf("/Count %d ", len(p.pages)))
	output.WriteString(">>\n")
	output.WriteString("endobj\n")
	objNum++

	// Font object (will be object 3)
	fontObjNum := objNum
	xrefPositions = append(xrefPositions, output.Len())
	output.WriteString(fmt.Sprintf("%d 0 obj\n", objNum))
	output.WriteString("<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>\n")
	output.WriteString("endobj\n")
	objNum++

	// Page objects and their content streams
	for i := range p.pages {
		// Page object
		xrefPositions = append(xrefPositions, output.Len())
		output.WriteString(fmt.Sprintf("%d 0 obj\n", objNum))
		output.WriteString("<< /Type /Page ")
		output.WriteString("/Parent 2 0 R ")
		output.WriteString(fmt.Sprintf("/MediaBox [0 0 %.2f %.2f] ", p.pageWidth, p.pageHeight))
		output.WriteString(fmt.Sprintf("/Contents %d 0 R ", objNum+len(p.pages)))
		output.WriteString(fmt.Sprintf("/Resources << /Font << /F1 %d 0 R >> >> ", fontObjNum))
		output.WriteString(">>\n")
		output.WriteString("endobj\n")

		// Update content stream position
		p.pages[i] = objNum + len(p.pages)
		objNum++
	}

	// Content streams
	for _, contentNum := range p.pages {
		xrefPositions = append(xrefPositions, output.Len())
		output.WriteString(fmt.Sprintf("%d 0 obj\n", contentNum))

		if contentNum <= len(p.objects) {
			output.Write(p.objects[contentNum-1])
		}

		output.WriteString("\nendobj\n")
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
