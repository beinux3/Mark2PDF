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
		c.renderInlineElements(elem.Children, c.pdf.GetFontSize("h1"))
		c.pdf.addSpace(5)

	case "h2":
		c.pdf.addSpace(8)
		c.renderInlineElements(elem.Children, c.pdf.GetFontSize("h2"))
		c.pdf.addSpace(4)

	case "h3":
		c.pdf.addSpace(6)
		c.renderInlineElements(elem.Children, c.pdf.GetFontSize("h3"))
		c.pdf.addSpace(3)

	case "h4":
		c.pdf.addSpace(5)
		c.renderInlineElements(elem.Children, c.pdf.GetFontSize("h4"))
		c.pdf.addSpace(2)

	case "h5":
		c.pdf.addSpace(4)
		c.renderInlineElements(elem.Children, c.pdf.GetFontSize("h5"))
		c.pdf.addSpace(2)

	case "h6":
		c.pdf.addSpace(3)
		c.renderInlineElements(elem.Children, c.pdf.GetFontSize("h6"))
		c.pdf.addSpace(2)

	case "p":
		c.renderInlineElements(elem.Children, c.pdf.GetFontSize("normal"))
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
				c.renderInlineElementsWithPrefix("  - ", elem.ItemChildren[i], c.pdf.GetFontSize("normal"))
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
				c.renderInlineElementsWithPrefix(prefix, elem.ItemChildren[i], c.pdf.GetFontSize("normal"))
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
				c.renderInlineElementsWithPrefix(checkbox, elem.ItemChildren[i], c.pdf.GetFontSize("normal"))
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
		c.pdf.addSpace(15)
		c.pdf.writeLine(c.pdf.pageWidth - c.pdf.margin*2)
		c.pdf.addSpace(15)
	}

	return nil
}

// renderInlineElementsWithPrefix renderizza un prefisso seguito da elementi inline
func (c *Converter) renderInlineElementsWithPrefix(prefix string, elements []InlineElement, baseFontSize float64) {
	// Build the complete text with formatting markers
	parts := []TextPart{}

	// Add prefix with regular font
	if prefix != "" {
		parts = append(parts, TextPart{Text: prefix, Font: "F1"})
	}

	// Add inline elements with their respective fonts
	for _, elem := range elements {
		var fontName string
		var text string

		switch elem.Type {
		case "text":
			fontName = "F1"
			text = elem.Content
		case "bold":
			fontName = "F2"
			text = elem.Content
		case "italic":
			fontName = "F3"
			text = elem.Content
		case "code":
			fontName = "F4"
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

		parts = append(parts, TextPart{Text: text, Font: fontName})
	}

	// Render with word wrapping
	c.writeMultiStyleTextWrapped(parts, baseFontSize)
}

// writeMultiStyleTextWrapped scrive testo multi-stile con word wrapping
func (c *Converter) writeMultiStyleTextWrapped(parts []TextPart, fontSize float64) {
	maxWidth := c.pdf.pageWidth - c.pdf.margin*2
	avgCharWidth := fontSize * 0.5

	currentLine := []TextPart{}
	currentWidth := 0.0

	for _, part := range parts {
		// Split text into words
		words := strings.Fields(part.Text)

		for i, word := range words {
			wordWidth := float64(len(word)) * avgCharWidth
			spaceWidth := avgCharWidth * 0.5

			// Add space before word (except for first word or after prefix)
			if i > 0 || (len(currentLine) > 0 && currentLine[len(currentLine)-1].Text != "") {
				testWidth := currentWidth + spaceWidth + wordWidth
				if testWidth > maxWidth && len(currentLine) > 0 {
					// Write current line
					c.pdf.writeMultiStyleText(currentLine, fontSize)
					currentLine = []TextPart{}
					currentWidth = 0
				} else {
					// Add space to last part if same font, otherwise create new part
					if len(currentLine) > 0 && currentLine[len(currentLine)-1].Font == part.Font {
						currentLine[len(currentLine)-1].Text += " "
						currentWidth += spaceWidth
					} else {
						currentLine = append(currentLine, TextPart{Text: " ", Font: part.Font})
						currentWidth += spaceWidth
					}
				}
			}

			// Check if word fits on current line
			if currentWidth+wordWidth > maxWidth && len(currentLine) > 0 {
				// Write current line and start new one
				c.pdf.writeMultiStyleText(currentLine, fontSize)
				currentLine = []TextPart{{Text: word, Font: part.Font}}
				currentWidth = wordWidth
			} else {
				// Add word to current line
				if len(currentLine) > 0 && currentLine[len(currentLine)-1].Font == part.Font {
					// Merge with previous part if same font
					currentLine[len(currentLine)-1].Text += word
				} else {
					currentLine = append(currentLine, TextPart{Text: word, Font: part.Font})
				}
				currentWidth += wordWidth
			}
		}
	}

	// Write remaining line
	if len(currentLine) > 0 {
		c.pdf.writeMultiStyleText(currentLine, fontSize)
	}
}

// renderInlineElements renderizza una lista di elementi inline con formattazione
func (c *Converter) renderInlineElements(elements []InlineElement, baseFontSize float64) {
	c.renderInlineElementsWithPrefix("", elements, baseFontSize)
}

// renderTable renderizza una tabella con bordi e celle
func (c *Converter) renderTable(elem MarkdownElement) {
	if len(elem.TableRows) == 0 {
		return
	}

	fontSize := c.pdf.GetFontSize("normal")
	cellPadding := 10.0 // Increased padding for better spacing
	rowHeight := fontSize * 2.5

	// Calculate column widths based on content
	maxWidth := c.pdf.pageWidth - c.pdf.margin*2
	numCols := len(elem.TableRows[0])

	// Calculate max width needed for each column (content only, without padding)
	colContentWidths := make([]float64, numCols)
	avgCharWidth := fontSize * 0.6 // More conservative estimate for character width

	for _, row := range elem.TableRows {
		for j, cell := range row {
			if j < numCols {
				contentWidth := float64(len(cell)) * avgCharWidth
				if contentWidth > colContentWidths[j] {
					colContentWidths[j] = contentWidth
				}
			}
		}
	}

	// Add padding to get total column widths
	colWidths := make([]float64, numCols)
	totalWidth := 0.0
	for j := range colContentWidths {
		colWidths[j] = colContentWidths[j] + cellPadding*2
		totalWidth += colWidths[j]
	}

	// Adjust if total width exceeds page width
	// Scale only the content width, not the padding
	if totalWidth > maxWidth {
		availableContentWidth := maxWidth - float64(numCols)*cellPadding*2
		totalContentWidth := 0.0
		for _, w := range colContentWidths {
			totalContentWidth += w
		}

		if totalContentWidth > 0 {
			scale := availableContentWidth / totalContentWidth
			for j := range colWidths {
				colWidths[j] = colContentWidths[j]*scale + cellPadding*2
			}
		}
		totalWidth = maxWidth
	}

	// Render each row
	for rowIdx, row := range elem.TableRows {
		startY := c.pdf.yPosition
		isHeader := rowIdx == 0

		// Draw cells for this row
		xPos := c.pdf.margin
		for colIdx, cell := range row {
			if colIdx >= numCols {
				break
			}

			cellWidth := colWidths[colIdx]

			// Draw cell border (thicker for header)
			if isHeader {
				c.pdf.drawThickRect(xPos, startY, cellWidth, rowHeight)
			} else {
				c.pdf.drawRect(xPos, startY, cellWidth, rowHeight)
			}

			// Truncate text if too long for cell
			maxCellChars := int((cellWidth - cellPadding*2) / avgCharWidth)
			displayText := cell
			if len(cell) > maxCellChars && maxCellChars > 3 {
				displayText = cell[:maxCellChars-3] + "..."
			}

			// Draw cell text (centered vertically in cell)
			// Calculate Y position: start from top of cell, go down by half row height, adjust for font baseline
			textY := startY - rowHeight/2 - fontSize/3
			c.pdf.writeTextAt(displayText, xPos+cellPadding, textY, fontSize, isHeader)

			xPos += cellWidth
		}

		// Move to next row
		c.pdf.yPosition = startY - rowHeight

		// Add extra spacing after header
		if isHeader {
			c.pdf.yPosition -= 2
		}
	}

	c.pdf.addSpace(5)
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
