package mark2pdf

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Converter converte Markdown in PDF
type Converter struct {
	pdf    *PDFWriter
	parser *MarkdownParser
}

// NewConverter crea un nuovo convertitore
func NewConverter(markdown string) *Converter {
	return &Converter{
		pdf:    NewPDFWriter(),
		parser: NewMarkdownParser(markdown),
	}
}

// Convert esegue la conversione e restituisce i byte del PDF
func (c *Converter) Convert() ([]byte, error) {
	elements := c.parser.Parse()

	for _, elem := range elements {
		if err := c.renderElement(elem); err != nil {
			return nil, err
		}
	}

	return c.pdf.Build()
}

// ConvertToFile converte e salva in un file
func (c *Converter) ConvertToFile(filename string) error {
	data, err := c.Convert()
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// ConvertToWriter converte e scrive su un Writer
func (c *Converter) ConvertToWriter(w io.Writer) error {
	data, err := c.Convert()
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

// renderElement renderizza un singolo elemento markdown
func (c *Converter) renderElement(elem MarkdownElement) error {
	switch elem.Type {
	case "h1":
		c.pdf.addSpace(10)
		c.renderInlineElements(elem.Children, c.pdf.GetFontSize("h1"), true)
		c.pdf.addSpace(5)

	case "h2":
		c.pdf.addSpace(8)
		c.renderInlineElements(elem.Children, c.pdf.GetFontSize("h2"), true)
		c.pdf.addSpace(4)

	case "h3":
		c.pdf.addSpace(6)
		c.renderInlineElements(elem.Children, c.pdf.GetFontSize("h3"), true)
		c.pdf.addSpace(3)

	case "h4":
		c.pdf.addSpace(5)
		c.renderInlineElements(elem.Children, c.pdf.GetFontSize("h4"), true)
		c.pdf.addSpace(2)

	case "h5":
		c.pdf.addSpace(4)
		c.renderInlineElements(elem.Children, c.pdf.GetFontSize("h5"), true)
		c.pdf.addSpace(2)

	case "h6":
		c.pdf.addSpace(3)
		c.renderInlineElements(elem.Children, c.pdf.GetFontSize("h6"), true)
		c.pdf.addSpace(2)

	case "p":
		c.renderInlineElements(elem.Children, c.pdf.GetFontSize("normal"), false)
		c.pdf.addSpace(8)

	case "code":
		c.pdf.addSpace(5)
		if elem.Language != "" {
			c.pdf.writeText("Code ("+elem.Language+"):", c.pdf.GetFontSize("normal"), false)
		}
		lines := strings.Split(elem.Content, "\n")
		for _, line := range lines {
			c.pdf.writeText("  "+line, c.pdf.GetFontSize("code"), false)
		}
		c.pdf.addSpace(5)

	case "list":
		c.pdf.addSpace(3)
		for i, item := range elem.Items {
			// Use inline children if available, otherwise use raw item
			if i < len(elem.ItemChildren) && len(elem.ItemChildren[i]) > 0 {
				c.pdf.writeText("  - ", c.pdf.GetFontSize("normal"), false)
				c.renderInlineElements(elem.ItemChildren[i], c.pdf.GetFontSize("normal"), false)
			} else {
				text := "  - " + item
				c.writeWrappedText(text, c.pdf.GetFontSize("normal"), false)
			}
		}
		c.pdf.addSpace(5)

	case "ordered-list":
		c.pdf.addSpace(3)
		for i, item := range elem.Items {
			prefix := fmt.Sprintf("  %d. ", i+1)
			if i < len(elem.ItemChildren) && len(elem.ItemChildren[i]) > 0 {
				c.pdf.writeText(prefix, c.pdf.GetFontSize("normal"), false)
				c.renderInlineElements(elem.ItemChildren[i], c.pdf.GetFontSize("normal"), false)
			} else {
				text := prefix + item
				c.writeWrappedText(text, c.pdf.GetFontSize("normal"), false)
			}
		}
		c.pdf.addSpace(5)

	case "task-list":
		c.pdf.addSpace(3)
		for i, item := range elem.Items {
			// Extract checkbox status
			checkbox := "  [ ] "
			if strings.HasPrefix(item, "[x]") {
				checkbox = "  [x] "
			}

			if i < len(elem.ItemChildren) && len(elem.ItemChildren[i]) > 0 {
				c.pdf.writeText(checkbox, c.pdf.GetFontSize("normal"), false)
				c.renderInlineElements(elem.ItemChildren[i], c.pdf.GetFontSize("normal"), false)
			} else {
				c.writeWrappedText("  "+item, c.pdf.GetFontSize("normal"), false)
			}
		}
		c.pdf.addSpace(5)

	case "blockquote":
		c.pdf.addSpace(5)
		text := "  | " + elem.Content
		c.writeWrappedText(text, c.pdf.GetFontSize("normal"), false)
		c.pdf.addSpace(5)

	case "table":
		c.pdf.addSpace(5)
		c.renderTable(elem)
		c.pdf.addSpace(5)

	case "hr":
		c.pdf.addSpace(5)
		c.pdf.writeLine(c.pdf.pageWidth - c.pdf.margin*2)
		c.pdf.addSpace(5)
	}

	return nil
}

// renderInlineElements renderizza una lista di elementi inline con formattazione
func (c *Converter) renderInlineElements(elements []InlineElement, baseFontSize float64, isBold bool) {
	if len(elements) == 0 {
		return
	}

	// Calculate text width to handle wrapping
	xOffset := 0.0
	avgCharWidth := baseFontSize * 0.5
	maxWidth := c.pdf.pageWidth - c.pdf.margin*2
	currentLineWidth := 0.0

	for _, elem := range elements {
		var fontName string
		var text string

		switch elem.Type {
		case "text":
			fontName = "F1" // Regular
			text = elem.Content
		case "bold":
			fontName = "F2" // Bold
			text = elem.Content
		case "italic":
			fontName = "F3" // Italic
			text = elem.Content
		case "code":
			fontName = "F4" // Monospace
			text = elem.Content
		case "link":
			fontName = "F1"
			text = elem.Content + " (" + elem.URL + ")"
		case "image":
			fontName = "F1"
			text = "[Image: " + elem.Alt + "]"
		case "strikethrough":
			fontName = "F1"
			text = elem.Content
		default:
			continue
		}

		// Estimate text width
		textWidth := float64(len(text)) * avgCharWidth

		// Check if we need to wrap
		if currentLineWidth+textWidth > maxWidth {
			// Move to next line
			c.pdf.yPosition -= baseFontSize * 1.5
			currentLineWidth = 0
			xOffset = 0
		}

		// Write the text inline
		c.pdf.writeInlineText(text, baseFontSize, fontName, xOffset)
		xOffset += textWidth
		currentLineWidth += textWidth
	}

	// Move to next line after finishing all elements
	c.pdf.yPosition -= baseFontSize * 1.5
}

// renderTable renderizza una tabella
func (c *Converter) renderTable(elem MarkdownElement) {
	if len(elem.TableRows) == 0 {
		return
	}

	// Calculate column widths
	maxWidth := c.pdf.pageWidth - c.pdf.margin*2
	numCols := len(elem.TableRows[0])
	_ = maxWidth / float64(numCols) // colWidth for future use

	// Render header row
	if len(elem.TableRows) > 0 {
		headerText := ""
		for i, cell := range elem.TableRows[0] {
			if i > 0 {
				headerText += " | "
			}
			headerText += cell
		}
		c.pdf.writeText(headerText, c.pdf.GetFontSize("normal"), true)
		c.pdf.writeLine(maxWidth)
	}

	// Render data rows
	for i := 1; i < len(elem.TableRows); i++ {
		rowText := ""
		for j, cell := range elem.TableRows[i] {
			if j > 0 {
				rowText += " | "
			}
			rowText += cell
		}
		c.writeWrappedText(rowText, c.pdf.GetFontSize("normal"), false)
	}
}

// writeWrappedText scrive testo con word wrapping
func (c *Converter) writeWrappedText(text string, fontSize float64, isBold bool) {
	maxWidth := c.pdf.pageWidth - c.pdf.margin*2
	avgCharWidth := fontSize * 0.5 // Approssimazione della larghezza media del carattere

	maxCharsPerLine := int(maxWidth / avgCharWidth)
	if maxCharsPerLine < 10 {
		maxCharsPerLine = 10
	}

	words := strings.Fields(text)
	currentLine := ""

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		if len(testLine) <= maxCharsPerLine {
			currentLine = testLine
		} else {
			// Scrive la linea corrente e inizia una nuova
			if currentLine != "" {
				c.pdf.writeText(currentLine, fontSize, isBold)
			}
			currentLine = word
		}
	}

	// Scrive l'ultima linea
	if currentLine != "" {
		c.pdf.writeText(currentLine, fontSize, isBold)
	}
}

// ConvertString è una funzione helper per conversioni veloci
func ConvertString(markdown string) ([]byte, error) {
	converter := NewConverter(markdown)
	return converter.Convert()
}

// ConvertFile converte un file markdown in PDF
func ConvertFile(inputFile, outputFile string) error {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("errore lettura file: %w", err)
	}

	converter := NewConverter(string(data))
	return converter.ConvertToFile(outputFile)
}
