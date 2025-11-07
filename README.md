# Mark2PDF

[![Go Reference](https://pkg.go.dev/badge/github.com/beinux3/Mark2PDF.svg)](https://pkg.go.dev/github.com/beinux3/Mark2PDF)
[![Go Report Card](https://goreportcard.com/badge/github.com/beinux3/Mark2PDF)](https://goreportcard.com/report/github.com/beinux3/Mark2PDF)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A pure Go library for converting Markdown documents to PDF without any external dependencies.

## Features

- **Zero external dependencies**: Completely native Go implementation
- **Full Markdown parser**: Supports headers, lists, code blocks, blockquotes, tables, and more
- **Rich inline formatting**: Bold, italic, inline code, links, images, and strikethrough
- **Color support**: Named colors, RGB, and hex color codes for text
- **Advanced table rendering**: Bordered tables with automatic column sizing and padding
- **Native PDF generator**: Creates valid PDFs conforming to PDF 1.4 standard
- **Word wrapping**: Automatic text wrapping for long paragraphs and inline elements
- **Multiple fonts**: Helvetica (regular, bold, oblique) and Courier for code
- **Simple API**: Intuitive and easy-to-use interface
- **Lightweight**: Efficient and performant code

## Installation

```bash
go get github.com/beinux3/Mark2PDF
```

## Quick Start

```go
package main

import (
    "log"
    "github.com/beinux3/Mark2PDF"
)

func main() {
    markdown := `# Hello World

This is an example of Markdown to PDF conversion.

## Features

- Easy to use
- No dependencies
- 100% Go
`

    // Convert Markdown string to PDF bytes
    pdfBytes, err := mark2pdf.ConvertString(markdown)
    if err != nil {
        log.Fatal(err)
    }

    // Or save directly to a file
    err = mark2pdf.ConvertFile("input.md", "output.pdf")
    if err != nil {
        log.Fatal(err)
    }
}
```

## API Reference

### ConvertString

Converts a Markdown string to PDF bytes:

```go
pdfBytes, err := mark2pdf.ConvertString(markdownString)
```

### ConvertFile

Converts a Markdown file to a PDF file:

```go
err := mark2pdf.ConvertFile("input.md", "output.pdf")
```

### Custom Converter

For more control:

```go
converter := mark2pdf.NewConverter(markdownString)

// Convert to file
err := converter.ConvertToFile("output.pdf")

// Or get bytes
pdfBytes, err := converter.Convert()

// Or write to an io.Writer
err := converter.ConvertToWriter(writer)
```

## Supported Markdown Elements

### Headers

All six levels of headers are supported with appropriate font sizes:

```markdown
# H1 (24pt)
## H2 (20pt)
### H3 (16pt)
#### H4 (14pt)
##### H5 (12pt)
###### H6 (11pt)
```

### Paragraphs

Regular text with automatic word wrapping to fit page width.

### Inline Formatting

Rich text formatting within paragraphs, headers, and lists:

```markdown
**bold text**
*italic text*
`inline code`
[link text](https://example.com)
![image alt text](https://example.com/image.png)
~~strikethrough~~
```

You can also combine formatting:
```markdown
**bold with *italic* inside**
**bold with `code`**
```

### Colors

Mark2PDF supports colored text using a simple syntax:

```markdown
{red}This text is red{/red}
{blue}This text is blue{/blue}
{color:green}This text is green{/color}
```

**Supported named colors:**
- `red`, `green`, `blue`, `yellow`, `cyan`, `magenta`, `orange`, `purple`, `gray`/`grey`, `black`, `white`

**Custom colors:**
```markdown
{color:rgb(255,100,50)}Custom RGB color{/color}
{color:#FF6347}Hex color (tomato){/color}
```

**Combining colors with other formatting:**
```markdown
**Bold with {red}red text{/red} inside**
{blue}Blue text with **bold** inside{/blue}
### {purple}Colored Header{/purple}
```

**Nested formatting (colors can contain formatting):**
```markdown
{blue}**Bold inside blue**{/blue}
{red}*Italic inside red*{/red}
{green}Text with **bold** and *italic*{/green}
```

### Unordered Lists

With full support for inline formatting:

```markdown
- Item 1 with **bold**
- Item 2 with *italic*
- Item 3 with `code`
```

### Ordered Lists

Automatically numbered lists:

```markdown
1. First item with **formatting**
2. Second item with *emphasis*
3. Third item with `code`
```

### Task Lists

GitHub-style task lists with checkboxes:

```markdown
- [ ] Incomplete task with **bold**
- [x] Completed task with *italic*
- [ ] Another task with `code`
```

### Code Blocks

Syntax-highlighted code blocks with language specification:

````markdown
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

```python
def hello():
    print("Hello, World!")
```
````

### Blockquotes

Quote blocks with inline formatting support:

```markdown
> This is a blockquote with **bold** and *italic*.
> It can span multiple lines.
```

### Tables

Markdown tables with automatic column sizing and borders:

```markdown
| Header 1 | Header 2 | Header 3 |
|----------|----------|----------|
| Cell 1   | Cell 2   | Cell 3   |
| Data 1   | Data 2   | Data 3   |
```

Features:
- Automatic column width calculation based on content
- Proportional scaling when table exceeds page width
- Thicker borders for header row
- Proper cell padding (10pt)
- Text truncation with "..." for overflow

### Horizontal Rules

Visual separators with ample spacing:

```markdown
---
***
___
```

## Project Structure

```
Mark2PDF/
├── mark2pdf.go      # Main API and converter
├── markdown.go      # Markdown parser with color support
├── pdf.go           # PDF generator with RGB colors
├── color_test.go    # Unit tests for color functionality
├── examples/        # Usage examples
│   ├── examples.go  # Example generator
│   └── README.md    # Examples documentation
├── go.mod
├── LICENSE
├── CHANGELOG.md
├── COLOR_SUPPORT.md
└── README.md
```

## How It Works

1. **Parsing**: The Markdown parser analyzes the document and breaks it down into structured elements (headers, paragraphs, lists, etc.)
2. **Rendering**: Each element is rendered in the PDF using basic PDF primitives
3. **Generation**: The PDF generator creates a valid PDF document conforming to PDF 1.4 standard

## Technical Details

### PDF Generation

- **PDF Version**: 1.4 specification
- **Page Size**: A4 (595.28 × 841.89 points)
- **Margins**: 50 points on all sides
- **Compression**: zlib (FlateDecode) for content streams
- **Fonts**:
  - F1: Helvetica (regular text)
  - F2: Helvetica-Bold (bold text, table headers)
  - F3: Helvetica-Oblique (italic text)
  - F4: Courier (code blocks and inline code)

### Font Sizes

- H1: 24pt
- H2: 20pt
- H3: 16pt
- H4: 14pt
- H5: 12pt
- H6: 11pt
- Normal text: 10pt
- Code: 9pt

### Word Wrapping

The library implements intelligent word wrapping for:
- Paragraphs with mixed inline formatting
- List items with inline formatting
- Long text that exceeds page width

Text is split on word boundaries and distributed across multiple lines while preserving formatting (bold, italic, code) throughout.

## Limitations

Some advanced features are not yet implemented:

- Image embedding (images are displayed as text references)
- Custom fonts (limited to standard PDF fonts)
- Nested lists
- Custom page sizes
- Headers and footers

## Integration with Your Go Project

### Step 1: Install the package

In your Go project directory, run:

```bash
go get github.com/beinux3/Mark2PDF
```

This will add the package to your `go.mod` file.

### Step 2: Import and use

```go
package main

import (
    "github.com/beinux3/Mark2PDF"
    "log"
)

func main() {
    // Your Markdown content
    md := "# My Document\n\nThis is a paragraph."

    // Convert to PDF
    if err := mark2pdf.ConvertFile("input.md", "output.pdf"); err != nil {
        log.Fatal(err)
    }
}
```

### Step 3: Build your application

```bash
go build
```

The `mark2pdf` package will be automatically downloaded and included in your build.

## Examples

Complete examples can be found in the `examples/` directory:

```bash
cd examples
go run examples.go
```

This will generate 4 PDF files demonstrating all features:

1. **1_basic.pdf** - Headers, lists, code blocks, tables, blockquotes
2. **2_inline_formatting.pdf** - Bold, italic, code, links, strikethrough
3. **3_colors.pdf** - Named colors, RGB, hex, nested formatting
4. **4_complete_demo.pdf** - Complete showcase of all features

See [`examples/README.md`](examples/README.md) for detailed documentation.

### Command Line Tool

Build and use the command-line tool:

```bash
# Build
make build

# Convert a Markdown file
./bin/mark2pdf -input document.md -output document.pdf

# Show version
./bin/mark2pdf -version

# Show help
./bin/mark2pdf -help
```

## Contributing

Contributions are welcome! Feel free to open issues or pull requests.

### Development

1. Clone the repository:
```bash
git clone https://github.com/beinux3/Mark2PDF.git
cd mark2pdf
```

2. Make your changes

3. Run tests (when available):
```bash
go test ./...
```

4. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) file for details

## Author

Created as an example of a pure Go library without external dependencies.