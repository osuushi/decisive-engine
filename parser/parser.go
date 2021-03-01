package parser

import (
	"encoding/hex"
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/image/colornames"
)

// Parse a row format.
// Example: "Title: @title{15} Description: @desc"
// @ prefixed implies an interpolated field.
// {} indicates padding width
// If no padding width is given, that field expands to fill space
//   if multiple fields have no padding value, they will divide up
//   the free space
// @@ is used to print a literal @ (cannot appear in field key)
// Two literals cannot be adjacent or have their own padding value
// All spaces are included in literals.
// Extra spaces between adjacent fields are ignored; use padding instead

type Node struct {
	// If false, this is a literal string
	IsField bool
	// For fields, the value is the key (no @ included)
	// For literals, it is the literal value (with escaped @s unescaped)
	Value string
	// The amount of padding (for fields only)
	// zero means fill
	Padding int
	// Formatting; nil for literals
	FieldFormatting *FieldFormatting
}

type FieldFormatting struct {
	// Color for field node; nil means default
	Color *color.RGBA
	// Styling
	Bold, Italic, Underline bool
}

func Parse(s string) ([]Node, error) {
	if len(s) == 0 {
		return make([]Node, 0), nil
	}

	firstNode, remainder, err := parseNode(s)
	if err != nil {
		return nil, err
	}
	nodes := []Node{*firstNode}

	remainingNodes, err := Parse(remainder)
	if err != nil {
		return nil, err
	}

	return append(nodes, remainingNodes...), nil
}

func parseNode(s string) (*Node, string, error) {
	if s[0] == '@' && len(s) > 1 && s[1] != '@' {
		return parseFieldNode(s)
	} else {
		return parseLiteralNode(s)
	}
}

// Try to decode a hex color, or return nil if it's not valid
func decodeHexColor(s string) *color.RGBA {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return nil
	}
	if len(bytes) != 3 {
		return nil
	}
	return &color.RGBA{bytes[0], bytes[1], bytes[2], 0xFF}
}

func (self *Node) parseFieldSpecifier(s string) error {
	parts := strings.Split(s, ".")
	value, formats := parts[0], parts[1:]
	formatting := FieldFormatting{}

	for _, format := range formats {
		switch format {
		case "bold":
			formatting.Bold = true
		case "italic":
			formatting.Italic = true
		case "underline":
			formatting.Underline = true
		default:
			alreadyHasColor := formatting.Color != nil
			addedNewColor := false
			// See if this is a crayon name
			color, ok := colornames.Map[format]
			if ok {
				formatting.Color = &color
				addedNewColor = true
			} else if colorPtr := decodeHexColor(format); colorPtr != nil {
				formatting.Color = colorPtr
				addedNewColor = true
			} else {
				return fmt.Errorf("Invalid field format: %s in %s", format, s)
			}
			if alreadyHasColor && addedNewColor {
				return fmt.Errorf("Cannot add second color %s in %s", format, s)
			}
		}
	}

	self.Value = value
	self.FieldFormatting = &formatting
	return nil
}

func parseFieldNode(s string) (*Node, string, error) {
	specifier := ""     // the text of the field not including any padding term
	remainder := "\x00" // null to indicate no special remainder
	var i int
	node := Node{IsField: true}

scan:
	for i = 1; i < len(s); i++ { // skip leading @
		r := s[i]
		switch r {
		case '{':
			if i >= len(s)-2 {
				return nil, "", fmt.Errorf("Unexpected end of string; expected padding followed by }")
			}
			// Chop off curly and parse
			remainder = s[i+1:]
			var padding int
			var err error
			padding, remainder, err = parseFieldPadding(remainder)
			if err != nil {
				return nil, "", err
			}
			node.Padding = padding
			break scan
		case ' ', '@':
			break scan
		default:
			// Not a special character. Note that the character may still be invalid,
			// in which case we will find that in *Node.parseFieldSpecifier before
			// returning the node.
			specifier += string(r)
		}
	}

	err := node.parseFieldSpecifier(specifier)
	if err != nil {
		return nil, "", err
	}

	if remainder == "\x00" {
		// No special remainder; slice the input string from i
		remainder = s[i:]
	}
	return &node, remainder, nil
}

func parseFieldPadding(s string) (int, string, error) {
	result := ""
	for i, r := range s {
		// Handle close brace case
		if r == '}' {
			if result == "" {
				return 0, "", fmt.Errorf("Unexpected empty padding. Expected digits.")
			}
			padding, _ := strconv.Atoi(result)
			return padding, s[i+1:], nil
		}

		if unicode.IsDigit(r) {
			result += string(r)
		} else {
			return 0, "", fmt.Errorf("Unexpected character '%s'; expected digit", string(r))
		}
	}
	return 0, "", fmt.Errorf("Unexpected end of input in padding string")
}

func parseLiteralNode(s string) (*Node, string, error) {
	result := ""
	var i int
	node := Node{IsField: false}
	for i = 0; i < len(s); i++ {
		r := s[i]
		switch r {
		case '@':
			// Look ahead to see if this is an escape
			if i == len(s)-1 {
				return nil, "", fmt.Errorf("Unexpected end of input in literal")
			} else if s[i+1] == '@' {
				// Escaped atsign
				result += string(r)
				// Skip next character
				i++
			} else {
				// Not an escape; field started; stop
				node.Value = result
				return &node, s[i:], nil
			}
		default:
			result += string(r)
		}
	}
	// Making it here means we're at the end of the input
	node.Value = result
	return &node, "", nil
}
