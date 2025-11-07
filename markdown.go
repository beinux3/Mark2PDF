package mark2pdf

import (
	"regexp"
	"strings"
)

// MarkdownElement rappresenta un elemento del markdown parsato
type MarkdownElement struct {
	Type         string            // "h1", "h2", "h3", "h4", "h5", "h6", "p", "code", "list", "hr", "blockquote", "table"
	Content      string            // Il contenuto testuale
	Level        int               // Per liste e headers
	Items        []string          // Per liste (raw content)
	ItemChildren [][]InlineElement // Inline elements per ogni item della lista
	Ordered      bool              // Se la lista è ordinata
	Language     string            // Per blocchi di codice
	TableRows    [][]string        // Per tabelle
	TableAlign   []string          // Allineamento colonne tabella
	Children     []InlineElement   // Elementi inline (bold, italic, code, link)
}

// InlineElement rappresenta elementi inline nel testo
type InlineElement struct {
	Type    string // "text", "bold", "italic", "code", "link", "image", "strikethrough"
	Content string
	URL     string // Per link e immagini
	Alt     string // Per immagini
}

// MarkdownParser parsea il markdown in elementi
type MarkdownParser struct {
	lines []string
}

// NewMarkdownParser crea un nuovo parser
func NewMarkdownParser(markdown string) *MarkdownParser {
	lines := strings.Split(markdown, "\n")
	return &MarkdownParser{lines: lines}
}

// Parse parsea il markdown e restituisce una lista di elementi
func (mp *MarkdownParser) Parse() []MarkdownElement {
	elements := make([]MarkdownElement, 0)
	i := 0

	for i < len(mp.lines) {
		line := mp.lines[i]
		trimmed := strings.TrimSpace(line)

		// Linea vuota
		if trimmed == "" {
			i++
			continue
		}

		// Headers ATX style (# ## ### etc)
		if strings.HasPrefix(trimmed, "#") {
			elem, consumed := mp.parseHeader(i)
			if consumed > 0 {
				elements = append(elements, elem)
				i += consumed
				continue
			}
		}

		// Headers Setext style (underlined with = or -)
		if i+1 < len(mp.lines) {
			nextLine := strings.TrimSpace(mp.lines[i+1])
			if isSetextHeader(nextLine) {
				level := 1
				if strings.HasPrefix(nextLine, "-") {
					level = 2
				}
				elements = append(elements, MarkdownElement{
					Type:     "h" + string(rune('0'+level)),
					Content:  trimmed,
					Level:    level,
					Children: mp.parseInline(trimmed),
				})
				i += 2
				continue
			}
		}

		// Horizontal rule
		if isHorizontalRule(trimmed) {
			elements = append(elements, MarkdownElement{
				Type: "hr",
			})
			i++
			continue
		}

		// Code block (fenced with ``` or ~~~)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			elem, consumed := mp.parseCodeBlock(i)
			elements = append(elements, elem)
			i += consumed
			continue
		}

		// Code block (indented with 4 spaces or tab)
		if len(line) > 0 && (strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t")) {
			elem, consumed := mp.parseIndentedCodeBlock(i)
			elements = append(elements, elem)
			i += consumed
			continue
		}

		// Blockquote
		if strings.HasPrefix(trimmed, ">") {
			elem, consumed := mp.parseBlockquote(i)
			elements = append(elements, elem)
			i += consumed
			continue
		}

		// Table
		if i+1 < len(mp.lines) && isTableSeparator(strings.TrimSpace(mp.lines[i+1])) {
			elem, consumed := mp.parseTable(i)
			elements = append(elements, elem)
			i += consumed
			continue
		}

		// Unordered list
		if isUnorderedListItem(trimmed) {
			elem, consumed := mp.parseList(i, false)
			elements = append(elements, elem)
			i += consumed
			continue
		}

		// Ordered list
		if isOrderedListItem(trimmed) {
			elem, consumed := mp.parseList(i, true)
			elements = append(elements, elem)
			i += consumed
			continue
		}

		// Task list
		if isTaskListItem(trimmed) {
			elem, consumed := mp.parseTaskList(i)
			elements = append(elements, elem)
			i += consumed
			continue
		}

		// Paragraph - raggruppa linee consecutive non vuote
		elem, consumed := mp.parseParagraph(i)
		elements = append(elements, elem)
		i += consumed
	}

	return elements
}

// parseHeader parsea un header ATX style
func (mp *MarkdownParser) parseHeader(startIdx int) (MarkdownElement, int) {
	line := strings.TrimSpace(mp.lines[startIdx])
	level := 0

	for _, ch := range line {
		if ch == '#' {
			level++
		} else {
			break
		}
	}

	if level > 0 && level <= 6 {
		content := strings.TrimSpace(line[level:])
		// Remove trailing #
		content = strings.TrimRight(content, "#")
		content = strings.TrimSpace(content)

		return MarkdownElement{
			Type:     "h" + string(rune('0'+level)),
			Content:  content,
			Level:    level,
			Children: mp.parseInline(content),
		}, 1
	}

	return MarkdownElement{}, 0
}

// parseCodeBlock parsea un blocco di codice fenced
func (mp *MarkdownParser) parseCodeBlock(startIdx int) (MarkdownElement, int) {
	line := strings.TrimSpace(mp.lines[startIdx])
	fence := "```"
	if strings.HasPrefix(line, "~~~") {
		fence = "~~~"
	}

	// Extract language
	language := strings.TrimSpace(strings.TrimPrefix(line, fence))

	codeLines := make([]string, 0)
	i := startIdx + 1

	for i < len(mp.lines) {
		if strings.HasPrefix(strings.TrimSpace(mp.lines[i]), fence) {
			i++
			break
		}
		codeLines = append(codeLines, mp.lines[i])
		i++
	}

	return MarkdownElement{
		Type:     "code",
		Content:  strings.Join(codeLines, "\n"),
		Language: language,
	}, i - startIdx
}

// parseIndentedCodeBlock parsea un blocco di codice indentato
func (mp *MarkdownParser) parseIndentedCodeBlock(startIdx int) (MarkdownElement, int) {
	codeLines := make([]string, 0)
	i := startIdx

	for i < len(mp.lines) {
		line := mp.lines[i]
		if len(line) == 0 || strings.TrimSpace(line) == "" {
			// Empty line in code block
			codeLines = append(codeLines, "")
			i++
			continue
		}

		if strings.HasPrefix(line, "    ") {
			codeLines = append(codeLines, line[4:])
			i++
		} else if strings.HasPrefix(line, "\t") {
			codeLines = append(codeLines, line[1:])
			i++
		} else {
			break
		}
	}

	return MarkdownElement{
		Type:    "code",
		Content: strings.Join(codeLines, "\n"),
	}, i - startIdx
}

// parseBlockquote parsea una blockquote
func (mp *MarkdownParser) parseBlockquote(startIdx int) (MarkdownElement, int) {
	quoteLines := make([]string, 0)
	i := startIdx

	for i < len(mp.lines) {
		line := strings.TrimSpace(mp.lines[i])
		if strings.HasPrefix(line, ">") {
			content := strings.TrimSpace(strings.TrimPrefix(line, ">"))
			quoteLines = append(quoteLines, content)
			i++
		} else if line == "" && i+1 < len(mp.lines) && strings.HasPrefix(strings.TrimSpace(mp.lines[i+1]), ">") {
			quoteLines = append(quoteLines, "")
			i++
		} else {
			break
		}
	}

	content := strings.Join(quoteLines, " ")
	return MarkdownElement{
		Type:     "blockquote",
		Content:  content,
		Children: mp.parseInline(content),
	}, i - startIdx
}

// parseList parsea una lista (ordinata o non ordinata)
func (mp *MarkdownParser) parseList(startIdx int, ordered bool) (MarkdownElement, int) {
	items := make([]string, 0)
	itemChildren := make([][]InlineElement, 0)
	i := startIdx

	for i < len(mp.lines) {
		line := strings.TrimSpace(mp.lines[i])

		if line == "" {
			if i+1 < len(mp.lines) {
				nextLine := strings.TrimSpace(mp.lines[i+1])
				if (ordered && isOrderedListItem(nextLine)) || (!ordered && isUnorderedListItem(nextLine)) {
					i++
					continue
				}
			}
			break
		}

		if ordered && isOrderedListItem(line) {
			dotIdx := strings.Index(line, ".")
			if dotIdx > 0 && dotIdx < len(line)-1 {
				content := strings.TrimSpace(line[dotIdx+1:])
				items = append(items, content)
				itemChildren = append(itemChildren, mp.parseInline(content))
			}
			i++
		} else if !ordered && isUnorderedListItem(line) {
			// Remove marker (-, *, +)
			content := line[1:]
			if len(content) > 0 && content[0] == ' ' {
				content = content[1:]
			}
			content = strings.TrimSpace(content)
			items = append(items, content)
			itemChildren = append(itemChildren, mp.parseInline(content))
			i++
		} else {
			break
		}
	}

	listType := "list"
	if ordered {
		listType = "ordered-list"
	}

	return MarkdownElement{
		Type:         listType,
		Items:        items,
		ItemChildren: itemChildren,
		Ordered:      ordered,
	}, i - startIdx
}

// parseTaskList parsea una task list
func (mp *MarkdownParser) parseTaskList(startIdx int) (MarkdownElement, int) {
	items := make([]string, 0)
	itemChildren := make([][]InlineElement, 0)
	i := startIdx

	for i < len(mp.lines) {
		line := strings.TrimSpace(mp.lines[i])
		if line == "" || !isTaskListItem(line) {
			break
		}

		// Extract task text (after [ ] or [x])
		taskRegex := regexp.MustCompile(`^[-*+]\s+\[([ xX])\]\s+(.+)$`)
		matches := taskRegex.FindStringSubmatch(line)
		if len(matches) == 3 {
			checked := strings.ToLower(matches[1]) == "x"
			text := matches[2]
			displayText := ""
			if checked {
				displayText = "[x] " + text
			} else {
				displayText = "[ ] " + text
			}
			items = append(items, displayText)
			itemChildren = append(itemChildren, mp.parseInline(text))
		}
		i++
	}

	return MarkdownElement{
		Type:         "task-list",
		Items:        items,
		ItemChildren: itemChildren,
	}, i - startIdx
}

// parseTable parsea una tabella
func (mp *MarkdownParser) parseTable(startIdx int) (MarkdownElement, int) {
	rows := make([][]string, 0)
	align := make([]string, 0)

	// Header row
	headerLine := strings.TrimSpace(mp.lines[startIdx])
	headerCells := parseTableRow(headerLine)
	rows = append(rows, headerCells)

	// Separator row
	if startIdx+1 < len(mp.lines) {
		sepLine := strings.TrimSpace(mp.lines[startIdx+1])
		align = parseTableAlignment(sepLine)
	}

	// Data rows
	i := startIdx + 2
	for i < len(mp.lines) {
		line := strings.TrimSpace(mp.lines[i])
		if line == "" || !strings.Contains(line, "|") {
			break
		}
		cells := parseTableRow(line)
		rows = append(rows, cells)
		i++
	}

	return MarkdownElement{
		Type:       "table",
		TableRows:  rows,
		TableAlign: align,
	}, i - startIdx
}

// parseParagraph parsea un paragrafo
func (mp *MarkdownParser) parseParagraph(startIdx int) (MarkdownElement, int) {
	paragraphLines := make([]string, 0)
	i := startIdx

	for i < len(mp.lines) {
		line := strings.TrimSpace(mp.lines[i])
		if line == "" {
			break
		}

		// Check if next line is a special element
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "```") ||
			strings.HasPrefix(line, "~~~") || isUnorderedListItem(line) ||
			strings.HasPrefix(line, "> ") || isHorizontalRule(line) ||
			isOrderedListItem(line) || isTaskListItem(line) {
			break
		}

		paragraphLines = append(paragraphLines, line)
		i++
	}

	content := strings.Join(paragraphLines, " ")
	return MarkdownElement{
		Type:     "p",
		Content:  content,
		Children: mp.parseInline(content),
	}, i - startIdx
}

// parseInline parsea elementi inline (bold, italic, code, link, etc)
func (mp *MarkdownParser) parseInline(text string) []InlineElement {
	elements := make([]InlineElement, 0)
	current := ""
	i := 0

	for i < len(text) {
		// Bold **text** or __text__
		if i+1 < len(text) && (text[i:i+2] == "**" || text[i:i+2] == "__") {
			if current != "" {
				elements = append(elements, InlineElement{Type: "text", Content: current})
				current = ""
			}

			delimiter := text[i : i+2]
			end := strings.Index(text[i+2:], delimiter)
			if end != -1 {
				boldText := text[i+2 : i+2+end]
				elements = append(elements, InlineElement{Type: "bold", Content: boldText})
				i += 2 + end + 2
				continue
			}
		}

		// Italic *text* or _text_
		if text[i] == '*' || text[i] == '_' {
			if current != "" {
				elements = append(elements, InlineElement{Type: "text", Content: current})
				current = ""
			}

			delimiter := string(text[i])
			end := strings.Index(text[i+1:], delimiter)
			if end != -1 && end > 0 {
				italicText := text[i+1 : i+1+end]
				// Check it's not part of bold
				if !(i > 0 && text[i-1] == text[i]) && !(i+1+end+1 < len(text) && text[i+1+end+1] == text[i]) {
					elements = append(elements, InlineElement{Type: "italic", Content: italicText})
					i += 1 + end + 1
					continue
				}
			}
		}

		// Strikethrough ~~text~~
		if i+1 < len(text) && text[i:i+2] == "~~" {
			if current != "" {
				elements = append(elements, InlineElement{Type: "text", Content: current})
				current = ""
			}

			end := strings.Index(text[i+2:], "~~")
			if end != -1 {
				strikeText := text[i+2 : i+2+end]
				elements = append(elements, InlineElement{Type: "strikethrough", Content: strikeText})
				i += 2 + end + 2
				continue
			}
		}

		// Inline code `text`
		if text[i] == '`' {
			if current != "" {
				elements = append(elements, InlineElement{Type: "text", Content: current})
				current = ""
			}

			end := strings.Index(text[i+1:], "`")
			if end != -1 {
				codeText := text[i+1 : i+1+end]
				elements = append(elements, InlineElement{Type: "code", Content: codeText})
				i += 1 + end + 1
				continue
			}
		}

		// Link [text](url) or Image ![alt](url)
		if text[i] == '[' || (text[i] == '!' && i+1 < len(text) && text[i+1] == '[') {
			if current != "" {
				elements = append(elements, InlineElement{Type: "text", Content: current})
				current = ""
			}

			isImage := text[i] == '!'
			startPos := i
			if isImage {
				startPos++
			}

			closeBracket := strings.Index(text[startPos:], "](")
			if closeBracket != -1 {
				closeParen := strings.Index(text[startPos+closeBracket+2:], ")")
				if closeParen != -1 {
					linkText := text[startPos+1 : startPos+closeBracket]
					url := text[startPos+closeBracket+2 : startPos+closeBracket+2+closeParen]

					if isImage {
						elements = append(elements, InlineElement{
							Type: "image",
							Alt:  linkText,
							URL:  url,
						})
						i = startPos + closeBracket + 2 + closeParen + 1
					} else {
						elements = append(elements, InlineElement{
							Type:    "link",
							Content: linkText,
							URL:     url,
						})
						i = startPos + closeBracket + 2 + closeParen + 1
					}
					continue
				}
			}
		}

		current += string(text[i])
		i++
	}

	if current != "" {
		elements = append(elements, InlineElement{Type: "text", Content: current})
	}

	return elements
}

// Helper functions

func isSetextHeader(line string) bool {
	if len(line) < 1 {
		return false
	}
	return (strings.Trim(line, "=") == "" && strings.Contains(line, "=")) ||
		(strings.Trim(line, "-") == "" && strings.Contains(line, "-"))
}

func isHorizontalRule(line string) bool {
	line = strings.TrimSpace(line)
	if len(line) < 3 {
		return false
	}

	// Remove spaces
	cleaned := strings.ReplaceAll(line, " ", "")

	// Check for ---, ***, ___
	if len(cleaned) < 3 {
		return false
	}

	char := cleaned[0]
	if char != '-' && char != '*' && char != '_' {
		return false
	}

	for _, c := range cleaned {
		if c != rune(char) {
			return false
		}
	}

	return len(cleaned) >= 3
}

func isUnorderedListItem(line string) bool {
	line = strings.TrimSpace(line)
	return len(line) >= 2 && (line[0] == '-' || line[0] == '*' || line[0] == '+') &&
		(line[1] == ' ' || line[1] == '\t')
}

func isOrderedListItem(line string) bool {
	line = strings.TrimSpace(line)
	if len(line) < 3 {
		return false
	}

	dotIdx := strings.Index(line, ".")
	if dotIdx <= 0 || dotIdx >= len(line)-1 {
		return false
	}

	numberPart := line[:dotIdx]
	for _, ch := range numberPart {
		if ch < '0' || ch > '9' {
			return false
		}
	}

	return line[dotIdx+1] == ' ' || line[dotIdx+1] == '\t'
}

func isTaskListItem(line string) bool {
	line = strings.TrimSpace(line)
	taskRegex := regexp.MustCompile(`^[-*+]\s+\[([ xX])\]\s+.+$`)
	return taskRegex.MatchString(line)
}

func isTableSeparator(line string) bool {
	if !strings.Contains(line, "|") {
		return false
	}
	// Remove spaces and check if it contains only |, -, and :
	cleaned := strings.ReplaceAll(line, " ", "")
	for _, ch := range cleaned {
		if ch != '|' && ch != '-' && ch != ':' {
			return false
		}
	}
	return strings.Count(cleaned, "-") >= 1
}

func parseTableRow(line string) []string {
	// Remove leading and trailing |
	line = strings.Trim(strings.TrimSpace(line), "|")
	cells := strings.Split(line, "|")

	for i := range cells {
		cells[i] = strings.TrimSpace(cells[i])
	}

	return cells
}

func parseTableAlignment(line string) []string {
	cells := parseTableRow(line)
	align := make([]string, len(cells))

	for i, cell := range cells {
		hasLeft := strings.HasPrefix(cell, ":")
		hasRight := strings.HasSuffix(cell, ":")

		if hasLeft && hasRight {
			align[i] = "center"
		} else if hasRight {
			align[i] = "right"
		} else {
			align[i] = "left"
		}
	}

	return align
}
