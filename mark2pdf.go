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
		c.pdf.writeText(elem.Content, c.pdf.GetFontSize("h1"), true)
		c.pdf.addSpace(5)

	case "h2":
		c.pdf.addSpace(8)
		c.pdf.writeText(elem.Content, c.pdf.GetFontSize("h2"), true)
		c.pdf.addSpace(4)

	case "h3":
		c.pdf.addSpace(6)
		c.pdf.writeText(elem.Content, c.pdf.GetFontSize("h3"), true)
		c.pdf.addSpace(3)

	case "h4":
		c.pdf.addSpace(5)
		c.pdf.writeText(elem.Content, c.pdf.GetFontSize("h4"), true)
		c.pdf.addSpace(2)

	case "h5":
		c.pdf.addSpace(4)
		c.pdf.writeText(elem.Content, c.pdf.GetFontSize("h5"), true)
		c.pdf.addSpace(2)

	case "h6":
		c.pdf.addSpace(3)
		c.pdf.writeText(elem.Content, c.pdf.GetFontSize("h6"), true)
		c.pdf.addSpace(2)

	case "p":
		content := stripInlineFormatting(elem.Content)
		c.writeWrappedText(content, c.pdf.GetFontSize("normal"), false)
		c.pdf.addSpace(8)

	case "code":
		c.pdf.addSpace(5)
		lines := strings.Split(elem.Content, "\n")
		for _, line := range lines {
			c.pdf.writeText(line, c.pdf.GetFontSize("code"), false)
		}
		c.pdf.addSpace(5)

	case "list":
		c.pdf.addSpace(3)
		for _, item := range elem.Items {
			text := "  - " + stripInlineFormatting(item)
			c.writeWrappedText(text, c.pdf.GetFontSize("normal"), false)
		}
		c.pdf.addSpace(5)

	case "ordered-list":
		c.pdf.addSpace(3)
		for i, item := range elem.Items {
			text := fmt.Sprintf("  %d. %s", i+1, stripInlineFormatting(item))
			c.writeWrappedText(text, c.pdf.GetFontSize("normal"), false)
		}
		c.pdf.addSpace(5)

	case "blockquote":
		c.pdf.addSpace(5)
		content := "  | " + stripInlineFormatting(elem.Content)
		c.writeWrappedText(content, c.pdf.GetFontSize("normal"), false)
		c.pdf.addSpace(5)

	case "hr":
		c.pdf.addSpace(5)
		c.pdf.writeLine(c.pdf.pageWidth - c.pdf.margin*2)
		c.pdf.addSpace(5)
	}

	return nil
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
