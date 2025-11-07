package main

import (
	"fmt"
	"log"

	"github.com/beinux3/Mark2PDF"
)

func main() {
	// Example 1: Convert from string
	markdown := `# Welcome to Mark2PDF

## Introduction

This is a **Go** library for converting Markdown to PDF without external dependencies.

### Features

- Native Markdown parser
- Pure Go PDF generator
- No external dependencies
- Easy to use

### Formatting Examples

#### Unordered Lists

- First item
- Second item
- Third item

#### Ordered Lists

1. First step
2. Second step
3. Third step

#### Code Blocks

` + "```" + `go
package main

func main() {
    fmt.Println("Hello, World!")
}
` + "```" + `

#### Blockquotes

> This is a sample blockquote.
> It can span multiple lines.

---

## Conclusion

Mark2PDF makes Markdown to PDF conversion easy!
`

	// Convert to PDF
	pdfBytes, err := mark2pdf.ConvertString(markdown)
	if err != nil {
		log.Fatal("Error during conversion:", err)
	}

	fmt.Printf("PDF created successfully! Size: %d bytes\n", len(pdfBytes))

	// Example 2: Convert from string to file
	converter := mark2pdf.NewConverter(markdown)
	err = converter.ConvertToFile("example.pdf")
	if err != nil {
		log.Fatal("Error:", err)
	}

	fmt.Println("File example.pdf created!")
}
