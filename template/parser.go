package template

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
// {N} indicates width
// If no width is given, that field expands to fill space
//   if multiple fields have no width value, they will divide up
//   the free space
// @@ is used to print a literal @ (cannot appear in field key)
// Two literals cannot be adjacent or have their own width value
// All spaces are included in literals.
// Extra spaces between adjacent fields are ignored; use width instead

func Parse(s string) (Template, error) {
	if len(s) == 0 {
		return make(Template, 0), nil
	}

	firstNode, remainder, err := parseNode(s)
	if err != nil {
		return nil, err
	}
	nodes := Template{*firstNode}

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
	updateBoolOnce := func(name string, boolPtr *bool) error {
		if *boolPtr {
			return fmt.Errorf("Cannot specify %s more than once", name)
		}
		*boolPtr = true
		return nil
	}

	for _, format := range formats {
		switch {
		case format == "bold":
			err := updateBoolOnce("bold", &formatting.Bold)
			if err != nil {
				return err
			}
		case format == "italic":
			err := updateBoolOnce("italic", &formatting.Italic)
			if err != nil {
				return err
			}
		case format == "underline":
			err := updateBoolOnce("underline", &formatting.Underline)
			if err != nil {
				return err
			}
		case format == "wrap":
			err := updateBoolOnce("wrap", &formatting.Wrap)
			if err != nil {
				return err
			}
		case isAlignmentName(format):
			if formatting.Alignment != AlignmentDefault {
				return fmt.Errorf("Cannot specify more than one alignment in %s", s)
			}
			formatting.Alignment = AlignmentsByName[format]
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
	specifier := ""     // the text of the field not including any width term
	remainder := "\x00" // null to indicate no special remainder
	var i int
	node := Node{IsField: true}

scan:
	for i = 1; i < len(s); i++ { // skip leading @
		r := s[i]
		switch r {
		case '{':
			if i >= len(s)-2 {
				return nil, "", fmt.Errorf("Unexpected end of string; expected width followed by }")
			}
			// Chop off curly and parse
			remainder = s[i+1:]
			var width int
			var err error
			width, remainder, err = parseFieldWidth(remainder)
			if err != nil {
				return nil, "", err
			}
			node.Width = width
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

func parseFieldWidth(s string) (int, string, error) {
	result := ""
	for i, r := range s {
		// Handle close brace case
		if r == '}' {
			if result == "" {
				return 0, "", fmt.Errorf("Unexpected empty width. Expected digits.")
			}
			width, _ := strconv.Atoi(result)
			return width, s[i+1:], nil
		}

		if unicode.IsDigit(r) {
			result += string(r)
		} else {
			return 0, "", fmt.Errorf("Unexpected character '%s'; expected digit", string(r))
		}
	}
	return 0, "", fmt.Errorf("Unexpected end of input in width string")
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
