package main

import (
	"fmt"
	"log"

	"github.com/beinux3/Mark2PDF"
)

func main() {
	fmt.Println("Starting Mark2PDF test suite...")
	fmt.Println("=====================================")

	// Test 1: Simple document
	test1Simple()

	// Test 2: Inline formatting
	test2InlineFormatting()

	// Test 3: Lists with inline formatting
	test3Lists()

	// Test 4: Code blocks
	test4CodeBlocks()

	// Test 5: Blockquotes
	test5Blockquotes()

	// Test 6: Tables
	test6Tables()

	// Test 7: Horizontal rules
	test7HorizontalRules()

	// Test 8: Comprehensive document
	test8Comprehensive()

	fmt.Println("=====================================")
	fmt.Println("All tests completed successfully!")
}

func test1Simple() {
	fmt.Println("\n[Test 1] Simple document...")

	markdown := `# Simple Document

This is a very simple Markdown document.

It contains just a header and some paragraphs.

## Second Header

Another paragraph here.`

	converter := mark2pdf.NewConverter(markdown)
	err := converter.ConvertToFile("test1_simple.pdf")
	if err != nil {
		log.Fatal("Error in test 1:", err)
	}

	fmt.Println("  ✓ Created: test1_simple.pdf")
}

func test2InlineFormatting() {
	fmt.Println("\n[Test 2] Inline formatting...")

	markdown := `# Inline Formatting Test

This document tests all inline formatting options.

## Text Styles

This paragraph contains **bold text**, *italic text*, and ` + "`" + `inline code` + "`" + `.

You can also use ~~strikethrough~~ text.

## Combined Formatting

Here we have **bold and *italic* together**.

You can also have **bold with ` + "`" + `code` + "`" + `** inside.

## Links and Images

Check out this [link to GitHub](https://github.com).

Here's an image: ![Alt text](https://example.com/image.png)`

	converter := mark2pdf.NewConverter(markdown)
	err := converter.ConvertToFile("test2_inline_formatting.pdf")
	if err != nil {
		log.Fatal("Error in test 2:", err)
	}

	fmt.Println("  ✓ Created: test2_inline_formatting.pdf")
}

func test3Lists() {
	fmt.Println("\n[Test 3] Lists with inline formatting...")

	markdown := `# Lists Test

## Unordered Lists

- First item with **bold**
- Second item with *italic*
- Third item with ` + "`" + `code` + "`" + `
- Fourth item with [a link](https://example.com)
- Fifth item with **bold**, *italic*, and ` + "`" + `code` + "`" + ` combined

## Ordered Lists

1. First step with **important** keyword
2. Second step with *emphasis*
3. Third step with ` + "`" + `command` + "`" + `
4. Fourth step with [documentation link](https://docs.example.com)

## Task Lists

- [ ] Incomplete task with **bold**
- [x] Completed task with *italic*
- [ ] Another task with ` + "`" + `inline code` + "`" + `
- [x] Task with [link](https://example.com)`

	converter := mark2pdf.NewConverter(markdown)
	err := converter.ConvertToFile("test3_lists.pdf")
	if err != nil {
		log.Fatal("Error in test 3:", err)
	}

	fmt.Println("  ✓ Created: test3_lists.pdf")
}

func test4CodeBlocks() {
	fmt.Println("\n[Test 4] Code blocks...")

	markdown := `# Code Blocks Test

## Go Code

` + "```go" + `
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
    result := add(5, 3)
    fmt.Printf("Result: %d\n", result)
}

func add(a, b int) int {
    return a + b
}
` + "```" + `

## Python Code

` + "```python" + `
def hello(name):
    """Greet someone by name."""
    print(f"Hello, {name}!")
    return True

if __name__ == "__main__":
    hello("World")
    numbers = [1, 2, 3, 4, 5]
    total = sum(numbers)
    print(f"Total: {total}")
` + "```" + `

## Plain Code Block

` + "```" + `
This is a plain code block
without language specification.
It can contain any text.
` + "```" + ``

	converter := mark2pdf.NewConverter(markdown)
	err := converter.ConvertToFile("test4_code_blocks.pdf")
	if err != nil {
		log.Fatal("Error in test 4:", err)
	}

	fmt.Println("  ✓ Created: test4_code_blocks.pdf")
}

func test5Blockquotes() {
	fmt.Println("\n[Test 5] Blockquotes...")

	markdown := `# Blockquotes Test

## Simple Blockquote

> This is a simple blockquote.
> It can span multiple lines.

## Blockquote with Formatting

> This blockquote contains **bold text** and *italic text*.
> It also has ` + "`" + `inline code` + "`" + ` and [a link](https://example.com).

## Multiple Blockquotes

> First quote here.

Some normal text in between.

> Second quote here with more content.
> This one is also multi-line.`

	converter := mark2pdf.NewConverter(markdown)
	err := converter.ConvertToFile("test5_blockquotes.pdf")
	if err != nil {
		log.Fatal("Error in test 5:", err)
	}

	fmt.Println("  ✓ Created: test5_blockquotes.pdf")
}

func test6Tables() {
	fmt.Println("\n[Test 6] Tables...")

	markdown := `# Tables Test

## Simple Table

| Name | Age | City |
|------|-----|------|
| John | 25  | NYC  |
| Jane | 30  | LA   |
| Bob  | 35  | SF   |

## Table with Alignment

| Left Aligned | Center Aligned | Right Aligned |
|:-------------|:--------------:|--------------:|
| Left         | Center         | Right         |
| Text         | More Text      | Even More     |
| Data         | Information    | Content       |

## Another Table

| Feature | Supported | Notes |
|---------|-----------|-------|
| Headers | Yes | H1-H6 |
| Lists | Yes | Ordered and unordered |
| Code | Yes | Inline and blocks |
| Tables | Yes | With alignment |`

	converter := mark2pdf.NewConverter(markdown)
	err := converter.ConvertToFile("test6_tables.pdf")
	if err != nil {
		log.Fatal("Error in test 6:", err)
	}

	fmt.Println("  ✓ Created: test6_tables.pdf")
}

func test7HorizontalRules() {
	fmt.Println("\n[Test 7] Horizontal rules...")

	markdown := `# Horizontal Rules Test

This is the first section.

---

This is the second section after a horizontal rule.

***

This is the third section after another type of rule.

___

This is the fourth section.

## Using Rules for Separation

Rules are useful for separating different topics.

---

Like this!`

	converter := mark2pdf.NewConverter(markdown)
	err := converter.ConvertToFile("test7_horizontal_rules.pdf")
	if err != nil {
		log.Fatal("Error in test 7:", err)
	}

	fmt.Println("  ✓ Created: test7_horizontal_rules.pdf")
}

func test8Comprehensive() {
	fmt.Println("\n[Test 8] Comprehensive document...")

	markdown := `# Comprehensive Markdown Test

This document demonstrates **all** supported features in Mark2PDF.

## Headers

### Level 3 Header
#### Level 4 Header
##### Level 5 Header
###### Level 6 Header

## Inline Formatting

This paragraph contains **bold text**, *italic text*, ` + "`" + `inline code` + "`" + `, and ~~strikethrough text~~.

You can also combine **bold and *italic* together**.

## Links and Images

Here's a [link to GitHub](https://github.com).

Here's an image reference: ![Alt text](https://example.com/image.png)

## Lists

### Unordered Lists with Formatting

- Item with **bold text**
- Item with *italic text*
- Item with ` + "`" + `inline code` + "`" + `
- Item with [a link](https://example.com)

### Ordered Lists

1. First step with **bold**
2. Second step with *italic*
3. Third step with ` + "`" + `code` + "`" + `
4. Fourth step complete

### Task Lists

- [ ] Incomplete task with **bold**
- [x] Completed task with *italic*
- [ ] Another task with ` + "`" + `code` + "`" + `

## Code Blocks

` + "```go" + `
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
` + "```" + `

## Blockquotes

> This is a blockquote with **formatting**.
> It can span multiple lines.

## Tables

| Header 1 | Header 2 | Header 3 |
|----------|----------|----------|
| Row 1 Col 1 | Row 1 Col 2 | Row 1 Col 3 |
| Row 2 Col 1 | Row 2 Col 2 | Row 2 Col 3 |

## Horizontal Rules

---

## Mixed Content

Here's a paragraph with **bold**, *italic*, and ` + "`" + `code` + "`" + ` all together. It also has a [link](https://example.com) and continues with more text to demonstrate word wrapping in the PDF output.

---

## Conclusion

This document demonstrates all the Markdown features currently supported by the Mark2PDF library. The PDF should render all these elements correctly with proper formatting and spacing.`

	converter := mark2pdf.NewConverter(markdown)
	err := converter.ConvertToFile("test8_comprehensive.pdf")
	if err != nil {
		log.Fatal("Error in test 8:", err)
	}

	fmt.Println("  ✓ Created: test8_comprehensive.pdf")
}
