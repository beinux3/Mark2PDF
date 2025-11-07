# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- **Color Support**: Comprehensive color support for text in PDF output
  - Named colors: red, green, blue, yellow, cyan, magenta, orange, purple, gray, black, white
  - Custom RGB colors: `{color:rgb(255,100,50)}text{/color}`
  - Hexadecimal colors: `{color:#FF6347}text{/color}`
  - Alternative syntax: `{red}text{/red}` or `{color:red}text{/color}`
  - Colors work with all inline formatting (bold, italic, code)
  - Colors can be used in headers, paragraphs, lists, and **tables**
  - Added `Color` struct with RGB support (0.0-1.0 range)
  - Added `NewColor(r, g, b int)` and `NewColorFloat(r, g, b float64)` constructors
  - Extended `InlineElement` to support color information
  - Extended `TextPart` to include color
  - Updated PDF writer to render colored text using PDF `rg` operator
  - Added `writeMultiStyleTextAt()` for positioned multi-style text rendering
  - Table cells now support full inline formatting including colors

### Changed
- Updated README with comprehensive color documentation
- Added color examples to examples directory
- Enhanced inline element parsing to recognize color tags
- Extended `MarkdownElement` with `TableCellsInline` field for inline elements in table cells
- Modified `parseTable()` to parse inline elements for each cell
- Updated `renderTable()` to render inline-formatted cell content
- RGB color validation now enforces 0-255 range

### Fixed
- Table cells now properly render colored text and other inline formatting
- Invalid RGB values (outside 0-255 range) are now rejected
- **Nested formatting inside colors** - Colors can now contain bold, italic, and other formatting
  - Example: `{blue}**Bold inside blue**{/blue}` now works correctly
  - Recursive parsing of inline elements within colored text
  - Proper inheritance of colors in nested structures
- **Color rendering in word-wrapped text** - Fixed critical bug where colors were lost during word wrapping
  - Issue: `TextPart.Color` field was not being copied when creating new text parts in `writeMultiStyleTextWrapped()`
  - Fixed in [mark2pdf.go:320, 330, 338](mark2pdf.go) by ensuring `Color` field is preserved during text wrapping operations
  - Now all colored text renders correctly in PDFs regardless of line length

### Technical Details
- Colors are rendered using PDF RGB color space
- RGB values normalized to 0.0-1.0 for PDF compliance
- Color parser supports multiple formats (names, rgb(), hex) with validation
- Full integration with existing word-wrapping and multi-style text rendering
- Table rendering fully supports inline elements with proper positioning
- **Recursive inline parsing**: Added `Children []InlineElement` field to `InlineElement`
- **New function**: `convertInlineToTextParts()` handles recursive conversion with color inheritance
- Nested formatting properly combines font styles and colors

## Previous Versions

See git commit history for details on previous releases.
