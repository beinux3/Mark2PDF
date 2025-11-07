package mark2pdf

import (
	"strings"
)

// MarkdownElement rappresenta un elemento del markdown parsato
type MarkdownElement struct {
	Type    string   // "h1", "h2", "h3", "h4", "h5", "h6", "p", "code", "list", "hr", "blockquote"
	Content string   // Il contenuto testuale
	Level   int      // Per liste e headers
	Items   []string // Per liste
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
			level := 0
			for _, ch := range trimmed {
				if ch == '#' {
					level++
				} else {
					break
				}
			}
			if level > 0 && level <= 6 {
				content := strings.TrimSpace(trimmed[level:])
				elements = append(elements, MarkdownElement{
					Type:    "h" + string(rune('0'+level)),
					Content: content,
					Level:   level,
				})
				i++
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

		// Code block (fenced with ```)
		if strings.HasPrefix(trimmed, "```") {
			codeLines := make([]string, 0)
			i++ // Skip opening ```
			for i < len(mp.lines) {
				if strings.HasPrefix(strings.TrimSpace(mp.lines[i]), "```") {
					i++ // Skip closing ```
					break
				}
				codeLines = append(codeLines, mp.lines[i])
				i++
			}
			elements = append(elements, MarkdownElement{
				Type:    "code",
				Content: strings.Join(codeLines, "\n"),
			})
			continue
		}

		// Blockquote
		if strings.HasPrefix(trimmed, ">") {
			quoteLines := make([]string, 0)
			for i < len(mp.lines) {
				l := strings.TrimSpace(mp.lines[i])
				if strings.HasPrefix(l, ">") {
					quoteLines = append(quoteLines, strings.TrimSpace(l[1:]))
					i++
				} else if l == "" {
					i++
				} else {
					break
				}
			}
			elements = append(elements, MarkdownElement{
				Type:    "blockquote",
				Content: strings.Join(quoteLines, " "),
			})
			continue
		}

		// Unordered list
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") || strings.HasPrefix(trimmed, "+ ") {
			items := make([]string, 0)
			for i < len(mp.lines) {
				l := strings.TrimSpace(mp.lines[i])
				if strings.HasPrefix(l, "- ") || strings.HasPrefix(l, "* ") || strings.HasPrefix(l, "+ ") {
					items = append(items, strings.TrimSpace(l[2:]))
					i++
				} else if l == "" {
					i++
					break
				} else {
					break
				}
			}
			elements = append(elements, MarkdownElement{
				Type:  "list",
				Items: items,
			})
			continue
		}

		// Ordered list
		if isOrderedListItem(trimmed) {
			items := make([]string, 0)
			for i < len(mp.lines) {
				l := strings.TrimSpace(mp.lines[i])
				if isOrderedListItem(l) {
					// Remove number and dot
					dotIdx := strings.Index(l, ".")
					if dotIdx > 0 {
						items = append(items, strings.TrimSpace(l[dotIdx+1:]))
					}
					i++
				} else if l == "" {
					i++
					break
				} else {
					break
				}
			}
			elements = append(elements, MarkdownElement{
				Type:  "ordered-list",
				Items: items,
			})
			continue
		}

		// Paragraph - raggruppa linee consecutive non vuote
		paragraphLines := make([]string, 0)
		for i < len(mp.lines) {
			l := strings.TrimSpace(mp.lines[i])
			if l == "" {
				break
			}
			// Check if next line is a special element
			if strings.HasPrefix(l, "#") || strings.HasPrefix(l, "```") ||
				strings.HasPrefix(l, "- ") || strings.HasPrefix(l, "* ") ||
				strings.HasPrefix(l, "> ") || isHorizontalRule(l) ||
				isOrderedListItem(l) {
				break
			}
			paragraphLines = append(paragraphLines, l)
			i++
		}

		if len(paragraphLines) > 0 {
			elements = append(elements, MarkdownElement{
				Type:    "p",
				Content: strings.Join(paragraphLines, " "),
			})
		}
	}

	return elements
}

// isHorizontalRule verifica se una linea è un separatore orizzontale
func isHorizontalRule(line string) bool {
	line = strings.TrimSpace(line)
	if len(line) < 3 {
		return false
	}

	// Conta caratteri -, *, _
	dashCount := strings.Count(line, "-")
	starCount := strings.Count(line, "*")
	underCount := strings.Count(line, "_")

	// Remove spaces
	cleaned := strings.ReplaceAll(line, " ", "")

	if dashCount >= 3 && cleaned == strings.Repeat("-", dashCount) {
		return true
	}
	if starCount >= 3 && cleaned == strings.Repeat("*", starCount) {
		return true
	}
	if underCount >= 3 && cleaned == strings.Repeat("_", underCount) {
		return true
	}

	return false
}

// isOrderedListItem verifica se una linea inizia con un numero seguito da un punto
func isOrderedListItem(line string) bool {
	line = strings.TrimSpace(line)
	if len(line) < 3 {
		return false
	}

	dotIdx := strings.Index(line, ".")
	if dotIdx <= 0 {
		return false
	}

	numberPart := line[:dotIdx]
	for _, ch := range numberPart {
		if ch < '0' || ch > '9' {
			return false
		}
	}

	return true
}

// stripInlineFormatting rimuove la formattazione inline (**, __, `, etc)
// In una implementazione completa, questo gestirebbe bold, italic, code inline, link
func stripInlineFormatting(text string) string {
	// Remove bold **text**
	text = strings.ReplaceAll(text, "**", "")
	text = strings.ReplaceAll(text, "__", "")

	// Remove italic *text* (attenzione a non rimuovere liste)
	// Implementazione semplificata

	// Remove inline code `text`
	for strings.Contains(text, "`") {
		start := strings.Index(text, "`")
		end := strings.Index(text[start+1:], "`")
		if end == -1 {
			break
		}
		text = text[:start] + text[start+1:start+1+end] + text[start+1+end+1:]
	}

	// Remove links [text](url)
	for strings.Contains(text, "](") {
		closeBracketPos := strings.Index(text, "](")
		if closeBracketPos == -1 {
			break
		}

		linkStart := strings.LastIndex(text[:closeBracketPos], "[")
		if linkStart == -1 {
			break
		}

		linkEnd := strings.Index(text[closeBracketPos:], ")")
		if linkEnd == -1 {
			break
		}
		linkEnd += closeBracketPos

		linkText := text[linkStart+1 : closeBracketPos]
		text = text[:linkStart] + linkText + text[linkEnd+1:]
	}

	return text
}
