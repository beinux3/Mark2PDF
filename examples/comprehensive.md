# Comprehensive Markdown Test

This document tests all supported Markdown features in Mark2PDF.

## Headers

### Level 3 Header
#### Level 4 Header
##### Level 5 Header
###### Level 6 Header

## Inline Formatting

This paragraph contains **bold text**, *italic text*, `inline code`, and ~~strikethrough text~~.

You can also combine **bold and *italic* together**.

## Links and Images

Here's a [link to GitHub](https://github.com).

Here's an image reference: ![Alt text](https://example.com/image.png)

## Lists

### Unordered Lists

- First item
- Second item
- Third item
  - Nested item (not fully supported yet)
- Fourth item

### Ordered Lists

1. First step
2. Second step
3. Third step
4. Fourth step

### Task Lists

- [ ] Incomplete task
- [x] Completed task
- [ ] Another incomplete task

## Code Blocks

### Fenced code block with language

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

### Fenced code block without language

```
Simple code block
No syntax highlighting
```

### Python example

```python
def hello():
    print("Hello from Python!")
    return True
```

## Blockquotes

> This is a blockquote.
> It can span multiple lines.

> Another blockquote here.

## Horizontal Rules

---

Text after first rule.

***

Text after second rule.

___

Text after third rule.

## Tables

| Header 1 | Header 2 | Header 3 |
|----------|----------|----------|
| Row 1 Col 1 | Row 1 Col 2 | Row 1 Col 3 |
| Row 2 Col 1 | Row 2 Col 2 | Row 2 Col 3 |

| Left Aligned | Center Aligned | Right Aligned |
|:-------------|:---------------:|--------------:|
| Left | Center | Right |
| Data | More Data | Even More |

## Mixed Content

Here's a paragraph with **bold**, *italic*, and `code` all together. It also has a [link](https://example.com) and continues with more text.

### Lists with Inline Formatting

- Item with **bold text**
- Item with *italic text*
- Item with `inline code`
- Item with [a link](https://example.com)

### Code with Description

Here's how to define a function in Go:

```go
func add(a, b int) int {
    return a + b
}
```

The function above adds two integers and returns the result.

## Conclusion

This document demonstrates all the Markdown features currently supported by Mark2PDF library.
