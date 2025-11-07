# Color Support in Mark2PDF

## Overview

Mark2PDF now supports colored text in PDF output using a simple and intuitive syntax. Colors can be applied to any inline text and work seamlessly with other formatting options like bold, italic, and code.

## Syntax

### Named Colors

Use predefined color names:

```markdown
{red}This text is red{/red}
{blue}This text is blue{/blue}
{green}This text is green{/green}
```

### Alternative Syntax

```markdown
{color:red}This text is red{/color}
{color:blue}This text is blue{/color}
```

### RGB Colors

Specify custom colors using RGB values (0-255):

```markdown
{color:rgb(255,100,50)}Custom orange color{/color}
```

### Hexadecimal Colors

Use standard hex color codes:

```markdown
{color:#FF6347}Tomato red{/color}
{color:#4169E1}Royal blue{/color}
```

## Supported Colors

### Predefined Named Colors

- `black` - RGB(0, 0, 0)
- `red` - RGB(255, 0, 0)
- `green` - RGB(0, 128, 0)
- `blue` - RGB(0, 0, 255)
- `yellow` - RGB(255, 255, 0)
- `cyan` - RGB(0, 255, 255)
- `magenta` - RGB(255, 0, 255)
- `orange` - RGB(255, 165, 0)
- `purple` - RGB(128, 0, 128)
- `gray` / `grey` - RGB(128, 128, 128)
- `white` - RGB(255, 255, 255)

## Usage Examples

### In Paragraphs

```markdown
This paragraph has {red}red text{/red} and {blue}blue text{/blue}.
```

### In Headers

```markdown
### {purple}Colored Header{/purple}
### Mixed {red}Red{/red} and {green}Green{/green}
```

### In Lists

```markdown
- {red}Red item{/red}
- {green}Green item{/green}
- {blue}Blue item{/blue}
```

### Combined with Formatting

```markdown
**Bold with {red}red text{/red} inside**
{blue}Blue text with *italic*{/blue}
{green}Green with `code`{/green}
```

## Programmatic Usage

### Using the API

```go
package main

import (
    "github.com/beinux3/Mark2PDF"
)

func main() {
    markdown := `# Colored Document

    {red}Important:{/red} This is a colored document.

    {color:rgb(255,165,0)}Custom orange text{/color}
    `

    // Convert to PDF
    mark2pdf.ConvertFile("input.md", "output.pdf")
}
```

### Working with Color Types

```go
// Create colors programmatically
redColor := mark2pdf.NewColor(255, 0, 0)
customColor := mark2pdf.NewColorFloat(0.5, 0.7, 0.9)

// Use predefined colors
greenColor := mark2pdf.ColorGreen
blueColor := mark2pdf.ColorBlue
```

## Technical Implementation

### PDF Color Space

Colors are rendered using PDF's RGB color space with the `rg` operator. RGB values are normalized to the 0.0-1.0 range as required by the PDF specification.

### Color Structure

```go
type Color struct {
    R float64 // 0.0 - 1.0
    G float64 // 0.0 - 1.0
    B float64 // 0.0 - 1.0
}
```

### Color Parsing

The parser recognizes three color formats:

1. **Named colors**: Direct color names (case-insensitive)
2. **RGB format**: `rgb(r,g,b)` where r,g,b are 0-255
3. **Hex format**: `#RRGGBB` standard hex notation

### Integration

Colors are fully integrated with:
- Word wrapping
- Multi-style text rendering
- All inline formatting elements
- Headers, paragraphs, lists, and tables

## Testing

Run the color tests:

```bash
cd examples
go run test_colors.go
```

Run unit tests:

```bash
go test -v
```

## Examples

See the `examples/` directory for complete examples:

- `test_colors.md` - Comprehensive color syntax reference
- `test_colors.go` - Programmatic color usage examples
- `demo_colors.md` - Full-featured color demonstration

## Limitations

- Colors do not work inside code blocks (by design)
- Colors are not visible in standard Markdown viewers
- Color tags are case-sensitive for the tag itself but not for color names

## Future Enhancements

Potential future additions:
- Background colors
- Color themes
- CMYK color space support
- Color gradients
- Transparency/opacity support
