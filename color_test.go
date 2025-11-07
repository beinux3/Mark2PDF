package mark2pdf

import (
	"testing"
)

func TestParseColorName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *Color
	}{
		{"red color", "red", &ColorRed},
		{"blue color", "blue", &ColorBlue},
		{"green color", "green", &ColorGreen},
		{"case insensitive", "RED", &ColorRed},
		{"with spaces", " blue ", &ColorBlue},
		{"rgb format", "rgb(255,0,0)", &Color{R: 1.0, G: 0.0, B: 0.0}},
		{"rgb with spaces", "rgb(255, 128, 0)", &Color{R: 1.0, G: 128.0 / 255.0, B: 0.0}},
		{"hex format", "#FF0000", &Color{R: 1.0, G: 0.0, B: 0.0}},
		{"hex lowercase", "#ff0000", &Color{R: 1.0, G: 0.0, B: 0.0}},
		{"invalid color", "notacolor", nil},
		{"invalid rgb", "rgb(300,0,0)", nil},
		{"invalid hex", "#GG0000", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseColorName(tt.input)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("Expected nil, got %v", result)
				}
			} else {
				if result == nil {
					t.Errorf("Expected color, got nil")
					return
				}
				// Compare with small epsilon for floating point
				epsilon := 0.01
				if abs(result.R-tt.expected.R) > epsilon ||
					abs(result.G-tt.expected.G) > epsilon ||
					abs(result.B-tt.expected.B) > epsilon {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func TestColorInlineElements(t *testing.T) {
	parser := NewMarkdownParser("{red}red text{/red} and {blue}blue text{/blue}")
	elements := parser.Parse()

	if len(elements) != 1 {
		t.Fatalf("Expected 1 element, got %d", len(elements))
	}

	if elements[0].Type != "p" {
		t.Errorf("Expected paragraph, got %s", elements[0].Type)
	}

	children := elements[0].Children
	if len(children) != 3 { // red + text + blue
		t.Fatalf("Expected 3 children, got %d", len(children))
	}

	// First should be colored red
	if children[0].Type != "color" {
		t.Errorf("Expected color element, got %s", children[0].Type)
	}
	if children[0].Content != "red text" {
		t.Errorf("Expected 'red text', got '%s'", children[0].Content)
	}
	if children[0].Color == nil || children[0].Color.R != 1.0 {
		t.Error("Expected red color")
	}

	// Second should be plain text " and "
	if children[1].Type != "text" {
		t.Errorf("Expected text element, got %s", children[1].Type)
	}

	// Third should be colored blue
	if children[2].Type != "color" {
		t.Errorf("Expected color element, got %s", children[2].Type)
	}
	if children[2].Content != "blue text" {
		t.Errorf("Expected 'blue text', got '%s'", children[2].Content)
	}
	if children[2].Color == nil || children[2].Color.B != 1.0 {
		t.Error("Expected blue color")
	}
}

func TestNewColor(t *testing.T) {
	color := NewColor(255, 128, 0)

	if color.R != 1.0 {
		t.Errorf("Expected R=1.0, got %f", color.R)
	}
	if abs(color.G-0.502) > 0.01 {
		t.Errorf("Expected G≈0.502, got %f", color.G)
	}
	if color.B != 0.0 {
		t.Errorf("Expected B=0.0, got %f", color.B)
	}
}

func TestColorWithFormatting(t *testing.T) {
	parser := NewMarkdownParser("{red}red text{/red} with **bold**")
	elements := parser.Parse()

	if len(elements) != 1 {
		t.Fatalf("Expected 1 element, got %d", len(elements))
	}

	children := elements[0].Children
	if len(children) == 0 {
		t.Fatal("Expected children elements")
	}

	// Should have color and bold elements
	hasColor := false
	hasBold := false
	for _, child := range children {
		if child.Type == "color" && child.Color != nil {
			hasColor = true
		}
		if child.Type == "bold" {
			hasBold = true
		}
	}

	if !hasColor {
		t.Error("Expected to find colored text")
	}
	if !hasBold {
		t.Error("Expected to find bold text")
	}
}
