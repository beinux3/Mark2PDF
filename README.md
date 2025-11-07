# Mark2PDF

[![Go Reference](https://pkg.go.dev/badge/github.com/beinux3/Mark2PDF.svg)](https://pkg.go.dev/github.com/beinux3/Mark2PDF)
[![Go Report Card](https://goreportcard.com/badge/github.com/beinux3/Mark2PDF)](https://goreportcard.com/report/github.com/beinux3/Mark2PDF)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A pure Go library for converting Markdown documents to PDF without any external dependencies.

## Features

- **Zero external dependencies**: Completely native Go implementation
- **Full Markdown parser**: Supports headers, lists, code blocks, blockquotes, and more
- **Native PDF generator**: Creates valid PDFs without external libraries
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

```markdown
# H1
## H2
### H3
#### H4
##### H5
###### H6
```

### Paragraphs

Regular text with automatic line wrapping.

### Unordered Lists

```markdown
- Item 1
- Item 2
- Item 3
```

### Ordered Lists

```markdown
1. First
2. Second
3. Third
```

### Code Blocks

````markdown
```go
package main
func main() {}
```
````

### Blockquotes

```markdown
> This is a blockquote
```

### Horizontal Rules

```markdown
---
***
___
```

## Project Structure

```
mark2pdf/
├── mark2pdf.go      # Main API and converter
├── markdown.go      # Markdown parser
├── pdf.go           # PDF generator
├── examples/        # Usage examples
│   ├── basic.go
│   └── test.md
├── go.mod
├── LICENSE
└── README.md
```

## How It Works

1. **Parsing**: The Markdown parser analyzes the document and breaks it down into structured elements (headers, paragraphs, lists, etc.)
2. **Rendering**: Each element is rendered in the PDF using basic PDF primitives
3. **Generation**: The PDF generator creates a valid PDF document conforming to PDF 1.4 standard

## Limitations

This is a basic implementation that supports the most common Markdown elements. Some advanced features are not yet implemented:

- Images
- Tables
- Complex inline formatting (bold, italic, links are partially supported)
- Custom fonts (uses only Helvetica)
- Custom colors

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
go run basic.go
```

This will create a sample PDF file demonstrating all supported features.

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