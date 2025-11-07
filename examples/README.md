# Mark2PDF Examples

This directory contains examples demonstrating all Mark2PDF features.

## Quick Start

Generate all example PDFs:

```bash
go run examples.go
```

This will create 4 PDF files demonstrating different features:

## Generated Examples

### 1. Basic Features (`1_basic.pdf`)
- Headers (all 6 levels)
- Paragraphs with word wrapping
- Lists (unordered, ordered, task lists)
- Code blocks with syntax highlighting
- Blockquotes
- Tables with borders
- Horizontal rules

### 2. Inline Formatting (`2_inline_formatting.pdf`)
- **Bold**, *italic*, `code`, ~~strikethrough~~
- [Links](https://example.com)
- Combined formatting styles
- Formatting in lists

### 3. Colors (`3_colors.pdf`)
- Named colors (red, green, blue, etc.)
- Custom RGB colors: `{color:rgb(255,105,180)}text{/color}`
- Hex colors: `{color:#8B4513}text{/color}`
- Nested formatting: `{blue}**bold inside blue**{/blue}`
- Colors in tables

### 4. Complete Demo (`4_complete_demo.pdf`)
- All features combined in one document
- Real-world usage examples
- Best practices demonstration

## Color Syntax

Mark2PDF supports colored text with simple syntax:

### Named Colors
```markdown
{red}Red text{/red}
{blue}Blue text{/blue}
```

Available colors: `red`, `green`, `blue`, `yellow`, `cyan`, `magenta`, `orange`, `purple`, `gray`, `black`, `white`

### Custom Colors
```markdown
{color:rgb(255,100,50)}Custom RGB{/color}
{color:#FF6347}Hex color{/color}
```

### Nested Formatting
```markdown
{blue}**Bold inside blue**{/blue}
{green}Text with **bold** and *italic*{/green}
```

## Creating Your Own Examples

Use the `mark2pdf` package in your Go code:

```go
package main

import "github.com/beinux3/Mark2PDF"

func main() {
    markdown := "# Hello\n\nThis is **bold** and {red}red{/red}."
    mark2pdf.ConvertString(markdown) // Returns []byte

    // Or save directly to file
    mark2pdf.ConvertFile("input.md", "output.pdf")
}
```

## Features Demonstrated

- ✅ All 6 header levels
- ✅ Text formatting (bold, italic, code, strikethrough)
- ✅ Colors (named, RGB, hex)
- ✅ Nested formatting (colors containing bold/italic)
- ✅ Lists (unordered, ordered, task lists)
- ✅ Code blocks with language tags
- ✅ Tables with borders and formatting
- ✅ Blockquotes
- ✅ Links
- ✅ Horizontal rules
- ✅ Word wrapping
- ✅ Multi-page documents

## Notes

- All PDFs are generated in PDF 1.4 format
- Page size is A4 (595.28 × 841.89 points)
- Margins are 50 points on all sides
- Text automatically wraps to fit page width
- Tables scale automatically to fit page width
