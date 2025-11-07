package main

import (
	"fmt"
	"log"

	"github.com/beinux3/Mark2PDF"
)

func main() {
	fmt.Println("Mark2PDF - Examples Generator")
	fmt.Println("==============================\n")

	examples := []struct {
		name     string
		markdown string
		output   string
	}{
		{
			name:   "Basic Features",
			output: "1_basic.pdf",
			markdown: `# Mark2PDF - Basic Features

## Headers
All six levels of headers with different sizes.

## Paragraphs
Regular paragraphs with automatic word wrapping. Text flows naturally across multiple lines.

## Lists

### Unordered Lists
- First item
- Second item
- Third item

### Ordered Lists
1. First item
2. Second item
3. Third item

### Task Lists
- [x] Completed task
- [ ] Pending task

## Code Blocks

` + "```go" + `
package main
import "fmt"
func main() {
    fmt.Println("Hello, World!")
}
` + "```" + `

## Blockquotes

> This is a blockquote.
> It can span multiple lines.

## Tables

| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
| Cell 1   | Cell 2   | Cell 3   |
| Data 1   | Data 2   | Data 3   |

## Horizontal Rules

---

That's it for basic features!
`,
		},
		{
			name:   "Inline Formatting",
			output: "2_inline_formatting.pdf",
			markdown: `# Inline Formatting

## Text Styles

- **Bold text**
- *Italic text*
- ` + "`Inline code`" + `
- ~~Strikethrough text~~
- [Link text](https://example.com)

## Combined Formatting

You can combine **bold with *italic*** or **bold with ` + "`code`" + `**.

Paragraphs can have **bold**, *italic*, ` + "`code`" + `, and [links](https://example.com) all mixed together.

## In Lists

- Item with **bold**
- Item with *italic*
- Item with ` + "`code`" + `
- Item with [link](https://example.com)
`,
		},
		{
			name:   "Colors",
			output: "3_colors.pdf",
			markdown: `# Color Support

## Named Colors

- {red}Red text{/red}
- {green}Green text{/green}
- {blue}Blue text{/blue}
- {yellow}Yellow text{/yellow}
- {cyan}Cyan text{/cyan}
- {magenta}Magenta text{/magenta}
- {orange}Orange text{/orange}
- {purple}Purple text{/purple}
- {gray}Gray text{/gray}

## Custom Colors

{color:rgb(255,105,180)}Hot pink using RGB{/color}

{color:#8B4513}Saddle brown using hex{/color}

## Colors with Formatting

{blue}**Bold inside blue**{/blue}

{red}*Italic inside red*{/red}

{green}Text with **bold** and *italic*{/green}

## Colors in Tables

| Status | Message | Notes |
|--------|---------|-------|
| {green}Success{/green} | All OK | Working |
| {red}Error{/red} | Failed | Check logs |
| {yellow}Warning{/yellow} | Review | Attention needed |
`,
		},
		{
			name:   "Complete Demo",
			output: "4_complete_demo.pdf",
			markdown: `# Mark2PDF Complete Demo

## Overview
This document showcases all Mark2PDF features.

---

## 1. Text Formatting

Regular text, **bold text**, *italic text*, and ` + "`inline code`" + `.

Combined: **bold with *italic*** and {blue}**colored bold**{/blue}.

## 2. Lists

### Unordered
- First item with **bold**
- {green}Colored item{/green}
- Item with ` + "`code`" + `

### Ordered
1. {red}**First**{/red}
2. *Second in italic*
3. Regular third

### Tasks
- [x] {green}Completed{/green}
- [ ] {red}Pending{/red}

## 3. Code

` + "```python" + `
def hello(name):
    print(f"Hello, {name}!")
    return True
` + "```" + `

## 4. Tables

| Feature | Status | Notes |
|---------|--------|-------|
| {blue}**Headers**{/blue} | {green}✓{/green} | 6 levels |
| {purple}*Colors*{/purple} | {green}✓{/green} | RGB & Hex |
| {orange}Tables{/orange} | {green}✓{/green} | With borders |

## 5. Blockquotes

> {blue}**Important:**{/blue} This is a blockquote with formatting.
> It supports all inline elements.

---

**That's all!** Visit [GitHub](https://github.com/beinux3/Mark2PDF) for more.
`,
		},
	}

	fmt.Println("Generating examples...\n")

	for _, ex := range examples {
		fmt.Printf("  • %s...", ex.name)
		converter := mark2pdf.NewConverter(ex.markdown)
		if err := converter.ConvertToFile(ex.output); err != nil {
			log.Printf("\n    ✗ Error: %v\n", err)
			continue
		}
		fmt.Printf(" ✓ %s\n", ex.output)
	}

	fmt.Println("\n✅ All examples generated successfully!")
	fmt.Println("\nGenerated PDFs:")
	fmt.Println("  1_basic.pdf            - Basic Markdown features")
	fmt.Println("  2_inline_formatting.pdf - Text formatting options")
	fmt.Println("  3_colors.pdf           - Color support examples")
	fmt.Println("  4_complete_demo.pdf    - Complete feature showcase")
}
